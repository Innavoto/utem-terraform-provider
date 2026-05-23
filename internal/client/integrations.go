package client

import (
	"context"
	"fmt"
)

type Integration struct {
	ID              int    `json:"id,omitempty"`
	Name            string `json:"name"`
	IntegrationType string `json:"integration_type"`
	IsEnabled       bool   `json:"is_enabled"`
	Config          map[string]interface{} `json:"config,omitempty"`
	WebhookURL      string `json:"webhook_url,omitempty"`
	Description     string `json:"description,omitempty"`
}

type IntegrationList struct {
	Items []Integration `json:"items"`
	Total int           `json:"total"`
}

func (c *Client) ListIntegrations(ctx context.Context) ([]Integration, error) {
	var result IntegrationList
	if err := c.Get(ctx, "/api/v1/integrations", &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (c *Client) GetIntegration(ctx context.Context, id int) (*Integration, error) {
	var result Integration
	if err := c.Get(ctx, fmt.Sprintf("/api/v1/integrations/%d", id), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateIntegration(ctx context.Context, input *Integration) (*Integration, error) {
	var result Integration
	if err := c.Post(ctx, "/api/v1/integrations", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateIntegration(ctx context.Context, id int, input *Integration) (*Integration, error) {
	var result Integration
	if err := c.Put(ctx, fmt.Sprintf("/api/v1/integrations/%d", id), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteIntegration(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("/api/v1/integrations/%d", id))
}
