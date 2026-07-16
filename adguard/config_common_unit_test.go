package adguard

import (
	"context"
	"testing"

	adgmodels "github.com/gmichels/adguard-client-go/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNormalizeDisabledDhcpStatus(t *testing.T) {
	status := adgmodels.DhcpStatus{
		Enabled:       false,
		InterfaceName: "eth0",
		V4: adgmodels.DhcpConfigV4{
			GatewayIp:     "192.0.2.1",
			SubnetMask:    "255.255.255.0",
			RangeStart:    "192.0.2.10",
			RangeEnd:      "192.0.2.20",
			LeaseDuration: 123,
		},
		V6: adgmodels.DhcpConfigV6{
			RangeStart:    "2001:db8::10",
			LeaseDuration: 456,
		},
		StaticLeases: []adgmodels.DhcpStaticLease{{
			Mac:      "00:11:22:33:44:55",
			Ip:       "192.0.2.30",
			Hostname: "stale",
		}},
	}

	normalized := normalizeDisabledDhcpStatus(status)

	if normalized.InterfaceName != "" {
		t.Fatalf("interface was not normalized: %q", normalized.InterfaceName)
	}
	if normalized.V4.GatewayIp != "" || normalized.V4.RangeStart != "" || normalized.V4.LeaseDuration != uint64(CONFIG_DHCP_V4_LEASE_DURATION) {
		t.Fatalf("IPv4 settings were not normalized: %#v", normalized.V4)
	}
	if normalized.V6.RangeStart != "" || normalized.V6.LeaseDuration != uint64(CONFIG_DHCP_V6_LEASE_DURATION) {
		t.Fatalf("IPv6 settings were not normalized: %#v", normalized.V6)
	}
	if normalized.StaticLeases != nil {
		t.Fatalf("static leases were not normalized: %#v", normalized.StaticLeases)
	}
}

func TestNormalizeDisabledDhcpStatusPreservesEnabledConfig(t *testing.T) {
	status := adgmodels.DhcpStatus{
		Enabled:       true,
		InterfaceName: "eth0",
		V4:            adgmodels.DhcpConfigV4{GatewayIp: "192.0.2.1"},
	}

	normalized := normalizeDisabledDhcpStatus(status)
	if normalized.InterfaceName != status.InterfaceName || normalized.V4.GatewayIp != status.V4.GatewayIp {
		t.Fatalf("enabled DHCP config changed: %#v", normalized)
	}
}

func TestShouldSetDhcpConfig(t *testing.T) {
	disabled := types.ObjectValueMust(dhcpConfigModel{}.attrTypes(), dhcpConfigModel{}.defaultObject())
	enabledValues := dhcpConfigModel{}.defaultObject()
	enabledValues["enabled"] = types.BoolValue(true)
	enabled := types.ObjectValueMust(dhcpConfigModel{}.attrTypes(), enabledValues)

	if shouldSetDhcpConfig(types.BoolValue(false), disabled) {
		t.Fatal("disabled DHCP must not be written when current and planned states are disabled")
	}
	if !shouldSetDhcpConfig(types.BoolValue(true), disabled) {
		t.Fatal("enabling DHCP must write its configuration")
	}
	if !shouldSetDhcpConfig(types.BoolValue(false), enabled) {
		t.Fatal("disabling DHCP must write its configuration")
	}
	if !shouldSetDhcpConfig(types.BoolValue(false), types.ObjectNull(dhcpConfigModel{}.attrTypes())) {
		t.Fatal("unknown prior DHCP state must fail closed by writing the planned configuration")
	}
}

func TestMutableRuntimeStringsRemainComputedDuringPlan(t *testing.T) {
	var response resource.SchemaResponse
	(&configResource{}).Schema(context.Background(), resource.SchemaRequest{}, &response)

	assertComputed := func(name string, attribute schema.Attribute) {
		t.Helper()
		stringAttribute, ok := attribute.(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s is not a string attribute", name)
		}
		if !stringAttribute.Computed {
			t.Fatalf("%s is not computed", name)
		}
		if len(stringAttribute.PlanModifiers) != 0 {
			t.Fatalf("%s must be unknown during updates so apply can refresh it", name)
		}
	}

	assertComputed("last_updated", response.Schema.Attributes["last_updated"])
	tlsAttribute := response.Schema.Attributes["tls"].(schema.SingleNestedAttribute)
	for _, name := range []string{"issuer", "key_type", "not_after", "not_before", "subject", "warning_validation"} {
		assertComputed("tls."+name, tlsAttribute.Attributes[name])
	}
}

func TestNormalizeTlsTimestamp(t *testing.T) {
	if got := normalizeTlsTimestamp("0001-01-01T00:00:00Z"); got != "" {
		t.Fatalf("zero timestamp was not normalized: %q", got)
	}
	const timestamp = "2026-07-16T12:34:56Z"
	if got := normalizeTlsTimestamp(timestamp); got != timestamp {
		t.Fatalf("real timestamp changed: %q", got)
	}
}
