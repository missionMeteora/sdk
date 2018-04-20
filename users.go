package sdk

import "context"

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
