package sdk

import (
	"encoding/json"
)

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
	BaseCPM float64 `json:"baseCpm"`
	MaxCPM  float64 `json:"maxCpm"`
}

func (AppAdvBidding) Name() string { return "advancedBidding" }
func (AppAdvBidding) app()         {}

// DefaultApps are the default apps to be added to a new empty campaign.
var DefaultApps = []App{
	&AppPacing{Status: true},
}

// SetApp is a little helper func to add apps to campaigns.
func SetApp(cmp *Campaign, name string, app App) error {
	j, err := json.Marshal(app)
	if err != nil {
		return err
	}
	if cmp.Apps == nil {
		cmp.Apps = map[string]*json.RawMessage{}
	}
	cmp.Apps[name] = (*json.RawMessage)(&j)
	return nil
}
