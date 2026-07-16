package adguard

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRewriteDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
data "adguardhome_rewrite" "test" {
	domain = "example.org"
	answer = "1.2.3.4"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.adguardhome_rewrite.test", "domain", "example.org"),
					resource.TestCheckResourceAttr("data.adguardhome_rewrite.test", "answer", "1.2.3.4"),
					resource.TestCheckResourceAttr("data.adguardhome_rewrite.test", "enabled", "true"),

					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.adguardhome_rewrite.test", "id", "placeholder"),
				),
			},
		},
	})
}
