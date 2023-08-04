package kindsys

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnstructuredJSON(t *testing.T) {
	simple := &UnstructuredResource{}
	simple.StaticMeta.Group = "ext.something.grafana.com"
	simple.StaticMeta.Version = "v1-1"
	simple.StaticMeta.Kind = "Example"
	simple.StaticMeta.Name = "test"
	simple.StaticMeta.Namespace = "default"
	simple.CommonMeta.CreatedBy = "ryan"
	simple.CommonMeta.UpdatedBy = "ryan"
	simple.CommonMeta.Origin = &ResourceOriginInfo{
		Name:      "file",
		Path:      "path/to/file",
		Key:       "hash",
		Timestamp: ptr(time.Date(2020, time.January, 1, 1, 10, 30, 0, time.UTC)),
	}
	simple.CommonMeta.Finalizers = []string{"a", "b", "c"}
	simple.CommonMeta.ExtraFields = map[string]any{
		"deletionGracePeriodSeconds": 30, // unknown meta fields
	}

	// "deletionGracePeriodSeconds": 30,
	// "annotations": {
	//   "grafana.com/tags": "a,b,c-d",
	//   "grafana.com/title": "A title here"
	// },
	// "finalizers": [
	//   "a",
	//   "b",
	//   "c"
	// ],
	simple.CommonMeta.CreationTimestamp = time.Date(2020, time.January, 21, 1, 10, 30, 0, time.UTC)
	simple.CommonMeta.UpdateTimestamp = time.Date(2022, time.January, 21, 1, 10, 30, 0, time.UTC)
	simple.Spec = map[string]any{
		"hello":  "world",
		"number": 1.234,
		"int":    25,
	}
	simple.Status = map[string]any{
		"hello": "world",
	}

	out, err := json.MarshalIndent(simple, "", "  ")
	require.NoError(t, err)
	fmt.Printf("%s\n", string(out))
	require.JSONEq(t, `{
		"apiVersion": "ext.something.grafana.com/v1-1",
		"kind": "Example",
		"metadata": {
		  "name": "test",
		  "namespace": "default",
		  "annotations": {
			"grafana.com/createdBy": "ryan",
			"grafana.com/updatedBy": "ryan",
			"grafana.com/origin.name": "file",
			"grafana.com/origin.path": "path/to/file",
			"grafana.com/origin.key": "hash",
			"grafana.com/origin.timestamp": "2020-01-01T01:10:30Z",
			"grafana.com/updatedTimestamp": "2022-01-21T01:10:30Z"
		  },
		  "creationTimestamp": "2020-01-21T01:10:30Z",
		  "finalizers": [
			"a",
			"b",
			"c"
		  ],
		  "deletionGracePeriodSeconds": 30
		},
		"spec": {
		  "hello": "world",
		  "int": 25,
		  "number": 1.234
		},
		"status": {
		  "hello": "world"
		}
	  }`, string(out))

	copy := &UnstructuredResource{}
	json.Unmarshal(out, copy)
	require.NoError(t, err)

	after, err := json.MarshalIndent(copy, "", "  ")
	require.NoError(t, err)

	//	fmt.Printf("\nAFTER:\n\n%s\n", string(after))
	require.JSONEq(t, string(out), string(after))
}
