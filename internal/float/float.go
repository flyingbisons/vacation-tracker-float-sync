package float

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/syumai/workers/cloudflare/fetch"
)

const ApiUrl = "https://api.float.com/v3"

const (
	TimeTypePaidID = 275796
	TimeTypeFreeID = 323532
	TimeTypeSickID = 275797
)

type Client struct {
	apiKey string
	apiUrl string
	http   *fetch.Client
}

func NewClient(apiKey, apiUrl string) *Client {
	return &Client{
		apiKey: apiKey,
		apiUrl: apiUrl,
		http:   fetch.NewClient(),
	}
}

func (c *Client) prepareRequest(ctx context.Context, method string, url string, body io.Reader) (*fetch.Request, error) {
	req, err := fetch.NewRequest(ctx, method, fmt.Sprintf("%s%s", c.apiUrl, url), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("User-Agent", "vt-float-flying-bisons/0.1 (pawel@flyingbisons.com)")

	return req, nil
}

func (c *Client) GetTimeOffTypes(ctx context.Context) ([]TimeOffType, error) {
	req, err := c.prepareRequest(ctx, "GET", "/timeoff-types?active=1", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err.Error())
	}
	resp, err := c.http.Do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting time off types: %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("error closing response body", slog.String("error", err.Error()))
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code: %s", resp.Status)
	}

	var timeOffTypes []TimeOffType
	err = json.NewDecoder(resp.Body).Decode(&timeOffTypes)
	if err != nil {
		return nil, fmt.Errorf("error decoding time off types: %s", err.Error())
	}

	return timeOffTypes, nil
}

func (c *Client) AddTimeOff(ctx context.Context, off TimeOff) (TimeOff, error) {
	var timeOff TimeOff

	body, err := json.Marshal(off)
	if err != nil {
		return timeOff, fmt.Errorf("error marshalling time off request: %s", err.Error())
	}

	req, err := c.prepareRequest(ctx, "POST", "/timeoffs", bytes.NewBuffer(body))
	if err != nil {
		return timeOff, fmt.Errorf("error creating time off request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req, nil)
	if err != nil {
		return timeOff, fmt.Errorf("error : %s", err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("error closing response body", slog.String("error", err.Error()))
		}
	}(resp.Body)

	if resp.StatusCode != 201 {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return timeOff, fmt.Errorf("error reading response body: %s", err.Error())
		}
		slog.Debug(string(body))
		return timeOff, fmt.Errorf("invalid status code: %s", resp.Status)
	}
	slog.Info(string(body))
	err = json.NewDecoder(resp.Body).Decode(&timeOff)
	if err != nil {
		return timeOff, fmt.Errorf("error decoding response body: %s", err.Error())
	}

	return timeOff, nil
}

func (c *Client) FindEmployeeByEmail(ctx context.Context, email string) (Employee, error) {
	var employee Employee
	req, err := c.prepareRequest(ctx, "GET", fmt.Sprintf("/people?email=%s", email), nil)
	if err != nil {
		return employee, fmt.Errorf("error creating time off request: %s", err.Error())
	}
	resp, err := c.http.Do(req, nil)
	if err != nil {
		return employee, fmt.Errorf("error : %s", err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("error closing response body", slog.String("error", err.Error()))
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return employee, fmt.Errorf("invalid status code: %s", resp.Status)
	}

	var employees []Employee
	err = json.NewDecoder(resp.Body).Decode(&employees)
	if err != nil {
		return employee, fmt.Errorf("error decoding response body: %s", err.Error())
	}

	if len(employees) == 0 {
		return employee, fmt.Errorf("no employee found %s", email)
	}

	return employees[0], nil
}
