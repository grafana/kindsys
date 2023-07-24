package codegen

import (
	"testing"

	"github.com/grafana/codejen"
	"github.com/matryer/is"
)

func TestPrefixer(t *testing.T) {
	is := is.New(t)
	inputFile := codejen.NewFile("some.file", []byte("with content"))

	resultFile, err := Prefixer("/the/prefix")(*inputFile)
	is.NoErr(err)

	is.Equal(resultFile.RelativePath, "/the/prefix/some.file")
	is.True(resultFile.Data != nil)
}
