package codegen

import (
	"testing"
)

func TestJsonSchemaJenny_NoParams(t *testing.T) {
	test := NewGenTest(t, GenTestConfig{
		OutputDir: "testdata/codegen/output/folder_JsonSchemaJenny_NoParams",
	})

	test.RunOneToOneFromModule(
		"testdata/codegen/schemas/folder",
		LatestJenny("root", JsonSchemaJenny{}),
	)
}
