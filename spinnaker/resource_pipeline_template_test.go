package spinnaker

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSpinnakerTemplate_build(t *testing.T) {
	rName := "build-image-template"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:            testAccSpinnakerTemplate_build(rName),
				ResourceName:      "spinnaker_pipeline_template.test1",
				ImportState:       true,
				ImportStateId:     rName,
				ImportStateVerify: false,
			},
		},
	})
}

/*
func TestAccSpinnakerTemplate_nondefault(t *testing.T) {
	resourceName := "spinnaker_pipeline_template.test2"
	rName := acctest.RandomWithPrefix("tf-acc-test")
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
					resource.TestCheckResourceAttr(resourceName, "platform_health_only", "true"),
					resource.TestCheckResourceAttr(resourceName, "platform_health_only_show_override", "true"),
				),
			},
		},
	})
}

func testAccCheckPipelineTemplateExists(n string) resource.TestCheckFunc {
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
*/

func testAccSpinnakerTemplate_build(rName string) string {
	return fmt.Sprintf(`
resource "spinnaker_pipeline_template" "test1" {
	name  = %q
}
`, rName)
}

/*
func testAccSpinnakerApplication_nondefault(rName string) string {
	return fmt.Sprintf(`
resource "spinnaker_application" "test2" {
	name                               = %q
	email                              = "acceptance@test.com"
	description                        = "My application"
	platform_health_only               = true
	platform_health_only_show_override = true
}
`, rName)
}
*/
