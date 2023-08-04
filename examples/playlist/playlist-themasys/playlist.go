package playlistthemasys

import (
	"embed"

	"github.com/grafana/kindsys/pkg/themasys"
)

//go:embed *.cue
var packageFS embed.FS

var loaded *themasys.ThemaCoreKind

func GetPlaylistKind() (*themasys.ThemaCoreKind, error) {
	if loaded == nil {
		cue, err := packageFS.ReadFile("playlist.cue")
		if err != nil {
			return nil, err
		}
		k, err := themasys.NewCoreResourceKind(cue)
		if err != nil {
			return nil, err
		}
		loaded = k
	}
	return loaded, nil
}
