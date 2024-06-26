//nolint:revive
package segments

// The following file has been autogenerated. Please avoid any changes!
import (
	"errors"

	model0 "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/segments"
	lrmodel1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/segments"
	lrmodel0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client2 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra/segments"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type StructValueClientContext utl.ClientContext

func NewDhcpStaticBindingConfigsClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *StructValueClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewDhcpStaticBindingConfigsClient(connector)

	case utl.Global:
		client = client1.NewDhcpStaticBindingConfigsClient(connector)

	case utl.Multitenancy:
		client = client2.NewDhcpStaticBindingConfigsClient(connector)

	default:
		return nil
	}
	return &StructValueClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c StructValueClientContext) Get(segmentIdParam string, bindingIdParam string) (*model0.StructValue, error) {
	var obj *model0.StructValue
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.DhcpStaticBindingConfigsClient)
		obj, err = client.Get(segmentIdParam, bindingIdParam)
		if err != nil {
			return obj, err
		}

	case utl.Global:
		client := c.Client.(client1.DhcpStaticBindingConfigsClient)
		obj, err = client.Get(segmentIdParam, bindingIdParam)
		if err != nil {
			return obj, err
		}

	case utl.Multitenancy:
		client := c.Client.(client2.DhcpStaticBindingConfigsClient)
		obj, err = client.Get(utl.DefaultOrgID, c.ProjectID, segmentIdParam, bindingIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c StructValueClientContext) Delete(segmentIdParam string, bindingIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.DhcpStaticBindingConfigsClient)
		err = client.Delete(segmentIdParam, bindingIdParam)

	case utl.Global:
		client := c.Client.(client1.DhcpStaticBindingConfigsClient)
		err = client.Delete(segmentIdParam, bindingIdParam)

	case utl.Multitenancy:
		client := c.Client.(client2.DhcpStaticBindingConfigsClient)
		err = client.Delete(utl.DefaultOrgID, c.ProjectID, segmentIdParam, bindingIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c StructValueClientContext) Patch(segmentIdParam string, bindingIdParam string, dhcpStaticBindingConfigParam *model0.StructValue) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.DhcpStaticBindingConfigsClient)
		err = client.Patch(segmentIdParam, bindingIdParam, dhcpStaticBindingConfigParam)

	case utl.Global:
		client := c.Client.(client1.DhcpStaticBindingConfigsClient)
		err = client.Patch(segmentIdParam, bindingIdParam, dhcpStaticBindingConfigParam)

	case utl.Multitenancy:
		client := c.Client.(client2.DhcpStaticBindingConfigsClient)
		err = client.Patch(utl.DefaultOrgID, c.ProjectID, segmentIdParam, bindingIdParam, dhcpStaticBindingConfigParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c StructValueClientContext) Update(segmentIdParam string, bindingIdParam string, dhcpStaticBindingConfigParam *model0.StructValue) (*model0.StructValue, error) {
	var err error
	var obj *model0.StructValue

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.DhcpStaticBindingConfigsClient)
		obj, err = client.Update(segmentIdParam, bindingIdParam, dhcpStaticBindingConfigParam)

	case utl.Global:
		client := c.Client.(client1.DhcpStaticBindingConfigsClient)
		obj, err = client.Update(segmentIdParam, bindingIdParam, dhcpStaticBindingConfigParam)

	case utl.Multitenancy:
		client := c.Client.(client2.DhcpStaticBindingConfigsClient)
		obj, err = client.Update(utl.DefaultOrgID, c.ProjectID, segmentIdParam, bindingIdParam, dhcpStaticBindingConfigParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c StructValueClientContext) List(segmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (lrmodel0.DhcpStaticBindingConfigListResult, error) {
	var err error
	var obj lrmodel0.DhcpStaticBindingConfigListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.DhcpStaticBindingConfigsClient)
		obj, err = client.List(segmentIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	case utl.Global:
		client := c.Client.(client1.DhcpStaticBindingConfigsClient)
		gmObj, err := client.List(segmentIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, lrmodel1.DhcpStaticBindingConfigListResultBindingType(), lrmodel0.DhcpStaticBindingConfigListResultBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(lrmodel0.DhcpStaticBindingConfigListResult)

	case utl.Multitenancy:
		client := c.Client.(client2.DhcpStaticBindingConfigsClient)
		obj, err = client.List(utl.DefaultOrgID, c.ProjectID, segmentIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}
