/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package sdn

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

// Client is an interface for accessing the Proxmox ACME API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to a full cluster ACME API path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/sdn/%s", path)
}

// Zones returns a client for managing the cluster's SDN zones.
func (c *Client) Zones() *zones.Client {
	return &zones.Client{Client: c.Client}
}
