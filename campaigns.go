package sdk

import (
	"context"
	"encoding/json"
	"strings"
)

const (
	CampaignMinimumBudget = 1
	CampaignMaximumBudget = 300

	CampaignMaximumImpBudget = 100000
)

// Campaign is a QoL alias for campaigns.Campaign
type Campaign = struct {
	ID      string `json:"id"`
	OwnerID string `json:"ownerID"`

	Active   bool `json:"active"`
	Archived bool `json:"archived,omitempty"`

	Name      string  `json:"name"`
	Notes     string  `json:"notes,omitempty"`
	Budget    float64 `json:"budget"`
	ImpBudget uint32  `json:"impBudget"`

	OptBucket uint8 `json:"optBucket"`

	Created   int64 `json:"created"`
	Scheduled bool  `json:"scheduled"`
	Start     int64 `json:"start"`
	End       int64 `json:"end"`

	Segments []string `json:"segments,omitempty"`
	Adgroups []string `json:"adGroups,omitempty"`

	Apps map[string]*json.RawMessage `json:"apps,omitempty"`

	Searches []string `json:"searches,omitempty"`
}

// CreateCampaign will create a campaign for a given user ID
func (c *Client) CreateCampaign(ctx context.Context, uid string, cmp *Campaign) (string, error) {
	return c.createCampaign(ctx, uid, cmp, false)
}

// CreateDraftCampaign will create a draft campaign for a given user ID
func (c *Client) CreateDraftCampaign(ctx context.Context, uid string, cmp *Campaign) (string, error) {
	return c.createCampaign(ctx, uid, cmp, true)
}

func (c *Client) createCampaign(ctx context.Context, uid string, cmp *Campaign, drafts bool) (cid string, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var (
		resp idOrDataResp

		ep = "campaigns/byAdv/"
	)

	if drafts {
		ep = "campaignDrafts/"
	}

	for _, app := range DefaultApps {
		if _, ok := cmp.Apps[app.Name()]; !ok {
			SetApp(cmp, app.Name(), app)
		}
	}

	if err = c.rawPost(ctx, ep+uid, cmp, &resp); err != nil {
		return
	}

	cid = resp.String()
	return
}

// GetCampaign will get a campaign by campaign id
func (c *Client) GetCampaign(ctx context.Context, cid string) (*Campaign, error) {
	return c.getCampaign(ctx, cid, false)
}

// GetDraftCampaign will get a campaign by campaign id
func (c *Client) GetDraftCampaign(ctx context.Context, cid string) (*Campaign, error) {
	return c.getCampaign(ctx, cid, true)
}

func (c *Client) getCampaign(ctx context.Context, cid string, drafts bool) (cmp *Campaign, err error) {
	if cid == "" {
		err = ErrMissingCID
		return
	}

	var (
		resp struct {
			Data *Campaign `json:"data"`
		}
		ep = "campaigns/byCID/"
	)

	if drafts {
		ep = "campaignDraft/"
	}

	if err = c.rawGet(ctx, ep+cid, &resp); err != nil {
		return
	}

	cmp = resp.Data
	return
}

// UpdateCampaign will update a campaign
func (c *Client) UpdateCampaign(ctx context.Context, cmp *Campaign) error {
	return c.updateCampaign(ctx, cmp, false)
}

// UpdateDraftCampaign will update a draft campaign
func (c *Client) UpdateDraftCampaign(ctx context.Context, cmp *Campaign) error {
	return c.updateCampaign(ctx, cmp, true)
}

// updateCampaign will update a campaign
func (c *Client) updateCampaign(ctx context.Context, cmp *Campaign, drafts bool) (err error) {
	var (
		ep = "campaigns/byCID/"
	)

	if drafts {
		ep = "campaignDraft/"
	}

	return c.rawPut(ctx, ep+cmp.ID, cmp, nil)
}

// DeleteCampaign will delete a campaign by it's ID
func (c *Client) DeleteCampaign(ctx context.Context, cid string) error {
	return c.deleteCampaign(ctx, cid, false)
}

// DeleteDraftCampaign will delete a draft campaign by it's ID
func (c *Client) DeleteDraftCampaign(ctx context.Context, cid string) error {
	return c.deleteCampaign(ctx, cid, true)
}

// deleteCampaign will delete a campaign by it's ID
func (c *Client) deleteCampaign(ctx context.Context, cid string, drafts bool) (err error) {
	if cid == "" {
		err = ErrMissingCID
		return
	}

	ep := "campaigns/byCID/"

	if drafts {
		ep = "campaignDraft/"
	}

	return c.rawDelete(ctx, ep+cid, nil)
}

// ListCampaigns will list all the campaigns for a given user id
func (c *Client) ListCampaigns(ctx context.Context, uid string) (map[string]*Campaign, error) {
	return c.listCampaigns(ctx, uid, false)
}

// ListDraftCampaigns will list all the draft campaigns for a given user id
func (c *Client) ListDraftCampaigns(ctx context.Context, uid string) (map[string]*Campaign, error) {
	return c.listCampaigns(ctx, uid, true)
}

func (c *Client) listCampaigns(ctx context.Context, uid string, drafts bool) (cmps map[string]*Campaign, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var (
		resp []*Campaign
		ep   = "campaignsList/"
	)

	if drafts {
		ep = "campaignDrafts/"
	}

	if err = c.rawGet(ctx, ep+uid, &resp); err != nil || len(resp) == 0 {
		return
	}

	cmps = make(map[string]*Campaign, len(resp))
	for _, cmp := range resp {
		cmp.OwnerID = uid
		cmps[cmp.ID] = cmp
	}

	return
}

type CreateFullCampaignRequest struct {
	Campaign *Campaign          `json:"campaign,omitempty"` // required
	Ads      []*CreateAdRequest `json:"ads,omitempty"`      // required

	CampaignApps []App `json:"campaignApps,omitempty"`

	Segments         []*Segment          `json:"segments,omitempty"`         // optional
	ProximitySegment []*ProximitySegment `json:"proximitySegment,omitempty"` // optional

	IsDraft bool `json:"isDraft"`

	// this is only used for rollback
	adIDs []string
}

func (req *CreateFullCampaignRequest) validate() (err error) {
	if req == nil {
		return ErrRequestIsNil
	}

	if err = validateCampaign(req.Campaign); err != nil {
		return
	}

	if len(req.Ads) == 0 {
		return ErrMissingAds
	}

	for _, ad := range req.Ads {
		if err = ad.validate(); err != nil {
			return
		}
	}

	for _, seg := range req.Segments {
		if seg.Name == "" {
			return ErrMissingSegName
		}
	}

	for _, seg := range req.ProximitySegment {
		if seg.Name == "" {
			return ErrMissingProxSegName
		}

		if len(seg.Locations) == 0 {
			return ErrMissingLocations
		}
	}

	return
}

// deletes any partially created resources
func (req *CreateFullCampaignRequest) Rollback(c *Client) {
	cmp := req.Campaign
	ctx := context.Background()

	if cmp.ID != "" {
		c.deleteCampaign(ctx, cmp.ID, req.IsDraft)
	}

	for _, segID := range cmp.Segments {
		if strings.HasPrefix(segID, "px_") {
			c.DeleteProximitySegment(ctx, segID)
		} else {
			c.DeleteSegment(ctx, segID)
		}
	}

	for _, adID := range req.adIDs {
		c.DeleteAd(ctx, adID)
	}

	for _, agID := range cmp.Adgroups {
		c.DeleteAdGroup(ctx, agID)
	}
}

// CreateFullCampaign takes full control of the passed request, reusing it can cause races and/or crashes.
// A new ADGroup will be created for this request and filled in with the ads, and appended to cmp.Adgroups.
func (c *Client) CreateFullCampaign(ctx context.Context, uid string, req *CreateFullCampaignRequest) (cmp *Campaign, err error) {
	if uid == "" {
		return nil, err
	}

	if err = req.validate(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			cmp = nil
			// should this be ran in a goroutine?
			req.Rollback(c)
		}
	}()

	cmp = req.Campaign

	cmp.OwnerID = uid

	var agID string
	if agID, err = c.CreateAdGroup(ctx, uid, cmp.Name); err != nil {
		return
	}

	cmp.Adgroups = append(cmp.Adgroups, agID)

	var newAd *Ad
	for _, ad := range req.Ads {
		ad.GroupID = agID
		if newAd, err = c.CreateAd(ctx, uid, ad); err != nil {
			return
		}

		req.adIDs = append(req.adIDs, newAd.ID)
	}

	var segID string
	for _, seg := range req.Segments {
		if segID, err = c.CreateSegment(ctx, uid, seg); err != nil {
			return
		}

		cmp.Segments = append(cmp.Segments, segID)
	}

	for _, seg := range req.ProximitySegment {
		if segID, err = c.CreateProximitySegment(ctx, uid, seg); err != nil {
			return
		}

		cmp.Segments = append(cmp.Segments, segID)
	}

	for _, app := range req.CampaignApps {
		SetApp(cmp, app.Name(), app)
	}

	if cmp.ID, err = c.createCampaign(ctx, uid, cmp, req.IsDraft); err != nil {
		return
	}

	return cmp, nil
}

// can't specify methods on aliased types
func validateCampaign(c *Campaign) (err error) {
	if len(c.Name) == 0 {
		return ErrInvalidName
	}

	if c.Budget < CampaignMinimumBudget && c.ImpBudget == 0 {
		return ErrInvalidBudget
	}

	if c.ImpBudget > CampaignMaximumBudget {
		return ErrInvalidImpBudget
	}

	if c.Scheduled && c.Start >= c.End {
		return ErrInvalidSchedule
	}

	return nil
}
