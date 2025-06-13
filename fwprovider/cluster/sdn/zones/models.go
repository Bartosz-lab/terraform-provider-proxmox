/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package sdn_zones

import (
	"context"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// INFO: Proxmox API provides additional type of SDN zone: "faucet"
// But I don't see any documentation about it.
//
// There are also "bridge-disable-mac-learning" and "dp-id" attributes
// But I also don't know how to use them.

type sdnZoneResourceModel struct {
	// Base attributes
	Name       types.String        `tfsdk:"name"`
	MTU        types.Int32         `tfsdk:"mtu"`
	Nodes      types.List          `tfsdk:"nodes"`
	IPAM       types.String        `tfsdk:"ipam"`
	DNS        types.String        `tfsdk:"dns"`
	ReverseDNS types.String        `tfsdk:"reversedns"`
	DNSZone    types.String        `tfsdk:"dnszone"`
	Simple     *sdnZoneSimpleModel `tfsdk:"simple"`
	VLAN       *sdnZoneVlanModel   `tfsdk:"vlan"`
	VXLAN      *sdnZoneVxlanModel  `tfsdk:"vxlan"`
	QinQ       *sdnZoneQinQModel   `tfsdk:"qinq"`
	EVPN       *sdnZoneEvpnModel   `tfsdk:"evpn"`
}

type sdnZoneSimpleModel struct {
	AutomaticDHCP types.String `tfsdk:"dhcp"`
}

type sdnZoneVlanModel struct {
	Bridge types.String `tfsdk:"bridge"`
}

type sdnZoneVxlanModel struct {
	Peers types.List  `tfsdk:"peers"`
	Port  types.Int32 `tfsdk:"port"`
}

type sdnZoneQinQModel struct {
	Bridge       types.String `tfsdk:"bridge"`
	Tag          types.Int32  `tfsdk:"tag"`
	VlanProtocol types.String `tfsdk:"vlan_protocol"`
}

type sdnZoneEvpnModel struct {
	Controller              types.String `tfsdk:"controller"`
	VrfVxlan                types.Int32  `tfsdk:"vrf_vxlan"`
	Mac                     types.String `tfsdk:"mac"`
	Exitnodes               types.List   `tfsdk:"exitnodes"`
	ExitnodesPrimary        types.String `tfsdk:"exitnodes_primary"`
	ExitnodesLocalRouting   types.Bool   `tfsdk:"exitnodes_local_routing"`
	AdvertiseSubnets        types.Bool   `tfsdk:"advertise_subnets"`
	DisableArpNdSuppression types.Bool   `tfsdk:"disable_arp_nd_suppression"`
	RtImport                types.String `tfsdk:"rt_import"`
}

// RemoveAllAttributes resets all attributes except the name.
func (m *sdnZoneResourceModel) RemoveAllAttributes() {
	*m = sdnZoneResourceModel{
		Name:  m.Name,
		Nodes: types.ListNull(types.StringType),
	}
}

// exportToSdnZoneBody converts the resource model to a SDN zone body for API requests.
func (m *sdnZoneResourceModel) exportToSdnZoneBody(ctx context.Context, diags *diag.Diagnostics) *zones.SdnZoneBody {
	result := &zones.SdnZoneBody{
		Name:       m.Name.ValueString(),
		Mtu:        m.MTU.ValueInt32Pointer(),
		Nodes:      convertListToString(m.Nodes, ctx, diags),
		Ipam:       m.IPAM.ValueStringPointer(),
		Dns:        m.DNS.ValueStringPointer(),
		Reversedns: m.ReverseDNS.ValueStringPointer(),
		Dnszone:    m.DNSZone.ValueStringPointer(),
	}

	var zoneType string
	if m.Simple != nil {
		zoneType = "simple"
		result.Dhcp = m.Simple.AutomaticDHCP.ValueStringPointer()

	} else if m.VLAN != nil {
		zoneType = "vlan"
		result.Bridge = m.VLAN.Bridge.ValueStringPointer()

	} else if m.VXLAN != nil {
		zoneType = "vxlan"
		result.Peers = convertListToString(m.VXLAN.Peers, ctx, diags)
		result.VxlanPort = m.VXLAN.Port.ValueInt32Pointer()

	} else if m.QinQ != nil {
		zoneType = "qinq"
		result.Bridge = m.QinQ.Bridge.ValueStringPointer()
		result.Tag = m.QinQ.Tag.ValueInt32Pointer()
		result.VlanProtocol = m.QinQ.VlanProtocol.ValueStringPointer()

	} else if m.EVPN != nil {
		zoneType = "evpn"
		result.Controller = m.EVPN.Controller.ValueStringPointer()
		result.VrfVxlan = m.EVPN.VrfVxlan.ValueInt32Pointer()
		result.Mac = m.EVPN.Mac.ValueStringPointer()
		result.Exitnodes = convertListToString(m.EVPN.Exitnodes, ctx, diags)
		result.ExitnodesPrimary = m.EVPN.ExitnodesPrimary.ValueStringPointer()
		result.ExitnodesLocalRouting = m.EVPN.ExitnodesLocalRouting.ValueBoolPointer()
		result.AdvertiseSubnets = m.EVPN.AdvertiseSubnets.ValueBoolPointer()
		result.DisableArpNdSuppression = m.EVPN.DisableArpNdSuppression.ValueBoolPointer()
		result.RtImport = m.EVPN.RtImport.ValueStringPointer()
	}

	result.Type = &zoneType

	return result
}

// importFromSdnZoneBody populates the resource model from a SDN zone body.
func (m *sdnZoneResourceModel) importFromSdnZoneBody(ctx context.Context, body *zones.SdnZoneBody, diags *diag.Diagnostics) {
	m.Name = types.StringValue(body.Name)
	m.MTU = types.Int32PointerValue(body.Mtu)
	m.Nodes = convertStringToList(body.Nodes, ctx, diags)
	m.IPAM = types.StringPointerValue(body.Ipam)
	m.DNS = types.StringPointerValue(body.Dns)
	m.ReverseDNS = types.StringPointerValue(body.Reversedns)
	m.DNSZone = types.StringPointerValue(body.Dnszone)

	switch *body.Type {
	case "simple":
		m.Simple = &sdnZoneSimpleModel{
			AutomaticDHCP: types.StringPointerValue(body.Dhcp),
		}
	case "vlan":
		m.VLAN = &sdnZoneVlanModel{
			Bridge: types.StringPointerValue(body.Bridge),
		}
	case "vxlan":
		m.VXLAN = &sdnZoneVxlanModel{
			Peers: convertStringToList(body.Peers, ctx, diags),
			Port:  types.Int32PointerValue(body.VxlanPort),
		}
	case "qinq":
		m.QinQ = &sdnZoneQinQModel{
			Bridge:       types.StringPointerValue(body.Bridge),
			Tag:          types.Int32PointerValue(body.Tag),
			VlanProtocol: types.StringPointerValue(body.VlanProtocol),
		}
	case "evpn":
		m.EVPN = &sdnZoneEvpnModel{
			Controller:              types.StringPointerValue(body.Controller),
			VrfVxlan:                types.Int32PointerValue(body.VrfVxlan),
			Mac:                     types.StringPointerValue(body.Mac),
			Exitnodes:               convertStringToList(body.Exitnodes, ctx, diags),
			ExitnodesPrimary:        types.StringPointerValue(body.ExitnodesPrimary),
			ExitnodesLocalRouting:   types.BoolPointerValue(body.ExitnodesLocalRouting),
			AdvertiseSubnets:        types.BoolPointerValue(body.AdvertiseSubnets),
			DisableArpNdSuppression: types.BoolPointerValue(body.DisableArpNdSuppression),
			RtImport:                types.StringPointerValue(body.RtImport),
		}
	default:
		diags.AddError(
			"Invalid SDN Zone Type",
			"SDN Zone type is not recognized: "+*body.Type,
		)
		return
	}
}

func (s *sdnZoneResourceModel) exportToUpdateBody(ctx context.Context, diags *diag.Diagnostics) *zones.SdnZoneBody {
	body := s.exportToSdnZoneBody(ctx, diags)

	// Add to delete_tab any fields that are unset in the request body.
	var deleteTab []string

	if body.Mtu == nil {
		deleteTab = append(deleteTab, "mtu")
	}
	if body.Nodes == nil {
		deleteTab = append(deleteTab, "nodes")
	}
	if body.Ipam == nil {
		deleteTab = append(deleteTab, "ipam")
	}
	if body.Dns == nil {
		deleteTab = append(deleteTab, "dns")
	}
	if body.Reversedns == nil {
		deleteTab = append(deleteTab, "reversedns")
	}
	if body.Dnszone == nil {
		deleteTab = append(deleteTab, "dnszone")
	}

	switch *body.Type {
	case "simple":
		if body.Dhcp == nil {
			deleteTab = append(deleteTab, "dhcp")
		}
	case "vlan":
		if body.Bridge == nil {
			deleteTab = append(deleteTab, "bridge")
		}
	case "vxlan":
		if body.Peers == nil {
			deleteTab = append(deleteTab, "peers")
		}
		if body.VxlanPort == nil {
			deleteTab = append(deleteTab, "vxlan-port")
		}
	case "qinq":
		if body.Bridge == nil {
			deleteTab = append(deleteTab, "bridge")
		}
		if body.Tag == nil {
			deleteTab = append(deleteTab, "tag")
		}
		if body.VlanProtocol == nil {
			deleteTab = append(deleteTab, "vlan-protocol")
		}
	case "evpn":
		if body.Controller == nil {
			deleteTab = append(deleteTab, "controller")
		}
		if body.VrfVxlan == nil {
			deleteTab = append(deleteTab, "vrf-vxlan")
		}
		if body.Mac == nil {
			deleteTab = append(deleteTab, "mac")
		}
		if body.Exitnodes == nil {
			deleteTab = append(deleteTab, "exitnodes")
		}
		if body.ExitnodesPrimary == nil {
			deleteTab = append(deleteTab, "exitnodes-primary")
		}
		if body.ExitnodesLocalRouting == nil {
			deleteTab = append(deleteTab, "exitnodes-local-routing")
		}
		if body.AdvertiseSubnets == nil {
			deleteTab = append(deleteTab, "advertise-subnets")
		}
		if body.DisableArpNdSuppression == nil {
			deleteTab = append(deleteTab, "disable-arp-nd-suppression")
		}
		if body.RtImport == nil {
			deleteTab = append(deleteTab, "rt-import")
		}
	}

	if len(deleteTab) > 0 {
		toDelete := strings.Join(deleteTab, ",")
		body.Delete = &toDelete
	}

	// Update requests don't accept the "type" field, so we remove it if present.
	body.Type = nil

	return body
}

// convertListToString converts a Terraform list to a comma-separated string.
func convertListToString(list types.List, ctx context.Context, diags *diag.Diagnostics) *string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	strs := make([]types.String, 0, len(list.Elements()))
	nodes_diags := list.ElementsAs(ctx, &strs, false)
	diags.Append(nodes_diags...)

	stringVals := make([]string, len(strs))
	for i, v := range strs {
		stringVals[i] = v.ValueString()
	}
	joined := strings.Join(stringVals, ",")
	return &joined
}

// convertStringToList converts a comma-separated string to a Terraform list.
func convertStringToList(value *string, ctx context.Context, diags *diag.Diagnostics) types.List {
	if value == nil || *value == "" {
		return types.ListNull(types.StringType)
	}

	parts := strings.Split(*value, ",")
	list, listDiags := types.ListValueFrom(ctx, types.StringType, parts)
	diags.Append(listDiags...)

	return list
}
