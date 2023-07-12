package tidyes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Helper struct {
	es *elasticsearch.Client
}

func NewHelper(es *elasticsearch.Client) *Helper {
	return &Helper{es: es}
}

func (h *Helper) Do(ctx context.Context, v any, req esapi.Request) error {
	res, err := req.Do(ctx, h.es)
	if err != nil {
		return err
	}

	defer closer(res.Body, &err)
	defer func() {
		// Drain the request, to make sure connection re-use works properly
		_, _ = io.Copy(io.Discard, res.Body)
	}()

	if res.IsError() {
		var errResult ErrorResult
		err = json.NewDecoder(res.Body).Decode(&errResult)
		if err != nil {
			return err
		}
		if errResult.Error.Type == "resource_already_exists_exception" {
			return nil
		}
		return fmt.Errorf("unexpected status %d: %#+v", res.StatusCode, errResult)
	}

	if v == nil {
		return nil
	}

	return json.NewDecoder(res.Body).Decode(v)
}

func closer(c io.Closer, in *error) {
	cerr := c.Close()
	if *in == nil {
		*in = cerr
	}
}

type RootCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type ErrorReason struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type FailedShard struct {
	Shard  int         `json:"shard"`
	Index  string      `json:"index"`
	Node   string      `json:"node"`
	Reason ErrorReason `json:"reason"`
}

type CausedByDetail struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type CausedBy struct {
	Type     string         `json:"type"`
	Reason   string         `json:"reason"`
	CausedBy CausedByDetail `json:"caused_by"`
}

type ErrorDetail struct {
	RootCause    []RootCause   `json:"root_cause"`
	Type         string        `json:"type"`
	Reason       string        `json:"reason"`
	Phase        string        `json:"phase"`
	Grouped      bool          `json:"grouped"`
	FailedShards []FailedShard `json:"failed_shards"`
	CausedBy     CausedBy      `json:"caused_by"`
}

type ErrorResult struct {
	Error  ErrorDetail `json:"error"`
	Status int         `json:"status"`
}
