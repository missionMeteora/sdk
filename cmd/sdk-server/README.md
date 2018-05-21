# sdk-server

**sdk-server** provides 1:1 mapping to the [Meteora SDK](https://godoc.org/github.com/missionMeteora/sdk) minus a few functions:

- RawRequest
- RawRequestCtx
- CurrentKey
- AsUser
- GetUserID
- GetUserAPIKey
- GetAPIVersion
- CreateAdFromFile
- ListAdsFilter
- GetAgencies`

## How to use it

`MethodFuncName(ctx, arg0, arg1) -> HTTP-METHOD /api/v1/funcName/arg0/arg1?key=your-agency-api-key`:

### Methods Mapping

```
GET /api/v1/userID -> GetUserID
GET /api/v1/listApps -> listApps

GET /api/v1/campaign/:cid -> GetCampaignz
GET /api/v1/draftCampaign/:cid -> GetDraftCampaign

GET /api/v1/heatmap/:uid -> GetHeatmap
GET /api/v1/adGroups/:uid -> ListAdGroups
GET /api/v1/ads/:uid -> ListAds
GET /api/v1/adsByAdGroup/:uid/:adGroupID -> ListAdsByAdGroup
GET /api/v1/advertisers/:agencyID -> ListAdvertisers
GET /api/v1/agencies -> ListAgencies
GET /api/v1/campaigns/:uid -> ListCampaigns
GET /api/v1/draftCampaigns/:uid -> ListDraftCampaigns
GET /api/v1/proximitySegments/:uid -> ListProximitySegments
GET /api/v1/segments/:uid -> ListSegments

POST /api/v1/campaign/:uid -> CreateCampaign
POST /api/v1/draftCampaign/:uid -> CreateDraftCampaign
POST /api/v1/fullCampaign/:uid -> CreateFullCampaign
POST /api/v1/upgradeCampaign/:uid/:draftCID -> UpgradeCampaign


POST /api/v1/ad/:uid -> CreateAd
POST /api/v1/adGroup/:uid -> CreateAdGroup
POST /api/v1/advertiser -> CreateAdvertiser
POST /api/v1/proximitySegment/:uid -> CreateProximitySegment
POST /api/v1/segment/:uid -> CreateSegment

PUT /api/v1/ad -> UpdateAd
PUT /api/v1/campaign -> UpdateCampaign
PUT /api/v1/draftCampaign -> UpdateDraftCampaign
PUT /api/v1/proximitySegment -> UpdateProximitySegment
PUT /api/v1/segment -> UpdateSegment


DELETE /api/v1/ad/:adID -> DeleteAd
DELETE /api/v1/adGroup/:adgroupID -> DeleteAdGroup
DELETE /api/v1/campaign/:cid -> DeleteCampaign
DELETE /api/v1/draftCampaign/:cid -> DeleteDraftCampaign
DELETE /api/v1/proximitySegment/:segID -> DeleteProximitySegment
DELETE /api/v1/segment/:segID -> DeleteSegment


// those 2 endpoints are limited total 100 requests per hour.
GET /adsReport/:uid/:start/:end
GET /campaignReport/:uid/:cid/:start/:end
```

### Notes

- For [CreateAd](https://godoc.org/github.com/missionMeteora/sdk#Client.CreateAd), [CreateAdRequest.AdImage](https://godoc.org/github.com/missionMeteora/sdk#CreateAdRequest) must be a fully encoded base64 image.

### Examples

#### [ListAdvertisers](https://godoc.org/github.com/missionMeteora/sdk#Client.ListAdvertisers)

```
➤ curl "https://rest.meteora.us/api/v1/advertisers/[your-agency-id]?apiKey=[your-meteora-api-key]" | jq

{
	"data": {
		"advertiser-id": {
			"id": "advertiser-id",
			"agencyID": "your-agency-id",
			"name": "advertiser name",
			"numCmps": 4,
			"status": true
		}
		"another-advertiser-id": {
			"id": "advertiser-id",
			"agencyID": "your-agency-id",
			"name": "advertiser name",
			"numCmps": 4,
			"status": true
		}
	},
	"code": 200,
	"success": true
}
```

#### [CreateSegment](https://godoc.org/github.com/missionMeteora/sdk#Client.CreateSegment) with [Segment](https://godoc.org/github.com/missionMeteora/sdk#Segment)

For any POST/GET requests that uses a struct, you need to match the required struct fields.

```
# Content-Type: application/json is required for any POST/PUT rquests.
➤ curl -H "Content-Type: application/json" -d '{"name": "segment name"}' 'http://localhost:8081/api/v1/segment/2?apiKey=382787ed-6c36-4ca0-a3f1-6bcf926fea7c'

{
	"data": "segment-id",
	"code":200,
	"success":true
}

```

#### [CreateAdGroup](https://godoc.org/github.com/missionMeteora/sdk#Client.CreateAdGroup)

```
# Content-Type: application/json is required for any POST/PUT rquests.
➤ curl -H "Content-Type: application/json" -d '{"name": "ad group name"}' 'http://localhost:8081/api/v1/adGroup/2?apiKey=382787ed-6c36-4ca0-a3f1-6bcf926fea7c'

{
	"data": "adgroup-id",
	"code":200,
	"success":true
}


#### [CreateFullCampaign](https://godoc.org/github.com/missionMeteora/sdk#Client.CreateFullCampaign)

```

{
	"campaign": {
		"active": false,
		"name": "SDK Test Full Campaign",
		"budget": 50,
		"impBudget": 0,
		"created": 0,
		"scheduled": false,
		"start": 0,
		"end": 0,
		"apps": {
			"advancedBidding": {
				"status": true,
				"baseCpm": 2,
				"maxCpm": 5
			},
			"searchRetargeting": {
				"status": true,
				"list": ["nike shoes", "adidas", "shiny shoes"]
			}
		}
	},

	"ads": [
		{
			"name": "sdkTestCampaign-1.png",
			"width": 300,
			"height": 250,
			"landingURL": "https://test.com",
			"adImage": "data:....,base64,"
		},
		{
			"name": "sdkTestCampaign-2.png",
			"width": 300,
			"height": 250,
			"landingURL": "https://test.com",
			"adImage": "data:....,base64,"
		}
	],

	"segments": [
		{
			"name": "Full Segment"
		}
	],

	"proximitySegment": [
		{
			"name": "Full Proximity Segment",
			"locations": [
				{
					"id": "",
					"label": "Starbucks",
					"type": "",
					"center": {
						"lat": 32.8826822,
						"lng": -97.39539739999998
					},
					"radius": 500
				}
			]
		}
	],

	"isDraft": true
}

```
