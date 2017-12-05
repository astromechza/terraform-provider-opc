package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCOrchestratedInstance_Basic(t *testing.T) {
	resName := "opc_compute_orchestrated_instance.test"
	ri := acctest.RandInt()
	config := testAccOrchestrationBasic(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
					resource.TestCheckResourceAttrSet(resName, "instance.0.id"),
				),
			},
		},
	})
}

func testAccCheckOrchestrationExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.Orchestrations()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_orchestration" {
			continue
		}

		input := compute.GetOrchestrationInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetOrchestration(&input); err != nil {
			return fmt.Errorf("Error retrieving state of Orchestration %s: %s", input.Name, err)
		}
	}

	return nil
}

func testAccCheckOrchestrationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.Orchestrations()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_orchestration" {
			continue
		}

		input := compute.GetOrchestrationInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetOrchestration(&input); err == nil {
			return fmt.Errorf("Orchestration %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccOrchestrationBasic(rInt int) string {
	return fmt.Sprintf(`
  resource "opc_compute_orchestrated_instance" "test" {
    name        = "test_orchestration-%d"
    desired_state = "active"
		instance {
			name = "acc-test-instance-%d"
			label = "TestAccOPCInstance_basic"
			shape = "oc3"
			image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
			instance_attributes = <<JSON
{
	"foo": "bar"
}
JSON
	  }
  }
  `, rInt, rInt)
}
