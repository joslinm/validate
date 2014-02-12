package validate

import (
	"fmt"
	"reflect"
)

func ConversionError(got interface{}, expected interface{}) error {
	return fmt.Errorf(`
    Got: %v
    Expecting: %v"`, reflect.TypeOf(got).Kind().String(), expected)
}
