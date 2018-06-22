package sdk

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PathDNA/ptk"
)

var quoteBytes = []byte{'"'}

// Ad is mostly a trimmed down version of api/server/types.Ad
type Ad struct {
	ID           string `json:"id,omitempty"`
	GroupID      string `json:"group,omitempty"`
	AdvertiserID string `json:"advertiserID,omitempty"`

	Active        bool   `json:"active,omitempty"`
	Clicks        int    `json:"clicks,omitempty"`
	Conv          int    `json:"conv,omitempty"`
	Width         int    `json:"width,omitempty"`
	Height        int    `json:"height,omitempty"`
	Imps          int    `json:"imps,omitempty"`
	AdType        string `json:"adType,omitempty"`
	ImageKey      string `json:"imageKey,omitempty"`
	ImageLocation string `json:"imageLocation,omitempty"`
	Name          string `json:"name,omitempty"`
	Size          string `json:"size,omitempty"`

	LandingURL string `json:"landingURL,omitempty"`

	ImpTracker   string `json:"impTrack,omitempty"`
	ClickTracker string `json:"clickTrack,omitempty"`
}

// CreateAd will create an Advertisement for a given user ID
func (c *Client) CreateAd(ctx context.Context, uid string, req *CreateAdRequest) (ad *Ad, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	if err = req.validate(); err != nil {
		return
	}

	mt := mime.TypeByExtension(filepath.Ext(req.Name))
	switch mt {
	case "image/jpeg", "image/gif", "image/png":
	default:
		err = fmt.Errorf("invalid image type from filename (%s): %s", req.Name, mt)
		return
	}

	pipeObj := ptk.M{
		"name":   req.Name,
		"width":  req.Width,
		"height": req.Height,
	}

	// this is needed to support the sdk-server,
	// otherwise there is no other way to create an ad.
	switch img := req.AdImage.(type) {
	case io.Reader:
		if !req.Encoded {
			pipeObj["data"] = ptk.Base64ToJSON(nil, mt, img)
		} else {
			pipeObj["data"] = io.MultiReader(bytes.NewReader(quoteBytes), img, bytes.NewReader(quoteBytes))
		}
	case string:
		if !strings.HasPrefix(img, "data:"+mt+";base64,") {
			err = fmt.Errorf("unexpected base64 data, the image data must start with: `data:%s;base64,`", mt)
			return
		}
		pipeObj["data"] = io.MultiReader(bytes.NewReader(quoteBytes), strings.NewReader(img), bytes.NewReader(quoteBytes))
	}

	imageReq := ptk.PipeJSONObject(pipeObj)

	var resp struct {
		ID       string `json:"id"`
		Location string `json:"location"`
	}

	if err = c.rawPost(ctx, "images/"+uid, imageReq, &resp); err != nil {
		return
	}

	ad = req.toAd(resp.ID, resp.Location)

	if err = c.rawPost(ctx, "ads/"+uid+"?incGroup=true", ad, &resp); err != nil {
		ad = nil
		return
	}

	ad.ID, ad.AdvertiserID = resp.ID, uid
	return
}

// CreateAdFromFile will create an Advertisement for a given user ID using a given filename
func (c *Client) CreateAdFromFile(ctx context.Context, uid string, req *CreateAdRequest, filename string) (ad *Ad, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var f *os.File
	if f, err = os.Open(filename); err != nil {
		return
	}
	defer f.Close()

	req.AdImage = f

	return c.CreateAd(ctx, uid, req)
}

// UpdateAd will update an ad
func (c *Client) UpdateAd(ctx context.Context, ad *Ad) error {
	return c.rawPut(ctx, "ads/"+ad.ID, ad, nil)
}

// DeleteAd will delete an ad by advertisement ID
func (c *Client) DeleteAd(ctx context.Context, adID string) (err error) {
	return c.rawDelete(ctx, "ads/"+adID, nil)
}

func (c *Client) ListAds(ctx context.Context, uid string) (ads map[string]*Ad, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var resp []struct {
		ID   string `json:"id"`
		Data *Ad    `json:"data"`
	}

	if err = c.rawGet(ctx, "adsList/"+uid, &resp); err != nil {
		return
	}

	ads = make(map[string]*Ad, len(resp))

	for _, v := range resp {
		ad := v.Data
		ad.ID, ad.AdvertiserID = v.ID, uid
		ads[ad.ID] = ad
	}

	return
}

// ListAdsByAdGroup will list all ads for a user and ad group
func (c *Client) ListAdsByAdGroup(ctx context.Context, uid, adGroupID string) (map[string]*Ad, error) {
	return c.ListAdsFilter(ctx, uid, func(ad *Ad) bool {
		return ad.GroupID == adGroupID
	})
}

// ListAdsFilter will list ads which pass a given filter
func (c *Client) ListAdsFilter(ctx context.Context, uid string, filterFn func(ad *Ad) bool) (ads map[string]*Ad, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	if ads, err = c.ListAds(ctx, uid); err != nil {
		return
	}

	for key, ad := range ads {
		if filterFn(ad) {
			continue
		}

		delete(ads, key)
	}

	return
}

// CreateAdRequest is the request needed to create an Ad
type CreateAdRequest struct {
	Name       string `json:"name,omitempty"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	LandingURL string `json:"landingURL,omitempty"`

	// Everything below here is optional
	GroupID      string `json:"group,omitempty"`
	ImpTracker   string `json:"impTrack,omitempty"`
	ClickTracker string `json:"clickTrack,omitempty"`
	Encoded      bool   `json:"encoded,omitempty"`

	AdImage interface{} `json:"adImage,omitempty"`
}

func (r *CreateAdRequest) toAd(iid, iloc string) *Ad {
	return &Ad{
		Name:    r.Name,
		GroupID: r.GroupID,
		AdType:  "banner",
		Active:  true,

		Width:  r.Width,
		Height: r.Height,
		Size:   strconv.Itoa(r.Width) + "x" + strconv.Itoa(r.Height),

		ImageKey:      iid,
		ImageLocation: iloc,

		LandingURL:   r.LandingURL,
		ImpTracker:   r.ImpTracker,
		ClickTracker: r.ClickTracker,
	}
}

// Validate will validate a create ad request
func (r *CreateAdRequest) validate() error {
	if r == nil {
		return ErrRequestIsNil
	}

	if len(r.Name) == 0 {
		return ErrInvalidName
	}

	if len(r.LandingURL) == 0 {
		return ErrInvalidLanding
	}

	if !isAllowedAdSize(r.Width, r.Height) {
		return ErrInvalidAdSize
	}

	if r.AdImage == nil {
		return ErrMissingAdImage
	}
	return nil
}

func AllowedAdRectsString() string {
	out := make([]string, 0, len(allowedAdRects))
	for _, rect := range allowedAdRects {
		out = append(out, fmt.Sprintf("%dx%d", rect.Width, rect.Height))
	}
	return strings.Join(out, ",")
}

// Rects represent a unit of width and height
type Rect struct {
	Width  int
	Height int
}

var (
	allowedAdRects = []Rect{
		{300, 250},
		{300, 600},
		{300, 1050},
		{160, 600},
		{728, 90},
		{180, 150},
		{970, 90},
		{970, 250},
		{300, 50},
		{320, 50},
		{300, 100},
		{468, 60},
	}
)

func isAllowedAdSize(width, height int) bool {
	for _, r := range allowedAdRects {
		if r.Width == width && r.Height == height {
			return true
		}
	}

	return false
}
