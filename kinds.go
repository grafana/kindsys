package kindsys

import (
	"context"
	"io"
)

type KindInfo struct {
	// Organization controlled prefix
	Group string `json:"group"`

	// Must be unique within the group
	Kind string `json:"kind"`

	// Description of the purpose of this kind
	Description string `json:"description"`

	// This kind depends on composable types
	UsesComposition []string `json:"usesComposition,omitempty"`

	// Indicate where this kind is in the dev cycle
	Maturity Maturity `json:"maturity"`
}

type VersionInfo struct {
	// Must be vMajor-Minor-alpha
	// for k8s... we can't have "." and support aggregation api
	Version string `json:"version"`

	// The software version when this schema was released.
	// NOTE(1): the version must follow semantic versioning so that the order is deterministic
	// NOTE(2): panel plugin version is saved in dashboards.  This can be used
	// to find the appropriate schema
	SoftwareVersion string `json:"software"`

	// Human readable descriptions of the changes in this version
	Changelog []string `json:"changelog,omitempty"`

	// The YYYY-MM-DD this version was published (or empty if not yet published)
	Published string `json:"published,omitempty"`

	// JSONSchema hash
	Signature string `json:"signature,omitempty"` // ?? hash of the json schema
}

type MachineNames struct {
	// This is used in k8s URLs
	Plural string `json:"plural,omitempty"`

	// Used as an alias in the display
	Singular string `json:"singular,omitempty"`

	// Optional shorter names that can be matched in a CLI
	Short []string `json:"short,omitempty"`
}

type Kind interface {
	// Get general information about this kind
	GetKindInfo() KindInfo

	// Get the latest version
	CurrentVersion() string

	// Get all versions
	GetVersions() []VersionInfo

	// Return a JSON schema definition for the selected version
	// When composition is required, the slots will have an any node
	// TODO? include an option to have `AnyOf(latest known options)`
	GetJSONSchema(version string) (string, error)
}

// A kind that manages a k8s style resource
type ResourceKind interface {
	Kind

	// K8S style machine names for this kind
	GetMachineNames() MachineNames

	// Read data into a Resource, when strict is true, all validation rules will be checked
	// otherwise a resource will be created if possible, but all validation may not have been run
	Read(reader io.Reader, strict bool) (Resource, error)

	// Migrate from one object to another version
	// NOTE, this may require database calls
	Migrate(ctx context.Context, obj Resource, targetVersion string) (Resource, error)
}

// Function that can be used to migrate resources
type ResourceMigrator = func(ctx context.Context, obj Resource, targetVersion string) (Resource, error)

// Composable kinds are used within other resources -- they do not exist on their own
// NOTE: this interface is still an untested sketch.  We need to implement something real
// before settling on the real method signatures.  However, it is still useful to have the
// common container metadata so that we can produce consistent docs and clients.
type ComposableKind interface {
	Kind

	// eg: panel(options+fieldconfig) | transformation | dataquery | matcher
	GetComposableType() string

	// panel currently has Options + FieldConfig
	// TODO?? can we get rid of slots and just have two composable kinds in the same plugin?
	GetComposableSlots() []string

	// Given an object (at a version) check that it is valid
	Validate(obj any, sourceVersion string) error

	// Migrate from one version of the object to another
	Migrate(ctx context.Context, obj any, sourceVersion string, targetVersion string) (any, error)
}

// Not yet implemented, but will be required for kinds that require composition
type KindRegistry interface {
	// List the objects that can be saved as k8s style resources
	GetResourceKinds(ctx context.Context) []ResourceKind

	// Get composable kinds of a given type
	GetComposableKinds(ctx context.Context) []ComposableKind

	// Get a notification when a kind changes (eg plugin/service registered)
	Watch(ctx context.Context) chan Kind
}
