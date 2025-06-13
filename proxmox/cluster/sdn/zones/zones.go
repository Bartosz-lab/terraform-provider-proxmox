/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zones

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// List returns a list of SDN zones in the Proxmox cluster.
func (c *Client) List(ctx context.Context) ([]*SdnZoneBody, error) {
	resBody := &SdnZoneListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing SDN zones: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody.Data, nil
}

// Get retrieves a single SDN zone based on its identifier.
func (c *Client) Get(ctx context.Context, zone string) (*SdnZoneBody, error) {
	resBody := &SdnZoneGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(url.PathEscape(zone)), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading SDN zone: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Create creates a new SDN zone.
func (c *Client) Create(ctx context.Context, data *SdnZoneBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), data, nil)
	if err != nil {
		return fmt.Errorf("error creating SDN zone: %w", err)
	}

	return nil
}

// Update updates an existing SDN zone.
func (c *Client) Update(ctx context.Context, zone string, data *SdnZoneBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(url.PathEscape(zone)), data, nil)
	if err != nil {
		return fmt.Errorf("error updating SDN zone: %w", err)
	}

	return nil
}

// Delete removes an SDN zone.
func (c *Client) Delete(ctx context.Context, zone string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(url.PathEscape(zone)), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting SDN zone: %w", err)
	}

	return nil
}
