package sdk_test

import (
	"bytes"
	"context"
	"flag"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/missionMeteora/ssync"

	TU "github.com/PathDNA/testutils"
	"github.com/missionMeteora/sdk"
)

const (
	localAPI   = "http://localhost:8080/"
	defaultUID = "3"
)

var adminKey string

func init() {
	flag.StringVar(&adminKey, "key", "382787ed-6c36-4ca0-a3f1-6bcf926fea7c", "local api admin key")
	flag.Parse()

	log.SetFlags(log.Lshortfile)
}

var ctx = context.Background()

func TestClient(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)

	ver, err := c.GetAPIVersion(ctx)
	TU.FatalIf(t, err)

	id, err := c.GetUserID(ctx)
	TU.FatalIf(t, err)

	t.Logf("api version: %v, key's userID: %v", ver, id)

	if err = c.RawRequest("GET", "does-not-exist", nil, nil); err == nil || !strings.HasPrefix(err.Error(), "error(404") {
		t.Fatalf("expected error 404, got : %v", err)
	}
}

func TestAdGroups(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)

	c2, err := c.AsUser(ctx, "2") // meteora agency
	TU.FatalIf(t, err)

	name := "sdk test:" + time.Now().String()

	nag, err := c2.CreateAdGroup(ctx, defaultUID, name)
	TU.FatalIf(t, err)

	ags, err := c2.ListAdGroups(ctx, defaultUID)
	TU.FatalIf(t, err)

	found := false
	for _, ag := range ags {
		if found = ag.Name == name && ag.ID == nag; found {
			break
		}
	}

	if !found {
		t.Errorf("couldn't find group (%s: %s): %+v", nag, name, ags)
	}

	TU.FatalIf(t, c2.DeleteAdGroup(ctx, nag))
}

func TestAds(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)

	c2, err := c.AsUser(ctx, defaultUID)
	TU.FatalIf(t, err)

	ad, err := c2.CreateAd(ctx, defaultUID, &sdk.CreateAdRequest{
		Name: "ad.png", GroupID: "1", Width: 300, Height: 250,
		LandingURL: "https://test.com", AdImage: dummyPNG(300, 250),
	})
	TU.FatalIf(t, err)

	t.Logf("New AD: %+v", ad)

	ad.Active = false
	TU.FatalIf(t, c2.UpdateAd(ctx, ad))

	ads, err := c2.ListAds(ctx, defaultUID)
	TU.FatalIf(t, err)

	if ad := ads[ad.ID]; ad == nil || ad.Active {
		t.Fatalf("expected to find the ad in the ads list, but we didn't :(\n%+v", ad)
	}

	as, err := c.GetAdsReport(ctx, defaultUID, sdk.DateToTime("2017-01-02"), sdk.DateToTime("2018-02-22"))
	t.Log(as, err)
	TU.FatalIf(t, c2.DeleteAd(ctx, ad.ID))
	TU.FatalIf(t, err)
}

func dummyPNG(w, h int) io.Reader {
	var (
		img = image.NewRGBA(image.Rect(0, 0, w, h))
		buf bytes.Buffer
	)

	png.Encode(&buf, img)

	return &buf
}

func TestSegments(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)

	c2, err := c.AsUser(ctx, defaultUID)
	TU.FatalIf(t, err)

	seg := &sdk.Segment{
		Name: "sdk test",
	}

	id, err := c2.CreateSegment(ctx, defaultUID, seg)
	TU.FatalIf(t, err)

	t.Logf("New SegmentID: %s", id)

	seg.Name = "RENAMED"
	TU.FatalIf(t, c2.UpdateSegment(ctx, seg))

	segs, err := c2.ListSegments(ctx, defaultUID)
	TU.FatalIf(t, err)

	if seg, ok := segs[id]; !ok || seg.Name != "RENAMED" {
		t.Fatalf("expected to find the new segment, but we didn't :( %s", id)
	}

	TU.FatalIf(t, c2.DeleteSegment(ctx, id))
}

func TestProxSegments(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)

	c2, err := c.AsUser(ctx, defaultUID)
	TU.FatalIf(t, err)

	seg := dummyProxSeg("sdk test (proximity)")

	id, err := c2.CreateProximitySegment(ctx, defaultUID, seg)
	seg.ID = id

	TU.FatalIf(t, err)

	t.Logf("New ProximitySegmentID: %s", id)

	seg.Name = "RENAMED"
	TU.FatalIf(t, c2.UpdateProximitySegment(ctx, seg))

	segs, err := c2.ListProximitySegments(ctx, defaultUID)
	TU.FatalIf(t, err)

	if seg, ok := segs[id]; !ok || seg.Name != "RENAMED" {
		t.Fatalf("expected to find the new segment, but we didn't :( %s", id)
	}

	TU.FatalIf(t, c2.DeleteProximitySegment(ctx, id))
}

func TestLists(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)

	ags, err := c.ListAgencies(ctx)
	TU.FatalIf(t, err)

	if ags["2"] == nil {
		t.Fatalf("couldn't find meteora agency (id: 2): %+v", ags)
	}

	advs, err := c.ListAdvertisers(ctx, "2")
	TU.FatalIf(t, err)

	if advs["3"] == nil {
		t.Fatalf("couldn't find meteora advertiser (id: 3): %+v", advs)
	}
}

func TestCreateAdvertiser(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)
	ts := time.Now().Format("20060102150405")

	req := &sdk.CreateAdvertiserRequest{
		Name:     "Test Adv (" + ts + ")",
		AgencyID: "2",
		Email:    ts + "@test.org",
	}

	t.Log(req)

	uid, err := c.CreateAdvertiser(ctx, req)
	TU.FatalIf(t, err)

	t.Logf("new advertiser id: %s", uid)
}

func TestHeatmaps(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)

	hm, err := c.GetHeatmap(ctx, defaultUID)
	TU.FatalIf(t, err)
	if len(hm.AdsList) == 0 {
		t.Fatalf("invalid response: %+v", hm)
	}
}

func TestSsync(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)
	sc, err := ssync.NewClient("", "bank", os.Getenv("SSYNC_ADDR"))
	TU.FatalIf(t, err)
	d := time.Date(2018, 9, 27, 0, 0, 0, 0, time.UTC)
	_, err = c.Receipts(ctx, sc, d, "2720", "12591")
	TU.FatalIf(t, err)
}

func dummyProxSeg(name string) *sdk.ProximitySegment {
	return &sdk.ProximitySegment{
		Name: name,
		Locations: []*sdk.Location{
			{
				Label:  "Starbucks",
				Center: sdk.Coords{Latitude: 32.8826822, Longitude: -97.39539739999998},
			},
		},
	}
}
