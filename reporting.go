package sdk

import (
	"context"
	"time"
)

var (
	AllTime = time.Unix(-1, 0)
)

// CampaignReport represents a basic campaign report
type CampaignReport struct {
	ID      string              `json:"id"`
	Imp     int64               `json:"imp,omitempty"`
	Clicks  int64               `json:"clicks,omitempty"`
	Spent   float64             `json:"spent,omitempty"`
	Period  uint32              `json:"reportingPeriod,omitempty"`
	Domains map[string]ImpClick `json:"domains,omitempty"`
	Visits  []Visit             `json:"visits,omtiempty"` // est visits please see terms
}

/// Visit is a single visit details
type Visit struct {
	CampaignID  string `json:"cid,omitempty"`
	ProximityID string `json:"pxID,omitempty"`
	Name        string `json:"pxName,omitempty"`
	StoreName   string `json:"storeName,omitempty"`
	StoreID     string `json:"storeID,omitempty"`
	TS          int64  `json:"ts,omitempty"`
}

// ImpClick has the numbers of imps and clicks for a domain
type ImpClick struct {
	Imp    int64 `json:"imp"`
	Clicks int64 `json:"clicks"`
}

// AdReport represents a basic ad report
type AdReport struct {
	ID     string  `json:"id"`
	Imp    int64   `json:"imp,omitempty"`
	Clicks int64   `json:"clicks,omitempty"`
	Spent  float64 `json:"spent,omitempty"`
}

// Click is a click action
type Click = struct {
	X int `json:"x"`
	Y int `json:"y"`
	// TTC represents time to click (in milliseconds)
	TTC int `json:"ttc,omitempty"`

	UUID string `json:"uuid"`
	AdID string `json:"adID"`
	TS   int64  `json:"ts"`
}

// GetCampaignReport will generate a report for a given campaign ID and date range
func (c *Client) GetCampaignReport(ctx context.Context, uid, cid string, start, end time.Time) (cs *CampaignReport, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	if cid == "" {
		err = ErrMissingCID
		return
	}

	var s, e string
	if s, e, err = getStartEnd(start, end); err != nil {
		return
	}

	ep := "r/campaignStats/" + uid + "/" + cid + "/" + s + "-" + e

	if err = c.rawGet(ctx, ep, &cs); err != nil {
		return
	}

	cs.ID = cid
	return
}

// GetAdsReport will generate a advertisements report for a given user ID and date range
func (c *Client) GetAdsReport(ctx context.Context, uid string, start, end time.Time) (rp map[string]*AdReport, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var s, e string
	if s, e, err = getStartEnd(start, end); err != nil {
		return
	}

	ep := "adsStats/" + uid + "/" + s + "/" + e

	if err = c.rawGet(ctx, ep, &rp); err != nil {
		return
	}

	for id, r := range rp {
		r.ID = id
	}

	return
}

// AdsListItem is a copied type from API for compatibility with the Heatmap JS lib
type AdsListItem struct {
	ID   string `json:"id"`
	Data *Ad    `json:"data"`
}

//Heatmap contains the raw Heatmap data from API
type Heatmap struct {
	AdsList []*AdsListItem `json:"adsList"`
	Clicks  []*Click       `json:"clicks"`
}

// Ads don't actually contain an ID, so the sdk version of it adds an ID to the Ad to make it easier to parse in clients.
// it.Data is the actual Ad returned from API.
func (hm *Heatmap) fillIDs() {
	if hm == nil {
		return
	}

	for _, it := range hm.AdsList {
		it.Data.ID = it.ID
	}
}

// AllAds returns a map of all the ads keyed by their IDs.
func (hm *Heatmap) AllAds() (out map[string]*Ad) {
	if hm == nil {
		return
	}

	out = map[string]*Ad{}
	for _, it := range hm.AdsList {
		out[it.ID] = it.Data
	}

	return
}

// ClicksByAds returns a map of a slice of clicks by Ad ID.
func (hm *Heatmap) ClicksByAds() (out map[string][]*Click) {
	if hm == nil {
		return
	}

	out = map[string][]*Click{}
	for _, c := range hm.Clicks {
		out[c.AdID] = append(out[c.AdID], c)
	}

	return
}

// GetHeatmap will return the heatmaps belonging to a user ID
func (c *Client) GetHeatmap(ctx context.Context, uid string) (out *Heatmap, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var resp struct {
		Data *Heatmap
	}
	if err = c.rawGet(ctx, "r/heatmap/alltime/"+uid, &resp); err != nil {
		return
	}

	out = resp.Data
	out.fillIDs()

	return
}
