package kindsys

import "github.com/grafana/thema"

// ResourceMigration is a migration function that transforms a resource from one
// of its schema versions to another.
type ResourceMigration struct {
	// The schema version to which Fn migrates resources
	To thema.SyntacticVersion
	// The schema version from which Fn migrates resources
	From thema.SyntacticVersion
	// The migration function. It is guaranteed that the provided
	// [UnstructuredResource] will be a valid instance of the From schema. The
	// returned [UnstructuredResource] is verified to be a valid instance of the To
	// schema.
	Fn func(resource *UnstructuredResource) (*UnstructuredResource, error)
}

// toLenses wraps the provided [ResourceMigration]s into a slice of imperative lenses that thema can execute.
func toLenses(rms ...ResourceMigration) []thema.ImperativeLens {
	lenses := make([]thema.ImperativeLens, 0, len(rms))
	for _, rm := range rms {
		lenses = append(lenses, thema.ImperativeLens{
			To:   rm.To,
			From: rm.From,
			Mapper: func(inst *thema.Instance, sch thema.Schema) (*thema.Instance, error) {
				data, ur := inst.Underlying(), &UnstructuredResource{}
				data.Decode(ur)
				urm, err := rm.Fn(ur)
				if err != nil {
					return nil, err
				}
				return sch.Validate(data.Context().Encode(urm))
			},
		})
	}
	return lenses
}
