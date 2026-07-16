package adguard

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestProviderMetadataUsesAdguardhomePrefix(t *testing.T) {
	var response provider.MetadataResponse
	(&adguardProvider{}).Metadata(context.Background(), provider.MetadataRequest{}, &response)
	if response.TypeName != "adguardhome" {
		t.Fatalf("unexpected provider type name: %q", response.TypeName)
	}
}

func TestEnvironmentValuePrefersAdguardhomeAndFallsBackToLegacy(t *testing.T) {
	t.Setenv("ADGUARDHOME_TEST_VALUE", "preferred")
	t.Setenv("ADGUARD_TEST_VALUE", "legacy")
	if got := environmentValue("ADGUARDHOME_TEST_VALUE", "ADGUARD_TEST_VALUE"); got != "preferred" {
		t.Fatalf("preferred environment value not selected: %q", got)
	}

	t.Setenv("ADGUARDHOME_TEST_VALUE", "")
	if got := environmentValue("ADGUARDHOME_TEST_VALUE", "ADGUARD_TEST_VALUE"); got != "legacy" {
		t.Fatalf("legacy environment fallback not selected: %q", got)
	}
}

func TestSecretSchemaAttributesAreSensitive(t *testing.T) {
	var providerResponse provider.SchemaResponse
	(&adguardProvider{}).Schema(context.Background(), provider.SchemaRequest{}, &providerResponse)
	password := providerResponse.Schema.Attributes["password"].(providerschema.StringAttribute)
	if !password.Sensitive {
		t.Fatal("provider password must be sensitive")
	}

	var resourceResponse resource.SchemaResponse
	(&configResource{}).Schema(context.Background(), resource.SchemaRequest{}, &resourceResponse)
	tls := resourceResponse.Schema.Attributes["tls"].(resourceschema.SingleNestedAttribute)
	privateKey := tls.Attributes["private_key"].(resourceschema.StringAttribute)
	if !privateKey.Sensitive {
		t.Fatal("TLS private key must be sensitive")
	}
}
