package codegen

import (
	"encoding/json"
	"strings"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/codejen"
	"github.com/grafana/thema/encoding/jsonschema"
)

type JsonnetSchemaJenny struct{}

func (j JsonnetSchemaJenny) JennyName() string {
	return "JsonnetImportsCoreJenny"
}

func (j JsonnetSchemaJenny) Generate(sfg SchemaForGen) (*codejen.File, error) {
	// TODO allow using name instead of machine name in thema generator
	ast, err := jsonschema.GenerateSchema(sfg.Schema)
	if err != nil {
		return nil, err
	}
	ctx := cuecontext.New()
	str, err := json.MarshalIndent(ctx.BuildFile(ast), "", "  ")
	if err != nil {
		return nil, err
	}

	// @TODO we should be receiving a name without schema interface type so that we don't
	// need to strip it with a hack like this:
	name := jsonnetFixKindName(sfg.Schema.Lineage().Name())
	return codejen.NewFile(name+"_types_gen.json", []byte(str), j), nil
}

func jsonnetFixKindName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "panelcfg", "")
	name = strings.ReplaceAll(name, "dataquery", "")
	return name
}
