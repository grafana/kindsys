package playlist

import (
	"embed"

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
