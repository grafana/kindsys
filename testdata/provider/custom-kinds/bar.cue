package provider

customKinds: Bar: {
	name:        "Bar"
	group:       "bar"
	maturity:    "experimental"
	description: "Lorem ipsum..."

	lineage: schemas: [{
		version: [0, 0]
		schema: {
			spec: {
				// the foo
				foo: string
			} @cuetsy(kind="interface")
		}
	}]
}
