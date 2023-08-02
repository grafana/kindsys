package codegen

import (
	"encoding/json"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/codejen"
	"github.com/grafana/thema/encoding/jsonschema"
)

type JsonSchemaJenny struct{}

func (j JsonSchemaJenny) JennyName() string {
	return "JsonSchemaJenny"
}

func (j JsonSchemaJenny) Generate(sfg SchemaForGen) (*codejen.File, error) {
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

	return codejen.NewFile(sfg.Schema.Lineage().Name()+"_types_gen.json", str, j), nil
}
