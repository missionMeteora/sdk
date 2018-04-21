package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
