package codegen

import (
	"testing"
)

func TestJsonnetSchemaJenny_NoParams(t *testing.T) {
	test := NewGenTest(t, GenTestConfig{
		OutputDir: "testdata/codegen/output/folder_JsonnetSchemaJenny_NoParams",
	})

	test.RunOneToOneFromModule(
		"testdata/codegen/schemas/folder",
		LatestJenny("root", JsonnetSchemaJenny{}),
	)
}
