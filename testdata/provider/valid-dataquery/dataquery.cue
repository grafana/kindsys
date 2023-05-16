package provider

// import (
// 	"github.com/grafana/kindsys"
// )

// kindsys.Provider

name:    "grafana-prometheus-datasource"
version: "1.0.0"

composableKinds: DataQuery: {
	name:     "PrometheusQuery"
	maturity: "experimental"

	lineage: {
		name: "PrometheusQuery"
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
