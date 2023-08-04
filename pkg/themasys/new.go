package themasys

import (
	"bytes"
	"fmt"
	"io"

	"github.com/grafana/kindsys"
	"github.com/grafana/kindsys/pkg/themasys/encoding"
	"github.com/grafana/thema"
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

func (k *ThemaCoreKind) CurrentVersion() string {
	return k.kind.CurrentVersion().String()
}

func (k *ThemaCoreKind) GetVersions() []kindsys.VersionInfo {
	versions := []kindsys.VersionInfo{}
	for _, schema := range k.kind.Lineage().All() {
		versions = append(versions, kindsys.VersionInfo{
			Version: schema.Version().String(),
		})
	}
	return versions
}

func (k *ThemaCoreKind) GetJSONSchema(version string) (string, error) {
	for _, schema := range k.kind.Lineage().All() {
		if version == schema.Version().String() {
			return "", fmt.Errorf("TODO... convert to JSONSchema")
		}
	}
	return "", fmt.Errorf("unknown version")
}

func (k *ThemaCoreKind) Read(reader io.Reader, strict bool) (kindsys.Resource, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	if strict {
		// ?? is this necessary, or part of the FromBytes below?
		err := k.kind.Validate(buf.Bytes(), &encoding.KubernetesJSONDecoder{})
		if err != nil {
			return nil, err
		}
	}

	return k.kind.FromBytes(buf.Bytes(), &encoding.KubernetesJSONDecoder{})
}

func (k *ThemaCoreKind) Migrate(obj kindsys.Resource, targetVersion string) (kindsys.Resource, error) {
	return nil, fmt.Errorf("TODO implement version migration")
}
