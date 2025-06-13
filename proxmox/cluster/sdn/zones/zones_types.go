/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zones

// SdnZonesListResponseBody contains the body from a SDN zones list response.
type SdnZoneListResponseBody struct {
	Data []*SdnZoneBody `json:"data,omitempty"`
}

// SdnZoneGetResponseData contains the data from a SDN zone get response.
type SdnZoneGetResponseBody struct {
	Data *SdnZoneBody `json:"data,omitempty"`
}

// SdnZoneBody represents the body of a SDN zone in Proxmox.
// Documented in: https://pve.proxmox.com/pve-docs/api-viewer/#/cluster/sdn/zones
type SdnZoneBody struct {
	Name string `json:"zone" url:"zone"`

	Type                     *string `json:"type,omitempty" url:"type,omitempty"`     // Should be omitted only with update requests.
	Delete                   *string `json:"delete,omitempty" url:"delete,omitempty"` // Should be used only with update requests.
	AdvertiseSubnets         *bool   `json:"advertise-subnets,omitempty" url:"advertise-subnets,omitempty"`
	Bridge                   *string `json:"bridge,omitempty" url:"bridge,omitempty"`
	BridgeDisableMacLearning *bool   `json:"bridge-disable-mac-learning,omitempty" url:"bridge-disable-mac-learning,omitempty"`
	Controller               *string `json:"controller,omitempty" url:"controller,omitempty"`
	Dhcp                     *string `json:"dhcp,omitempty" url:"dhcp,omitempty"`
	DisableArpNdSuppression  *bool   `json:"disable-arp-nd-suppression,omitempty" url:"disable-arp-nd-suppression,omitempty"`
	Dns                      *string `json:"dns,omitempty" url:"dns,omitempty"`
	Dnszone                  *string `json:"dnszone,omitempty" url:"dnszone,omitempty"`
	DpID                     *int32  `json:"dp-id,omitempty" url:"dp-id,omitempty"`
	Exitnodes                *string `json:"exitnodes,omitempty" url:"exitnodes,omitempty"`
	ExitnodesLocalRouting    *bool   `json:"exitnodes-local-routing,omitempty" url:"exitnodes-local-routing,omitempty"`
	ExitnodesPrimary         *string `json:"exitnodes-primary,omitempty" url:"exitnodes-primary,omitempty"`
	Ipam                     *string `json:"ipam,omitempty" url:"ipam,omitempty"`
	Mac                      *string `json:"mac,omitempty" url:"mac,omitempty"`
	Mtu                      *int32  `json:"mtu,omitempty" url:"mtu,omitempty"`
	Nodes                    *string `json:"nodes,omitempty" url:"nodes,omitempty"`
	Peers                    *string `json:"peers,omitempty" url:"peers,omitempty"`
	Reversedns               *string `json:"reversedns,omitempty" url:"reversedns,omitempty"`
	RtImport                 *string `json:"rt-import,omitempty" url:"rt-import,omitempty"`
	Tag                      *int32  `json:"tag,omitempty" url:"tag,omitempty"`
	VlanProtocol             *string `json:"vlan-protocol,omitempty" url:"vlan-protocol,omitempty"`
	VrfVxlan                 *int32  `json:"vrf-vxlan,omitempty" url:"vrf-vxlan,omitempty"`
	VxlanPort                *int32  `json:"vxlan-port,omitempty" url:"vxlan-port,omitempty"`
}
