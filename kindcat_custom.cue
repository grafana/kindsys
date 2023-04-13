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
    finalizers: [string]
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

	// appID is the unique identifier of the app which is the owner of this Custom kind.
	// TODO: should this just be pluginID? Is there a world where an app and plugin may not share the same ID?
	// The practical reason to have it different is to allow the appID to be shorter than the pluginID as the group
	// is generated from the appID, and <kind name>+<group> cannot exceed 63 characters in kubernetes.
	appID: =~"^([a-z][a-z0-9-]{0,61}[a-z0-9])$"

	lineage: { 
		name: S.machineName 
		// If the crd trait is defined, the schemas in the lineage must follow the format:
		// {
		//     "metadata": CommonMetadata & {...string}
		//     "spec": {...}
		//     "status": {...}
		// }
		if S.crd != _|_ {
			joinSchema: {
				metadata: CommonMetadata & {
					// All extensions to this metadata need to have string values (for APIServer encoding-to-annotations purposes)
					// Can't use this as it's not yet enforced CUE:
					//...string
					// Have to do this gnarly regex instead
					[!~"^(uid|creationTimestamp|deletionTimestamp|finalizers|resourceVersion|labels|updateTimestamp|createdBy|updatedBy)$"]: string
				}
				spec: {...}
				status: {...}
			}
		}
	}
	lineageIsGroup: false

	// crd contains properties specific to converting this kind to a Kubernetes CRD.
	// Unlike in Core, crd is optional and is used as a signaling mechanism for whether the kind is intended to be registered as a Kubernetes CRD 
	// and/or a resource in a compatible API server. When present, additional structure is enforced on the kind's lineage's schemas.
	// When absent, a lineage's schema has no restrictions as it is assumed that a CRD or similar resource type will not be generated from it.
	// 
	// TODO: rather than `crd`, should this trait be something more generic, as it really indicates more if a resource should be available in a
	// kubernetes-compatible APIServer, not specifically as CRD (though that _is_ an implementation)
	crd?: {
		// group is used as the CRD group name in the GVK.
		group: strings.ToLower(S.appID) + ".apps.grafana.com"

		// scope determines whether resources of this kind exist globally ("Cluster") or
		// within Kubernetes namespaces.
		scope: "Cluster" | *"Namespaced"

		// deepCopy determines whether a generic implementation of copying should be
		// generated, or a passthrough call to a Go function.
		//   deepCopy: *"generic" | "passthrough"
	}

	// codegen contains properties specific to generating code using tooling
	codegen: {
		// frontend indicates whether front-end TypeScript code should be generated for this kind's schema
		frontend: bool | *true
		// backend indicates whether back-end Go code should be generated for this kind's schema
		backend: bool | *true
	}
}
