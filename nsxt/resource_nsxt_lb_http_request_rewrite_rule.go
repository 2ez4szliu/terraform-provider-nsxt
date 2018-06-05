/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/loadbalancer"
	"log"
	"net/http"
)

func resourceNsxtLbHTTPRequestRewriteRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLbHTTPRequestRewriteRuleCreate,
		Read:   resourceNsxtLbHTTPRequestRewriteRuleRead,
		Update: resourceNsxtLbHTTPRequestRewriteRuleUpdate,
		Delete: resourceNsxtLbHTTPRequestRewriteRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"match_strategy": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Strategy when multiple match conditions are specified in one rule (ANY vs ALL)",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL", "ANY"}, false),
				Default:      "ALL",
			},
			"header_condition":        getLbRuleHTTPHeaderConditionSchema(),
			"body_condition":          getLbRuleHTTPRequestBodyConditionSchema(),
			"method_condition":        getLbRuleHTTPRequestMethodConditionSchema(),
			"cookie_condition":        getLbRuleHTTPHeaderConditionSchema(),
			"version_condition":       getLbRuleHTTPVersionConditionSchema(),
			"uri_condition":           getLbRuleHTTPRequestURIConditionSchema(),
			"uri_arguments_condition": getLbRuleHTTPRequestURIArgumentsConditionSchema(),
			"ip_condition":            getLbRuleIPConditionSchema(),
			"tcp_condition":           getLbRuleTCPConditionSchema(),

			"uri_rewrite_action":    getLbRuleURIRewriteActionSchema(),
			"header_rewrite_action": getLbRuleHeaderRewriteActionSchema(),
		},
	}
}

func getLbRuleInverseSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Whether to reverse match result of this condition",
		Optional:    true,
		Default:     false,
	}
}

func getLbRuleCaseSensitiveSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Description: "If true, case is significant in condition matching",
		Optional:    true,
		Default:     false,
	}
}

func getLbRuleMatchTypeSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Match type (STARTS_WITH, ENDS_WITH, EQUALS, CONTAINS, REGEX)",
		ValidateFunc: validation.StringInSlice([]string{"STARTS_WITH", "ENDS_WITH", "EQUALS", "CONTAINS", "REGEX"}, false),
		Required:     true,
	}
}

func getLbRuleHTTPRequestBodyConditionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Rule condition based on http request body",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"inverse": getLbRuleInverseSchema(),
				"value": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"case_sensitive": getLbRuleCaseSensitiveSchema(),
				"match_type":     getLbRuleMatchTypeSchema(),
			},
		},
	}
}

func getLbRuleHTTPHeaderConditionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Rule condition based on http header",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"inverse": getLbRuleInverseSchema(),
				"name": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"value": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"case_sensitive": getLbRuleCaseSensitiveSchema(),
				"match_type":     getLbRuleMatchTypeSchema(),
			},
		},
	}
}

func getLbRuleHTTPRequestMethodConditionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Rule condition based on http request method",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"inverse": getLbRuleInverseSchema(),
				"method": &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"GET", "OPTIONS", "POST", "HEAD", "PUT"}, false),
				},
			},
		},
	}
}

func getLbRuleHTTPVersionConditionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Rule condition based on http request version",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"inverse": getLbRuleInverseSchema(),
				"version": &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"HTTP_VERSION_1_0", "HTTP_VERSION_1_1", "HTTP_VERSION_2_0"}, false),
				},
			},
		},
	}
}

func getLbRuleHTTPRequestURIConditionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Rule condition based on http request URI",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"inverse": getLbRuleInverseSchema(),
				"uri": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func getLbRuleHTTPRequestURIArgumentsConditionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Rule condition based on http request URI arguments (query string)",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"inverse": getLbRuleInverseSchema(),
				"uri_arguments": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func getLbRuleIPConditionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Rule condition based on request IP settings",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"inverse": getLbRuleInverseSchema(),
				"source_address": &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validateSingleIP(),
				},
			},
		},
	}
}

func getLbRuleTCPConditionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Rule condition based on request TCP settings",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"inverse": getLbRuleInverseSchema(),
				"source_port": &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validateSinglePort(),
				},
			},
		},
	}
}

func getLbRuleURIRewriteActionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Uri to replace original URI in outgoing request",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"uri": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"uri_arguments": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func getLbRuleHeaderRewriteActionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Header to replace original header in outgoing request",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"value": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func initLbHTTPRuleMatchCondition(data map[string]interface{}, conditionType string) loadbalancer.LbRuleCondition {
	condition := loadbalancer.LbRuleCondition{
		Inverse:   data["inverse"].(bool),
		Type_:     conditionType,
		MatchType: data["match_type"].(string),
	}

	condition.CaseSensitive = new(bool)
	*condition.CaseSensitive = data["case_sensitive"].(bool)
	return condition
}

func getLbRuleHTTPRequestConditionsFromSchema(d *schema.ResourceData) []loadbalancer.LbRuleCondition {
	var conditionList []loadbalancer.LbRuleCondition
	conditions := d.Get("header_condition").(*schema.Set).List()
	for _, condition := range conditions {
		data := condition.(map[string]interface{})
		elem := initLbHTTPRuleMatchCondition(data, "LbHTTPRequestHeaderCondition")
		elem.HeaderName = data["name"].(string)
		elem.HeaderValue = data["value"].(string)

		conditionList = append(conditionList, elem)
	}

	conditions = d.Get("cookie_condition").(*schema.Set).List()
	for _, condition := range conditions {
		data := condition.(map[string]interface{})
		elem := initLbHTTPRuleMatchCondition(data, "LbHTTPRequestCookieCondition")
		elem.CookieName = data["name"].(string)
		elem.CookieValue = data["value"].(string)

		conditionList = append(conditionList, elem)
	}

	conditions = d.Get("body_condition").(*schema.Set).List()
	for _, condition := range conditions {
		data := condition.(map[string]interface{})
		elem := initLbHTTPRuleMatchCondition(data, "LbHTTPRequestBodyCondition")
		elem.BodyValue = data["value"].(string)

		conditionList = append(conditionList, elem)
	}

	conditions = d.Get("method_condition").(*schema.Set).List()
	for _, condition := range conditions {
		data := condition.(map[string]interface{})
		elem := loadbalancer.LbRuleCondition{
			Inverse: data["inverse"].(bool),
			Type_:   "LbHTTPRequestMethodCondition",
			Method:  data["method"].(string),
		}

		conditionList = append(conditionList, elem)
	}

	conditions = d.Get("version_condition").(*schema.Set).List()
	for _, condition := range conditions {
		data := condition.(map[string]interface{})
		elem := loadbalancer.LbRuleCondition{
			Inverse: data["inverse"].(bool),
			Type_:   "LbHTTPRequestVersionCondition",
			Version: data["version"].(string),
		}

		conditionList = append(conditionList, elem)
	}

	conditions = d.Get("uri_condition").(*schema.Set).List()
	for _, condition := range conditions {
		data := condition.(map[string]interface{})
		elem := initLbHTTPRuleMatchCondition(data, "LbHTTPRequestUriCondition")
		elem.Uri = data["uri"].(string)

		conditionList = append(conditionList, elem)
	}

	conditions = d.Get("uri_arguments").(*schema.Set).List()
	for _, condition := range conditions {
		data := condition.(map[string]interface{})
		elem := initLbHTTPRuleMatchCondition(data, "LbHTTPRequestUriArgumentsCondition")
		elem.UriArguments = data["uri_arguments"].(string)

		conditionList = append(conditionList, elem)
	}

	conditions = d.Get("ip_condition").(*schema.Set).List()
	for _, condition := range conditions {
		data := condition.(map[string]interface{})
		elem := loadbalancer.LbRuleCondition{
			Inverse:       data["inverse"].(bool),
			Type_:         "LbIpHeaderCondition",
			SourceAddress: data["source_address"].(string),
		}

		conditionList = append(conditionList, elem)
	}

	conditions = d.Get("tcp_condition").(*schema.Set).List()
	for _, condition := range conditions {
		data := condition.(map[string]interface{})
		elem := loadbalancer.LbRuleCondition{
			Inverse:    data["inverse"].(bool),
			Type_:      "LbTcpHeaderCondition",
			SourcePort: data["source_port"].(string),
		}

		conditionList = append(conditionList, elem)
	}

	return conditionList
}

func setLbRuleHTTPRequestConditionsInSchema(d *schema.ResourceData, conditions []loadbalancer.LbRuleCondition) error {
	var headerConditionList []map[string]interface{}
	var cookieConditionList []map[string]interface{}
	var bodyConditionList []map[string]interface{}
	var methodConditionList []map[string]interface{}
	var versionConditionList []map[string]interface{}
	var uriConditionList []map[string]interface{}
	var uriArgumentsConditionList []map[string]interface{}
	var ipConditionList []map[string]interface{}
	var tcpConditionList []map[string]interface{}

	for _, condition := range conditions {
		elem := make(map[string]interface{})

		if condition.Type_ == "LbHTTPRequestHeaderCondition" {
			elem["name"] = condition.HeaderName
			elem["value"] = condition.HeaderValue
			elem["inverse"] = condition.Inverse
			elem["match_type"] = condition.MatchType
			elem["case_sensitive"] = *condition.CaseSensitive
			headerConditionList = append(headerConditionList, elem)
		}

		if condition.Type_ == "LbHTTPRequestCookieCondition" {
			elem["name"] = condition.CookieName
			elem["value"] = condition.CookieValue
			elem["inverse"] = condition.Inverse
			elem["match_type"] = condition.MatchType
			elem["case_sensitive"] = *condition.CaseSensitive
			cookieConditionList = append(headerConditionList, elem)
		}

		if condition.Type_ == "LbHTTPRequestBodyCondition" {
			elem["value"] = condition.BodyValue
			elem["inverse"] = condition.Inverse
			elem["match_type"] = condition.MatchType
			elem["case_sensitive"] = *condition.CaseSensitive
			bodyConditionList = append(bodyConditionList, elem)
		}

		if condition.Type_ == "LbHTTPRequestMethodCondition" {
			elem["method"] = condition.Method
			elem["inverse"] = condition.Inverse
			methodConditionList = append(methodConditionList, elem)
		}

		if condition.Type_ == "LbHTTPRequestVersionCondition" {
			elem["version"] = condition.Version
			elem["inverse"] = condition.Inverse
			versionConditionList = append(versionConditionList, elem)
		}

		if condition.Type_ == "LbHTTPRequestUriCondition" {
			elem["uri"] = condition.Uri
			elem["inverse"] = condition.Inverse
			elem["match_type"] = condition.MatchType
			elem["case_sensitive"] = *condition.CaseSensitive
			uriConditionList = append(uriConditionList, elem)
		}

		if condition.Type_ == "LbHTTPRequestUriArgumentsCondition" {
			elem["uri_arguments"] = condition.UriArguments
			elem["inverse"] = condition.Inverse
			elem["match_type"] = condition.MatchType
			elem["case_sensitive"] = *condition.CaseSensitive
			uriArgumentsConditionList = append(uriArgumentsConditionList, elem)
		}

		if condition.Type_ == "LbIpHeaderCondition" {
			elem["source_address"] = condition.SourceAddress
			elem["inverse"] = condition.Inverse
			ipConditionList = append(ipConditionList, elem)
		}

		if condition.Type_ == "LbTcpHeaderCondition" {
			elem["source_port"] = condition.SourcePort
			elem["inverse"] = condition.Inverse
			tcpConditionList = append(tcpConditionList, elem)
		}

		err := d.Set("header_condition", headerConditionList)
		if err != nil {
			return err
		}

		err = d.Set("cookie_condition", cookieConditionList)
		if err != nil {
			return err
		}

		err = d.Set("body_condition", bodyConditionList)
		if err != nil {
			return err
		}

		err = d.Set("method_condition", methodConditionList)
		if err != nil {
			return err
		}

		err = d.Set("version_condition", versionConditionList)
		if err != nil {
			return err
		}

		err = d.Set("uri_condition", uriConditionList)
		if err != nil {
			return err
		}

		err = d.Set("uri_arguments_condition", uriArgumentsConditionList)
		if err != nil {
			return err
		}

		err = d.Set("ip_condition", ipConditionList)
		if err != nil {
			return err
		}

		err = d.Set("tcp_condition", tcpConditionList)
		if err != nil {
			return err
		}

	}
	return nil
}

func getLbRuleRewriteActionsFromSchema(d *schema.ResourceData) []loadbalancer.LbRuleAction {
	var actionList []loadbalancer.LbRuleAction
	actions := d.Get("header_rewrite_action").(*schema.Set).List()
	for _, action := range actions {
		data := action.(map[string]interface{})
		elem := loadbalancer.LbRuleAction{
			Type_:       "LbHTTPRequestHeaderRewriteAction",
			HeaderName:  data["name"].(string),
			HeaderValue: data["value"].(string),
		}

		actionList = append(actionList, elem)
	}

	actions = d.Get("uri_rewrite_action").(*schema.Set).List()
	for _, action := range actions {
		data := action.(map[string]interface{})
		elem := loadbalancer.LbRuleAction{
			Type_:        "LbHTTPRequestUriRewriteAction",
			Uri:          data["uri"].(string),
			UriArguments: data["uri_arguments"].(string),
		}

		actionList = append(actionList, elem)
	}

	return actionList
}

func setLbRuleRewriteActionsInSchema(d *schema.ResourceData, actions []loadbalancer.LbRuleAction) error {
	var uriActionList []map[string]string
	var headerActionList []map[string]string

	for _, action := range actions {
		elem := make(map[string]string)
		if action.Type_ == "LbHTTPRequestHeaderRewriteAction" {
			elem["name"] = action.HeaderName
			elem["value"] = action.HeaderValue
			headerActionList = append(headerActionList, elem)
		}

		if action.Type_ == "LbHTTPRequestUriRewriteAction" {
			elem["uri"] = action.Uri
			elem["uri_arguments"] = action.UriArguments
			uriActionList = append(uriActionList, elem)
		}

		err := d.Set("header_rewrite_action", headerActionList)
		if err != nil {
			return err
		}

		err = d.Set("uri_rewrite_action", uriActionList)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceNsxtLbHTTPRequestRewriteRuleCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	matchConditions := getLbRuleHTTPRequestConditionsFromSchema(d)
	actions := getLbRuleRewriteActionsFromSchema(d)
	matchStrategy := d.Get("match_strategy").(string)
	phase := "HTTP_REQUEST_REWRITE"

	lbRule := loadbalancer.LbRule{
		Description:     description,
		DisplayName:     displayName,
		Tags:            tags,
		Actions:         actions,
		MatchConditions: matchConditions,
		MatchStrategy:   matchStrategy,
		Phase:           phase,
	}

	lbRule, resp, err := nsxClient.ServicesApi.CreateLoadBalancerRule(nsxClient.Context, lbRule)

	if err != nil {
		return fmt.Errorf("Error during LoadBalancerRule create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during LoadBalancerRule create: %v", resp.StatusCode)
	}
	d.SetId(lbRule.Id)

	return resourceNsxtLbHTTPRequestRewriteRuleRead(d, m)
}

func resourceNsxtLbHTTPRequestRewriteRuleRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	lbRule, resp, err := nsxClient.ServicesApi.ReadLoadBalancerRule(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LoadBalancerRule read: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LoadBalancerRule %s not found", id)
		d.SetId("")
		return nil
	}
	d.Set("revision", lbRule.Revision)
	d.Set("description", lbRule.Description)
	d.Set("display_name", lbRule.DisplayName)
	setTagsInSchema(d, lbRule.Tags)
	setLbRuleHTTPRequestConditionsInSchema(d, lbRule.MatchConditions)
	setLbRuleRewriteActionsInSchema(d, lbRule.Actions)
	d.Set("match_strategy", lbRule.MatchStrategy)

	return nil
}

func resourceNsxtLbHTTPRequestRewriteRuleUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int32(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	matchConditions := getLbRuleHTTPRequestConditionsFromSchema(d)
	actions := getLbRuleRewriteActionsFromSchema(d)
	matchStrategy := d.Get("match_strategy").(string)
	phase := "HTTP_REQUEST_REWRITE"

	lbRule := loadbalancer.LbRule{
		Revision:        revision,
		Description:     description,
		DisplayName:     displayName,
		MatchStrategy:   matchStrategy,
		Phase:           phase,
		Actions:         actions,
		MatchConditions: matchConditions,
		Tags:            tags,
	}

	lbRule, resp, err := nsxClient.ServicesApi.UpdateLoadBalancerRule(nsxClient.Context, id, lbRule)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LoadBalancerRule update: %v", err)
	}

	return resourceNsxtLbHTTPRequestRewriteRuleRead(d, m)
}

func resourceNsxtLbHTTPRequestRewriteRuleDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.ServicesApi.DeleteLoadBalancerRule(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LoadBalancerRule delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LoadBalancerRule %s not found", id)
		d.SetId("")
	}
	return nil
}
