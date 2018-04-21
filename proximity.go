package sdk

import (
	"context"
	"strconv"
)

// Segment is a segment by proximity location
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

// Location represents a specified location
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

// Coords are a set of coordinates
type Coords = struct {
	// Note, the json keys for these values are setup to match google's coordinate object
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

// CreateProxSegment will create a proximity segment for a given user ID
func (c *Client) CreateProximitySegment(ctx context.Context, uid string, seg *ProximitySegment) (id string, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var resp idOrDataResp

	seg.OwnerID = uid
	seg.IDCounter = len(seg.Locations)

	// auto-fill location ids if needed
	for i, loc := range seg.Locations {
		loc.ID = strconv.Itoa(i)
		if loc.Type == "" && loc.Center.Latitude != 0 && loc.Center.Longitude != 0 {
			loc.Type = "circle"
		}
	}

	if err = c.rawPost(ctx, "segments/proximity/byKey/"+uid, seg, &resp); err != nil {
		return
	}

	id = resp.String()
	seg.ID = id

	return
}

// UpdateProxSegment will update a proximity segment by it's ID
func (c *Client) UpdateProximitySegment(ctx context.Context, seg *ProximitySegment) (err error) {
	if seg == nil || seg.ID == "" {
		err = ErrMissingSegID
		return
	}

	return c.rawPut(ctx, "segments/proximity/byKey/"+seg.ID, seg, nil)
}

// DeleteProxSegment will delete a proximity segment by it's ID
func (c *Client) DeleteProximitySegment(ctx context.Context, segID string) (err error) {
	if segID == "" {
		err = ErrMissingSegID
		return
	}

	return c.rawDelete(ctx, "segments/proximity/byKey/"+segID, nil)
}

// ListProxSegments will list the proximity segments belonging to a given user ID
func (c *Client) ListProximitySegments(ctx context.Context, uid string) (segs map[string]*ProximitySegment, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var resp struct {
		Data map[string]*ProximitySegment `json:"data"`
	}

	if err = c.rawGet(ctx, "segments/proximity/byOwner/"+uid, &resp); err != nil {
		return
	}

	segs = resp.Data

	for id, seg := range segs {
		seg.ID = id
	}

	return
}
