package validate

import (
	"fmt"
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

	// callbacks
	Customs  []CustomCallback
	Prepares []PrepareCallback
	Alters   []AlterCallback

	// using these since max / min get initialized to 0
	DidSetMin bool
	DidSetMax bool
}

/* ****** */
/* Public */
/* ****** */

func (rule *Rule) Process(input interface{}) (bool, []error) {
	// ret values
	var ok bool
	var errors []error

	// route return values by type
	switch rule.Type {
	case Int:
		fallthrough
	case Float:
		fallthrough
	case Number:
		ok, errors = rule.evalNumber(input)
		break
	}

	return ok, errors
}

/* * * * * * * * * * * * *
  Type Eval Functions
* * * * * * * * * * * * */
func (rule *Rule) evalNumber(val interface{}) (bool, []error) {
	var ok bool
	var errors []error

	switch val.(type) {
	case string:
		fmt.Println("Got String input.. converting..")
		float, err := strconv.ParseFloat(val.(string), 64)
		if err == nil {
			fmt.Println("Converted to float")
			ok, errors = rule.evalFloat(float)
		} else {
			integer, err := strconv.ParseInt(val.(string), 2, 64)
			if err != nil {
				fmt.Println("Converted to int")
				ok, errors = rule.evalInt(int(integer))
			}
		}
		break
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
		ok, err := rule.evalNumericMin(val)
		if !ok {
			fmt.Printf("Min failed")
			errors = append(errors, err)
			allOk = false
		}
	}
	if rule.DidSetMax {
		if ok, err := rule.evalNumericMax(val); !ok {
			fmt.Printf("Max failed")
			errors = append(errors, err)
			allOk = false
		}
	}

	return allOk, errors
}

/* * * * * * * * * * * * *
  Rule Eval Functions
* * * * * * * * * * * * */
func (rule *Rule) evalNumericMin(val float64) (bool, error) {
	ok := true
	var err error

	fmt.Printf("Comparing %v > %v", val, rule.Min)
	if val < rule.Min {
		err = fmt.Errorf("Input(%v) is LESS THAN(<) Minimum(%v)", val, rule.Min)
		ok = false
	}

	return ok, err
}

func (rule *Rule) evalNumericMax(val float64) (bool, error) {
	ok := true
	var err error

	if val > rule.Max {
		err = fmt.Errorf("Input(%v) is GREATER THAN(>) Minimum(%v)", val, rule.Max)
		ok = false
	}

	return ok, err
}

/* * * * * * * * * * * * *
  Helper Functions
* * * * * * * * * * * * */
