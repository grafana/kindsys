package kindsys

import (
	"errors"
	"fmt"
	"io/fs"

	"cuelang.org/go/cue"
	"github.com/grafana/thema"
)

const providerPackageName = "provider"

// Provider is a structure that can hold a set of core, composable and
// custom kinds together with provider metadata.
type Provider struct {
	// Name the unique name of the provider.
	Name string

	// Version the version of the provider.
	Version string

	// CoreKinds list of core kinds that this provider provides.
	CoreKinds map[string]Core

	// ComposableKinds list of composable kinds that this provider provides.
	ComposableKinds map[string]Composable

	// CustomKinds list of custom kinds that this provider provides.
	CustomKinds map[string]Custom

	// V is the cue.Value containing the entire provider definition.
	V cue.Value
}

// LoadProvider takes a virtual filesystem and checks that it contains a valid
// set of files that statically define a Provider.
//
// If any .cue files exist in the provider package, these will also be loaded and
// validated according to the [Provider] specification. This includes the validation
// of any core, composable or custom kinds and their contained lineages,
// via [thema.BindLineage].
//
// This function parses exactly one provider. It does not descend into
// subdirectories to search for additional .cue files.
//
// [Provider]: https://github.com/grafana/kindsys/blob/main/provider.cue
func LoadProvider(fsys fs.FS, rt *thema.Runtime) (*Provider, error) {
	ctx := cueContext()

	if fsys == nil {
		return nil, errors.New("fsys cannot be nil")
	}
	if rt == nil {
		rt = thema.NewRuntime(ctx)
	}

	bi, err := LoadInstance("", providerPackageName, fsys)
	if err != nil || bi.Err != nil {
		if err == nil {
			err = bi.Err
		}
		return nil, fmt.Errorf("failed to load instance: %w", err)
	}

	return BindProvider(rt, ctx.BuildInstance(bi))
}
