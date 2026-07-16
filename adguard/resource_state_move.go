package adguard

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

const (
	legacyProviderAddress  = "registry.opentofu.org/gmichels/adguard"
	currentProviderAddress = "registry.opentofu.org/erikbpf/adguardhome"
	legacyProviderType     = "adguard"
)

func legacyStateMover(sourceTypeName string, sourceSchema schema.Schema, newModel func() any) resource.StateMover {
	return resource.StateMover{
		SourceSchema: &sourceSchema,
		StateMover: func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
			if (req.SourceProviderAddress != legacyProviderAddress &&
				req.SourceProviderAddress != currentProviderAddress &&
				req.SourceProviderAddress != legacyProviderType) ||
				req.SourceTypeName != sourceTypeName ||
				req.SourceSchemaVersion != 0 {
				return
			}

			if req.SourceState == nil {
				resp.Diagnostics.AddError(
					"Unable to Move Legacy Resource State",
					"The legacy resource state does not match its version 0 schema.",
				)
				return
			}

			model := newModel()
			resp.Diagnostics.Append(req.SourceState.Get(ctx, model)...)
			if resp.Diagnostics.HasError() {
				return
			}

			resp.Diagnostics.Append(resp.TargetState.Set(ctx, model)...)
			if !resp.Diagnostics.HasError() {
				resp.TargetPrivate = req.SourcePrivate
			}
		},
	}
}

func resourceSchema(ctx context.Context, r resource.Resource) schema.Schema {
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	return resp.Schema
}
