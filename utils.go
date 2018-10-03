package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

type errorResp struct {
	Error   interface{} `json:"error"`
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
}

type respOrError struct {
	V interface{}
}

func (roe *respOrError) Handle(s int, r io.Reader) (err error) {
	if s >= 200 && s < 400 {
		if roe.V == nil {
			return nil
		}

		return json.NewDecoder(r).Decode(roe.V)
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(r); err != nil || buf.Len() == 0 {
		return
	}

	if ch := buf.Bytes()[0]; ch != '{' && ch != '[' {
		return fmt.Errorf("error(%d): %s", s, buf.Bytes())
	}

	var er errorResp
	// using bytes.NewReader so the buffer doesn't get reset
	if err = json.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(&er); err != nil {
		return
	}

	switch v := er.Error.(type) {
	case bool:
	case string:
		er.Message = v

	default:
		return fmt.Errorf("unexpected response from api: %s", buf.Bytes())
	}

	if er.Message != "" {
		return fmt.Errorf("api error: %s", er.Message)
	}

	return fmt.Errorf("multiple API errors: %q", er.Errors)
}

type idOrDataResp struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

func (id *idOrDataResp) String() string {
	if id.ID != "" {
		return id.ID
	}

	return id.Data
}

func getStartEnd(start, end time.Time) (s, e string, err error) {
	if start.IsZero() || end.IsZero() {
		err = ErrDateRange
		return
	}

	if start == AllTime {
		s = time.Now().UTC().AddDate(-1, 0, 0).Format("2006-01-02")
	} else {
		s = start.UTC().Format("2006-01-02")
	}

	if end == AllTime {
		e = time.Now().UTC().Format("2006-01-02")
	} else {
		e = end.UTC().Format("2006-01-02")
	}

	return
}

func MidnightToMidnight(t time.Time) (start, end time.Time) {
	start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end = start.AddDate(0, 0, 1).Add(-time.Second)
	return
}

func DateToTime(date string) time.Time {
	if date == "-1" {
		return AllTime
	}

	if t, err := time.Parse(`20060102`, date); err == nil {
		return t.UTC()
	}

	if t, err := time.Parse(`2006-01-02`, date); err == nil {
		return t.UTC()
	}

	if t, err := time.Parse(`02 Jan 06`, date); err == nil {
		return t.UTC()
	}

	if u, err := strconv.ParseInt(date, 10, 64); err == nil {
		if len(date) > 10 { // js timestamps are in MS
			u /= 1000
		}
		return time.Unix(u, 0)
	}

	return time.Time{}
}

func isNumber(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func verifyUserCampaign(ctx context.Context, c *Client, uid, cid string) bool {
	if cid != "" && cid != "-1" {
		cmps, _ := c.ListCampaigns(ctx, uid)
		return cmps[cid] != nil
	}

	_, err := c.AsUser(ctx, uid)
	return err == nil
}
