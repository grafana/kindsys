package kindsys

// Core specifies the kind category for core-defined arbitrary types.
// Familiar types and functional resources in Grafana, such as dashboards
// and datasources, are represented as core kinds.
Core: S=close({
	_sharedKind
	_rootKind

	lineage: { name: S.machineName, joinSchema: _crdSchema }
	lineageIsGroup: false

	// crd contains properties specific to converting this kind to a Kubernetes CRD.
	crd: {
		// group is used as the CRD group name in the GVK.
		group: "\(S.machineName).core.grafana.com"

		// scope determines whether resources of this kind exist globally ("Cluster") or
		// within Kubernetes namespaces.
		scope: "Cluster" | *"Namespaced"

		// dummySchema determines whether a dummy OpenAPI schema - where the schema is
		// simply an empty, open object - should be generated for the kind.
		//
		// It is a goal that this option eventually be force dto false. Only set to
		// true when Grafana's code generators produce OpenAPI that is rejected by
		// Kubernetes' CRD validation.
		dummySchema: bool | *false

		// deepCopy determines whether a generic implementation of copying should be
		// generated, or a passthrough call to a Go function.
		//   deepCopy: *"generic" | "passthrough"
	}
})
