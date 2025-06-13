/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package sdn_zones

import (
	"context"
	"fmt"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &sdnZoneResource{}
	_ resource.ResourceWithConfigure = &sdnZoneResource{}
)

// NewSdnZoneResource creates a new instance of the sdn zone resource.
// It is a helper function to simplify the provider implementation.
func NewSdnZoneResource() resource.Resource {
	return &sdnZoneResource{}
}

type sdnZoneResource struct {
	client proxmox.Client
}

// Metadata returns the resource type name.
func (r *sdnZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdn_zone"
}

// Schema defines the schema for the resource.
func (r *sdnZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	recreatemodifier := objectplanmodifier.RequiresReplaceIf(
		func(ctx context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifier.RequiresReplaceIfFuncResponse) {
			if req.StateValue.IsNull() != req.PlanValue.IsNull() {
				resp.RequiresReplace = true
			}
		},
		"Changes of the SDN zone type require a resource replacement",
		"Changes of the SDN zone type require a resource replacement",
	)

	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox SDN zone.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the SDN zone.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mtu": schema.Int32Attribute{
				Description: "MTU",
				Optional:    true,
			},
			"nodes": schema.ListAttribute{
				Description: "List of nodes that are part of the SDN zone.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"ipam": schema.StringAttribute{
				Description: "IPAM name",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("pve"),
			},
			"dns": schema.StringAttribute{
				Description: "DNS api server",
				Optional:    true,
			},
			"reversedns": schema.StringAttribute{
				Description: "Reverse DNS api server",
				Optional:    true,
			},
			"dnszone": schema.StringAttribute{
				Description: "DNS zone name",
				Optional:    true,
			},
			"simple": schema.SingleNestedAttribute{
				Description: "Simple SDN zone configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"dhcp": schema.StringAttribute{
						Description: "Enable automatic DHCP.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("dnsmasq"),
						},
					},
				},
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRoot("simple"),
						path.MatchRoot("vlan"),
						path.MatchRoot("vxlan"),
						path.MatchRoot("qinq"),
						path.MatchRoot("evpn"),
					),
				},
				PlanModifiers: []planmodifier.Object{
					recreatemodifier,
				},
			},
			"vlan": schema.SingleNestedAttribute{
				Description: "VLAN SDN zone configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"bridge": schema.StringAttribute{
						Description: "Bridge to use for the VLAN zone.",
						Required:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					recreatemodifier,
				},
			},
			"vxlan": schema.SingleNestedAttribute{
				Description: "VXLAN SDN zone configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"peers": schema.ListAttribute{
						Description: "List of peer nodes for the VXLAN zone.",
						Required:    true,
						ElementType: types.StringType,
					},
					"port": schema.Int32Attribute{
						Description: "Vxlan tunnel udp port.",
						Optional:    true,
						Computed:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					recreatemodifier,
				},
			},
			"qinq": schema.SingleNestedAttribute{
				Description: "QinQ SDN zone configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"bridge": schema.StringAttribute{
						Description: "Bridge to use for the QinQ zone.",
						Required:    true,
					},
					"tag": schema.Int32Attribute{
						Description: "VLAN tag for the QinQ zone.",
						Required:    true,
					},
					"vlan_protocol": schema.StringAttribute{
						Description: "VLAN protocol for the QinQ zone.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("802.1q"),
						Validators: []validator.String{
							stringvalidator.OneOf("802.1q", "802.1ad"),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					recreatemodifier,
				},
			},
			"evpn": schema.SingleNestedAttribute{
				Description: "EVPN SDN zone configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"controller": schema.StringAttribute{
						Description: "EVPN controller address.",
						Required:    true,
					},
					"vfr_vxlan": schema.Int32Attribute{
						Description: "VRF VXLAN ID for the EVPN zone.",
						Required:    true,
					},
					"mac": schema.StringAttribute{
						Description: "Anycast logical router mac address.",
						Optional:    true,
						Computed:    true,
					},
					"exitnodes": schema.ListAttribute{
						Description: "List of exit nodes for the EVPN zone.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
					},
					"exitnodes_primary": schema.StringAttribute{
						Description: "Primary exit node for the EVPN zone.",
						Optional:    true,
						Computed:    true,
					},
					"exitnodes_local_routing": schema.BoolAttribute{
						Description: "Enable local routing for exit nodes.",
						Optional:    true,
						Computed:    true,
					},
					"advertise_subnets": schema.BoolAttribute{
						Description: "Advertise subnets to exit nodes.",
						Optional:    true,
						Computed:    true,
					},
					"disable_arp_nd_suppression": schema.BoolAttribute{
						Description: "Disable ipv4 arp && ipv6 neighbour discovery suppression",
						Optional:    true,
						Computed:    true,
					},
					"rt_import": schema.StringAttribute{
						Description: "Route target import.",
						Optional:    true,
						Computed:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					recreatemodifier,
				},
			},
		},
	}
}

func (r *sdnZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource but got: %T", req.ProviderData),
		)
		return
	}

	r.client = cfg.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *sdnZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sdnZoneResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Cluster().SDN().Zones().Create(ctx, plan.exportToSdnZoneBody(ctx, &resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SDN Zone",
			fmt.Sprintf("Failed to list SDN zones: %s", err),
		)
		return
	}

	r.read(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// read fetches the current state of the resource from the Proxmox API and updates the model.
func (r *sdnZoneResource) read(ctx context.Context, model *sdnZoneResourceModel, diags *diag.Diagnostics) {
	zone, err := r.client.Cluster().SDN().Zones().Get(ctx, model.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			diags.AddWarning(
				"SDN Zone Not Found",
				fmt.Sprintf("SDN zone %s does not exist, setting to empty state", model.Name.ValueString()),
			)
			model.RemoveAllAttributes()
		} else {
			diags.AddError(
				"Error listing SDN Zones",
				fmt.Sprintf("Failed to list SDN zones: %s", err),
			)
		}
		return
	}

	model.importFromSdnZoneBody(ctx, zone, diags)
}

// Read refreshes the Terraform state with the latest data.
func (r *sdnZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sdnZoneResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.read(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sdnZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sdnZoneResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Cluster().SDN().Zones().Update(ctx, plan.Name.ValueString(), plan.exportToUpdateBody(ctx, &resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SDN Zone",
			fmt.Sprintf("Failed to update SDN zone %s: %s", plan.Name.ValueString(), err),
		)
		return
	}

	r.read(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sdnZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sdnZoneResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Cluster().SDN().Zones().Delete(ctx, state.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			resp.Diagnostics.AddWarning(
				"SDN Zone Not Found",
				fmt.Sprintf("SDN zone %s does not exist, skipping deletion", state.Name.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Deleting SDN Zone",
				fmt.Sprintf("Failed to delete SDN zone %s: %s", state.Name.ValueString(), err),
			)
		}
		return
	}
}
