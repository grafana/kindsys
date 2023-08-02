package folder

import (
	"time"
)

// Defines values for StatusOperatorStateState.
const (
	StatusOperatorStateStateFailed     StatusOperatorStateState = "failed"
	StatusOperatorStateStateInProgress StatusOperatorStateState = "in_progress"
	StatusOperatorStateStateSuccess    StatusOperatorStateState = "success"
)

// Folder defines model for Folder.
type Folder struct {
	// metadata contains embedded CommonMetadata and can be extended with custom string fields
	// TODO: use CommonMetadata instead of redefining here; currently needs to be defined here
	// without external reference as using the CommonMetadata reference breaks thema codegen.
	Metadata struct {
		CreatedBy         string     `json:"createdBy"`
		CreationTimestamp time.Time  `json:"creationTimestamp"`
		DeletionTimestamp *time.Time `json:"deletionTimestamp,omitempty"`

		// extraFields is reserved for any fields that are pulled from the API server metadata but do not have concrete fields in the CUE metadata
		ExtraFields     map[string]any    `json:"extraFields"`
		Finalizers      []string          `json:"finalizers"`
		Labels          map[string]string `json:"labels"`
		ResourceVersion string            `json:"resourceVersion"`
		Uid             string            `json:"uid"`
		UpdateTimestamp time.Time         `json:"updateTimestamp"`
		UpdatedBy       string            `json:"updatedBy"`
	} `json:"metadata"`
	Spec struct {
		// Description of the folder.
		Description *string `json:"description,omitempty"`

		// UID of the parent folder.
		Parent *string `json:"parent,omitempty"`

		// Folder title
		Title string `json:"title"`

		// Unique folder id. (will be k8s name)
		Uid string `json:"uid"`
	} `json:"spec"`
	Status struct {
		// additionalFields is reserved for future use
		AdditionalFields map[string]any `json:"additionalFields,omitempty"`

		// operatorStates is a map of operator ID to operator state evaluations.
		// Any operator which consumes this kind SHOULD add its state evaluation information to this field.
		OperatorStates map[string]StatusOperatorState `json:"operatorStates,omitempty"`
	} `json:"status"`
}

// _kubeObjectMetadata is metadata found in a kubernetes object's metadata field.
// It is not exhaustive and only includes fields which may be relevant to a kind's implementation,
// As it is also intended to be generic enough to function with any API Server.
type KubeObjectMetadata struct {
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	DeletionTimestamp *time.Time        `json:"deletionTimestamp,omitempty"`
	Finalizers        []string          `json:"finalizers"`
	Labels            map[string]string `json:"labels"`
	ResourceVersion   string            `json:"resourceVersion"`
	Uid               string            `json:"uid"`
}

// StatusOperatorState defines model for status.#OperatorState.
type StatusOperatorState struct {
	// descriptiveState is an optional more descriptive state field which has no requirements on format
	DescriptiveState *string `json:"descriptiveState,omitempty"`

	// details contains any extra information that is operator-specific
	Details map[string]any `json:"details,omitempty"`

	// lastEvaluation is the ResourceVersion last evaluated
	LastEvaluation string `json:"lastEvaluation"`

	// state describes the state of the lastEvaluation.
	// It is limited to three possible states for machine evaluation.
	State StatusOperatorStateState `json:"state"`
}

// StatusOperatorStateState state describes the state of the lastEvaluation.
// It is limited to three possible states for machine evaluation.
type StatusOperatorStateState string
