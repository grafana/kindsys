export interface spec {
  /**
   * Description of the folder.
   */
  description?: string;
  /**
   * UID of the parent folder.
   */
  parent?: string;
  /**
   * Folder title
   */
  title: string;
  /**
   * Unique folder id. (will be k8s name)
   */
  uid: string;
}

export interface Folder {
  /**
   * metadata contains embedded CommonMetadata and can be extended with custom string fields
   * TODO: use CommonMetadata instead of redefining here; currently needs to be defined here
   * without external reference as using the CommonMetadata reference breaks thema codegen.
   */
  metadata: {
    updateTimestamp: string;
    createdBy: string;
    updatedBy: string;
    uid: string;
    creationTimestamp: string;
    deletionTimestamp?: string;
    finalizers: Array<string>;
    resourceVersion: string;
    /**
     * extraFields is reserved for any fields that are pulled from the API server metadata but do not have concrete fields in the CUE metadata
     */
    extraFields: Record<string, unknown>;
    labels: Record<string, string>;
  };
  spec: {
    /**
     * Unique folder id. (will be k8s name)
     */
    uid: string;
    /**
     * UID of the parent folder.
     */
    parent?: string;
    /**
     * Folder title
     */
    title: string;
    /**
     * Description of the folder.
     */
    description?: string;
  };
  status: {
    /**
     * operatorStates is a map of operator ID to operator state evaluations.
     * Any operator which consumes this kind SHOULD add its state evaluation information to this field.
     */
    operatorStates?: Record<string, {
  /**
   * lastEvaluation is the ResourceVersion last evaluated
   */
  lastEvaluation: string,
  /**
   * state describes the state of the lastEvaluation.
   * It is limited to three possible states for machine evaluation.
   */
  state: ('success' | 'in_progress' | 'failed'),
  /**
   * descriptiveState is an optional more descriptive state field which has no requirements on format
   */
  descriptiveState?: string,
  /**
   * details contains any extra information that is operator-specific
   */
  details?: Record<string, unknown>,
}>;
    /**
     * additionalFields is reserved for future use
     */
    additionalFields?: Record<string, unknown>;
  };
}
