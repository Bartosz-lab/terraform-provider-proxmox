/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zones

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is an interface for accessing the Proxmox SDN zones management API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to a full cluster SDN zones API path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/sdn/zones/%s", path)
}
