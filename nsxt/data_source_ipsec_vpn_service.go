/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNsxtPolicyIPSecVpnService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIPSecVpnServiceRead,

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"gateway_path": getPolicyGatewayPathSchema(),
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enable/Disable IPSec VPN service.",
				Optional:    true,
				Default:     true,
			},
			"ha_sync": {
				Type:        schema.TypeBool,
				Description: "Enable/Disable IPSec VPN service HA state sync.",
				Optional:    true,
				Default:     true,
			},
			"ike_loglevel": {
				Type:         schema.TypeString,
				Description:  "Log level for internet key exchange (IKE).",
				ValidateFunc: validation.StringInSlice(IPSecVpnServiceIkeLogLevelTypes, false),
				Optional:     true,
			},
			"bypass_rules": getIPSecVPNRulesSchema(),
		},
	}
}

func dataSourceNsxtPolicyIPSecVpnServiceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	_, err := policyDataSourceResourceRead(d, connector, isPolicyGlobalManager(m), "IPSecVpnService", nil)
	if err != nil {
		return err
	}

	return nil
}
