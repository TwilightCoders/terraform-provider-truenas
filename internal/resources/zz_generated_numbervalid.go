// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Integer attributes are exposed as arbitrary-precision numbers rather than int64: TrueNAS
// returns values that do not fit in 64 bits (an X.509 serial is up to 20 octets). The framework's
// number validators cover OneOf and NoneOf only, so the bounds and integrality checks that
// Int64Attribute gave for free are reimplemented here against big.Float.

// numberAtLeast rejects values below min.
type numberAtLeast struct{ min float64 }

func (v numberAtLeast) Description(context.Context) string {
	return fmt.Sprintf("value must be at least %v", v.min)
}

func (v numberAtLeast) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (v numberAtLeast) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if req.ConfigValue.ValueBigFloat().Cmp(big.NewFloat(v.min)) < 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Attribute Value",
			fmt.Sprintf("Attribute %s %s, got: %v", req.Path, v.Description(ctx), req.ConfigValue.ValueBigFloat()))
	}
}

// numberAtMost rejects values above max.
type numberAtMost struct{ max float64 }

func (v numberAtMost) Description(context.Context) string {
	return fmt.Sprintf("value must be at most %v", v.max)
}

func (v numberAtMost) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (v numberAtMost) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if req.ConfigValue.ValueBigFloat().Cmp(big.NewFloat(v.max)) > 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Attribute Value",
			fmt.Sprintf("Attribute %s %s, got: %v", req.Path, v.Description(ctx), req.ConfigValue.ValueBigFloat()))
	}
}

// numberIsInteger rejects fractional values, which Int64Attribute used to reject by its type.
type numberIsInteger struct{}

func (numberIsInteger) Description(context.Context) string { return "value must be a whole number" }

func (v numberIsInteger) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (v numberIsInteger) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if f := req.ConfigValue.ValueBigFloat(); f != nil && !f.IsInt() {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Attribute Value",
			fmt.Sprintf("Attribute %s %s, got: %v", req.Path, v.Description(ctx), f))
	}
}
