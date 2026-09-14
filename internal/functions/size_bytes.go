// Package functions implements provider-defined functions.
package functions

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = SizeBytes{}

// SizeBytes converts a human-readable size such as "10G" or "1.5 TiB" to bytes.
type SizeBytes struct{}

// NewSizeBytes returns the size_bytes function.
func NewSizeBytes() function.Function { return SizeBytes{} }

func (SizeBytes) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "size_bytes"
}

func (SizeBytes) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Converts a size such as \"10G\" or \"1.5 TiB\" to bytes.",
		MarkdownDescription: "Units are binary, as ZFS uses them: `K`, `M`, `G`, `T` and `P` (with optional `iB` or `B`) " +
			"multiply by 1024, 1024², and so on. A bare number is bytes. The result must be a whole number of bytes.",
		Parameters: []function.Parameter{
			function.StringParameter{Name: "size", MarkdownDescription: "Size to convert, e.g. `\"500M\"`."},
		},
		Return: function.Int64Return{},
	}
}

func (SizeBytes) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var size string
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &size))
	if resp.Error != nil {
		return
	}
	n, err := ParseSize(size)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())
		return
	}
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, n))
}

var sizePattern = regexp.MustCompile(`^\s*([0-9]+(?:\.[0-9]+)?)\s*([kmgtp]?)(?:i?b)?\s*$`)

// ParseSize parses a binary size string into bytes.
func ParseSize(s string) (int64, error) {
	m := sizePattern.FindStringSubmatch(strings.ToLower(s))
	if m == nil {
		return 0, fmt.Errorf("%q is not a size; use a number with an optional K, M, G, T or P suffix", s)
	}
	value, _, err := big.ParseFloat(m[1], 10, 128, big.ToNearestEven)
	if err != nil {
		return 0, fmt.Errorf("%q: %w", s, err)
	}
	exponent := strings.Index("kmgtp", m[2]) + 1
	if m[2] == "" {
		exponent = 0
	}
	value.Mul(value, new(big.Float).SetFloat64(math.Pow(1024, float64(exponent))))
	if !value.IsInt() {
		return 0, fmt.Errorf("%q is not a whole number of bytes", s)
	}
	n, acc := value.Int64()
	if acc != big.Exact {
		return 0, fmt.Errorf("%q is too large", s)
	}
	return n, nil
}
