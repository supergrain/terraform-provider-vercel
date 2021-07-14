package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/project_domain"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVercelProjectDomain(t *testing.T) {
	projectDomainName, _ := uuid.GenerateUUID()
	projectID, _ := uuid.GenerateUUID()
	redirect := "redirect.supergrain.com"
	var (

		// Holds the project domain fetched from vercel when we create it at the beginning
		actualProjectDomainAfterCreation project_domain.ProjectDomain

		// Changing the redirect should not result in the recreation of the project domain, so we expect to have the same name.

		// Used everywhere else
		actualProjectDomain project_domain.ProjectDomain
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckVercelProjectDomainDestroy(projectID, projectDomainName),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVercelProjectDomainConfig(projectDomainName, projectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectDomainStateHasValues(
						"vercel_project_domain.new", project_domain.ProjectDomain{Name: projectDomainName, ProjectID: projectID},
					),
					testAccCheckVercelProjectDomainExists("vercel_project_domain.new", &actualProjectDomainAfterCreation),
					testAccCheckActualProjectDomainHasValues(&actualProjectDomainAfterCreation, &project_domain.ProjectDomain{
						Name:      projectDomainName,
						ProjectID: projectID,
					}),
				),
			},
			{
				Config: testAccCheckVercelProjectDomainConfigWithRedirect(projectDomainName, projectID, redirect),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVercelProjectDomainExists("vercel_project.new", &actualProjectDomain),
					testAccCheckProjectDomainStateHasValues(
						"vercel_project_domain.new", project_domain.ProjectDomain{
							Name:      projectDomainName,
							ProjectID: projectID,
							Redirect:  redirect,
						},
					),
					testAccCheckActualProjectDomainHasValues(&actualProjectDomain, &project_domain.ProjectDomain{
						Name:      projectDomainName,
						ProjectID: projectID,
						Redirect:  redirect,
					},
					),
				),
			},
		},
	})
}

// Combines multiple `resource.TestCheckResourceAttr` calls
func testAccCheckProjectDomainStateHasValues(name string, want project_domain.ProjectDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tests := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(
				name, "name", want.Name),
			resource.TestCheckResourceAttr(
				name, "project_id", want.ProjectID),
			resource.TestCheckResourceAttr(
				name, "redirect", want.Redirect),
		}

		for _, test := range tests {
			err := test(s)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckActualProjectDomainHasValues(actual *project_domain.ProjectDomain, want *project_domain.ProjectDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if actual.Name != want.Name {
			return fmt.Errorf("name does not match, expected: %s, got: %s", want.Name, actual.Name)
		}
		if actual.Name == "" {
			return fmt.Errorf("Name is empty")
		}

		if actual.ProjectID != want.ProjectID {
			return fmt.Errorf("project_id does not match: expected: %s, got: %s", want.ProjectID, actual.ProjectID)
		}

		if actual.ProjectID == "" {
			return fmt.Errorf("ProjectID is empty")
		}

		return nil
	}
}

// Test whether the project was destroyed properly and finishes the job if necessary
func testAccCheckVercelProjectDomainDestroy(projectID, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := vercel.New(os.Getenv("VERCEL_TOKEN"))

		for _, rs := range s.RootModule().Resources {
			if rs.Type != name {
				continue
			}

			domainName := rs.Primary.Attributes["name"]
			if domainName == "" {
				return fmt.Errorf("No name set")
			}
			projectID := rs.Primary.Attributes["project_id"]
			if projectID == "" {
				return fmt.Errorf("No project id set")
			}

			projectDomain, err := client.ProjectDomain.Read(projectID, domainName, "")
			if err == nil {
				message := "Project domain was not deleted from vercel during terraform destroy."
				deleteErr := client.ProjectDomain.Delete(projectDomain.ProjectID, projectDomain.Name, "")
				if deleteErr != nil {
					return fmt.Errorf(message+" Automated removal did not succeed. Please manually remove @%s. Error: %w", projectDomain.Name, err)
				}
				return fmt.Errorf(message + " It was removed now.")
			}

		}
		return nil
	}
}

func testAccCheckVercelProjectDomainConfig(name, projectID string) string {
	return fmt.Sprintf(`
	resource "vercel_project_domain" "new" {
		name = "%s"
		project_id = "%s"
	}
	`, name, projectID)
}

func testAccCheckVercelProjectDomainConfigWithRedirect(name, projectID, redirect string) string {
	return fmt.Sprintf(`
	resource "vercel_project_domain" "new" {
		name = "%s"
		project_id = "%s"
		redirect = "%s"
	}
	`, name, projectID, redirect)
}

func testAccCheckVercelProjectDomainExists(n string, actual *project_domain.ProjectDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s in %+v", n, s.RootModule().Resources)
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			return fmt.Errorf("No name set")
		}
		projectID := rs.Primary.Attributes["project_id"]
		if projectID == "" {
			return fmt.Errorf("No project id set")
		}

		projectDomain, err := vercel.New(os.Getenv("VERCEL_TOKEN")).ProjectDomain.Read(projectID, name, "")
		if err != nil {
			return err
		}
		*actual = projectDomain
		return nil
	}
}
