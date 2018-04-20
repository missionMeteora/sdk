package sdk

import (
	"encoding/json"
)

type Pacing = struct {
	Status bool `json:"status"`
}

// DefaultApps are the default apps to be added to a new empty campaign if cmp.Apps == nil.
var DefaultApps = map[string]interface{}{
	"pacing": &Pacing{Status: true},
}

// SetApp is a little helper func to add apps to campaigns.
func SetApp(cmp *Campaign, name string, app interface{}) error {
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
