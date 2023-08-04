package playlist

import (
	"embed"
	"io/fs"

	"github.com/grafana/kindsys"
	"github.com/grafana/kindsys/pkg/santhoshsys"
	"github.com/grafana/kindsys/pkg/themasys"
)

//go:embed *.cue schemas/*.json
var packageFS embed.FS

func GetThemaKind() (*themasys.ThemaCoreKind, error) {
	cue, err := packageFS.ReadFile("playlist.cue")
	if err != nil {
		return nil, err
	}
	return themasys.NewCoreResourceKind(cue)
}

func GetSanthoshKind() (kindsys.ResourceKind, error) {
	schemas, err := fs.Sub(packageFS, "schemas")
	if err != nil {
		return nil, err
	}

	return santhoshsys.NewResourceKind(schemas)
}
