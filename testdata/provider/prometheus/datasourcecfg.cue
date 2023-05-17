package provider

composableKinds: DataSourceCfg: {
	name:     "PrometheusDataSourceCfg"
	maturity: "experimental"

	lineage: {
		name: "PrometheusDataSourceCfg"
		schemas: [{
			version: [0, 0]
			schema: {
				Options: {
					foo: string
				}
				SecureOptions: {
					bar: string
				}
			}
		}]
	}
}
