package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/syumai/workers/cloudflare"
	"github.com/syumai/workers/cloudflare/fetch"
)

type Datetime struct {
	Now time.Time `json:"datetime"`
}

// getDatetime returns current time in the given timezone
// using http://worldtimeapi.org/ because Cloudflare Workers (https://github.com/syumai/workers) don't support time.Time
func getDatetime(ctx context.Context) (Datetime, error) {
	client := fetch.NewClient()
	timezone := cloudflare.Getenv(ctx, "TIMEZONE")
	req, err := fetch.NewRequest(ctx, "GET", fmt.Sprintf("client/api/timezone/%s", timezone), nil)
	if err != nil {
		return Datetime{}, err
	}
	res, err := client.Do(req, nil)
	if err != nil {
		return Datetime{}, err
	}
	defer res.Body.Close()
	var datetime Datetime
	err = json.NewDecoder(res.Body).Decode(&datetime)
	if err != nil {
		return Datetime{}, err
	}

	return datetime, nil
}
