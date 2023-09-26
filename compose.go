package kindsys

import "fmt"

// ErrNoSlotInKind indicates that it was attempted to [Kind.Compose] into a
// [Slot] for a [Kind] that defines no such slot.
type ErrNoSlotInKind struct {
	Slot Slot
	Kind Kind
}

func (e *ErrNoSlotInKind) Error() string {
	return fmt.Sprintf("no slot named %s in kind %s", e.Slot.Name, e.Kind.Name())
}

type ErrKindDoesNotImplementInterface struct {
	Kind      Composable
	Interface SchemaInterface
}

func (e *ErrKindDoesNotImplementInterface) Error() string {
	return fmt.Sprintf("Composable kind %s does not implement schema interface %s", e.Kind.Name(), e.Interface.name)
}
