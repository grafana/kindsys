package kindsys

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadProvider(t *testing.T) {
	t.Run("load core-kinds provider", func(t *testing.T) {
		p, err := loadProviderTestCase(t, "core-kinds")
		require.NoError(t, err)
		require.NotNil(t, p)
		require.NotNil(t, p.V)

		require.Equal(t, "core-kinds", p.Name)
		require.Equal(t, "1.0.1", p.Version)
		require.Len(t, p.CoreKinds, 2)
		require.Len(t, p.ComposableKinds, 0)
		require.Len(t, p.CustomKinds, 0)

		fooKind := p.CoreKinds["Foo"]
		require.NotNil(t, fooKind)
		require.Equal(t, "Foo", fooKind.Name())

		barKind := p.CoreKinds["Bar"]
		require.NotNil(t, barKind)
		require.Equal(t, "Bar", barKind.Name())
	})

	t.Run("load custom-kinds provider", func(t *testing.T) {
		p, err := loadProviderTestCase(t, "custom-kinds")
		require.NoError(t, err)
		require.NotNil(t, p)
		require.NotNil(t, p.V)

		require.Equal(t, "custom-kinds", p.Name)
		require.Equal(t, "1.0.2", p.Version)
		require.Len(t, p.CoreKinds, 0)
		require.Len(t, p.ComposableKinds, 0)
		require.Len(t, p.CustomKinds, 2)

		fooKind := p.CustomKinds["Foo"]
		require.NotNil(t, fooKind)
		require.Equal(t, "Foo", fooKind.Name())

		barKind := p.CustomKinds["Bar"]
		require.NotNil(t, barKind)
		require.Equal(t, "Bar", barKind.Name())
	})

	t.Run("load prometheus provider", func(t *testing.T) {
		p, err := loadProviderTestCase(t, "prometheus")
		require.NoError(t, err)
		require.NotNil(t, p)
		require.NotNil(t, p.V)

		require.Equal(t, "prometheus", p.Name)
		require.Equal(t, "1.0.0", p.Version)
		require.Len(t, p.CoreKinds, 0)
		require.Len(t, p.ComposableKinds, 2)
		require.Len(t, p.CustomKinds, 0)

		dqKind := p.ComposableKinds["DataQuery"]
		require.NotNil(t, dqKind)
		require.Equal(t, "PrometheusQuery", dqKind.Name())
		require.Equal(t, "DataQuery", dqKind.Def().Properties.SchemaInterface)

		dscKind := p.ComposableKinds["DataSourceCfg"]
		require.NotNil(t, dscKind)
		require.Equal(t, "PrometheusDataSourceCfg", dscKind.Name())
		require.Equal(t, "DataSourceCfg", dscKind.Def().Properties.SchemaInterface)

		require.NotNil(t, p.Metadata)

		metadata, ok := p.Metadata.(map[string]interface{})
		require.Equal(t, true, ok)
		require.Equal(t, "10.0.0", metadata["grafanaVersion"])
	})

	t.Run("load timeseries provider", func(t *testing.T) {
		p, err := loadProviderTestCase(t, "timeseries")
		require.NoError(t, err)
		require.NotNil(t, p)
		require.NotNil(t, p.V)

		require.Equal(t, "timeseries", p.Name)
		require.Equal(t, "2.0.0", p.Version)
		require.Len(t, p.CoreKinds, 0)
		require.Len(t, p.ComposableKinds, 1)
		require.Len(t, p.CustomKinds, 0)

		pcKind := p.ComposableKinds["PanelCfg"]
		require.NotNil(t, pcKind)
		require.Equal(t, "TimeseriesPanelCfg", pcKind.Name())
		require.Equal(t, "PanelCfg", pcKind.Def().Properties.SchemaInterface)
	})
}

func loadProviderTestCase(t *testing.T, testcase string) (*Provider, error) {
	t.Helper()
	fsys := os.DirFS(fmt.Sprintf("./testdata/provider/%s", testcase))
	require.NotNil(t, fsys)
	return LoadProvider(fsys, nil)
}
