package kindsys

import "github.com/grafana/thema"

// CommonProperties contains the metadata common to all categories of kinds.
type CommonProperties struct {
	Name              string   `json:"name"`
	PluralName        string   `json:"pluralName"`
	MachineName       string   `json:"machineName"`
	PluralMachineName string   `json:"pluralMachineName"`
	LineageIsGroup    bool     `json:"lineageIsGroup"`
	Maturity          Maturity `json:"maturity"`
	Description       string   `json:"description,omitempty"`
}

// CoreProperties represents the static properties in the definition of a
// Core kind that are representable with basic Go types. This
// excludes Thema schemas.
//
// When .cue file(s) containing a Core definition is loaded through the standard
// [LoadCoreKindDef], func, it is fully validated and populated according to all
// rules specified in CUE for Core kinds.
type CoreProperties struct {
	CommonProperties
	CurrentVersion thema.SyntacticVersion `json:"currentVersion"`
	CRD            struct {
		Group       string `json:"group"`
		Scope       string `json:"scope"`
		DummySchema bool   `json:"dummySchema"`
	} `json:"crd"`
	Slots map[string]Slot `json:"slots"`
}

func (m CoreProperties) _private() {}
func (m CoreProperties) Common() CommonProperties {
	return m.CommonProperties
}

// CustomProperties represents the static properties in the definition of a
// Custom kind that are representable with basic Go types. This
// excludes Thema schemas.
type CustomProperties struct {
	CommonProperties
	CurrentVersion thema.SyntacticVersion `json:"currentVersion"`
	IsCRD          bool                   `json:"isCRD"`
	Group          string                 `json:"group"`
	CRD            struct {
		Group         string  `json:"group"`
		Scope         string  `json:"scope"`
		GroupOverride *string `json:"groupOverride"`
	} `json:"crd"`
	Codegen struct {
		Frontend bool `json:"frontend"`
		Backend  bool `json:"backend"`
	} `json:"codegen"`
	Slots map[string]Slot `json:"slots"`
}

// Slot describes a single composition slot defined in a Core or Custom kind.
type Slot struct {
	// Name is the string that uniquely identifies this slot within the kind.
	Name string `json:"name"`

	// SchemaInterface is the string name of the schema interface that this slot
	// accepts.
	SchemaInterface string `json:"schemaInterface"`
}

func (m CustomProperties) _private() {}
func (m CustomProperties) Common() CommonProperties {
	return m.CommonProperties
}

// ComposableProperties represents the static properties in the definition of a
// Composable kind that are representable with basic Go types. This
// excludes Thema schemas.
type ComposableProperties struct {
	CommonProperties
	CurrentVersion  thema.SyntacticVersion `json:"currentVersion"`
	SchemaInterface string                 `json:"schemaInterface"`
}

func (m ComposableProperties) _private() {}
func (m ComposableProperties) Common() CommonProperties {
	return m.CommonProperties
}

// SomeKindProperties is an interface type to abstract over the different kind
// property struct types: [CoreProperties], [CustomProperties],
// [ComposableProperties].
//
// It is the traditional interface counterpart to the generic type constraint
// KindProperties.
type SomeKindProperties interface {
	_private()
	Common() CommonProperties
}

// KindProperties is a type parameter that comprises the base possible set of
// kind metadata configurations.
type KindProperties interface {
	CoreProperties | CustomProperties | ComposableProperties
}
