package playlist

import (
	"bytes"
	"os"
	"testing"

	playlistthemasys "github.com/grafana/kindsys/examples/playlist/playlist-themasys"
	"github.com/stretchr/testify/require"
)

func TestThemaVersion(t *testing.T) {
	// This only supports
	k, err := playlistthemasys.GetPlaylistKind()
	require.NoError(t, err)

	// read valid file
	raw, err := os.ReadFile("testdata/valid-v0-0.json")
	require.NoError(t, err)

	obj, err := k.Read(bytes.NewReader(raw), true)
	require.NoError(t, err)
	common := obj.CommonMetadata()
	require.Equal(t, "me", common.CreatedBy)
	require.Equal(t, "you", common.UpdatedBy)
}
