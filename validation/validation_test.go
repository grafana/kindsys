package validation

import (
	"cuelang.org/go/cue/cuecontext"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnsureNoExportedKindName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "Must return an error if schema contains #Kind field",
			value:   `#Kind: "foo" | "bar"`,
			wantErr: true,
		},
		{
			name:    "Must return an error if schema contains Kind field",
			value:   `Kind: "foo" | "bar"`,
			wantErr: true,
		},
		{
			name:    "Must return an error if schema contains #kind field",
			value:   `#kind: "foo" | "bar"`,
			wantErr: true,
		},
		{
			name:    "Must return an error if schema contains kind field",
			value:   `kind: "foo" | "bar"`,
			wantErr: true,
		},
		{
			name:    "Must not return an error",
			value:   `#Bar: "foo" | "bar"`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := cuecontext.New()
			valStr := "schemas: [{ " + tt.value + " }]"
			val := ctx.CompileString(valStr)

			err := EnsureNoExportedKindName(val)
			assert.True(t, (err != nil) == tt.wantErr, tt.name)
		})
	}
}
