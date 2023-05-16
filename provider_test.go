package kindsys

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadProvider(t *testing.T) {
	// ctx := cuecontext.New()
	// ksVal := CUEFramework(ctx)
	// require.Nil(t, ksVal)

	tcs := []struct {
		testcase string
	}{
		{
			testcase: "valid-panelcfg",
		},
		{
			testcase: "valid-dataquery",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.testcase, func(t *testing.T) {
			f := os.DirFS(fmt.Sprintf("./testdata/provider/%s", tc.testcase))
			require.NotNil(t, f)
			bi, err := LoadInstance(fmt.Sprintf("testdata/provider/%s", tc.testcase), "provider", f)
			require.NoError(t, err)
			require.NotNil(t, bi)

			// spew.Dump(bi.Files)

			val := cueContext().BuildInstance(bi)
			require.NoError(t, val.Err())
			require.True(t, val.Exists())

			fmt.Println(val)
		})
	}

	require.Fail(t, "")
}
