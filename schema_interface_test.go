package kindsys_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grafana/thema"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/kindsys"
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
			v := ctx.CompileBytes(data.cue)
			require.NoError(t, v.Validate())

			k, err := verify(ctx, v)
			if err != nil {
				fixErr := strings.Trim(err.Error(), "\n")
				assert.Equal(t, fixErr, data.error)
			} else if data.output != "" {
				json, err := json.MarshalIndent(k.Props(), "", "    ")
				require.NoError(t, err)
				assert.Equal(t, data.output, string(json))
			}
		})
	}
}

func verify(ctx *cue.Context, v cue.Value) (kindsys.Kind, error) {
	rt := thema.NewRuntime(ctx)
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

		return kindsys.BindCore(rt, def)
	case "custom":
		def, err := kindsys.ToDef[kindsys.CustomProperties](v)
		if err != nil {
			return nil, err
		}

		return kindsys.BindCustom(rt, def)
	case "composable":
		def, err := kindsys.ToDef[kindsys.ComposableProperties](v)
		if err != nil {
			return nil, err
		}

		return kindsys.BindComposable(rt, def)
	}

	return nil, errors.New(fmt.Sprintf("unknown package: %s. Only core, custom and composable package names are allowed", pkg))
}

type testData struct {
	name   string
	cue    []byte
	output string
	error  string
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
	var output string
	for _, f := range archive.Files {
		if f.Name == "error" {
			err = strings.TrimSuffix(string(f.Data), "\n")
		}

		if f.Name == "output" {
			output = strings.TrimSuffix(string(f.Data), "\n")
		}
	}

	return testData{
		name:   name,
		cue:    archive.Files[0].Data,
		output: output,
		error:  err,
	}
}
