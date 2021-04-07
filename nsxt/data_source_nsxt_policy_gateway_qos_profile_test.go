/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestAccDataSourceNsxtPolicyGatewayQosProfile_basic(t *testing.T) {
	name := getAccTestDataSourceName()
	testResourceName := "data.nsxt_policy_gateway_qos_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyGatewayQosProfileDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyGatewayQosProfileCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyGatewayQosProfileReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyGatewayQosProfileCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	displayName := name
	description := name
	obj := model.GatewayQosProfile{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	id := newUUID()

	if testAccIsGlobalManager() {
		gmObj, convErr := convertModelBindingType(obj, model.GatewayQosProfileBindingType(), gm_model.GatewayQosProfileBindingType())
		if convErr != nil {
			return convErr
		}
		client := gm_infra.NewGatewayQosProfilesClient(connector)
		err = client.Patch(id, gmObj.(gm_model.GatewayQosProfile), nil)
	} else {
		client := infra.NewGatewayQosProfilesClient(connector)
		err = client.Patch(id, obj, nil)
	}

	if err != nil {
		return handleCreateError("GatewayQosProfile", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyGatewayQosProfileDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	// Find the object by name and delete it
	if testAccIsGlobalManager() {
		objID, err := testGetObjIDByName(name, "GatewayQosProfile")
		if err == nil {
			client := gm_infra.NewGatewayQosProfilesClient(connector)
			err := client.Delete(objID, nil)
			if err != nil {
				return handleDeleteError("GatewayQosProfile", objID, err)
			}
			return nil
		}
	} else {
		client := infra.NewGatewayQosProfilesClient(connector)
		// Find the object by name
		objList, err := client.List(nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleListError("GatewayQosProfile", err)
		}
		for _, objInList := range objList.Results {
			if *objInList.DisplayName == name {
				err := client.Delete(*objInList.Id, nil)
				if err != nil {
					return handleDeleteError("GatewayQosProfile", *objInList.Id, err)
				}
				return nil
			}
		}
	}
	return fmt.Errorf("Error while deleting GatewayQosProfile '%s': resource not found", name)
}

func testAccNsxtPolicyGatewayQosProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_gateway_qos_profile" "test" {
  display_name = "%s"
}`, name)
}
