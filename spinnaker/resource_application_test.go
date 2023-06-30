package spinnaker

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSpinnakerApplication_basic(t *testing.T) {
	//resourceName := "spinnaker_application.test"
	rName := "docta"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:            testAccSpinnakerApplication_basic(rName),
				ResourceName:      "spinnaker_application.test1",
				ImportState:       true,
				ImportStateId:     rName,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccSpinnakerApplication_nondefault(t *testing.T) {
	resourceName := "spinnaker_application.test2"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	emailUpdate := "changed@test.com"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSpinnakerApplication_nondefault(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "email", "acceptance@test.com"),
					resource.TestCheckResourceAttr(resourceName, "description", "My application"),
					resource.TestCheckResourceAttr(resourceName, "port", "8080"),
				),
			},
			{
				Config: testAccSpinnakerApplication_update(rName, emailUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPipelineTemplateConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "email", emailUpdate),
				),
			},
		},
	})
}

func testAccCheckApplicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Application Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Application ID is set")
		}
		client := testAccProvider.Meta().(gateConfig).client
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, resp, err := client.ApplicationControllerApi.GetApplicationUsingGET(client.Context, rs.Primary.ID, nil)

			if resp != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					return resource.RetryableError(fmt.Errorf("application does not exit"))
				} else if resp.StatusCode != http.StatusOK {
					return resource.NonRetryableError(fmt.Errorf("encountered an error getting application, status code: %d", resp.StatusCode))
				}
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Unable to find Application after retries: %s", err)
		}
		return nil
	}
}

func testAccSpinnakerApplication_basic(rName string) string {
	return fmt.Sprintf(`
resource "spinnaker_application" "test1" {
	name  = %q
	email = "acceptance@test.com"
}
`, rName)
}

func testAccSpinnakerApplication_nondefault(rName string) string {
	return fmt.Sprintf(`
resource "spinnaker_application" "test2" {
	name        = %q
	email       = "acceptance@test.com"
	description = "My application"
	port        = 8080
}
`, rName)
}

func testAccSpinnakerApplication_update(rName string, email string) string {
	return fmt.Sprintf(`
resource "spinnaker_application" "test2" {
	name        = %q
	email       = %q
	description = "My application"
	port        = 8080
}
`, rName, email)
}
