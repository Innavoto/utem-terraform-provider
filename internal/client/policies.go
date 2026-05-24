package client

import (
	"context"
	"fmt"
	"net/url"
)

type Policy struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Category       string `json:"category"`
	Severity       string `json:"severity"`
	IsEnabled      bool   `json:"is_enabled"`
	ResourceType   string `json:"resource_type,omitempty"`
	RegoPolicy     string `json:"rego_policy,omitempty"`
	CustomerFacing bool   `json:"customer_facing"`
}

type PolicyList struct {
	Items []Policy `json:"items"`
	Total int      `json:"total"`
}

func (c *Client) ListPolicies(ctx context.Context, category string) ([]Policy, error) {
	path := "/api/v1/policies"
	if category != "" {
		params := url.Values{}
		params.Set("category", category)
		path += "?" + params.Encode()
	}
	var result PolicyList
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (c *Client) GetPolicy(ctx context.Context, id string) (*Policy, error) {
	var result Policy
	if err := c.Get(ctx, fmt.Sprintf("/api/v1/policies/%s", id), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreatePolicy(ctx context.Context, input *Policy) (*Policy, error) {
	var result Policy
	if err := c.Post(ctx, "/api/v1/policies", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdatePolicy(ctx context.Context, id string, input *Policy) (*Policy, error) {
	var result Policy
	if err := c.Put(ctx, fmt.Sprintf("/api/v1/policies/%s", id), input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeletePolicy(ctx context.Context, id string) error {
	return c.Delete(ctx, fmt.Sprintf("/api/v1/policies/%s", id))
}
