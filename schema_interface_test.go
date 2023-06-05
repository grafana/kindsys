package kindsys_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
	"github.com/rogpeppe/go-internal/txtar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidSchemas(t *testing.T) {
	type tc struct {
		name string
		path string
	}

	var testCases []tc

	err := filepath.Walk("testdata", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.LastIndex(info.Name(), ".txtar") != -1 {
			testCases = append(testCases, tc{
				path: path,
				name: info.Name(),
			})
		}

		return nil
	})
	require.NoError(t, err)

	ctx := cuecontext.New()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := os.ReadFile(tc.path)
			require.NoError(t, err)

			data := getData(t, b)
			v := ctx.CompileBytes(data.CUE)
			require.NoError(t, v.Validate())

			_, err = getKind(ctx, v)
			if err != nil {
				fixErr := strings.Trim(err.Error(), "\n")
				assert.Equal(t, fixErr, data.Error)
			}
		})
	}
}

func getKind(ctx *cue.Context, v cue.Value) (kindsys.Kind, error) {
	instance := v.BuildInstance()
	if instance == nil {
		return nil, errors.New("cannot build instance")
	}

	pkg := instance.Files[0].PackageName()
	switch pkg {
	case "core":
		def, err := kindsys.ToDef[kindsys.CoreProperties](v)
		if err != nil {
			return nil, err
		}

		return kindsys.BindCore(thema.NewRuntime(ctx), def)
	case "custom":
		def, err := kindsys.ToDef[kindsys.CustomProperties](v)
		if err != nil {
			return nil, err
		}

		return kindsys.BindCustom(thema.NewRuntime(ctx), def)
	case "composable":
		def, err := kindsys.ToDef[kindsys.ComposableProperties](v)
		if err != nil {
			return nil, err
		}

		return kindsys.BindComposable(thema.NewRuntime(ctx), def)
	}

	return nil, errors.New(fmt.Sprintf("unknown package: %s", pkg))
}

type testData struct {
	Name  string
	CUE   []byte
	Error string
}

func getData(t *testing.T, b []byte) testData {
	archive := txtar.Parse(b)
	if len(archive.Files) < 1 {
		t.Fatal("It should include at least a cue file")
	}

	name := archive.Files[0].Name
	if !strings.HasSuffix(name, "cue") {
		t.Fatal("First argument should be a cue file")
	}

	var err string
	if len(archive.Files) > 1 {
		if archive.Files[1].Name != "error.out" {
			t.Fatal("Second argument should be error.out file")
		}
		err = strings.TrimSuffix(string(archive.Files[1].Data), "\n")
	}

	return testData{
		Name:  name,
		CUE:   archive.Files[0].Data,
		Error: err,
	}
}
