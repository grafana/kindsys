package kindsys

import (
	"errors"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/grafana/thema"
)

func BindProvider(rt *thema.Runtime, val cue.Value) (*Provider, error) {
	p := Provider{
		CoreKinds:       map[string]Core{},
		ComposableKinds: map[string]Composable{},
		CustomKinds:     map[string]Custom{},
	}

	if val.Err() != nil {
		return nil, fmt.Errorf("could not bind provider due to cue value error: %w", val.Err())
	}

	if !val.Exists() {
		return nil, errors.New("could not bind provider due to missing cue.Value")
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

	name, err := nameVal.String()
	if err != nil {
		return nil, fmt.Errorf("invalid provider name: %w", err)
	}
	p.Name = name

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
		it, err := coreKindsVal.Fields()
		if err != nil {
			return nil, fmt.Errorf("coreKinds is not a struct: %w", err)
		}

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
		it, err := composableKindsVal.Fields()
		if err != nil {
			return nil, fmt.Errorf("composableKinds is not a struct: %w", err)
		}

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
		it, err := customKindsVal.Fields()
		if err != nil {
			return nil, fmt.Errorf("customKinds is not a struct: %w", err)
		}

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
