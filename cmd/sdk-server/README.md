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
GET /api/v1/campaign/:cid -> GetCampaign
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

POST /api/v1/ad/:uid -> CreateAd
POST /api/v1/adGroup/:uid/:name -> CreateAdGroup
POST /api/v1/advertiser -> CreateAdvertiser
POST /api/v1/campaign/:uid -> CreateCampaign
POST /api/v1/draftCampaign/:uid -> CreateDraftCampaign
POST /api/v1/fullCampaign/:uid -> CreateFullCampaign
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

### Examples

#### [ListAdvertisers](https://godoc.org/github.com/missionMeteora/sdk#Client.ListAdvertisers)

```json
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
