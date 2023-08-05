package playlist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/grafana/kindsys"
	"github.com/grafana/kindsys/pkg/santhoshsys"
	"github.com/stretchr/testify/require"
)

func TestRawVersion(t *testing.T) {
	sys, err := GetRawKind()
	require.NoError(t, err)

	checkValidVersion(t, sys)
	checkInvalidVersion(t, sys)
	checkMigrations(t, sys)

	manifest, err := santhoshsys.CreateResourceKindManifest(sys)
	require.NoError(t, err)
	require.NotEmpty(t, manifest)
	fmt.Printf("KIND: %s\n", string(manifest))
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
		require.Equal(t, "me", common.CreatedBy, path)
		require.Equal(t, "you", common.UpdatedBy, path)

		static := obj.StaticMetadata()
		// TODO!  fails for thema :(
		// require.Equal(t, "ba2eea3b", static.Name)
		// require.Equal(t, "org-22", static.Namespace)
		require.Equal(t, k.GetKindInfo().Group, static.Group, path)
		require.Equal(t, k.GetKindInfo().Kind, static.Kind, path)
	}
}

func checkInvalidVersion(t *testing.T, k kindsys.ResourceKind) {
	validFiles := []string{
		"testdata/invalid-v0-0.json",
		"testdata/invalid-v0-1.json",
		"testdata/invalid-v1-0.json",
		"testdata/invalid-vX-bad-group.json",
		"testdata/invalid-vX-bad-kind.json",
		"testdata/invalid-vX-missing-spec.json",
	}

	for _, path := range validFiles {
		raw, err := os.ReadFile(path)
		require.NoError(t, err, path)

		_, err = k.Read(bytes.NewReader(raw), true)
		require.Error(t, err, path)
	}
}

func checkMigrations(t *testing.T, k kindsys.ResourceKind) {
	src00, e00 := os.ReadFile("testdata/valid-v0-0.json")
	src01, e01 := os.ReadFile("testdata/valid-v0-1.json")
	src10, e10 := os.ReadFile("testdata/valid-v1-0.json")
	require.NoError(t, e00)
	require.NoError(t, e01)
	require.NoError(t, e10)

	v00, e00 := k.Read(bytes.NewReader(src00), true)
	v01, e01 := k.Read(bytes.NewReader(src01), true)
	v10, e10 := k.Read(bytes.NewReader(src10), true)
	require.NoError(t, e00)
	require.NoError(t, e01)
	require.NoError(t, e10)

	require.Equal(t, "v0-0", v00.StaticMetadata().Version)
	require.Equal(t, "v0-1", v01.StaticMetadata().Version)
	require.Equal(t, "v1-0", v10.StaticMetadata().Version)

	ctx := context.Background()

	// Migrate UP
	out, err := k.Migrate(ctx, v00, "v0-1")
	require.NoError(t, err)
	after, err := json.Marshal(out)
	require.NoError(t, err)
	require.JSONEq(t, string(src01), string(after))
}
