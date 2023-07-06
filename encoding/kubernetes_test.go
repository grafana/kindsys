package encoding

import (
	"encoding/json"
	"fmt"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	testGrafanaSpec = struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}{
		Foo: "foobar",
		Bar: 128,
	}
	testCommonMetadata = commonMetadata{
		UID:             "foo",
		ResourceVersion: "12345",
		Labels: map[string]string{
			"label": "value",
		},
		Finalizers:        []string{"f1"},
		CreatedBy:         "me",
		UpdatedBy:         "you",
		CreationTimestamp: time.Now().Truncate(time.Second).UTC(),
		UpdateTimestamp:   time.Date(2023, time.July, 6, 3, 8, 1, 0, time.UTC),
		ExtraFields: map[string]any{
			"generation": 123,
			"annotations": map[string]string{
				annotationPrefix + "createdBy":       "me",
				annotationPrefix + "updatedBy":       "you",
				annotationPrefix + "updateTimestamp": time.Date(2023, time.July, 6, 3, 8, 1, 0, time.UTC).Format(time.RFC3339),
				annotationPrefix + "field1":          testCustomMetadata.Field1,
				annotationPrefix + "field2":          testCustomMetadata.Field2,
			},
		},
	}
	testCustomMetadata = struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2"`
	}{
		Field1: "foo",
		Field2: "bar",
	}
	testStatusSubresource = struct {
		State string `json:"state"`
	}{
		State: "active",
	}
	testKind    = "Test"
	testGroup   = "test.ext.grafana.com"
	testVersion = "v1-0"

	testGrafanaSpecJSONBytes, _       = json.Marshal(testGrafanaSpec)
	testCommonMetadataJSONBytes, _    = json.Marshal(testCommonMetadata)
	testCustomMetadataJSONBytes, _    = json.Marshal(testCustomMetadata)
	testStatusSubresourceJSONBytes, _ = json.Marshal(testStatusSubresource)

	testKubernetesBytes, _ = json.Marshal(struct {
		Kind       string            `json:"kind"`
		APIVersion string            `json:"apiVersion"`
		Metadata   metav1.ObjectMeta `json:"metadata"`
		Spec       any               `json:"spec"`
		Status     any               `json:"status"`
	}{
		Kind:       testKind,
		APIVersion: fmt.Sprintf("%s/%s", testGroup, testVersion),
		Spec:       testGrafanaSpec,
		Status:     testStatusSubresource,
		Metadata: metav1.ObjectMeta{
			UID:               types.UID(testCommonMetadata.UID),
			ResourceVersion:   testCommonMetadata.ResourceVersion,
			Generation:        123,
			CreationTimestamp: metav1.NewTime(testCommonMetadata.CreationTimestamp),
			Labels:            testCommonMetadata.Labels,
			Finalizers:        testCommonMetadata.Finalizers,
			Annotations: map[string]string{
				annotationPrefix + "createdBy":       testCommonMetadata.CreatedBy,
				annotationPrefix + "updatedBy":       testCommonMetadata.UpdatedBy,
				annotationPrefix + "updateTimestamp": testCommonMetadata.UpdateTimestamp.Format(time.RFC3339),
				annotationPrefix + "field1":          testCustomMetadata.Field1,
				annotationPrefix + "field2":          testCustomMetadata.Field2,
			},
		},
	})
)

func TestKubernetesJSONEncoder_Encode(t *testing.T) {
	emptyJSONErr := json.Unmarshal(nil, &struct{}{})

	tests := []struct {
		name          string
		grafanaBytes  GrafanaShapeBytes
		expectedBytes []byte
		expectedError error
	}{{
		name:          "nil metadata",
		grafanaBytes:  GrafanaShapeBytes{},
		expectedError: fmt.Errorf("unable to parse metadata: %w", emptyJSONErr),
	}, {
		name: "success",
		grafanaBytes: GrafanaShapeBytes{
			Kind:           testKind,
			Group:          testGroup,
			Version:        testVersion,
			Spec:           testGrafanaSpecJSONBytes,
			Metadata:       testCommonMetadataJSONBytes,
			CustomMetadata: testCustomMetadataJSONBytes,
			Subresources: map[string][]byte{
				"status": testStatusSubresourceJSONBytes,
			},
		},
		expectedBytes: testKubernetesBytes,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoder := KubernetesJSONEncoder{}
			res, err := encoder.Encode(test.grafanaBytes)
			if len(test.expectedBytes) > 0 {
				assert.JSONEq(t, string(test.expectedBytes), string(res))
			} else {
				assert.Empty(t, res)
			}
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestKubernetesJSONDecoder_Decode(t *testing.T) {
	tests := []struct {
		name        string
		bytes       []byte
		expected    GrafanaShapeBytes
		expectedErr error
	}{{
		name:  "success",
		bytes: testKubernetesBytes,
		expected: GrafanaShapeBytes{
			Kind:           testKind,
			Group:          testGroup,
			Version:        testVersion,
			Spec:           testGrafanaSpecJSONBytes,
			Metadata:       testCommonMetadataJSONBytes,
			CustomMetadata: testCustomMetadataJSONBytes,
			Subresources: map[string][]byte{
				"status": testStatusSubresourceJSONBytes,
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decoder := KubernetesJSONDecoder{}
			res, err := decoder.Decode(test.bytes)
			assert.Equal(t, test.expected, res)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
