//nolint:revive
package realizedstate

// The following file has been autogenerated. Please avoid any changes!
import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/realized_state"
	model1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/realized_state"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client2 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra/realized_state"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type RealizedEntityClientContext utl.ClientContext

func NewRealizedEntitiesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *RealizedEntityClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewRealizedEntitiesClient(connector)

	case utl.Global:
		client = client1.NewRealizedEntitiesClient(connector)

	case utl.Multitenancy:
		client = client2.NewRealizedEntitiesClient(connector)

	default:
		return nil
	}
	return &RealizedEntityClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c RealizedEntityClientContext) List(intentPathParam string, sitePathParam *string) (model0.GenericPolicyRealizedResourceListResult, error) {
	var err error
	var obj model0.GenericPolicyRealizedResourceListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.RealizedEntitiesClient)
		obj, err = client.List(intentPathParam, sitePathParam)

	case utl.Global:
		client := c.Client.(client1.RealizedEntitiesClient)
		gmObj, err := client.List(intentPathParam, sitePathParam)
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.GenericPolicyRealizedResourceListResultBindingType(), model0.GenericPolicyRealizedResourceListResultBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.GenericPolicyRealizedResourceListResult)

	case utl.Multitenancy:
		client := c.Client.(client2.RealizedEntitiesClient)
		obj, err = client.List(utl.DefaultOrgID, c.ProjectID, intentPathParam, sitePathParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}
