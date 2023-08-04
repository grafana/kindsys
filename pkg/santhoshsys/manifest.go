package santhoshsys

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/grafana/kindsys"
	"github.com/santhosh-tekuri/jsonschema"
)

// JUST USED while testing things -- with private package
func CreateResourceKindManifest(k kindsys.ResourceKind) ([]byte, error) {
	names := k.GetMachineNames()
	info := &manifest{
		KindInfo:     k.GetKindInfo(),
		Versions:     k.GetVersions(),
		MachineNames: &names,
	}
	return json.MarshalIndent(info, "", "  ")
}

// Internal type that describes what kind of thing we are looking at
type manifest struct {
	kindsys.KindInfo

	// Only valid for resource types
	MachineNames *kindsys.MachineNames `json:"machineName,omitempty"`

	// Only valid for composable types (panel|dataquery|transformer|matcher)
	ComposableType string `json:"composableType,omitempty"`

	// Only valid for composable types
	// ??? do we want/need multiple slots?  should each slot be a different type?
	// Currently: panel => Options | FieldConfig
	ComposableSlots []string `json:"slots,omitempty"`

	// List of version info
	Versions []kindsys.VersionInfo `json:"versions"`
}

var _ kindsys.Kind = &kindFromManifest{}

type kindFromManifest struct {
	info     kindsys.KindInfo
	current  kindsys.VersionInfo
	versions []kindsys.VersionInfo
	raw      map[string]string // raw
	parsed   map[string]*jsonschema.Schema
}

// Load all the schemas
// Currently fails if anything is invalid
func (m *kindFromManifest) init(sfs fs.FS) (*manifest, error) {
	filerc, err := sfs.Open("kind.json")
	if err != nil {
		return nil, fmt.Errorf("unable to find kind manifest")
	}

	buf := bytes.NewBuffer(nil)
	if true {
		defer filerc.Close()
		_, err = buf.ReadFrom(filerc)
		if err != nil {
			return nil, fmt.Errorf("error reading manifest")
		}
	}

	manifest := &manifest{}
	err = json.Unmarshal(buf.Bytes(), manifest)
	if err != nil {
		return manifest, fmt.Errorf("error parsing manifest")
	}
	m.info = manifest.KindInfo
	if m.info.Group == "" {
		return manifest, fmt.Errorf("missing group name")
	}
	if m.info.Kind == "" {
		return manifest, fmt.Errorf("missing kind name")
	}
	if m.info.Maturity <= kindsys.MaturityUnknown {
		return manifest, fmt.Errorf("unknown maturity") // ¯\_(ツ)_/¯
	}

	m.versions = make([]kindsys.VersionInfo, 0)
	m.raw = make(map[string]string)
	m.parsed = make(map[string]*jsonschema.Schema)
	hasher := sha256.New()

	// Load each schema version
	for _, v := range manifest.Versions {
		// TODO, make sure versions are sequential etc
		buf.Reset()
		filerc, err = sfs.Open(v.Version + ".json")
		if err != nil {
			return manifest, fmt.Errorf("error opening schema file: %s // %w", v.Version, err)
		}
		defer filerc.Close()
		_, err = buf.ReadFrom(filerc)
		if err != nil {
			return manifest, fmt.Errorf("error reading schema file: %s // %w", v.Version, err)
		}

		data := buf.Bytes()

		compiler := jsonschema.NewCompiler()
		compiler.AddResource(v.Version, bytes.NewReader(data))
		sch, err := compiler.Compile(v.Version)
		if err != nil {
			return manifest, fmt.Errorf("error parsing schema: %s // %w", v.Version, err)
		}
		m.parsed[v.Version] = sch
		m.raw[v.Version] = string(data)

		hasher.Reset()
		hasher.Write(data)
		hasher.Sum(nil)

		sum := fmt.Sprintf("%x", hasher.Sum(nil))
		if v.Signature != "" {
			if v.Signature != sum {
				return manifest, fmt.Errorf("signature changed for %s", v.Version)
			}
		} else {
			v.Signature = sum
		}
		m.versions = append(m.versions, v)
	}
	m.current = m.versions[len(m.versions)-1]
	return manifest, nil
}

func (m *kindFromManifest) GetKindInfo() kindsys.KindInfo {
	return m.info
}

func (m *kindFromManifest) CurrentVersion() string {
	return m.current.Version
}

func (m *kindFromManifest) GetVersions() []kindsys.VersionInfo {
	return m.versions
}

func (m *kindFromManifest) GetJSONSchema(version string) (string, error) {
	s, ok := m.raw[version]
	if !ok {
		return "", fmt.Errorf("unknown version")
	}
	return s, nil
}
