package provider

customKinds: Foo: {
	name:        "Foo"
	group:       "foo"
	maturity:    "merged"
	description: "Lorem ipsum"

	lineage: schemas: [{
		version: [0, 0]
		schema: {
			spec: {
				// the bar
				bar: string
			} @cuetsy(kind="interface")
		}
	}]
}
