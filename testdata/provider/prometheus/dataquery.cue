package provider

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
