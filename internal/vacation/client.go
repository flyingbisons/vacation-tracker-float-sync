package vacation

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/syumai/workers/cloudflare/fetch"
)

var (
	LeaveTypeFreeTime       = "161715c6-5e1f-48ff-8bdb-dc9ee028f061"
	LeaveTypePaidOffHalfDay = "3c3bfb46-c669-4a85-bf19-834b0ea375ac"
	LeaveTypePaidOff        = "ed58b4ca-f5b0-435c-b415-f60c48eed9a4"
)

const ApiUrl = "https://api.vacationtracker.io/v1"

type Client struct {
	apiKey      string
	apiUrl      string
	client      *fetch.Client
	currentTime time.Time
}

func NewClient(apiKey, apiUrl string, currentTime time.Time) *Client {
	return &Client{
		apiKey:      apiKey,
		apiUrl:      apiUrl,
		client:      fetch.NewClient(),
		currentTime: currentTime,
	}
}

func (c *Client) prepareRequest(ctx context.Context, method string, url string, body io.Reader) (*fetch.Request, error) {
	req, err := fetch.NewRequest(ctx, method, fmt.Sprintf("%s%s", c.apiUrl, url), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("User-Agent", "vt-float-flying-bisons/0.1 (pawel@flyingbisons.com)")

	return req, nil
}

type getLeaveRequestsPageParams struct {
	nextToken *string
	startDate string
	endDate   string
}

func (c *Client) getLeaveRequestsPage(ctx context.Context, params getLeaveRequestsPageParams) (Leavs, error) {
	var leaves Leavs
	token := ""
	if params.nextToken != nil {
		token = *params.nextToken
	}

	url := fmt.Sprintf("/leaves?limit=200&startDate=%s&endDate=%s&statuses=APPROVED&nextToken=%s", params.startDate, params.endDate, token)

	req, err := c.prepareRequest(ctx, "GET", url, nil)
	if err != nil {
		return leaves, fmt.Errorf("error creating request: %s", err.Error())
	}

	res, err := c.client.Do(req, nil)
	if err != nil {
		return leaves, fmt.Errorf("error getting time off types: %s", err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln("error closing response body:" + err.Error())
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return leaves, fmt.Errorf("invalid status code: %s", res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(&leaves); err != nil {
		return leaves, fmt.Errorf("error decoding time off types: %s", err.Error())
	}

	slog.Info("leaves", slog.Int("c", len(leaves.Data)))
	return leaves, err
}

func (c *Client) LeaveRequests(ctx context.Context) ([]Leave, error) {
	startDate, endDate := c.getDates()
	var leaves []Leave
	var nextToken *string

	for {
		pageLeaves, err := c.getLeaveRequestsPage(ctx, getLeaveRequestsPageParams{
			nextToken: nextToken,
			startDate: startDate,
			endDate:   endDate,
		})
		if err != nil {
			return leaves, err
		}

		for _, leave := range pageLeaves.Data {
			leaves = append(leaves, leave)
		}

		if pageLeaves.NextToken == nil {
			break
		}

		nextToken = pageLeaves.NextToken
	}
	return leaves, nil
}

func (c *Client) getUsersPage(ctx context.Context, nextToken *string) (map[string]string, *string, error) {
	pageUsers := make(map[string]string)
	token := ""
	if nextToken != nil {
		token = *nextToken
	}
	url := fmt.Sprintf("/users?limit=150&nextToken=%s", token)

	req, err := c.prepareRequest(ctx, "GET", url, nil)
	if err != nil {
		return pageUsers, nil, fmt.Errorf("error creating request: %s", err.Error())
	}
	res, err := c.client.Do(req, nil)
	if err != nil {
		return pageUsers, nil, fmt.Errorf("error getting pageUsers: %s", err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("error closing response body:" + err.Error())
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return pageUsers, nil, fmt.Errorf("get vacation tracker pageUsers - invalid status code: %s", res.Status)
	}

	var usersResponse UsersResponse
	if err := json.NewDecoder(res.Body).Decode(&usersResponse); err != nil {
		return pageUsers, nil, fmt.Errorf("error decoding pageUsers response: %s", err.Error())
	}

	for _, user := range usersResponse.Users {
		pageUsers[user.ID] = user.Email
	}

	return pageUsers, usersResponse.NextToken, nil
}

func (c *Client) Users(ctx context.Context) (map[string]string, error) {
	users := make(map[string]string)
	var nextToken *string

	for {
		pageUsers, newToken, err := c.getUsersPage(ctx, nextToken)
		if err != nil {
			return users, err
		}

		for k, v := range pageUsers {
			users[k] = v
		}

		if newToken == nil {
			break
		}

		nextToken = newToken
	}

	return users, nil
}
func (c *Client) getDates() (string, string) {
	endDate := c.currentTime.AddDate(0, 6, 0).Format("2006-01-02")
	startDate := c.currentTime.Format("2006-01-02")
	return startDate, endDate
}
