package sdk_test

import (
	"encoding/json"
	"strings"
	"testing"

	TU "github.com/PathDNA/testutils"
	"github.com/missionMeteora/sdk"
)

func TestListCampaigns(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)

	cmps, err := c.ListCampaigns(ctx, defaultUID)
	TU.FatalIf(t, err)

	for _, cmp := range cmps {
		t.Logf("Campaign: %+v", cmp)
		t.Log(c.GetCampaignReport(ctx, defaultUID, cmp.ID, sdk.DateToTime("20180101"), sdk.DateToTime("20180102")))
	}

	cmps, err = c.ListDraftCampaigns(ctx, defaultUID)
	TU.FatalIf(t, err)

	for _, cmp := range cmps {
		t.Logf("Draft: %+v", cmp)
	}

}

func TestCampaigns(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)

	c2, err := c.AsUser(ctx, "2") // meteora agency
	TU.FatalIf(t, err)

	var (
		agID, segID, psegID, cmpID, fullCmpID string

		ad *sdk.Ad
	)

	defer func() {
		if fullCmpID != "" {
			TU.FailIf(t, c2.DeleteCampaign(ctx, fullCmpID))
		}
		if cmpID != "" {
			TU.FailIf(t, c2.DeleteDraftCampaign(ctx, cmpID))
		}

		if psegID != "" {
			TU.FailIf(t, c2.DeleteProximitySegment(ctx, psegID))
		}
		if segID != "" {
			TU.FailIf(t, c2.DeleteSegment(ctx, segID))
		}
		if ad != nil {
			TU.FailIf(t, c2.DeleteAd(ctx, ad.ID))
		}
		if agID != "" {
			TU.FailIf(t, c2.DeleteAdGroup(ctx, agID))
		}
	}()

	agID, err = c2.CreateAdGroup(ctx, defaultUID, "SDK Test Campaign")
	TU.FatalIf(t, err)

	ad, err = c2.CreateAd(ctx, defaultUID, &sdk.CreateAdRequest{
		Name: "sdkTestCampaign.png", GroupID: agID, Width: 300, Height: 250,
		LandingURL: "https://test.com", AdImage: dummyPNG(300, 250),
	})
	TU.FatalIf(t, err)

	segID, err = c2.CreateSegment(ctx, defaultUID, &sdk.Segment{Name: "SDK Test Campaign"})
	TU.FatalIf(t, err)

	psegID, err = c2.CreateProximitySegment(ctx, defaultUID, dummyProxSeg("SDK Test Campaign (Prox)"))
	TU.FatalIf(t, err)

	cmpID, err = c2.CreateDraftCampaign(ctx, defaultUID, &sdk.Campaign{
		Name:     "SDK Test Campaign",
		Segments: []string{segID, psegID},
		Adgroups: []string{agID},
		Budget:   10,
	})
	TU.FatalIf(t, err)

	fullCmpID, err = c2.UpgradeCampaign(ctx, defaultUID, cmpID)
	if err != nil && strings.Contains(err.Error(), "Your card number is incorrect") { // expected
		err = nil
	}
	TU.FatalIf(t, err)
}

func TestCreateFullCampaign(t *testing.T) {
	c := sdk.NewWithAddr(localAPI, adminKey)

	req := &sdk.CreateFullCampaignRequest{
		Campaign: &sdk.Campaign{
			Name: "SDK Test Full Campaign",
			Apps: map[string]*json.RawMessage{
				"advancedBidding": sdk.RawMarshal(sdk.AppAdvBidding{
					Status: true, BaseCPM: 2, MaxCPM: 5,
				}),
				"searchRetargeting": sdk.RawMarshal(sdk.AppSearchRetargeting{
					Status: true,
					List:   []string{"nike shoes", "adidas", "shiny shoes"},
				}),
			},
			Budget: 50,
		},

		Ads: []*sdk.CreateAdRequest{
			{
				Name: "sdkTestCampaign-1.png", Width: 300, Height: 250,
				LandingURL: "https://test.com", AdImage: dummyPNG(300, 250),
			},
			{
				Name: "sdkTestCampaign-2.png", Width: 300, Height: 250,
				LandingURL: "https://test.com", AdImage: dummyPNG(300, 250),
			},
		},

		Segments: []*sdk.Segment{
			{Name: "Full Segment"},
		},

		ProximitySegment: []*sdk.ProximitySegment{
			dummyProxSeg("Full Proximity Segment"),
		},

		IsDraft: true,
	}

	t.Logf("%s", *sdk.RawMarshal(req))

	c2, err := c.AsUser(ctx, "3") // meteora user
	TU.FatalIf(t, err)

	cmp, err := c2.CreateFullCampaign(ctx, defaultUID, req)
	TU.FatalIf(t, err)

	defer req.Rollback(c2)

	t.Logf("%+v", cmp)
}
