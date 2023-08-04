package playlist

import (
	"bytes"
	"os"
	"testing"

	"github.com/grafana/kindsys"
	"github.com/stretchr/testify/require"
)

func TestThemaVersion(t *testing.T) {
	sys, err := GetThemaKind()
	require.NoError(t, err)

	checkValidVersion(t, sys)

	// Thema is not yet using the resource version to validate a payload
	// checkInvalidVersion(t, themasys)
}

func TestSanthoshVersion(t *testing.T) {
	sys, err := GetThemaKind()
	require.NoError(t, err)

	checkValidVersion(t, sys)
	checkInvalidVersion(t, sys)
}

func checkValidVersion(t *testing.T, k kindsys.ResourceKind) {
	validFiles := []string{
		"testdata/valid-v0-0.json",
		"testdata/valid-v0-1.json",
		"testdata/valid-v1-0.json",
	}

	for _, path := range validFiles {
		raw, err := os.ReadFile(path)
		require.NoError(t, err)

		obj, err := k.Read(bytes.NewReader(raw), true)
		require.NoError(t, err)
		common := obj.CommonMetadata()
		require.Equal(t, "me", common.CreatedBy)
		require.Equal(t, "you", common.UpdatedBy)
	}
}

func checkInvalidVersion(t *testing.T, k kindsys.ResourceKind) {
	validFiles := []string{
		"testdata/invalid-v0-0.json",
		"testdata/invalid-v0-1.json",
		"testdata/invalid-v1-0.json",
	}

	for _, path := range validFiles {
		raw, err := os.ReadFile(path)
		require.NoError(t, err)

		obj, err := k.Read(bytes.NewReader(raw), true)
		require.Error(t, err, path)

		// Should the read return an object if it can?
		if obj != nil {
			require.Equal(t, "playlists.ext.grafana.com", obj.StaticMetadata().Group)
		}
	}
}
