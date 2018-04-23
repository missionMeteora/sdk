# sdk [![GoDoc](https://godoc.org/github.com/missionMeteora/sdk?status.svg)](https://godoc.org/github.com/missionMeteora/sdk)

This provides a full SDK to access Meteora APIs from Go.

## Install

	go get -u github.com/missionMeteora/sdk

## Example

```go
req := &sdk.CreateFullCampaignRequest{
	// Required
	Campaign: &sdk.Campaign{
		Name:   "SDK Test Full Campaign",
		Budget: 50,
	},

	// Required
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

	// Optional
	CampaignApps: []sdk.App{
		&sdk.AppAdvBidding{BaseCPM: 2, MaxCPM: 5},
	},

	// Optional
	Segments: []*sdk.Segment{
		{Name: "Full Segment"},
	},

	// Optional
	ProximitySegment: []*sdk.ProximitySegment{
		dummyProxSeg("Full Proximity Segment"),
	},

	IsDraft: true,
}

cmp, err := c2.CreateFullCampaign(ctx, defaultUID, req)
// checkErr(err)
log.Printf("CampaignID: %s, Full Campaign: %+v", cmp.ID, cmp)
```
