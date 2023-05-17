package provider

coreKinds: Foo: {
  name:        "Foo"
  maturity:    "merged"
  description: "Lorem ipsum"

  lineage: schemas: [{
    version: [0, 0]
    schema: spec: {
      // the bar
      bar: string
    } @cuetsy(kind="interface")
  }]
}
