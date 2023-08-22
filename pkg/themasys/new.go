package themasys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/kindsys"
	"github.com/grafana/kindsys/pkg/themasys/encoding"
	"github.com/grafana/thema"
	"github.com/grafana/thema/encoding/jsonschema"
	"k8s.io/kube-openapi/pkg/common"
)

var _ kindsys.ResourceKind = &ThemaCoreKind{}

type ThemaCoreKind struct {
	kind Core
}

// Load a new kind based on a cue definition
func NewCoreResourceKind(cuefile []byte) (*ThemaCoreKind, error) {
	rt := thema.NewRuntime(ctx)
	cv := ctx.CompileBytes(cuefile)
	def, err := ToDef[CoreProperties](cv)
	if err != nil {
		return nil, err
	}
	k, err := BindCore(rt, def)
	if err != nil {
		return nil, err
	}
	return &ThemaCoreKind{kind: k}, nil
}

func (k *ThemaCoreKind) CoreKind() Core {
	return k.kind
}

func (k *ThemaCoreKind) GetMachineNames() kindsys.MachineNames {
	p := k.kind.Props()
	c := p.Common()
	return kindsys.MachineNames{
		Plural:   c.PluralName,
		Singular: c.MachineName,
	}
}

func (k *ThemaCoreKind) GetKindInfo() kindsys.KindInfo {
	p := k.kind.Props()
	c := p.Common()
	return kindsys.KindInfo{
		Group:       k.kind.Group(),
		Kind:        c.Name,
		Description: c.Description,
	}
}

func syntaxVersionToString(v thema.SyntacticVersion) string {
	return fmt.Sprintf("v%d-%d", v[1], v[1])
}

func (k *ThemaCoreKind) CurrentVersion() string {
	return k.kind.CurrentVersion().String()
}

func (k *ThemaCoreKind) GetVersions() []kindsys.VersionInfo {
	versions := []kindsys.VersionInfo{}
	// TODO??? this only gets the first version?
	for _, schema := range k.kind.Lineage().All() {
		versions = append(versions, kindsys.VersionInfo{
			Version: syntaxVersionToString(schema.Version()),
		})
	}
	return versions
}

// Converts the cue to JSON schema
func (k *ThemaCoreKind) GetOpenAPIDefinition(version string, ref common.ReferenceCallback) (common.OpenAPIDefinition, error) {
	api := common.OpenAPIDefinition{}

	for _, schema := range k.kind.Lineage().All() {
		if version == syntaxVersionToString(schema.Version()) {
			ast, err := jsonschema.GenerateSchema(schema)
			if err != nil {
				return api, err
			}
			ctx := cuecontext.New()
			out, err := json.MarshalIndent(ctx.BuildFile(ast), "", "  ")
			if err != nil {
				return api, err
			}

			// TODO... not quite the same!
			return kindsys.LoadOpenAPIDefinition(out)
		}
	}
	return api, fmt.Errorf("unknown version")
}

func (k *ThemaCoreKind) Read(reader io.Reader, strict bool) (kindsys.Resource, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	// if strict {
	// 	// ?? is this necessary, or part of the FromBytes below?
	// 	err := k.kind.Validate(buf.Bytes(), &encoding.KubernetesJSONDecoder{})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return k.kind.FromBytes(buf.Bytes(), &encoding.KubernetesJSONDecoder{})
}

func (k *ThemaCoreKind) Migrate(ctx context.Context, obj kindsys.Resource, targetVersion string) (kindsys.Resource, error) {
	return nil, fmt.Errorf("TODO implement version migration")
}
