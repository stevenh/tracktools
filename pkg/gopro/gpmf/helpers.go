package gpmf

import (
	"fmt"
	"strings"
)

// validateTypeDef validates that e has a type definition which matches size and typeDef.
func validateTypeDef(e *Element, typeDefs map[string]byte) error {
	f := e.friendlyName()

	td, ok := e.Metadata[friendlyName(KeyTypeDef)]
	if !ok {
		return fmt.Errorf("%s: missing type def", f)
	}

	t, ok := td.(string)
	if !ok {
		return fmt.Errorf("%s: unexpected type def type %T (expected string)", f, td)
	}

	types := make([]string, 0, len(typeDefs))
	for k, v := range typeDefs {
		if t == k {
			if e.Header.Size != v {
				return fmt.Errorf("%s: unexpected data size %d (expected %d)", f, e.Header.Size, v)
			}
			return nil
		}
		types = append(types, k)
	}

	return fmt.Errorf("%s: unexpected type def type %q (expected one of %s)", f, t, strings.Join(types, ","))
}
