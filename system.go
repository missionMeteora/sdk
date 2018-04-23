package sdk

import "context"

// GetAPIVersion will get the current API version
func (c *Client) GetAPIVersion(ctx context.Context) (ver string, err error) {
	var resp struct {
		Version string `json:"version"`
	}

	if err = c.rawGet(ctx, "version", &resp); err != nil {
		return
	}

	ver = resp.Version
	return
}

// GetUserID returns the current key owner's uid.
func (c *Client) GetUserID(ctx context.Context) (uid string, err error) {
	var resp idOrDataResp

	if err = c.rawGet(ctx, "userID", &resp); err != nil {
		return
	}

	uid = resp.String()
	return
}

// GetUserAPIKey will return a user's API key
func (c *Client) GetUserAPIKey(ctx context.Context, uid string) (apiKey string, err error) {
	if uid == "" {
		err = ErrMissingUID
		return
	}

	var resp struct {
		Key string `json:"key"`
	}

	if err = c.rawGet(ctx, "apiKey/"+uid, &resp); err != nil {
		return
	}

	apiKey = resp.Key
	return
}
