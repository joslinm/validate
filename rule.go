package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

// Types
const (
	Unknown = iota
	Int
	Float
	Number
	Bool
	String
	Time
)

//

// Type of callbacks to be used in a Rule
type AlterCallback func(value interface{}) interface{}
type PrepareCallback func(value interface{}) interface{}
type CustomCallback func(value interface{}) bool

// Rule encompasses a single validation rule for a parameter
type Rule struct {
	// validations
	Key      string
	Type     int
	Required bool
	Regex    string
	Message  string
	Min      float64
	Max      float64
	Before   time.Time
	After    time.Time
	In       []string

	// callbacks
	Customs  []CustomCallback
	Prepares []PrepareCallback
	Alters   []AlterCallback

	// using these since max / min get initialized to 0
	DidSetMin bool
	DidSetMax bool
}

func (rule *Rule) typeOkFor(input interface{}) bool {
	var ok bool

	switch rule.Type {
	case Int:
		_, ok = input.(int64)
		if !ok {
			_, ok = input.(int32)
		}
		break
	case Float:
		_, ok = input.(float64)
		if !ok {
			_, ok = input.(float32)
		}
		break
	case Number:
		_, ok = input.(float64)
		if !ok {
			_, ok = input.(float32)
		}
		if !ok {
			_, ok = input.(int64)
		}
		if !ok {
			_, ok = input.(int32)
		}
		break
	case String:
		_, ok = input.(string)
		break
	case Bool:
		_, ok = input.(bool)
		break
	}

	return ok
}

func (rule *Rule) convertString(input string) (interface{}, bool) {
	Log.Debug("convertString <- %v", input)

	var converted interface{}
	var ok bool

	switch rule.Type {
	case Bool:
		val, err := strconv.ParseBool(input)
		ok = err == nil
		converted = val
		Log.Debug("Tried to convert bool '%v' to '%v': succeeded? %v", input, converted, ok)
		break
	case Int:
		fallthrough
	case Float:
		fallthrough
	case Number:
		var num interface{}
		var t int
		num, t, ok = convertStringToNumber(input)
		if !ok {
			Log.Error("Could not convert string '%v' to number!", input)
		} else {
			switch t {
			case Int:
				converted = num.(int)
				break
			case Float:
				converted = num.(float64)
				break
			}
			Log.Debug("Converted %v to %v", input, converted)
		}
		break
	}

	Log.Debug("convertString -> %v, %v", converted, ok)
	return converted, ok
}

// Validates an input
func (rule *Rule) Process(input interface{}) (bool, []error) {
	// ret values
	var ok bool
	var errors []error

	// type check
	if !rule.typeOkFor(input) {
		// try to convert a string first
		_, ok := input.(string)
		if ok {
			input, ok = rule.convertString(input.(string))
		}
		if !ok {
			msg := fmt.Sprintf("Bad input type. Expecting type %v", rule.Type)
			Log.Warning(msg)
			errors = append(errors, fmt.Errorf(msg))
			return false, errors
		}
	}
	Log.Info("Input '%v' type is: %v", reflect.ValueOf(input), reflect.TypeOf(input))

	// route return values by type
	switch rule.Type {
	case Int:
		fallthrough
	case Float:
		fallthrough
	case Number:
		ok, errors = rule.evalNumber(input)
		break
	case String:
		ok, errors = rule.evalString(input.(string))
		break
	case Bool:
		ok, errors = rule.evalBoolean(input.(bool))
		break
	}

	Log.Debug("process(...) -> %v, %v", ok, errors)
	return ok, errors
}

/* * * * * * * * * * * * *
  Type Eval Functions
* * * * * * * * * * * * */

func (rule *Rule) evalBoolean(val bool) (bool, []error) {
	var errors []error
	return true, errors
}

func (rule *Rule) evalString(val string) (bool, []error) {
	allOk := true
	var errors []error

	Log.Debug("Length of regex %v (%v)", len(rule.Regex), rule.Regex)
	if len(rule.Regex) > 0 {
		ok, err := rule.evalRegex(val)
		if !ok {
			Log.Debug("Regex failed")
			errors = append(errors, err)
			allOk = false
		} else {
			Log.Debug("Regex succeeded")
		}
	}
	if len(rule.In) > 0 {
		ok, err := rule.evalIn(val)
		if !ok {
			Log.Debug("In failed")
			errors = append(errors, err)
			allOk = false
		} else {
			Log.Debug("In succeeded")
		}
	}

	return allOk, errors
}

func (rule *Rule) evalNumber(val interface{}) (bool, []error) {
	var ok bool
	var errors []error

	switch val.(type) {
	case int32:
		ok, errors = rule.evalInt(int(val.(int32)))
	case int64:
		ok, errors = rule.evalInt(int(val.(int64)))
	case int:
		ok, errors = rule.evalInt(val.(int))
	case float32:
		ok, errors = rule.evalFloat(val.(float64))
	case float64:
		ok, errors = rule.evalFloat(val.(float64))
	}

	return ok, errors
}

func (rule *Rule) evalInt(val int) (bool, []error) {
	return rule.evalFloat(float64(val))
}

func (rule *Rule) evalFloat(val float64) (bool, []error) {
	allOk := true
	var errors []error

	if rule.DidSetMin {
		ok, err := rule.evalMin(val)
		if !ok {
			errors = append(errors, err)
			allOk = false
		}
	}
	if rule.DidSetMax {
		if ok, err := rule.evalMax(val); !ok {
			errors = append(errors, err)
			allOk = false
		}
	}
	Log.Debug("evalFloat(...) -> %v, %v", allOk, errors)
	return allOk, errors
}

/* * * * * * * * * * * * *
  Rule Eval Functions
* * * * * * * * * * * * */
func (rule *Rule) evalIn(val string) (bool, error) {
	ok := false
	err := fmt.Errorf("[%v] not in %v", val, rule.In)

	Log.Debug("Looking up [%v] in %v", val, rule.In)
	for _, inVal := range rule.In {
		if inVal == val {
			ok = true
			err = nil
			break
		}
	}

	return ok, err

}

func (rule *Rule) evalRegex(val string) (bool, error) {
	ok := true
	var err error

	Log.Debug("Validating %v =~ %v", val, rule.Regex)
	expr, err := regexp.Compile(rule.Regex)
	if err == nil {
		// check regex
		if k := expr.MatchString(val); !k {
			Log.Debug("Failed regex")
			err = fmt.Errorf("[%v] did not match regex [%v]", val, rule.Regex)
			ok = false
		} else {
			Log.Debug("Passed regex")
		}
	}

	return ok, err
}

func (rule *Rule) evalMin(val float64) (bool, error) {
	ok := true
	var err error

	Log.Debug("Validating %v > %v...", val, rule.Min)
	if val < rule.Min {
		err = fmt.Errorf("Input(%v) < Minimum(%v)", val, rule.Min)
		ok = false
	}

	return ok, err
}

func (rule *Rule) evalMax(val float64) (bool, error) {
	ok := true
	var err error

	if val > rule.Max {
		err = fmt.Errorf("Input(%v) > Maximum(%v)", val, rule.Min)
		ok = false
	}

	return ok, err
}

/* * * * * * * * * * * * *
  Helper Functions
* * * * * * * * * * * * */

func convertStringToNumber(val string) (interface{}, int, bool) {
	Log.Debug("convertStringToNumber <- %v", val)

	var num interface{}
	var t = 0
	var ok = false

	// first try to convert to float
	for size := range []int{64, 32} { // try each size
		float, err := strconv.ParseFloat(val, size)
		if err == nil { // float success
			num, t, ok = float, Float, true
		} else { // try int
			integer, err := strconv.ParseInt(val, 2, size)
			if err == nil {
				num, t, ok = integer, Int, true
			}
		}
	}

	Log.Debug("convertStringToNumber -> %v, %v, %v", num, t, ok)
	return num, t, ok
}
