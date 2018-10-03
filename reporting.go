package sdk

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/PathDNA/ptk"

	"github.com/PathDNA/ssync"
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
		err = ErrBadUID
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
		err = ErrBadUID
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
		err = ErrBadUID
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

// the following funcs are used internally as the services they use are private

type Receipt struct {
	CID   string `json:"id"`              // Campaign id
	AdvID string `json:"advertiserId"`    // Advertiser id
	SegID string `json:"segID,omitempty"` // Segment id
	ImpID string `json:"impid"`           // Meteora impression id
	AdID  string `json:"adId"`            // Ad id
	UUID  string `json:"uuid"`            // User id for ad viewer

	Timestamp int64   `json:"timestamp"`         // Timestamp of win
	Amount    float64 `json:"amount,omitempty"`  // Cost of the impression (true value) [Deprecated]
	Credits   int64   `json:"credits,omitempty"` // Cost of the impression (true value * 1,000,000)

	Domain    string `json:"domain,omitempty"`  // Domain the ad was served on
	Services  string `json:"usedSvc,omitempty"` // Comma separated list of services used for the impression
	Inventory string `json:"invt,omitempty"`    // Platform used to view impression (IE desktop, mobile, mobile app)

	ExID    string `json:"svc"`              // Servicer (exchange) id
	ExImpID string `json:"eid,omitempty"`    // Exchange-provided impression id
	ExInfo  string `json:"exInfo,omitempty"` // Miscellaneous exchange information

	PxID string `json:"pxID,omitempty"` // Proximity target id
}

func (c *Client) Receipts(ctx context.Context, sc *ssync.Client, date time.Time, uid, cid string) (out []byte, err error) {
	if sc == nil {
		err = ErrInvalidSSyncClient
		return
	}

	const (
		tsAndSepLen = len("1538006719@")
		bufSize     = 1 << 17 // 128kb
	)

	if uid == "" || !isNumber(uid) || !verifyUserCampaign(ctx, c, uid, cid) {
		err = ErrBadUID
	}

	if cid != "" && !isNumber(cid) {
		err = ErrCampaignIsNil
		return
	}

	var (
		files    []*ssync.FileInfo
		basePath = filepath.Join("receipts", date.Format("2006/01/02"))
		sem      = ptk.NewSem(10)
		m        = []byte(getMatchString(uid, cid))
		mux      sync.Mutex
		buf      = bytes.NewBuffer(make([]byte, 0, bufSize))
		first    = true
	)

	defer sem.Close()

	if files, err = sc.ListFiles(ctx, basePath); err != nil {
		return
	}

	process := func(i int, f *ssync.FileInfo) {
		if err := sc.StreamFile(ctx, filepath.Join(basePath, f.Path), func(rd io.Reader) error {
			gz, err := gzip.NewReader(bufio.NewReaderSize(rd, bufSize/2))
			if err != nil {
				return err
			}
			defer gz.Close()
			br := bufio.NewScanner(gz)

			for br.Scan() {
				val := br.Bytes()
				// this is faster than decoding each record from json
				if !bytes.Contains(val, m) {
					continue
				}

				mux.Lock()
				if first {
					first = false
				} else {
					buf.WriteByte(',')
				}
				buf.Write(val[tsAndSepLen:])
				mux.Unlock()
			}
			return nil
		}); err != nil {
			log.Printf("error streaming file (%s/%s): %v", basePath, f.Path, err)
		}

		sem.Done()
	}

	buf.WriteByte('[')
	for i, f := range files {
		sem.Add(1)
		i, f := i, f
		go process(i, f)
	}

	sem.Wait()
	buf.WriteByte(']')

	return buf.Bytes(), nil
}

func getMatchString(uid, cid string) string {
	if cid == "" {
		return `"advertiserId":"` + uid + `"`
	}
	return `"id":"` + cid + `","advertiserId":"` + uid + `"`
}

func (c *Client) Clicks(ctx context.Context, clicksAddr string, date time.Time, uid, cid string) (out []byte, err error) {
	if clicksAddr == "" {
		err = ErrMissingClicksServer
		return
	}

	if cid == "" {
		cid = "-1"
	}

	if !verifyUserCampaign(ctx, c, uid, cid) {
		err = ErrBadUID
		return
	}

	var (
		start, end = MidnightToMidnight(date.UTC())

		u = fmt.Sprintf(clicksAddr, uid, cid, start.Unix(), end.Unix())
	)

	err = c.c.RequestCtx(ctx, "GET", "application/json", u, nil, func(r io.Reader) error {
		out, err = ioutil.ReadAll(r)
		return err
	})

	return
}

func (c *Client) Visits(ctx context.Context, visitsAddr string, date time.Time, uid, cid string) (out []byte, err error) {
	if visitsAddr == "" {
		err = ErrMissingVisitsServer
		return
	}

	if cid == "" {
		cid = "-1"
	}

	if !verifyUserCampaign(ctx, c, uid, cid) {
		err = ErrBadUID
		return
	}

	var (
		start, end = MidnightToMidnight(date.UTC())

		u = fmt.Sprintf(visitsAddr, cid, start.Unix(), end.Unix())
	)

	err = c.c.RequestCtx(ctx, "GET", "application/json", u, nil, func(r io.Reader) error {
		out, err = ioutil.ReadAll(r)
		return err
	})

	return
}
