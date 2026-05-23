package client

import (
	"context"
	"fmt"
)

type NotificationRule struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	IsEnabled   bool     `json:"is_enabled"`
	Severities  []string `json:"severities,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	ChannelType string   `json:"channel_type"`
	ChannelID   string   `json:"channel_id,omitempty"`
}

type NotificationRuleList struct {
	Items []NotificationRule `json:"items"`
	Total int                `json:"total"`
}

type Webhook struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	IsEnabled bool   `json:"is_enabled"`
	Secret    string `json:"secret,omitempty"`
	Events    []string `json:"events,omitempty"`
}

type WebhookList struct {
	Items []Webhook `json:"items"`
	Total int       `json:"total"`
}

func (c *Client) ListNotificationRules(ctx context.Context) ([]NotificationRule, error) {
	var result NotificationRuleList
	if err := c.Get(ctx, "/api/v1/notification-rules", &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (c *Client) GetNotificationRule(ctx context.Context, id int) (*NotificationRule, error) {
	var result NotificationRule
	if err := c.Get(ctx, fmt.Sprintf("/api/v1/notification-rules/%d", id), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateNotificationRule(ctx context.Context, input *NotificationRule) (*NotificationRule, error) {
	var result NotificationRule
	if err := c.Post(ctx, "/api/v1/notification-rules", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateNotificationRule(ctx context.Context, id int, input *NotificationRule) (*NotificationRule, error) {
	var result NotificationRule
	if err := c.Put(ctx, fmt.Sprintf("/api/v1/notification-rules/%d", id), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteNotificationRule(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("/api/v1/notification-rules/%d", id))
}

func (c *Client) ListWebhooks(ctx context.Context) ([]Webhook, error) {
	var result WebhookList
	if err := c.Get(ctx, "/api/v1/webhooks", &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (c *Client) GetWebhook(ctx context.Context, id int) (*Webhook, error) {
	var result Webhook
	if err := c.Get(ctx, fmt.Sprintf("/api/v1/webhooks/%d", id), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateWebhook(ctx context.Context, input *Webhook) (*Webhook, error) {
	var result Webhook
	if err := c.Post(ctx, "/api/v1/webhooks", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateWebhook(ctx context.Context, id int, input *Webhook) (*Webhook, error) {
	var result Webhook
	if err := c.Put(ctx, fmt.Sprintf("/api/v1/webhooks/%d", id), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteWebhook(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("/api/v1/webhooks/%d", id))
}
