package kindsys

import (
	"fmt"

	"github.com/grafana/thema"
)

// genericComposable is a general representation of a parsed and validated
// Composable kind.
type genericComposable struct {
	def   Def[ComposableProperties]
	lin   thema.Lineage
	schif SchemaInterface
}

func (k genericComposable) Maturity() Maturity {
	return k.def.Properties.Maturity
}

func (k genericComposable) CurrentVersion() thema.SyntacticVersion {
	return k.def.Properties.CurrentVersion
}

var _ Composable = genericComposable{}

func (k genericComposable) Props() SomeKindProperties {
	return k.def.Properties
}

func (k genericComposable) Name() string {
	return k.def.Properties.Name
}

func (k genericComposable) MachineName() string {
	return k.def.Properties.MachineName
}

func (k genericComposable) Def() Def[ComposableProperties] {
	return k.def
}

func (k genericComposable) Lineage() thema.Lineage {
	return k.lin
}

func (k genericComposable) Implements() SchemaInterface {
	return k.schif
}

// TODO docs
func BindComposable(rt *thema.Runtime, def Def[ComposableProperties], opts ...thema.BindOption) (Composable, error) {
	lin, err := def.Some().BindKindLineage(rt, opts...)
	if err != nil {
		return nil, err
	}

	schif, err := FindSchemaInterface(def.Properties.SchemaInterface)
	if err != nil {
		panic(fmt.Sprintf("unreachable - got %s as string name for schema interface which should have been rejected by declarative validation", def.Properties.SchemaInterface))
	}

	return genericComposable{
		def:   def,
		lin:   lin,
		schif: schif,
	}, nil
}
