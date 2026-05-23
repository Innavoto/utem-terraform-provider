package client

import (
	"context"
	"fmt"
)

type ScanSchedule struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	TargetHost  string `json:"target_host"`
	ScanType    string `json:"scan_type"`
	CronExpr    string `json:"cron_expression"`
	IsEnabled   bool   `json:"is_enabled"`
	Modules     []string `json:"modules,omitempty"`
}

type ScanScheduleList struct {
	Items []ScanSchedule `json:"items"`
	Total int            `json:"total"`
}

func (c *Client) ListScanSchedules(ctx context.Context) ([]ScanSchedule, error) {
	var result ScanScheduleList
	if err := c.Get(ctx, "/api/v1/scan-schedules", &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (c *Client) GetScanSchedule(ctx context.Context, id string) (*ScanSchedule, error) {
	var result ScanSchedule
	if err := c.Get(ctx, fmt.Sprintf("/api/v1/scan-schedules/%s", id), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateScanSchedule(ctx context.Context, input *ScanSchedule) (*ScanSchedule, error) {
	var result ScanSchedule
	if err := c.Post(ctx, "/api/v1/scan-schedules", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateScanSchedule(ctx context.Context, id string, input *ScanSchedule) (*ScanSchedule, error) {
	var result ScanSchedule
	if err := c.Put(ctx, fmt.Sprintf("/api/v1/scan-schedules/%s", id), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteScanSchedule(ctx context.Context, id string) error {
	return c.Delete(ctx, fmt.Sprintf("/api/v1/scan-schedules/%s", id))
}
