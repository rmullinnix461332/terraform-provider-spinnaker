package spinnaker

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func TestAccSpinnakerTemplate_nondefault(t *testing.T) {
	resourceName := "spinnaker_pipeline_template.test2"
	rName := "test-build-image-template"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSpinnakerTemplate_nondefault(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPipelineTemplateExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckPipelineTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Template Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Template ID is set")
		}
		client := testAccProvider.Meta().(gateConfig).client
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, resp, err := client.PipelineTemplatesControllerApi.GetUsingGET(client.Context, rs.Primary.ID)

			if resp != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					return resource.RetryableError(fmt.Errorf("template does not exit"))
				} else if resp.StatusCode != http.StatusOK {
					return resource.NonRetryableError(fmt.Errorf("encountered an error getting template, status code: %d", resp.StatusCode))
				}
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Unable to find template after retries: %s", err)
		}
		return nil
	}
}

func testAccSpinnakerTemplate_build(rName string) string {
	return fmt.Sprintf(`
resource "spinnaker_pipeline_template" "test1" {
	name  = %q
}
`, rName)
}

func testAccSpinnakerTemplate_nondefault(rName string) string {
	return fmt.Sprintf(`

locals {
	rendered_template = templatefile("${path.module}/test/test_pipeline_template.yaml", {})
}

resource "spinnaker_pipeline_template" "test2" {
	name                               = %q
	template                           = local.rendered_template
}
`, rName)
}
