package gpmf

import (
	"fmt"
)

// TypeDef a complex type definition.
type TypeDef string

func parseTypeDef(e *Element) error {
	s, ok := e.Data.(string)
	if !ok {
		return fmt.Errorf("type def: unexpected type %T (expected string)", e.Data)
	}

	d := TypeDef(s)

	e.Data = d
	e.parent.typeDef = s

	return nil
}
