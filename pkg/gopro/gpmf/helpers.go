package gpmf

import (
	"fmt"
)

// validateTypeDef validates that e has a type definition which matches size and typeDef.
func validateTypeDef(e *Element, size byte, typeDef string) error {
	f := e.friendlyName()
	if e.Header.Size != size {
		return fmt.Errorf("%s: unexpected data size %d (expected %d)", f, e.Header.Size, size)
	}

	td, ok := e.Metadata[friendlyName(KeyTypeDef)]
	if !ok {
		return fmt.Errorf("%s: missing type def", f)
	}

	t, ok := td.(string)
	if !ok {
		return fmt.Errorf("%s: unexpected type def type %T (expected string)", f, td)
	}

	if t != typeDef {
		return fmt.Errorf("%s: unexpected type def type %sq (expected %q)", f, t, faceTypeDef)
	}

	return nil
}
