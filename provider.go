package kindsys

import (
	"errors"
	"fmt"
	"io/fs"

	"cuelang.org/go/cue"
	"github.com/grafana/thema"
)

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

	bi, err := LoadInstance("", "", fsys)
	if err != nil || bi.Err != nil {
		if err == nil {
			err = bi.Err
		}
		return nil, fmt.Errorf("failed to load instance: %w", err)
	}

	p := Provider{
		CoreKinds:       map[string]Core{},
		ComposableKinds: map[string]Composable{},
		CustomKinds:     map[string]Custom{},
	}

	val := ctx.BuildInstance(bi)
	if val.Err() != nil {
		return nil, fmt.Errorf("failed to create a cue.Value from build.Instance: %w", err)
	}

	if !val.Exists() {
		return nil, errors.New("cue.Value doesn't exist")
	}

	providerVal := defaultFramework.LookupPath(cue.MakePath(cue.Str("Provider")))
	if !providerVal.Exists() {
		return nil, errors.New("provider not found in framework")
	}

	val = val.Unify(providerVal)
	p.V = val

	nameVal := val.LookupPath(cue.MakePath(cue.Str("name")))
	if !nameVal.Exists() {
		return nil, errors.New("provider name is required")
	}
	if nameVal.Err() != nil {
		return nil, fmt.Errorf("provider name error: %w", nameVal.Err())
	}

	p.Name, err = nameVal.String()
	if err != nil {
		return nil, fmt.Errorf("invalid provider name: %w", err)
	}

	versionVal := val.LookupPath(cue.MakePath(cue.Str("version")))
	if !versionVal.Exists() {
		return nil, errors.New("provider version is required")
	}
	if versionVal.Err() != nil {
		return nil, fmt.Errorf("provider version error: %w", versionVal.Err())
	}

	p.Version, err = versionVal.String()
	if err != nil {
		return nil, fmt.Errorf("invalid provider version: %w", err)
	}

	coreKindsVal := val.LookupPath(cue.MakePath(cue.Str("coreKinds")))
	if coreKindsVal.Exists() {
		s, err := coreKindsVal.Struct()
		if err != nil {
			return nil, fmt.Errorf("coreKinds is not a struct: %w", err)
		}

		it := s.Fields()
		for it.Next() {
			props, err := ToKindProps[CoreProperties](it.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to populate core kind: %w", err)
			}

			coreKind, err := BindCore(rt, Def[CoreProperties]{
				Properties: props,
				V:          it.Value(),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to bind core kind %q: %w", it.Label(), err)
			}

			p.CoreKinds[it.Label()] = coreKind
		}
	}

	composableKindsVal := val.LookupPath(cue.MakePath(cue.Str("composableKinds")))
	if composableKindsVal.Exists() {
		s, err := composableKindsVal.Struct()
		if err != nil {
			return nil, fmt.Errorf("composableKinds is not a struct: %w", err)
		}

		it := s.Fields()
		for it.Next() {
			props, err := ToKindProps[ComposableProperties](it.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to populate composable kind: %w", err)
			}

			composableKind, err := BindComposable(rt, Def[ComposableProperties]{
				Properties: props,
				V:          it.Value(),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to bind composable kind %q: %w", it.Label(), err)
			}

			p.ComposableKinds[it.Label()] = composableKind
		}
	}

	customKindsVal := val.LookupPath(cue.MakePath(cue.Str("customKinds")))
	if customKindsVal.Exists() {
		s, err := customKindsVal.Struct()
		if err != nil {
			return nil, fmt.Errorf("customKinds is not a struct: %w", err)
		}

		it := s.Fields()
		for it.Next() {
			props, err := ToKindProps[CustomProperties](it.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to populate custom kind: %w", err)
			}

			customKind, err := BindCustom(rt, Def[CustomProperties]{
				Properties: props,
				V:          it.Value(),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to bind custom kind %q: %w", it.Label(), err)
			}

			p.CustomKinds[it.Label()] = customKind
		}
	}

	return &p, nil
}
