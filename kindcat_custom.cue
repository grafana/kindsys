package kindsys

import (
	"strings"
	"time"	
)

// _kubeObjectMetadata is metadata found in a kubernetes object's metadata field.
// It is not exhaustive and only includes fields which may be relevant to a kind's implementation,
// As it is also intended to be generic enough to function with any API Server.
_kubeObjectMetadata: {
    uid: string
    creationTimestamp: string & time.Time
    deletionTimestamp?: string & time.Time
    finalizers: [...string]
    resourceVersion: string
    labels: {
        [string]: string
    }
}

// CommonMetadata is a combination of API Server metadata and additional metadata 
// intended to exist commonly across all kinds, but may have varying implementations as to its storage mechanism(s).
CommonMetadata: {
    _kubeObjectMetadata

    updateTimestamp: string & time.Time
    createdBy: string
    updatedBy: string

	// TODO: additional metadata fields?

	// extraFields is reserved for any fields that are pulled from the API server metadata but do not have concrete fields in the CUE metadata
	extraFields: {...}
}

// _crdSchema is the schema format for a CRD.
_crdSchema: {
	// metadata contains embedded CommonMetadata and can be extended with custom string fields
	// TODO: use CommonMetadata instead of redfining here; currently needs to be defined here 
	// without extenal reference as using the CommonMetadata reference breaks thema codegen.
	metadata: {
		_kubeObjectMetadata
		
		updateTimestamp: string & time.Time
		createdBy: string
		updatedBy: string

		// TODO: additional metadata fields?
		// Additional metadata can be added at any future point, as it is allowed to be constant across lineage versions

		// extraFields is reserved for any fields that are pulled from the API server metadata but do not have concrete fields in the CUE metadata
		extraFields: {...}
	} & {
		// All extensions to this metadata need to have string values (for APIServer encoding-to-annotations purposes)
		// Can't use this as it's not yet enforced CUE:
		//...string
		// Have to do this gnarly regex instead
		[!~"^(uid|creationTimestamp|deletionTimestamp|finalizers|resourceVersion|labels|updateTimestamp|createdBy|updatedBy|extraFields)$"]: string
	}
	spec: {...}
	status: {
		#OperatorState: {
			// lastEvaluation is the ResourceVersion last evaluated
			lastEvaluation: string
			// state describes the state of the lastEvaluation.
			// It is limited to three possible states for machine evaluation.
			state: "success" | "in_progress" | "failed"
			// descriptiveState is an optional more descriptive state field which has no requirements on format
			descriptiveState?: string
			// details contains any extra information that is operator-specific
			details?: {...}
		}
		// operatorStates is a map of operator ID to operator state evaluations.
		// Any operator which consumes this kind SHOULD add its state evaluation information to this field.
		operatorStates?: {
			[string]: #OperatorState
		}
		// additionalFields is reserved for future use
		additionalFields: {...}
		...
	}
}

// Custom specifies the kind category for plugin-defined arbitrary types.
// Custom kinds have the same purpose as Core kinds, differing only in
// that they are defined by external plugins rather than in Grafana core. As such,
// this specification is kept closely aligned with the Core kind.
//
// Grafana provides Kubernetes apiserver-shaped HTTP APIs for interacting with custom
// kinds - the same API patterns (and clients) used to interact with k8s CustomResources.
Custom: S={
	_sharedKind

	// pluginID is the unique identifier of the plugin which owns this Custom kind
	pluginID: =~"^([a-z][a-z0-9-]*[a-z0-9])$"

	// isCRD is true if the `crd` trait is present in the kind.
	isCRD: S.crd != _|_

	lineage: { 
		name: S.machineName
	}
	lineageIsGroup: false

	if isCRD {
		// If the crd trait is defined, the schemas in the lineage must follow the format:
		// {
		//     "metadata": CommonMetadata & {...string}
		//     "spec": {...}
		//     "status": {...}
		// }
		lineage: joinSchema: _crdSchema
	}

	// crd contains properties specific to converting this kind to a Kubernetes CRD.
	// Unlike in Core, crd is optional and is used as a signaling mechanism for whether the kind is intended to be registered as a Kubernetes CRD 
	// and/or a resource in a compatible API server. When present, additional structure is enforced on the kind's lineage's schemas.
	// When absent, a lineage's schema has no restrictions as it is assumed that a CRD or similar resource type will not be generated from it.
	// 
	// TODO: rather than `crd`, should this trait be something more generic, as it really indicates more if a resource should be available in a
	// kubernetes-compatible APIServer, not specifically as CRD (though that _is_ an implementation)
	crd?: {
		// groupOverride is an override that is used in the crd trait if present.
		// If left empty, plugin.id is used to generate the group name
		groupOverride?: =~"^([a-z][a-z0-9-]{0,32}[a-z0-9])$"

		// _computedGroups is a list of groups computed from information in the plugin trait.
		// The first element is always the "most correct" one to use.
		// This field could be inlined into `group`, but is separate for clarity.
		_computedGroups: [
			if S.crd.groupOverride != _|_ {
				strings.ToLower(S.crd.groupOverride) + ".apps.grafana.com",
			}
			strings.ToLower(strings.Replace(S.pluginID, "-","_",-1)) + ".apps.grafana.com"
		]

		// group is used as the CRD group name in the GVK.
		// It is computed from information in the plugin trait, using plugin.id unless groupName is specified.
		// The length of the computed group + the length of the name (plus 1) cannot exceed 63 characters for a valid CRD.
		// This length restriction is checked via _computedGroupKind
		group: _computedGroups[0] & =~"^([a-z][a-z0-9_.]{0,61}[a-z0-9])$"

		// _computedGroupKind checks the validity of the CRD kind + group
		_computedGroupKind: S.machineName + "." + group & =~"^([a-z][a-z0-9_.]{0,63}[a-z0-9])$"

		// scope determines whether resources of this kind exist globally ("Cluster") or
		// within Kubernetes namespaces.
		scope: "Cluster" | *"Namespaced"
	}

	// codegen contains properties specific to generating code using tooling
	codegen: {
		// frontend indicates whether front-end TypeScript code should be generated for this kind's schema
		frontend: bool | *true
		// backend indicates whether back-end Go code should be generated for this kind's schema
		backend: bool | *true
	}
}
