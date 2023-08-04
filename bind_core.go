package kindsys

import (
	"fmt"
	"sort"

	"github.com/grafana/thema"
)

// genericCore is a dynamically typed representation of a parsed and
// validated [Core] kind, implemented with thema.
type genericCore struct {
	def Def[CoreProperties]
	lin thema.Lineage

	// map of string name of slot to the currently composed contents of the slot
	composed map[string][]Composable
}

func (k genericCore) Validate(b []byte, codec Decoder) error {
	_, err := bytesToAnyInstance(k, b, codec)
	return err
}

func (k genericCore) CurrentVersion() thema.SyntacticVersion {
	return k.def.Properties.CurrentVersion
}

func (k genericCore) Group() string {
	return k.def.Properties.CRD.Group
}

func (k genericCore) FromBytes(b []byte, codec Decoder) (*UnstructuredResource, error) {
	inst, err := bytesToAnyInstance(k, b, codec)
	if err != nil {
		return nil, err
	}
	// we have a valid instance! decode into unstructured
	return grafanaShapeToUnstructured(k, inst)
}

var _ Core = genericCore{}

func (k genericCore) Props() SomeKindProperties {
	return k.def.Properties
}

func (k genericCore) Name() string {
	return k.def.Properties.Name
}

func (k genericCore) MachineName() string {
	return k.def.Properties.MachineName
}

func (k genericCore) Maturity() Maturity {
	return k.def.Properties.Maturity
}

func (k genericCore) Def() Def[CoreProperties] {
	return k.def
}

func (k genericCore) Lineage() thema.Lineage {
	return k.lin
}

func (k genericCore) Compose(slot Slot, kinds ...Composable) (Core, error) {
	// first, check that this kind supports this slot
	if k.def.Properties.Slots[slot.Name] != slot {
		return nil, &ErrNoSlotInKind{
			Slot: slot,
			Kind: k,
		}
	}

	schif, err := FindSchemaInterface(slot.SchemaInterface)
	if err != nil {
		panic(fmt.Sprintf("unreachable - slot was for nonexistent schema interface %s which should have been rejected at bind time", slot.SchemaInterface))
	}

	// then check that all provided kinds are implementors of the slot
	for _, kind := range kinds {
		if kind.Implements().Name() != schif.Name() {
			return nil, &ErrKindDoesNotImplementInterface{
				Kind:      kind,
				Interface: schif,
			}
		}
	}

	// Inputs look good. Make a copy with our built-up compose map
	com := make(map[string][]Composable)
	for k, v := range k.composed {
		com[k] = v
	}

	var all []Composable
	copy(all, com[slot.Name])
	all = append(all, kinds...)
	// Sort to ensure deterministic output of validation error messages, etc.
	sort.Slice(all, func(i, j int) bool {
		return all[i].Name() < all[j].Name()
	})
	com[slot.Name] = all

	return genericCore{
		def:      k.def,
		lin:      k.lin,
		composed: com,
	}, nil
}

// TODO docs
func BindCore(rt *thema.Runtime, def Def[CoreProperties], opts ...thema.BindOption) (Core, error) {
	lin, err := def.Some().BindKindLineage(rt, opts...)
	if err != nil {
		return nil, err
	}

	return genericCore{
		def: def,
		lin: lin,
	}, nil
}

// TODO docs
func BindCoreResource[R Resource](k Core) (TypedCore[R], error) {
	// TODO implement me
	panic("implement me")
}
