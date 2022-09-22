/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_domains "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/domains"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicySecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicySecurityPolicyCreate,
		Read:   resourceNsxtPolicySecurityPolicyRead,
		Update: resourceNsxtPolicySecurityPolicyUpdate,
		Delete: resourceNsxtPolicySecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtDomainResourceImporter,
		},
		Schema: getPolicySecurityPolicySchema(false),
	}
}

func getSecurityPolicyInDomain(id string, domainName string, connector *client.RestConnector, isGlobalManager bool) (model.SecurityPolicy, error) {
	if isGlobalManager {
		client := gm_domains.NewSecurityPoliciesClient(connector)
		gmObj, err := client.Get(domainName, id)
		if err != nil {
			return model.SecurityPolicy{}, err
		}
		rawObj, convErr := convertModelBindingType(gmObj, gm_model.SecurityPolicyBindingType(), model.SecurityPolicyBindingType())
		if convErr != nil {
			return model.SecurityPolicy{}, convErr
		}
		return rawObj.(model.SecurityPolicy), nil
	}
	client := domains.NewSecurityPoliciesClient(connector)
	return client.Get(domainName, id)

}

func resourceNsxtPolicySecurityPolicyExistsInDomain(id string, domainName string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {
	var err error
	if isGlobalManager {
		client := gm_domains.NewSecurityPoliciesClient(connector)
		_, err = client.Get(domainName, id)
	} else {
		client := domains.NewSecurityPoliciesClient(connector)
		_, err = client.Get(domainName, id)
	}

	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving Security Policy", err)
}

func resourceNsxtPolicySecurityPolicyExistsPartial(domainName string) func(id string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {
	return func(id string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {
		return resourceNsxtPolicySecurityPolicyExistsInDomain(id, domainName, connector, isGlobalManager)
	}
}

func policySecurityPolicyBuildAndPatch(d *schema.ResourceData, m interface{}, connector *client.RestConnector, isGlobalManager bool, id string) error {

	domain := d.Get("domain").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	category := d.Get("category").(string)
	comments := d.Get("comments").(string)
	locked := d.Get("locked").(bool)
	scope := getStringListFromSchemaSet(d, "scope")
	sequenceNumber := int64(d.Get("sequence_number").(int))
	stateful := d.Get("stateful").(bool)
	tcpStrict := d.Get("tcp_strict").(bool)
	revision := int64(d.Get("revision").(int))
	objType := "SecurityPolicy"

	obj := model.SecurityPolicy{
		Id:             &id,
		DisplayName:    &displayName,
		Description:    &description,
		Tags:           tags,
		Category:       &category,
		Comments:       &comments,
		Locked:         &locked,
		Scope:          scope,
		SequenceNumber: &sequenceNumber,
		Stateful:       &stateful,
		TcpStrict:      &tcpStrict,
		ResourceType:   &objType,
	}
	log.Printf("[INFO] Creating Security Policy with ID %s", id)

	if len(d.Id()) > 0 {
		// This is update flow
		obj.Revision = &revision
	}

	policyChildren, err := getUpdatedRuleChildren(d)
	if err != nil {
		return err
	}
	if len(policyChildren) > 0 {
		obj.Children = policyChildren
	}

	return securityPolicyInfraPatch(obj, domain, m)
}

func resourceNsxtPolicySecurityPolicyCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicySecurityPolicyExistsPartial(d.Get("domain").(string)))
	if err != nil {
		return err
	}

	err = policySecurityPolicyBuildAndPatch(d, m, connector, isPolicyGlobalManager(m), id)

	if err != nil {
		return handleCreateError("Security Policy", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicySecurityPolicyRead(d, m)
}

func resourceNsxtPolicySecurityPolicyRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	domainName := d.Get("domain").(string)
	if id == "" {
		return fmt.Errorf("Error obtaining Security Policy id")
	}
	var obj model.SecurityPolicy
	if isPolicyGlobalManager(m) {
		client := gm_domains.NewSecurityPoliciesClient(connector)
		gmObj, err := client.Get(domainName, id)
		if err != nil {
			return handleReadError(d, "SecurityPolicy", id, err)
		}
		rawObj, err := convertModelBindingType(gmObj, gm_model.SecurityPolicyBindingType(), model.SecurityPolicyBindingType())
		if err != nil {
			return err
		}
		obj = rawObj.(model.SecurityPolicy)
	} else {
		var err error
		client := domains.NewSecurityPoliciesClient(connector)
		obj, err = client.Get(domainName, id)
		if err != nil {
			return handleReadError(d, "SecurityPolicy", id, err)
		}
	}
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("domain", getDomainFromResourcePath(*obj.Path))
	d.Set("category", obj.Category)
	d.Set("comments", obj.Comments)
	d.Set("locked", obj.Locked)
	if len(obj.Scope) == 1 && obj.Scope[0] == "ANY" {
		d.Set("scope", nil)
	} else {
		d.Set("scope", obj.Scope)
	}
	d.Set("sequence_number", obj.SequenceNumber)
	d.Set("stateful", obj.Stateful)
	d.Set("tcp_strict", obj.TcpStrict)
	d.Set("revision", obj.Revision)
	return setPolicyRulesInSchema(d, obj.Rules)
}

func resourceNsxtPolicySecurityPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Security Policy id")
	}
	err := policySecurityPolicyBuildAndPatch(d, m, connector, isPolicyGlobalManager(m), id)
	if err != nil {
		return handleUpdateError("Security Policy", id, err)
	}

	return resourceNsxtPolicySecurityPolicyRead(d, m)
}

func resourceNsxtPolicySecurityPolicyDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Security Policy id")
	}

	connector := getPolicyConnector(m)
	var err error

	if isPolicyGlobalManager(m) {
		client := gm_domains.NewSecurityPoliciesClient(connector)
		err = client.Delete(d.Get("domain").(string), id)
	} else {
		client := domains.NewSecurityPoliciesClient(connector)
		err = client.Delete(d.Get("domain").(string), id)
	}

	if err != nil {
		return handleDeleteError("Security Policy", id, err)
	}

	return nil
}
