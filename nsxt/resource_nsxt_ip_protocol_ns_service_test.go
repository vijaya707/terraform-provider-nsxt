/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"net/http"
	"testing"
)

func TestAccResourceNsxtIpProtocolNsService_basic(t *testing.T) {
	serviceName := fmt.Sprintf("test-nsx-ip-protocol-service")
	updateServiceName := fmt.Sprintf("%s-update", serviceName)
	testResourceName := "nsxt_ip_protocol_ns_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXIpProtocolServiceCheckDestroy(state, serviceName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXIpProtocolServiceCreateTemplate(serviceName, 6),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXIpProtocolServiceExists(serviceName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", serviceName),
					resource.TestCheckResourceAttr(testResourceName, "description", "ip protocol service"),
					resource.TestCheckResourceAttr(testResourceName, "protocol", "6"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNSXIpProtocolServiceCreateTemplate(updateServiceName, 17),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXIpProtocolServiceExists(updateServiceName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateServiceName),
					resource.TestCheckResourceAttr(testResourceName, "description", "ip protocol service"),
					resource.TestCheckResourceAttr(testResourceName, "protocol", "17"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtIpProtocolNsService_importBasic(t *testing.T) {
	serviceName := fmt.Sprintf("test-nsx-ip-protocol-service")
	testResourceName := "nsxt_ip_protocol_ns_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXIpProtocolServiceCheckDestroy(state, serviceName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXIpProtocolServiceCreateTemplate(serviceName, 6),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNSXIpProtocolServiceExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX ip protocol service resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX ip protocol service resource ID not set in resources ")
		}

		service, responseCode, err := nsxClient.GroupingObjectsApi.ReadIpProtocolNSService(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving ip protocol service ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if ip protocol service %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if displayName == service.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX ip protocol ns service %s wasn't found", displayName)
	}
}

func testAccNSXIpProtocolServiceCheckDestroy(state *terraform.State, displayName string) error {
	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_ip_protocol_ns_service" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		service, responseCode, err := nsxClient.GroupingObjectsApi.ReadIpProtocolNSService(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving L4 ns service ID %s. Error: %v", resourceID, err)
		}

		if displayName == service.DisplayName {
			return fmt.Errorf("NSX L4 ns service %s still exists", displayName)
		}
	}
	return nil
}

func testAccNSXIpProtocolServiceCreateTemplate(serviceName string, protocol int) string {
	return fmt.Sprintf(`
resource "nsxt_ip_protocol_ns_service" "test" {
  description  = "ip protocol service"
  display_name = "%s"
  protocol     = "%d"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, serviceName, protocol)
}
