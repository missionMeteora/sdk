# sdk

Meteora API's SDK

## Download
```sh
$ go get -u github.com/missionMeteora/sdk
```


## Usage

`import "github.com/missionMeteora/sdk"`

```go
const (
	// Version is the current SDK version
	Version       = "v0.5"
	DefaultServer = "https://meteora.us"
)
```

```go
var (
	// ErrMissingUID is returned when a user ID is expected and missing
	ErrMissingUID = errors.New("missing user id")

	// ErrMissingCID is returned when a campaign ID is expected and missing
	ErrMissingCID = errors.New("missing campaign id")

	// ErrMissingSegID is returned when a segment ID is expected and missing
	ErrMissingSegID = errors.New("missing segment id")

	// ErrDateRange is returned for invalid or missing date ranges
	ErrDateRange = errors.New("bad or missing date range")

	// ErrInvalidName is returned when a name is invalid
	ErrInvalidName = errors.New("invalid name")

	// ErrInvalidLanding is returned when a landing url is invalid
	ErrInvalidLanding = errors.New("invalid landing url")

	// ErrInvalidAdSize is returned when an invalid ad size is provided
	ErrInvalidAdSize = fmt.Errorf("invalid ad size, the accepted sizes are: %s", allowedAdRectsString())
)
```
common errors

```go
var DefaultApps = map[string]interface{}{
	"pacing": &Pacing{Status: true},
}
```
DefaultApps are the default apps to be added to a new empty campaign if cmp.Apps
== nil.

#### func  SetApp

```go
func SetApp(cmp *Campaign, name string, app interface{}) error
```
SetApp is a little helper func to add apps to campaigns.

#### type Ad

```go
type Ad struct {
	ID           string `json:"id,omitempty"`
	GroupID      string `json:"group,omitempty"`
	AdvertiserID string `json:"advertiserID,omitempty"`

	Active        bool   `json:"active,omitempty"`
	Clicks        int    `json:"clicks,omitempty"`
	Conv          int    `json:"conv,omitempty"`
	Width         int    `json:"width,omitempty"`
	Height        int    `json:"height,omitempty"`
	Imps          int    `json:"imps,omitempty"`
	AdType        string `json:"adType,omitempty"`
	ImageKey      string `json:"imageKey,omitempty"`
	ImageLocation string `json:"imageLocation,omitempty"`
	Name          string `json:"name,omitempty"`
	Size          string `json:"size,omitempty"`

	LandingURL string `json:"landingURL,omitempty"`

	ImpTracker   string `json:"impTrack,omitempty"`
	ClickTracker string `json:"clickTrack,omitempty"`
}
```

Ad is mostly a trimmed down version of api/server/types.Ad

#### type AdGroup

```go
type AdGroup struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	AdCount int    `json:"adCount"`
}
```

AdGroup represents an ad group

#### type AdReport

```go
type AdReport struct {
	ID     string `json:"id"`
	Imp    int64  `json:"imp,omitempty"`
	Clicks int64  `json:"clicks,omitempty"`
}
```

AdReport represents a basic ad report

#### type AdsListItem

```go
type AdsListItem struct {
	ID   string `json:"id"`
	Data *Ad    `json:"data"`
}
```

AdsListItem is a copied type from API for compatibility with the Heatmap JS lib

#### type Advertiser

```go
type Advertiser struct {
	ID           string `json:"id"`
	AgencyID     string `json:"agencyID"`
	Name         string `json:"name"`
	NumCampaigns int    `json:"numCmps"`
	Status       bool   `json:"status"`
}
```

Advertiser represents an advertiser user

#### type Agency

```go
type Agency struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
```

Agency represents an agency user

#### type Campaign

```go
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
```

Campaign is a QoL alias for campaigns.Campaign

#### type CampaignReport

```go
type CampaignReport struct {
	ID            string              `json:"id"`
	Imp           int64               `json:"imp,omitempty"`
	Clicks        int64               `json:"clicks,omitempty"`
	Spent         float64             `json:"spent,omitempty"`
	Period        uint32              `json:"reportingPeriod,omitempty"`
	InStoreVisits uint32              `json:"inStoreVisits,omitempty"`
	Domains       map[string]ImpClick `json:"domains,omitempty"`
}
```

CampaignReport represents a basic campaign report

#### type Click

```go
type Click = struct {
	X int `json:"x"`
	Y int `json:"y"`
	// TTC represents time to click (in milliseconds)
	TTC int `json:"ttc,omitempty"`

	UUID string `json:"uuid"`
	AdID string `json:"adID"`
	TS   int64  `json:"ts"`
}
```

Click is a click action

#### type Client

```go
type Client struct {
}
```

Client is a helper structure who holds onto an apiKey

#### func  New

```go
func New(apiKey string) *Client
```
New returns a new instance with the default server addr and given apiKey.

#### func  NewWithAddr

```go
func NewWithAddr(apiAddr string, apiKey string) *Client
```
NewWithAddr returns a new instance of Client with the given api addr and key. If
testing locally on an unsecure https server, use
https+insecure://localhost:8080.

#### func (*Client) AsUser

```go
func (c *Client) AsUser(ctx context.Context, uid string) (nc *Client, err error)
```
AsUser fetches the given uid's api key and returns a new client using it or an
error.

#### func (*Client) CreateAd

```go
func (c *Client) CreateAd(ctx context.Context, uid string, req *CreateAdRequest, r io.Reader) (ad *Ad, err error)
```
CreateAd will create an Advertisement for a given user ID

#### func (*Client) CreateAdFromFile

```go
func (c *Client) CreateAdFromFile(ctx context.Context, uid string, req *CreateAdRequest, filename string) (ad *Ad, err error)
```
CreateAdFromFile will create an Advertisement for a given user ID using a given
filename

#### func (*Client) CreateAdGroup

```go
func (c *Client) CreateAdGroup(ctx context.Context, uid, name string) (id string, err error)
```
CreateAdGroup will create an adgroup by the owning user ID

#### func (*Client) CreateCampaign

```go
func (c *Client) CreateCampaign(ctx context.Context, uid string, cmp *Campaign) (string, error)
```
CreateCampaign will create a campaign for a given user ID

#### func (*Client) CreateDraftCampaign

```go
func (c *Client) CreateDraftCampaign(ctx context.Context, uid string, cmp *Campaign) (string, error)
```
CreateDraftCampaign will create a draft campaign for a given user ID

#### func (*Client) CreateProximitySegment

```go
func (c *Client) CreateProximitySegment(ctx context.Context, uid string, seg *ProximitySegment) (id string, err error)
```
CreateProxSegment will create a proximity segment for a given user ID

#### func (*Client) CreateSegment

```go
func (c *Client) CreateSegment(ctx context.Context, uid string, seg *Segment) (id string, err error)
```
CreateSegment will create a new segment for a given user ID

#### func (*Client) CurrentKey

```go
func (c *Client) CurrentKey() string
```
CurrentKey returns the API key used to initalize this client

#### func (*Client) DeleteAd

```go
func (c *Client) DeleteAd(ctx context.Context, adID string) (err error)
```
DeleteAd will delete an ad by advertisement ID

#### func (*Client) DeleteAdGroup

```go
func (c *Client) DeleteAdGroup(ctx context.Context, adgroupID string) (err error)
```
DeleteAdGroup will delete an adgroup by it's ID

#### func (*Client) DeleteCampaign

```go
func (c *Client) DeleteCampaign(ctx context.Context, cid string) error
```
DeleteCampaign will delete a campaign by it's ID

#### func (*Client) DeleteDraftCampaign

```go
func (c *Client) DeleteDraftCampaign(ctx context.Context, cid string) error
```
DeleteDraftCampaign will delete a draft campaign by it's ID

#### func (*Client) DeleteProximitySegment

```go
func (c *Client) DeleteProximitySegment(ctx context.Context, segID string) (err error)
```
DeleteProxSegment will delete a proximity segment by it's ID

#### func (*Client) DeleteSegment

```go
func (c *Client) DeleteSegment(ctx context.Context, segID string) (err error)
```
DeleteSegment will delete a segment by it's ID

#### func (*Client) GetAPIVersion

```go
func (c *Client) GetAPIVersion(ctx context.Context) (ver string, err error)
```
GetAPIVersion will get the current API version

#### func (*Client) GetAdsReport

```go
func (c *Client) GetAdsReport(ctx context.Context, uid, start, end string) (rp map[string]*AdReport, err error)
```
GetAdsReport will generate a advertisements report for a given user ID and date
range

#### func (*Client) GetCampaign

```go
func (c *Client) GetCampaign(ctx context.Context, cid string) (*Campaign, error)
```
GetCampaign will get a campaign by campaign id

#### func (*Client) GetCampaignReport

```go
func (c *Client) GetCampaignReport(ctx context.Context, uid, cid, start, end string) (cs *CampaignReport, err error)
```
GetCampaignReport will generate a report for a given campaign ID and date range

#### func (*Client) GetDraftCampaign

```go
func (c *Client) GetDraftCampaign(ctx context.Context, cid string) (*Campaign, error)
```
GetDraftCampaign will get a campaign by campaign id

#### func (*Client) GetHeatmap

```go
func (c *Client) GetHeatmap(ctx context.Context, uid string) (out *Heatmap, err error)
```
GetHeatmap will return the heatmaps belonging to a user ID

#### func (*Client) GetUserAPIKey

```go
func (c *Client) GetUserAPIKey(ctx context.Context, uid string) (apiKey string, err error)
```
GetUserAPIKey will return a user's API key

#### func (*Client) GetUserID

```go
func (c *Client) GetUserID(ctx context.Context) (uid string, err error)
```
GetUserID returns the current key owner's uid.

#### func (*Client) ListAdGroups

```go
func (c *Client) ListAdGroups(ctx context.Context, uid string) (ags map[string]*AdGroup, err error)
```
ListAdGroups will list all the adgroups for a given user id

#### func (*Client) ListAds

```go
func (c *Client) ListAds(ctx context.Context, uid string) (ads map[string]*Ad, err error)
```

#### func (*Client) ListAdsByAdGroup

```go
func (c *Client) ListAdsByAdGroup(ctx context.Context, uid, adGroupID string) (map[string]*Ad, error)
```
ListAdsByAdGroup will list all ads for a user and ad group

#### func (*Client) ListAdsFilter

```go
func (c *Client) ListAdsFilter(ctx context.Context, uid string, filterFn func(ad *Ad) bool) (ads map[string]*Ad, err error)
```
ListAdsFilter will list ads which pass a given filter

#### func (*Client) ListAdvertisers

```go
func (c *Client) ListAdvertisers(ctx context.Context, agencyID string) (out map[string]*Advertiser, err error)
```
ListAdvertisers will list the advertisers for a given agency

#### func (*Client) ListAgencies

```go
func (c *Client) ListAgencies(ctx context.Context) (out map[string]*Agency, err error)
```
ListAgencies will list the agencies

#### func (*Client) ListCampaigns

```go
func (c *Client) ListCampaigns(ctx context.Context, uid string) (map[string]*Campaign, error)
```
ListCampaigns will list all the campaigns for a given user id

#### func (*Client) ListDraftCampaigns

```go
func (c *Client) ListDraftCampaigns(ctx context.Context, uid string) (map[string]*Campaign, error)
```
ListDraftCampaigns will list all the draft campaigns for a given user id

#### func (*Client) ListProximitySegments

```go
func (c *Client) ListProximitySegments(ctx context.Context, uid string) (segs map[string]*ProximitySegment, err error)
```
ListProxSegments will list the proximity segments belonging to a given user ID

#### func (*Client) ListSegments

```go
func (c *Client) ListSegments(ctx context.Context, uid string) (segs map[string]*Segment, err error)
```
ListSegments will list the segments belonging to a given user ID

#### func (*Client) RawRequest

```go
func (c *Client) RawRequest(method, endpoint string, req, resp interface{}) (err error)
```
RawRequest is an alias for RawRequestCtx(context.Background(), method, endpoint,
req, resp)

#### func (*Client) RawRequestCtx

```go
func (c *Client) RawRequestCtx(ctx context.Context, method, endpoint string, req, resp interface{}) (err error)
```
RawRequestCtx allows a raw c.raw request to the given endpoint

#### func (*Client) UpdateAd

```go
func (c *Client) UpdateAd(ctx context.Context, ad *Ad) error
```
UpdateAd will update an ad

#### func (*Client) UpdateCampaign

```go
func (c *Client) UpdateCampaign(ctx context.Context, cmp *Campaign) error
```
UpdateCampaign will update a campaign

#### func (*Client) UpdateDraftCampaign

```go
func (c *Client) UpdateDraftCampaign(ctx context.Context, cmp *Campaign) error
```
UpdateDraftCampaign will update a draft campaign

#### func (*Client) UpdateProximitySegment

```go
func (c *Client) UpdateProximitySegment(ctx context.Context, seg *ProximitySegment) (err error)
```
UpdateProxSegment will update a proximity segment by it's ID

#### func (*Client) UpdateSegment

```go
func (c *Client) UpdateSegment(ctx context.Context, seg *Segment) (err error)
```
UpdateSegment will update a segment by it's ID

#### type Coords

```go
type Coords = struct {
	// Note, the json keys for these values are setup to match google's coordinate object
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}
```

Coords are a set of coordinates

#### type CreateAdRequest

```go
type CreateAdRequest struct {
	Name       string `json:"name,omitempty"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	LandingURL string `json:"landingURL,omitempty"`

	// Everything below here is optional
	GroupID      string `json:"group,omitempty"`
	ImpTracker   string `json:"impTrack,omitempty"`
	ClickTracker string `json:"clickTrack,omitempty"`
	Encoded      bool   `json:"encoded,omitempty"`
}
```

CreateAdRequest is the request needed to create an Ad

#### func (*CreateAdRequest) Validate

```go
func (r *CreateAdRequest) Validate() error
```
Validate will validate a create ad request

#### type Heatmap

```go
type Heatmap struct {
	AdsList []*AdsListItem `json:"adsList"`
	Clicks  []*Click       `json:"clicks"`
}
```

Heatmap contains the raw Heatmap data from API

#### func (*Heatmap) AllAds

```go
func (hm *Heatmap) AllAds() (out map[string]*Ad)
```
AllAds returns a map of all the ads keyed by their IDs.

#### func (*Heatmap) ClicksByAds

```go
func (hm *Heatmap) ClicksByAds() (out map[string][]*Click)
```
ClicksByAds returns a map of a slice of clicks by Ad ID.

#### type ImpClick

```go
type ImpClick struct {
	Imp    int64 `json:"imp"`
	Clicks int64 `json:"clicks"`
}
```

ImpClick has the numbers of imps and clicks for a domain

#### type Location

```go
type Location = struct {
	// ID of location
	ID string `json:"id"`
	// Label of the location
	Label string `json:"label"`
	// Type of location, can be "circle" or "polygon"
	Type string `json:"type"`
	// Center coordinates of the location
	Center Coords `json:"center"`
	// Radius in meters
	Radius float64 `json:"radius"`
	// Points for drawing on map, this is only used for polygons
	Points []Coords `json:"points"`
}
```

Location represents a specified location

#### type Pacing

```go
type Pacing = struct {
	Status bool `json:"status"`
}
```


#### type ProximitySegment

```go
type ProximitySegment = struct {
	ID      string `json:"id,omitempty"` // this is strictly used by the SDK
	Name    string `json:"name"`
	OwnerID string `json:"ownerID"`

	// Current location id
	IDCounter int `json:"idCounter"`

	// List of locations
	Locations []*Location `json:"locations"`
	// List of deletedLocations
	DeletedLocations []*Location `json:"deletedLocations,omitempty"`

	// Lookback period
	Lookback int16 `json:"lookback"`
	// Radius threshold for GPS padding
	RadiusThreshold uint16 `json:"radiusThreshold"`
}
```

Segment is a segment by proximity location

#### type Rects

```go
type Rects struct {
	Width  int
	Height int
}
```

Rects represent a unit of width and height

#### type Segment

```go
type Segment struct {
	SegmentID              string   `json:"segmentID,omitempty"`
	AdvertiserID           string   `json:"advertiserID,omitempty"`
	Active                 bool     `json:"active,omitempty"`
	Name                   string   `json:"name,omitempty"`
	TargetConsumers        string   `json:"targetConsumers,omitempty"`
	MinimumPageViews       int      `json:"minimumPageViews,omitempty"`
	VisitedSiteAtLeast     int      `json:"visitedSiteAtLeast,omitempty"`
	UniqueUsers            int      `json:"uniqueUsers,omitempty"`
	AvgPurchase            int      `json:"avgPurchase,omitempty"`
	PreviouslySpentBetween []int    `json:"previouslySpentBetween,omitempty"`
	LastVisitedBetween     []int    `json:"lastVisitedBetween,omitempty"`
	ViewedProductsBetween  []int    `json:"viewedProductsBetween,omitempty"`
	UrlsVisited            []string `json:"urlsVisited,omitempty"`
	SkusVisited            []string `json:"skusVisited,omitempty"`
	IncludeParams          []string `json:"includeParams,omitempty"`
	ExcludeParams          []string `json:"excludeParams,omitempty"`
}
```

Segment is a copy of api/internal/common.Segment
