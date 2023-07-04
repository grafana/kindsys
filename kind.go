package kindsys

import (
	"fmt"

	"github.com/grafana/thema"
)

// TODO docs
type Maturity string

const (
	MaturityMerged       Maturity = "merged"
	MaturityExperimental Maturity = "experimental"
	MaturityStable       Maturity = "stable"
	MaturityMature       Maturity = "mature"
)

func maturityIdx(m Maturity) int {
	// icky to do this globally, this is effectively setting a default
	if string(m) == "" {
		m = MaturityMerged
	}

	for i, ms := range maturityOrder {
		if m == ms {
			return i
		}
	}
	panic(fmt.Sprintf("unknown maturity milestone %s", m))
}

var maturityOrder = []Maturity{
	MaturityMerged,
	MaturityExperimental,
	MaturityStable,
	MaturityMature,
}

func (m Maturity) Less(om Maturity) bool {
	return maturityIdx(m) < maturityIdx(om)
}

func (m Maturity) String() string {
	return string(m)
}

// Kind is a runtime representation of a Grafana kind definition.
//
// Kind definitions are canonically written in CUE. Loading and validating such
// CUE definitions produces instances of this type. Kind, and its
// sub-interfaces, are the expected canonical way of working with kinds in Go.
//
// Kind has six sub-interfaces, all of which provide:
//
// - Access to the kind's defined meta-properties, such as `name`, `pluralName`, or `maturity`
// - Access to the schemas defined in the kind
// - Methods for certain key operations on the kind and object instances of its schemas
//
// Kind definitions are written in CUE. The meta-schema specifying how to write
// kind definitions are also written in CUE. See the files at the root of
// [the kindsys repository].
//
// There are three categories of kinds, each having its own sub-interface:
// [Core], [Custom], and [Composable]. All kind definitions are in exactly one
// category (a kind can't be both Core and Composable). Correspondingly, all
// instances of Kind also implement exactly one of these sub-interfaces.
//
// Conceptually, kinds are similar to class definitions in object-oriented
// programming. They define a particular type of object, and how instances of
// that object should be created. The object defined in a [Core] or [Custom] kind
// is called a [Resource]. TODO name for the associated object for composable kinds
//
// [Core], [Custom] and [Composable] all provide methods for unmarshaling []byte
// into an unstructured Go type, [UnstructuredResource], similar to how
// json.Unmarshal can use map[string]any as a universal fallback. Relying on
// this untyped approach is recommended for use cases that need to work
// generically on any Kind. This is especially because untyped Kinds are
// portable, and can be loaded at runtime in Go: the original CUE definition is
// sufficient to create instances of [Core], [Custom] or [Composable].
//
// However, when working with a specific, finite set of kinds, it is usually
// preferable to use the typed interfaces:
//
// - [Core] -> [TypedCore]
// - [Custom] -> [TypedCustom]
// - [Composable] -> [TypedComposable] (TODO not yet implemented)
//
// Each embeds the corresponding untyped interface, and takes a generic type
// parameter. The provided struct is verified to be assignable to the latest
// schema defined in the kind. (See [thema.BindType]) Additional methods are
// provided on Typed* variants that do the same as their untyped counterparts,
// but using the type given in the generic type parameter.
//
// [the kindsys repository]: https://github.com/grafana/kindsys
type Kind interface {
	// Name returns the kind's name, as defined in the name field of the kind definition.
	//
	// Note that this is the capitalized name of the kind. For other names and
	// common kind properties, see [Props.CommonProperties].
	Name() string

	// Props returns a [kindsys.SomeKindProps], representing the properties
	// of the kind as declared in the .cue source. The underlying type is
	// determined by the category of kind.
	Props() SomeKindProperties

	// CurrentVersion returns the version number of the schema that is considered
	// the 'current' version, usually the latest version. When initializing object
	// instances of this Kind, the current version is used by default.
	CurrentVersion() thema.SyntacticVersion

	// Lineage returns the kind's [thema.Lineage]. The lineage contains the full
	// history of object schemas associated with the kind.
	//
	// TODO hide thema away down on the Def
	Lineage() thema.Lineage
}

// Core is the untyped runtime representation of a Grafana core kind definition.
// It is one in a family of interfaces, see [Kind] for context.
//
// A Core kind provides untyped interactions with its corresponding [Resource]
// using [UnstructuredResource].
type Core interface {
	Kind

	// Group returns the kind's group, as defined in the group field of the kind definition.
	//
	// This is equivalent to the group of a Kubernetes CRD.
	Group() string

	// Def returns a wrapper around the underlying CUE value that represents the
	// loaded and validated kind definition.
	Def() Def[CoreProperties]

	// New initializes an object of this kind, represented as an
	// UnstructuredResource and populated with schema-specified defaults.
	New() UnstructuredResource
}

// Custom is the untyped runtime representation of a Grafana core kind definition.
// It is one in a family of interfaces, see [Kind] for context.
//
// A Custom kind provides untyped interactions with its corresponding [Resource]
// using [UnstructuredResource].
//
// Custom kinds are declared in Grafana extensions, rather than in Grafana core. It
// is likely that this distinction will go away in the future, leaving only
// Custom kinds.
type Custom interface {
	Kind

	// Group returns the kind's group, as defined in the group field of the kind definition.
	//
	// This is equivalent to the group field in a Kubernetes CRD.
	Group() string

	// Def returns a wrapper around the underlying CUE value that represents the
	// loaded and validated kind definition.
	Def() Def[CustomProperties]

	// New initializes an object of this kind, represented as an
	// UnstructuredResource and populated with schema-specified defaults.
	New() UnstructuredResource
}

// Composable is the untyped runtime representation of a Grafana core kind definition.
// It is one in a family of interfaces, see [Kind] for context.
//
// TODO sort out the type used for generic associated...objects? do we even need one?
type Composable interface {
	Kind

	// Def returns a wrapper around the underlying CUE value that represents the
	// loaded and validated kind definition.
	Def() Def[ComposableProperties]
}

// TypedCore is the typed runtime representation of a Grafana core kind definition.
// It is one in a family of interfaces, see [Kind] for context.
//
// A TypedCore provides typed interactions with the [Resource] type given as its
// generic type parameter. As it embeds [Core], untyped interaction is also available.
//
// A TypedCore is created by calling [BindCoreResource] on a [Core] with a
// Go type to which it is assignable (see [thema.BindType]).
type TypedCore[R Resource] interface {
	Core

	// NewTyped creates a new object of this kind, represented as the generic type R
	// and populated with schema-specified defaults.
	NewTyped() R
}

// TypedCustom is the typed runtime representation of a Grafana core kind definition.
// It is one in a family of interfaces, see [Kind] for context.
//
// A TypedCustom provides typed interactions with the [Resource] type given as its
// generic type parameter. As it embeds [Custom], untyped interaction is also available.
//
// A TypedCustom is created by calling [BindCustomResource] on a [Custom] with a
// Go type to which it is assignable (see [thema.BindType]).
type TypedCustom[R Resource] interface {
	Custom

	// NewTyped creates a new object of this kind, represented as the generic type R
	// and populated with schema-specified defaults.
	NewTyped() R
}
