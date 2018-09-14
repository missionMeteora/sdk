package sdk

import (
	"encoding/json"
)

// DefaultApps are the default apps to be added to a new empty campaign.
var DefaultApps = []App{
	&AppPacing{Status: true},
}

var AllApps = []App{
	AppAdvBidding{},
	AppSearchRetargeting{},
	AppDomainTargeting{},
	AppGeography{},
	UUIDTargeting{},
}

type App interface {
	Name() string
	app()
}

type AppPacing struct {
	Status bool `json:"status"`
}

func (AppPacing) Name() string { return "pacing" }
func (AppPacing) app()         {}

type AppAdvBidding struct {
	Status bool `json:"status"`

	BaseCPM float64 `json:"baseCpm"`
	MaxCPM  float64 `json:"maxCpm"`
}

func (AppAdvBidding) Name() string { return "advancedBidding" }
func (AppAdvBidding) app()         {}

type AppSearchRetargeting struct {
	Status bool `json:"status"`

	List []string `json:"list"`
}

func (AppSearchRetargeting) Name() string { return "searchRetargeting" }
func (AppSearchRetargeting) app()         {}

type AppDomainTargeting struct {
	Targeted []string `json:"targeted"`
	Banned   []string `json:"banned"`
}

func (AppDomainTargeting) Name() string { return "domainTargeting" }
func (AppDomainTargeting) app()         {}

type AppGeography struct {
	Status   bool     `json:"status"`
	Cities   []string `json:"cities"`
	States   []string `json:"states"`
	Zipcodes []string `json:"zipCodes"`
	DMAs     []string `json:"dmas"`
}

func (AppGeography) Name() string { return "geography" }
func (AppGeography) app()         {}

type UUIDTargeting struct {
	Whitelist []string `json:"whitelist"`
	Blacklist []string `json:"blacklist"`
	Status    bool     `json:"status"`
}

func (UUIDTargeting) Name() string { return "uuidTargeting" }
func (UUIDTargeting) app()         {}

// SetApp is a little helper func to add apps to campaigns.
func SetApp(cmp *Campaign, app App) {
	if cmp.Apps == nil {
		cmp.Apps = map[string]*json.RawMessage{}
	}
	cmp.Apps[app.Name()] = RawMarshal(app)
}

func GetApp(name string, val *json.RawMessage) (app App, err error) {
	for _, a := range AllApps {
		if a.Name() == name {
			app = a
			break
		}
	}

	if app == nil {
		return
	}

	if err = json.Unmarshal(*val, &app); err != nil {
		app = nil
	}

	return
}

func RawMarshal(v interface{}) *json.RawMessage {
	j, _ := json.MarshalIndent(v, "", "\t")
	return (*json.RawMessage)(&j)
}
