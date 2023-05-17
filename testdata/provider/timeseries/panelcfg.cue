package provider

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
