package sdk

import (
	"context"
	"encoding/json"
)

// Campaign is a QoL alias for campaigns.Campaign
type Campaign = struct {
	ID      string `json:"id"`
	OwnerID string `json:"-"`

	Active   bool `json:"active"`
	Archived bool `json:"archived,omitempty"`

	Name      string  `json:"name"`
	Notes     string  `json:"notes,omitempty"`
	Budget    float64 `json:"budget"`
	ImpBudget int64   `json:"impBudget"`

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
		resp struct {
			ID string `json:"data"`
		}

		ep = "campaigns/byAdv/"
	)

	if drafts {
		ep = "campaignDrafts/"
	}

	for name, app := range DefaultApps {
		if _, ok := cmp.Apps[name]; !ok {
			SetApp(cmp, name, app)
		}
	}

	if err = c.rawPost(ctx, ep+uid, cmp, &resp); err != nil {
		return
	}

	cid = resp.ID
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
