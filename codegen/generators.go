package codegen

import (
	"github.com/grafana/codejen"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
)

// OneToOne is a jenny that takes a single Kind and generates a single File.
//
// Most jennies should be OneToOne, as they best follow the single
// responsibility principle of "one input object, one output file".
type OneToOne codejen.OneToOne[kindsys.Kind]

// OneToMany is a jenny that generates multiple files from one kind.
//
// Examples:
//
// - ["github.com/grafana/kindsys/codegen/jenny".LatestMajorsOrXJenny]: A wrapper jenny that takes another Jenny which expects exactly one [thema.Schema], and runs it for the latest in each major.
// - A jenny that takes a Kind and generates one .go file for each top-level field in a Resource (typically Metadata, Spec, Status).
//
// Implement OneToMany jennies with care. Single-responsibility jennies are the
// most composable, and keeping to "one input object, one output file" is the
// simplest way of thinking about single responsibility. Most OneToMany jennies
// should be wrappers around [OneToOne] jennies.
type OneToMany codejen.OneToMany[kindsys.Kind]

// ManyToOne is a jenny that generates a single file from multiple kinds.
//
// These jennies are ideal for generating things like index files which include
// all known kinds.
type ManyToOne codejen.ManyToOne[kindsys.Kind]

// ManyToMany is a jenny that generates multiple files from multiple kinds.
//
// This jenny type generally only represents an entire jenny pipeline. It should
// almost never be implemented directly, for the same reasons as [OneToMany].
type ManyToMany codejen.ManyToMany[kindsys.Kind]

// ForLatestSchema returns a [SchemaForGen] for the latest schema in the
// provided [kindsys.Kind]'s lineage.
//
// TODO this will be replaced by thema-native constructs
func ForLatestSchema(k kindsys.Kind) SchemaForGen {
	comm := k.Props().Common()
	return SchemaForGen{
		Name:    comm.Name,
		Schema:  k.Lineage().Latest(),
		IsGroup: comm.LineageIsGroup,
	}
}

// SchemaForGen is an intermediate values type for jennies that holds both a [thema.Schema],
// and values relevant to generating the schema that should properly, eventually, be in
// thema itself.
//
// TODO this will be replaced by thema-native constructs
type SchemaForGen struct {
	// The PascalCase name of the schematized type.
	Name string
	// The schema to be rendered for the type itself.
	Schema thema.Schema
	// Whether the schema is grouped. See https://github.com/grafana/thema/issues/62
	IsGroup bool
}
