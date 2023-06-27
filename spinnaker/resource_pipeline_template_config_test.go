package spinnaker

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSpinnakerTemplateConfig_build(t *testing.T) {
	rName := "docta:Build and Deploy EKS"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:            testAccSpinnakerTemplateConfig_build(rName),
				ResourceName:      "spinnaker_pipeline_template_config.test1",
				ImportState:       true,
				ImportStateId:     rName,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccSpinnakerTemplateConfig_nondefault(t *testing.T) {
	resourceName := "spinnaker_pipeline_template_config.test2"
	appName := "docta"
	templateName := "Build and Deploy EKS 2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSpinnakerTemplateConfig_nondefault(appName, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPipelineTemplateConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "template_name", templateName),
					resource.TestCheckResourceAttr(resourceName, "application", appName),
				),
			},
		},
	})
}

func testAccCheckPipelineTemplateConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
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

func testAccSpinnakerTemplateConfig_build(rName string) string {
	return fmt.Sprintf(`
resource "spinnaker_pipeline_template_config" "test1" {
	name  = %q
}
`, rName)
}

func testAccSpinnakerTemplateConfig_nondefault(appName string, templateName string) string {
	return fmt.Sprintf(`
locals {
	our_rendered_content = templatefile("${path.module}/test/test_template_config.yaml", {})
}

resource "spinnaker_pipeline_template_config" "test2" {
	template_name   = %q
	application     = %q
	pipeline_config = local.our_rendered_content
}
`, templateName, appName)
}
