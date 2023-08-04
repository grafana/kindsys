package themasys

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestThemaResource(t *testing.T) {
	var testkind = `
name: "TestKind"
description: "Blammo!"
maturity: "experimental"
lineage: schemas: [{
	version: [0, 0]
	schema: {
		spec: aSpecField: int32
	}
}]
`

	var testresource = `
{
	"apiVersion": "core.grafana.com/v0",
	"kind": "TestKind",
	"metadata": {
		"name": "test",
		"namespace": "default",
		"annotations": {
			"grafana.com/createdBy": "me",
			"grafana.com/updatedBy": "you",
			"grafana.com/updateTimestamp": "2023-07-06T03:08:01Z"
		}
	},
	"spec": {
		"aSpecField": 42
	}
}`

	k, err := NewCoreResourceKind([]byte(testkind))
	require.NoError(t, err)
	require.Equal(t, "TestKind", k.GetKindInfo().Kind)
	require.Equal(t, "Blammo!", k.GetKindInfo().Description)

	res, err := k.Read(bytes.NewReader([]byte(testresource)), true)
	require.NoError(t, err)

	require.Equal(t, "me", res.CommonMetadata().CreatedBy)
	require.Equal(t, "you", res.CommonMetadata().UpdatedBy)
	require.Equal(t, "2023-07-06T03:08:01Z", res.CommonMetadata().UpdateTimestamp.Format(time.RFC3339))
}
