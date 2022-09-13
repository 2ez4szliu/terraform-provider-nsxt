/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var IPSecVpnServiceIkeLogLevelTypes = []string{
	model.IPSecVpnService_IKE_LOG_LEVEL_DEBUG,
	model.IPSecVpnService_IKE_LOG_LEVEL_INFO,
	model.IPSecVpnService_IKE_LOG_LEVEL_WARN,
	model.IPSecVpnService_IKE_LOG_LEVEL_ERROR,
	model.IPSecVpnService_IKE_LOG_LEVEL_EMERGENCY,
}

func resourceNsxtPolicyIPSecVpnService() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIPSecVpnServiceCreate,
		Read:   resourceNsxtPolicyIPSecVpnServiceRead,
		Update: resourceNsxtPolicyIPSecVpnServiceUpdate,
		Delete: resourceNsxtPolicyIPSecVpnServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

func getNsxtPolicyIPSecVpnServiceByID(connector *client.RestConnector, gwID string, isT0 bool, serviceID string, isGlobalManager bool) (model.IPSecVpnService, error) {
	if isT0 {
		client := tier_0s.NewIpsecVpnServicesClient(connector)
		return client.Get(gwID, serviceID)
	}
	client := tier_1s.NewIpsecVpnServicesClient(connector)
	return client.Get(gwID, serviceID)
}

func patchNsxtPolicyIPSecVpnService(connector *client.RestConnector, gwID string, ipSecVpnService model.IPSecVpnService, isT0 bool) error {
	id := *ipSecVpnService.Id
	if isT0 {
		client := tier_0s.NewIpsecVpnServicesClient(connector)
		return client.Patch(gwID, id, ipSecVpnService)
	}
	client := tier_1s.NewIpsecVpnServicesClient(connector)
	return client.Patch(gwID, id, ipSecVpnService)
}

func deleteNsxtPolicyIPSecVpnService(connector *client.RestConnector, gwID string, isT0 bool, id string) error {
	if isT0 {
		client := tier_0s.NewIpsecVpnServicesClient(connector)
		return client.Delete(gwID, id)
	}
	client := tier_1s.NewIpsecVpnServicesClient(connector)
	return client.Delete(gwID, id)
}

func resourceNsxtPolicyIPSecVpnServiceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnService ID")
	}
	gwPolicyPath := d.Get("gateway_path").(string)
	isT0, gwID := parseGatewayPolicyPath(gwPolicyPath)
	if gwID == "" {
		return fmt.Errorf("gateway_path is not valid")
	}
	obj, err := getNsxtPolicyIPSecVpnServiceByID(connector, gwID, isT0, id, isPolicyGlobalManager(m))
	if err != nil {
		return handleReadError(d, "IPSecVpnService", id, err)
	}
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("enabled", obj.Enabled)
	d.Set("ha_sync", obj.HaSync)

	if obj.BypassRules != nil {
		d.Set("bypass_rules", obj.BypassRules)
	}
	if obj.IkeLogLevel != nil {
		d.Set("ike_loglevel", obj.IkeLogLevel)
	}
	return nil
}

func resourceNsxtPolicyIPSecVpnServiceCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	gwPolicyPath := d.Get("gateway_path").(string)
	isT0, gwID := parseGatewayPolicyPath(gwPolicyPath)
	if gwID == "" {
		return fmt.Errorf("gateway_path is not valid")
	}
	isGlobalManager := isPolicyGlobalManager(m)
	id := d.Get("nsx_id").(string)
	if id == "" {
		id = newUUID()
	} else {
		_, err := getNsxtPolicyIPSecVpnServiceByID(connector, gwID, isT0, id, isGlobalManager)
		if err == nil {
			return fmt.Errorf("IPSecVpnService with nsx_id '%s' already exists", id)
		} else if !isNotFoundError(err) {
			return err
		}
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	ha_sync := d.Get("ha_sync").(bool)
	rules := getIPSecVPNRulesFromSchema(d)
	tags := getPolicyTagsFromSchema(d)

	ipSecVpnService := model.IPSecVpnService{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Enabled:     &enabled,
		HaSync:      &ha_sync,
		BypassRules: rules,
	}

	ike_loglevel := d.Get("ike_loglevel").(string)
	if ike_loglevel != "" {
		ipSecVpnService.IkeLogLevel = &ike_loglevel
	}

	log.Printf("[INFO] Creating IPSecVpnService with ID %s", id)
	err := patchNsxtPolicyIPSecVpnService(connector, gwID, ipSecVpnService, isT0)
	if err != nil {
		return handleCreateError("IPSecVpnService", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyIPSecVpnServiceRead(d, m)
}

func resourceNsxtPolicyIPSecVpnServiceUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSec VPN Service ID")
	}
	gwPolicyPath := d.Get("gateway_path").(string)
	isT0, gwID := parseGatewayPolicyPath(gwPolicyPath)
	if gwID == "" {
		return fmt.Errorf("gateway_path is not valid")
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	ha_sync := d.Get("ha_sync").(bool)
	rules := getIPSecVPNRulesFromSchema(d)
	tags := getPolicyTagsFromSchema(d)
	ipSecVpnService := model.IPSecVpnService{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Enabled:     &enabled,
		HaSync:      &ha_sync,
		BypassRules: rules,
	}

	ike_loglevel := d.Get("ike_loglevel").(string)
	if ike_loglevel != "" {
		ipSecVpnService.IkeLogLevel = &ike_loglevel
	}

	log.Printf("[INFO] Updating IPSecVpnService with ID %s", id)
	err := patchNsxtPolicyIPSecVpnService(connector, gwID, ipSecVpnService, isT0)
	if err != nil {
		return handleUpdateError("IPSecVpnService", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyIPSecVpnServiceRead(d, m)
}

func resourceNsxtPolicyIPSecVpnServiceDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSec VPN Service ID")
	}

	gwPolicyPath := d.Get("gateway_path").(string)
	isT0, gwID := parseGatewayPolicyPath(gwPolicyPath)
	if gwID == "" {
		return fmt.Errorf("gateway_path is not valid")
	}

	err := deleteNsxtPolicyIPSecVpnService(getPolicyConnector(m), gwID, isT0, id)
	if err != nil {
		return handleDeleteError("IPSecVpnService", id, err)
	}
	return nil
}
