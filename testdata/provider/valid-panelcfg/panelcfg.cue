package provider

// import (
// 	"github.com/grafana/kindsys"
// )

// kindsys.Provider

name:    "grafana-timeseries-panel"
version: "1.0.0"

composableKinds: PanelCfg: {
	name:     "TimeseriesPanelCfg"
	maturity: "experimental"

	lineage: {
		name: "TimeseriesPanelCfg"
		schemas: [{
			version: [0, 0]
			schema: {
				Options: {
					foo: string
				} @cuetsy(kind="interface")
			}
		}]
	}
}
