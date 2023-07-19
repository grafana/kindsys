package kindsys

import (
	"github.com/grafana/kindsys/encoding"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/grafana/thema"
)

func TestFromBytes(t *testing.T) {
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

	rt := thema.NewRuntime(ctx)

	cv := ctx.CompileString(testkind)
	def, err := ToDef[CoreProperties](cv)
	require.NoError(t, err)

	k, err := BindCore(rt, def)
	require.NoError(t, err)

	res, err := k.FromBytes([]byte(testresource), &encoding.KubernetesJSONDecoder{})
	require.NoError(t, err)

	require.Equal(t, "me", res.CommonMeta.CreatedBy)
	require.Equal(t, "you", res.CommonMeta.UpdatedBy)
	require.Equal(t, "2023-07-06T03:08:01Z", res.CommonMeta.UpdateTimestamp.Format(time.RFC3339))
}
