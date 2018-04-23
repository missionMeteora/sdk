package sdk

import (
	"context"
	"regexp"
)

// Agency represents an agency user
type Agency struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Advertiser represents an advertiser user
type Advertiser struct {
	ID           string `json:"id"`
	AgencyID     string `json:"agencyID"`
	Name         string `json:"name"`
	NumCampaigns int    `json:"numCmps"`
	Status       bool   `json:"status"`
}

// AsUser fetches the given uid's api key and returns a new client using it or an error.
func (c *Client) AsUser(ctx context.Context, uid string) (nc *Client, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var key string
	if key, err = c.GetUserAPIKey(ctx, uid); err != nil {
		return
	}

	nc = newClient(c.u, key)

	return
}

// ListAgencies will list the agencies
func (c *Client) ListAgencies(ctx context.Context) (out map[string]*Agency, err error) {
	var ags []*Agency
	if err = c.rawGet(ctx, "agenciesList", &ags); err != nil {
		return
	}

	out = make(map[string]*Agency, len(ags))
	for _, ag := range ags {
		out[ag.ID] = ag
	}

	return
}

// ListAdvertisers will list the advertisers for a given agency
func (c *Client) ListAdvertisers(ctx context.Context, agencyID string) (out map[string]*Advertiser, err error) {
	var advs []*Advertiser
	if err = c.rawGet(ctx, "advertiserList/"+agencyID, &advs); err != nil {
		return
	}

	out = make(map[string]*Advertiser, len(advs))
	for _, adv := range advs {
		adv.AgencyID = agencyID
		out[adv.ID] = adv
	}

	return
}

type CreateAdvertiserRequest struct {
	Name          string `json:"company"`       // required
	Email         string `json:"email"`         // required
	AgencyID      string `json:"agency"`        // Required
	AdvertiserFee int    `json:"advertiserFee"` // required

	Password        string `json:"pass"`        // optional, needed only if you want to login from the meteora dash
	PasswordConfirm string `json:"passConfirm"` // not needed, filled automatically if password is set
}

var emailRE = regexp.MustCompile(`.+@.+\.\w+`)

func (c *Client) CreateAdvertiser(ctx context.Context, req *CreateAdvertiserRequest) (uid string, err error) {
	if req == nil {
		err = ErrRequestIsNil
		return
	}
	if req.Name == "" {
		err = ErrInvalidName
		return
	}

	if !emailRE.MatchString(req.Email) {
		err = ErrInvalidEmail
		return
	}

	if req.AgencyID == "" {
		err = ErrInvalidAgencyID
		return
	}

	if req.Password != "" {
		if len(req.Password) < 8 {
			err = ErrInvalidPassword
			return
		}

		req.PasswordConfirm = req.Password
	}

	var resp idOrDataResp

	if err = c.rawPost(ctx, "signUp/advertiser/"+req.AgencyID, req, &resp); err != nil {
		return
	}

	uid = resp.String()

	return
}
