package sdk

import "context"

// Segment is a copy of api/internal/common.Segment
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

// CreateSegment will create a new segment for a given user ID
func (c *Client) CreateSegment(ctx context.Context, uid string, seg *Segment) (id string, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var resp struct {
		ID string `json:"id"`
	}

	seg.AdvertiserID = uid
	seg.UniqueUsers = 0

	// dashboard defaults
	if seg.TargetConsumers == "" {
		seg.TargetConsumers = "Anyone"
	}
	if len(seg.LastVisitedBetween) == 0 {
		seg.LastVisitedBetween = []int{0, 90}
	}
	if len(seg.PreviouslySpentBetween) == 0 {
		seg.PreviouslySpentBetween = []int{0, -1}
	}

	if err = c.rawPost(ctx, "segments/retargeting/"+uid, seg, &resp); err != nil {
		return
	}

	id = resp.ID
	seg.SegmentID = id

	return
}

// UpdateSegment will update a segment by it's ID
func (c *Client) UpdateSegment(ctx context.Context, seg *Segment) (err error) {
	if seg == nil || seg.SegmentID == "" {
		err = ErrMissingSegID
		return
	}

	return c.rawPut(ctx, "segments/retargeting/"+seg.SegmentID, seg, nil)
}

// DeleteSegment will delete a segment by it's ID
func (c *Client) DeleteSegment(ctx context.Context, segID string) (err error) {
	if segID == "" {
		err = ErrMissingSegID
		return
	}

	return c.rawDelete(ctx, "segments/retargeting/"+segID, nil)
}

// ListSegments will list the segments belonging to a given user ID
func (c *Client) ListSegments(ctx context.Context, uid string) (segs map[string]*Segment, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var resp []struct {
		ID   string   `json:"id"`
		Data *Segment `json:"data"`
	}
	if err = c.rawGet(ctx, "segmentsList/"+uid, &resp); err != nil {
		return
	}

	segs = make(map[string]*Segment, len(resp))
	for _, sd := range resp {
		s := sd.Data
		s.SegmentID = sd.ID
		s.AdvertiserID = uid
		segs[s.SegmentID] = s
	}

	return
}
