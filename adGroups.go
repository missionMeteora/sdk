package sdk

import (
	"context"
)

// AdGroup represents an ad group
type AdGroup struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	AdCount int    `json:"adCount"`
}

// CreateAdGroup will create an adgroup by the owning user ID
func (c *Client) CreateAdGroup(ctx context.Context, uid, name string) (id string, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var (
		req AdGroup

		resp idOrDataResp
	)

	// Set request name to match provided name
	req.Name = name
	// Attempt to post request
	if err = c.rawPost(ctx, "adGroups/"+uid, req, &resp); err != nil {
		return
	}

	// Set ID as the response ID
	id = resp.String()
	return
}

// ListAdGroups will list all the adgroups for a given user id
func (c *Client) ListAdGroups(ctx context.Context, uid string) (ags map[string]*AdGroup, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var resp []struct {
		ID   string   `json:"id"`
		Data *AdGroup `json:"data"`
	}

	if err = c.rawGet(ctx, "adGroupsList/"+uid, &resp); err != nil || len(resp) == 0 {
		return
	}

	ags = make(map[string]*AdGroup, len(resp))

	for _, ag := range resp {
		ag.Data.ID = ag.ID
		ags[ag.ID] = ag.Data
	}

	return
}

// DeleteAdGroup will delete an adgroup by it's ID
func (c *Client) DeleteAdGroup(ctx context.Context, adgroupID string) (err error) {
	return c.rawDelete(ctx, "adGroups/"+adgroupID, nil)
}
