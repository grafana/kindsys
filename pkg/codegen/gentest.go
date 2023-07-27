package codegen

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/codejen"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
	"github.com/stretchr/testify/require"
)

type GenTestConfig struct {
	OutputDir string
}

type GenTest struct {
	config GenTestConfig

	t            *testing.T
	themaRuntime *thema.Runtime
}

func NewGenTest(t *testing.T, config GenTestConfig) GenTest {
	return GenTest{
		config:       config,
		t:            t,
		themaRuntime: thema.NewRuntime(cuecontext.New()),
	}
}

func (genTest GenTest) RunManyToOneFromModule(modulePath string, jenny ManyToOne) {
	genTest.Run(func(t *testing.T) codejen.Files {
		req := require.New(t)

		kind, err := genTest.ModuleToCoreKind(modulePath)
		req.NoError(err)

		resultFile, err := jenny.Generate(kind)
		req.NoError(err)

		return codejen.Files{*resultFile}
	})
}

func (genTest GenTest) RunOneToOneFromModule(modulePath string, jenny OneToOne) {
	genTest.Run(func(t *testing.T) codejen.Files {
		req := require.New(t)

		kind, err := genTest.ModuleToCoreKind(modulePath)
		req.NoError(err)

		resultFile, err := jenny.Generate(kind)
		req.NoError(err)

		return codejen.Files{*resultFile}
	})
}

func (genTest GenTest) RunOneToManyFromModule(modulePath string, jenny OneToMany) {
	genTest.Run(func(t *testing.T) codejen.Files {
		req := require.New(t)

		kind, err := genTest.ModuleToCoreKind(modulePath)
		req.NoError(err)

		resultFiles, err := jenny.Generate(kind)
		req.NoError(err)

		return resultFiles
	})
}

func (genTest GenTest) Run(inner func(t *testing.T) codejen.Files) {
	req := require.New(genTest.t)

	updateOutputFiles := os.Getenv("KINDSYS_GEN_UPDATE_GOLDEN_FILES") != ""
	rootCodeJenFS := codejen.NewFS()

	generatedFiles := inner(genTest.t)
	for _, file := range generatedFiles {
		req.NoError(rootCodeJenFS.Add(file))
	}

	if updateOutputFiles {
		req.NoError(rootCodeJenFS.Write(context.Background(), genTest.config.OutputDir))
	} else {
		req.NoError(rootCodeJenFS.Verify(context.Background(), genTest.config.OutputDir))
	}
}

func (genTest GenTest) ModuleToCoreKind(modulePath string) (kindsys.Core, error) {
	overlayFS, err := dirToPrefixedFS(modulePath, "")
	if err != nil {
		return nil, err
	}

	cueInstance, err := kindsys.BuildInstance(genTest.themaRuntime.Context(), ".", "kind", overlayFS)
	if err != nil {
		return nil, fmt.Errorf("could not load kindsys instance: %w", err)
	}

	props, err := kindsys.ToKindProps[kindsys.CoreProperties](cueInstance)
	if err != nil {
		return nil, fmt.Errorf("could not convert cue value to kindsys props: %w", err)
	}

	kindDefinition := kindsys.Def[kindsys.CoreProperties]{
		V:          cueInstance,
		Properties: props,
	}

	boundKind, err := kindsys.BindCore(genTest.themaRuntime, kindDefinition)
	if err != nil {
		return nil, fmt.Errorf("could not bind kind definition to kind: %w", err)
	}

	return boundKind, nil
}

func dirToPrefixedFS(directory string, prefix string) (fs.FS, error) {
	dirHandle, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	commonFS := fstest.MapFS{}
	for _, file := range dirHandle {
		if file.IsDir() {
			continue
		}

		content, err := os.ReadFile(filepath.Join(directory, file.Name()))
		if err != nil {
			return nil, err
		}

		commonFS[filepath.Join(prefix, file.Name())] = &fstest.MapFile{Data: content}
	}

	return commonFS, nil
}
