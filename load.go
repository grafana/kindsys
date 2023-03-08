package kindsys

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"github.com/grafana/thema"
	"github.com/grafana/thema/load"
)

func LoadInstance(pkg, relpath string, fs fs.FS) (*build.Instance, error) {
	if pkg != "" {
		return load.InstanceWithThema(fs, relpath, load.Package(pkg))
	}
	return load.InstanceWithThema(fs, relpath)
}

func BuildInstance(ctx *cue.Context, relpath string, pkg string, overlay fs.FS) (cue.Value, error) {
	bi, err := LoadInstance(relpath, pkg, overlay)
	if err != nil {
		return cue.Value{}, err
	}

	if ctx == nil {
		return cue.Value{}, fmt.Errorf("nil cue context")
	}

	v := ctx.BuildInstance(bi)
	if v.Err() != nil {
		return v, fmt.Errorf("%s not a valid CUE instance: %w", relpath, v.Err())
	}
	return v, nil
}

// ToKindProps takes a cue.Value expected to represent a kind of the category
// specified by the type parameter and populates the Go type from the cue.Value.
func ToKindProps[T KindProperties](fw cue.Value, v cue.Value) (T, error) {
	props := new(T)
	if !v.Exists() {
		return *props, ErrValueNotExist
	}

	var kdef cue.Value
	anyprops := any(*props).(SomeKindProperties)
	switch anyprops.(type) {
	case CoreProperties:
		kdef = fw.LookupPath(cue.MakePath(cue.Str("Core")))
	case CustomProperties:
		kdef = fw.LookupPath(cue.MakePath(cue.Str("Custom")))
	case ComposableProperties:
		kdef = fw.LookupPath(cue.MakePath(cue.Str("Composable")))
	default:
		// unreachable so long as all the possibilities in KindProperties have switch branches
		panic("unreachable")
	}

	item := v.Unify(kdef)
	if item.Err() != nil {
		return *props, errors.Wrap(errors.Promote(ErrValueNotAKind, ""), item.Err())
	}

	if err := item.Decode(props); err != nil {
		// Should only be reachable if CUE and Go framework types have diverged
		panic(errors.Details(err, nil))
	}

	return *props, nil
}

// SomeDef represents a single kind definition, having been loaded and
// validated by a func such as [LoadCoreKindDef].
//
// The underlying type of the Properties field indicates the category of kind.
type SomeDef struct {
	// V is the cue.Value containing the entire Kind definition.
	V cue.Value
	// Properties contains the kind's declarative non-schema properties.
	Properties SomeKindProperties
}

// BindKindLineage binds the lineage for the kind definition.
//
// For kinds with a corresponding Go type, it is left to the caller to associate
// that Go type with the lineage returned from this function by a call to
// [thema.BindType].
func (def SomeDef) BindKindLineage(rt *thema.Runtime, opts ...thema.BindOption) (thema.Lineage, error) {
	if rt == nil {
		return &thema.UnaryLineage{}, fmt.Errorf("nil thema.Runtime")
	}
	return thema.BindLineage(def.V.LookupPath(cue.MakePath(cue.Str("lineage"))), rt, opts...)
}

// IsCore indicates whether the represented kind is a core kind.
func (def SomeDef) IsCore() bool {
	_, is := def.Properties.(CoreProperties)
	return is
}

// IsCustom indicates whether the represented kind is a custom kind.
func (def SomeDef) IsCustom() bool {
	_, is := def.Properties.(CustomProperties)
	return is
}

// IsComposable indicates whether the represented kind is a composable kind.
func (def SomeDef) IsComposable() bool {
	_, is := def.Properties.(ComposableProperties)
	return is
}

// Def represents a single kind definition, having been loaded and validated by
// a func such as [LoadCoreKindDef].
//
// Its type parameter indicates the category of kind.
//
// Thema lineages in the contained definition have not yet necessarily been
// validated.
type Def[T KindProperties] struct {
	// V is the cue.Value containing the entire Kind definition.
	V cue.Value
	// Properties contains the kind's declarative non-schema properties.
	Properties T
}

// Some converts the typed Def to the equivalent typeless SomeDef.
func (def Def[T]) Some() SomeDef {
	return SomeDef{
		V:          def.V,
		Properties: any(def.Properties).(SomeKindProperties),
	}
}

// LoadCoreKindDef loads and validates a core kind definition of the kind category
// indicated by the type parameter. On success, it returns a [Def] which
// contains the entire contents of the kind definition.
//
// declpath is the path to the directory containing the core kind definition,
// relative to the grafana/grafana root. For example, dashboards are in
// "kinds/dashboard".
//
// The .cue file bytes containing the core kind definition will be retrieved
// from the central embedded FS, [grafana.CueSchemaFS]. If desired (e.g. for
// testing), an optional fs.FS may be provided via the overlay parameter, which
// will be merged over [grafana.CueSchemaFS]. But in typical circumstances,
// overlay can and should be nil.
//
// This is a low-level function, primarily intended for use in code generation.
// For representations of core kinds that are useful in Go programs at runtime,
// see ["github.com/grafana/grafana/pkg/registry/corekind"].
func LoadCoreKindDef(defpath string, ctx *cue.Context, overlay fs.FS) (Def[CoreProperties], error) {
	none := Def[CoreProperties]{}
	vk, err := BuildInstance(ctx, defpath, "kind", overlay)
	if err != nil {
		return none, err
	}

	fw, err := LoadFrameworkCUE(ctx, defpath)
	if err != nil {
		return none, err
	}

	props, err := ToKindProps[CoreProperties](vk, fw)
	if err != nil {
		return none, err
	}

	return Def[CoreProperties]{
		V:          vk,
		Properties: props,
	}, nil
}

func LoadFrameworkCUE(ctx *cue.Context, fp string) (cue.Value, error) {
	v, err := BuildInstance(ctx, filepath.Join(fp, "kindsys"), "kindsys", nil)
	if err != nil {
		return v, err
	}

	if err = v.Validate(cue.Concrete(false), cue.All()); err != nil {
		return cue.Value{}, fmt.Errorf("kindsys framework loaded cue.Value has err: %w", err)
	}

	return v, nil
}
