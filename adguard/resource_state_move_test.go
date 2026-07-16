package adguard

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestLegacyStateMovesPreserveRequiredResourceState(t *testing.T) {
	ctx := context.Background()
	rules, diags := types.ListValueFrom(ctx, types.StringType, []string{"||secret.example^"})
	if diags.HasError() {
		t.Fatalf("building rules value: %v", diags)
	}

	tests := []struct {
		name       string
		resource   resource.ResourceWithMoveState
		sourceType string
		model      any
		result     func() any
	}{
		{
			name:       "rewrite",
			resource:   &rewriteResource{},
			sourceType: "adguard_rewrite",
			model: &rewriteResourceModel{
				ID: types.StringValue("example.org||192.0.2.1"), LastUpdated: types.StringValue("unchanged"),
				Domain: types.StringValue("example.org"), Answer: types.StringValue("192.0.2.1"), Enabled: types.BoolValue(true),
			},
			result: func() any { return &rewriteResourceModel{} },
		},
		{
			name:       "list_filter",
			resource:   &listFilterResource{},
			sourceType: "adguard_list_filter",
			model: &listFilterResourceModel{
				ID: types.StringValue("7"), Url: types.StringValue("https://example.org/filter.txt"), Name: types.StringValue("example"),
				LastUpdated: types.StringValue("unchanged"), RulesCount: types.Int64Value(42), Enabled: types.BoolValue(true), Whitelist: types.BoolValue(false),
			},
			result: func() any { return &listFilterResourceModel{} },
		},
		{
			name:       "user_rules",
			resource:   &userRulesResource{},
			sourceType: "adguard_user_rules",
			model:      &userRulesResourceModel{ID: types.StringValue("1"), Rules: rules, LastUpdated: types.StringValue("unchanged")},
			result:     func() any { return &userRulesResourceModel{} },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mover := test.resource.MoveState(ctx)[0]
			source := tfsdk.State{Schema: *mover.SourceSchema}
			if diags := source.Set(ctx, test.model); diags.HasError() {
				t.Fatalf("setting source state: %v", diags)
			}
			resp := resource.MoveStateResponse{TargetState: tfsdk.State{Schema: *mover.SourceSchema}}
			mover.StateMover(ctx, resource.MoveStateRequest{
				SourceProviderAddress: legacyProviderAddress,
				SourceTypeName:        test.sourceType,
				SourceSchemaVersion:   0,
				SourceState:           &source,
			}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("moving state: %v", resp.Diagnostics)
			}
			got := test.result()
			if diags := resp.TargetState.Get(ctx, got); diags.HasError() {
				t.Fatalf("reading target state: %v", diags)
			}
			if !reflect.DeepEqual(got, test.model) {
				t.Fatalf("state changed during move")
			}
		})
	}
}

func TestLegacyStateMoverRejectsUnsupportedSource(t *testing.T) {
	ctx := context.Background()
	mover := (&rewriteResource{}).MoveState(ctx)[0]

	tests := []resource.MoveStateRequest{
		{SourceProviderAddress: "registry.terraform.io/gmichels/adguard", SourceTypeName: "adguard_rewrite", SourceSchemaVersion: 0},
		{SourceProviderAddress: legacyProviderAddress, SourceTypeName: "adguard_client", SourceSchemaVersion: 0},
		{SourceProviderAddress: legacyProviderAddress, SourceTypeName: "adguard_rewrite", SourceSchemaVersion: 1},
	}
	for _, req := range tests {
		var resp resource.MoveStateResponse
		mover.StateMover(ctx, req, &resp)
		if len(resp.Diagnostics) != 0 || !reflect.DeepEqual(resp.TargetState, tfsdk.State{}) {
			t.Fatalf("unsupported source was handled: %#v", req)
		}
	}
}

func TestLegacyStateMoverAcceptsStateAfterProviderAddressReplacement(t *testing.T) {
	ctx := context.Background()
	mover := (&rewriteResource{}).MoveState(ctx)[0]
	var resp resource.MoveStateResponse
	mover.StateMover(ctx, resource.MoveStateRequest{
		SourceProviderAddress: currentProviderAddress,
		SourceTypeName:        "adguard_rewrite",
		SourceSchemaVersion:   0,
	}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("current provider address did not match after provider address replacement")
	}
}

func TestLegacyStateMoverAcceptsOpenTofuProviderTypeAddress(t *testing.T) {
	ctx := context.Background()
	mover := (&rewriteResource{}).MoveState(ctx)[0]
	var resp resource.MoveStateResponse
	mover.StateMover(ctx, resource.MoveStateRequest{
		SourceProviderAddress: legacyProviderType,
		SourceTypeName:        "adguard_rewrite",
		SourceSchemaVersion:   0,
	}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("OpenTofu provider type address did not match")
	}
}

func TestAllLegacyResourceTypesHaveStateMovers(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		resource   resource.ResourceWithMoveState
		sourceType string
	}{
		{&rewriteResource{}, "adguard_rewrite"},
		{&listFilterResource{}, "adguard_list_filter"},
		{&userRulesResource{}, "adguard_user_rules"},
		{&clientResource{}, "adguard_client"},
		{&configResource{}, "adguard_config"},
	}
	for _, test := range tests {
		movers := test.resource.MoveState(ctx)
		if len(movers) != 1 || movers[0].SourceSchema == nil {
			t.Fatalf("%s has no source schema", test.sourceType)
		}
		var resp resource.MoveStateResponse
		movers[0].StateMover(ctx, resource.MoveStateRequest{
			SourceProviderAddress: legacyProviderAddress,
			SourceTypeName:        test.sourceType,
			SourceSchemaVersion:   0,
		}, &resp)
		if !resp.Diagnostics.HasError() {
			t.Fatalf("%s did not match its legacy source", test.sourceType)
		}
	}
}
