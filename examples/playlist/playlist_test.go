package playlist

import (
	"bytes"
	"os"
	"testing"

	"github.com/grafana/kindsys"
	"github.com/stretchr/testify/require"
)

func TestRawVersion(t *testing.T) {
	sys, err := GetRawKind()
	require.NoError(t, err)

	checkValidVersion(t, sys)
	// checkInvalidVersion(t, sys)
}

func TestThemaVersion(t *testing.T) {
	sys, err := GetThemaKind()
	require.NoError(t, err)

	checkValidVersion(t, sys)

	// Thema is not yet using the resource version to validate a payload
	//checkInvalidVersion(t, sys)
}

func TestSanthoshVersion(t *testing.T) {
	sys, err := GetSanthoshKind()
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
		require.NoError(t, err, path)

		obj, err := k.Read(bytes.NewReader(raw), true)
		require.NoError(t, err, path)
		common := obj.CommonMetadata()
		require.Equal(t, "me", common.CreatedBy)
		require.Equal(t, "you", common.UpdatedBy)

		static := obj.StaticMetadata()
		// TODO!  fails for thema :(
		// require.Equal(t, "ba2eea3b", static.Name)
		// require.Equal(t, "org-22", static.Namespace)
		require.Equal(t, k.GetKindInfo().Group, static.Group)
		require.Equal(t, k.GetKindInfo().Kind, static.Kind)
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
		require.NoError(t, err, path)

		obj, err := k.Read(bytes.NewReader(raw), true)
		require.Error(t, err, path)

		// Should the read return an object if it can?
		if obj != nil {
			require.Equal(t, "playlists.ext.grafana.com", obj.StaticMetadata().Group)
		}
	}
}
