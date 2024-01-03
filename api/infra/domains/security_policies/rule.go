//nolint:revive
package securitypolicies

// The following file has been autogenerated. Please avoid any changes!
import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/domains/security_policies"
	model1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains/security_policies"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client2 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra/domains/security_policies"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RuleClientContext utl.ClientContext

func NewRulesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RuleClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewRulesClient(connector)

	case utl.Global:
		client = client1.NewRulesClient(connector)

	case utl.Multitenancy:
		client = client2.NewRulesClient(connector)

	default:
		return nil
	}
	return &RuleClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c RuleClientContext) Get(domainIdParam string, securityPolicyIdParam string, ruleIdParam string) (model0.Rule, error) {
	var obj model0.Rule
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RulesClient)
		obj, err = client.Get(domainIdParam, securityPolicyIdParam, ruleIdParam)
		if err != nil {
			return obj, err
		}

	case utl.Global:
		client := c.Client.(client1.RulesClient)
		gmObj, err1 := client.Get(domainIdParam, securityPolicyIdParam, ruleIdParam)
		if err1 != nil {
			return obj, err1
		}
		var rawObj interface{}
		rawObj, err = utl.ConvertModelBindingType(gmObj, model1.RuleBindingType(), model0.RuleBindingType())
		obj = rawObj.(model0.Rule)

	case utl.Multitenancy:
		client := c.Client.(client2.RulesClient)
		obj, err = client.Get(utl.DefaultOrgID, c.ProjectID, domainIdParam, securityPolicyIdParam, ruleIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c RuleClientContext) Delete(domainIdParam string, securityPolicyIdParam string, ruleIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RulesClient)
		err = client.Delete(domainIdParam, securityPolicyIdParam, ruleIdParam)

	case utl.Global:
		client := c.Client.(client1.RulesClient)
		err = client.Delete(domainIdParam, securityPolicyIdParam, ruleIdParam)

	case utl.Multitenancy:
		client := c.Client.(client2.RulesClient)
		err = client.Delete(utl.DefaultOrgID, c.ProjectID, domainIdParam, securityPolicyIdParam, ruleIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c RuleClientContext) Patch(domainIdParam string, securityPolicyIdParam string, ruleIdParam string, ruleParam model0.Rule) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RulesClient)
		err = client.Patch(domainIdParam, securityPolicyIdParam, ruleIdParam, ruleParam)

	case utl.Global:
		client := c.Client.(client1.RulesClient)
		gmObj, err1 := utl.ConvertModelBindingType(ruleParam, model0.RuleBindingType(), model1.RuleBindingType())
		if err1 != nil {
			return err1
		}
		err = client.Patch(domainIdParam, securityPolicyIdParam, ruleIdParam, gmObj.(model1.Rule))

	case utl.Multitenancy:
		client := c.Client.(client2.RulesClient)
		err = client.Patch(utl.DefaultOrgID, c.ProjectID, domainIdParam, securityPolicyIdParam, ruleIdParam, ruleParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c RuleClientContext) Update(domainIdParam string, securityPolicyIdParam string, ruleIdParam string, ruleParam model0.Rule) (model0.Rule, error) {
	var err error
	var obj model0.Rule

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RulesClient)
		obj, err = client.Update(domainIdParam, securityPolicyIdParam, ruleIdParam, ruleParam)

	case utl.Global:
		client := c.Client.(client1.RulesClient)
		gmObj, err := utl.ConvertModelBindingType(ruleParam, model0.RuleBindingType(), model1.RuleBindingType())
		if err != nil {
			return obj, err
		}
		gmObj, err = client.Update(domainIdParam, securityPolicyIdParam, ruleIdParam, gmObj.(model1.Rule))
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.RuleBindingType(), model0.RuleBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.Rule)

	case utl.Multitenancy:
		client := c.Client.(client2.RulesClient)
		obj, err = client.Update(utl.DefaultOrgID, c.ProjectID, domainIdParam, securityPolicyIdParam, ruleIdParam, ruleParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c RuleClientContext) List(domainIdParam string, securityPolicyIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.RuleListResult, error) {
	var err error
	var obj model0.RuleListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RulesClient)
		obj, err = client.List(domainIdParam, securityPolicyIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	case utl.Global:
		client := c.Client.(client1.RulesClient)
		gmObj, err := client.List(domainIdParam, securityPolicyIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.RuleListResultBindingType(), model0.RuleListResultBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.RuleListResult)

	case utl.Multitenancy:
		client := c.Client.(client2.RulesClient)
		obj, err = client.List(utl.DefaultOrgID, c.ProjectID, domainIdParam, securityPolicyIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}