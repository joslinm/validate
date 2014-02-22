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
  Before   *time.Time
  After    *time.Time
  In       []string

  // callbacks
  Customs  []CustomCallback
  Prepares []PrepareCallback
  Alters   []AlterCallback

  // using these since max / min get initialized to 0
  DidSetMin bool
  DidSetMax bool
}

// Validates an input
func (rule *Rule) Process(input interface{}) (interface{}, []error) {
  // ret values
  var ok bool
  var errors []error
  var retInput = input

  // type check
  retInput, ok = rule.TypeOkFor(input)
  if !ok { // failed type check
    msg := fmt.Sprintf("Bad input type. Expecting type %v. Got: %v", rule.Type, reflect.TypeOf(retInput))
    errors = append(errors, fmt.Errorf(msg))

    Log.Warning(msg)
    retInput = nil
  } else {
    Log.Info("Input '%v' type is: %v", reflect.ValueOf(retInput), reflect.TypeOf(retInput))

    // route return values by type
    switch rule.Type {
    case Int:
      fallthrough
    case Float:
      fallthrough
    case Number:
      ok, errors = rule.evalNumber(retInput)
      break
    case String:
      ok, errors = rule.evalString(retInput.(string))
      break
    case Bool:
      ok, errors = rule.evalBoolean(retInput.(bool))
      break
    case Time:
      ok, errors = rule.evalTime(retInput.(time.Time))
      break
    }
  }

  Log.Debug("process(...) -> %v, %v", ok, errors)
  return retInput, errors
}

/* * * * * * * * * * * * *
  Type Eval Functions
* * * * * * * * * * * * */

func (rule *Rule) evalTime(val time.Time) (bool, []error) {
  allOk := true
  var errors []error

  if rule.After != nil {
    ok, err := rule.evalAfter(val)
    if !ok {
      Log.Debug("Given time (%v) failed after test", val)
      errors = append(errors, err)
      allOk = false
    } else {
      Log.Debug("Given time (%v) > (%v) -- SUCCESS", val, *rule.After)
    }
  }
  if rule.Before != nil {
    ok, err := rule.evalBefore(val)
    if !ok {
      Log.Debug("Given time (%v) failed after test", val)
      errors = append(errors, err)
      allOk = false
    } else {
      Log.Debug("Given time (%v) < (%v) -- SUCCESS", val, *rule.Before)
    }
  }

  return allOk, errors
}

func (rule *Rule) evalBoolean(val bool) (bool, []error) {
  return true, []error{}
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
func (rule *Rule) evalBefore(val time.Time) (bool, error) {
  ok := true
  var err error

  if val.After(*rule.Before) {
    ok = false
    err = fmt.Errorf("[%v] is AFTER %v (expecting it to be BEFORE)", val, rule.Before)
  }

  return ok, err
}

func (rule *Rule) evalAfter(val time.Time) (bool, error) {
  ok := true
  var err error

  if val.Before(*rule.After) {
    ok = false
    err = fmt.Errorf("[%v] is BEFORE %v (expecting it to be AFTER)", val, rule.After)
  }

  return ok, err
}

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
func (rule *Rule) TypeOkFor(input interface{}) (interface{}, bool) {
  var ok bool

  switch rule.Type {
  case Int:
    _, ok = input.(int)
    if !ok {
      _, ok = input.(int32)
    }
    if !ok {
      _, ok = input.(int64)
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
      _, ok = input.(int)
    }
    if !ok {
      _, ok = input.(int64)
    }
    if !ok {
      _, ok = input.(int32)
    }
    if !ok {
      Log.Warning("Could not convert %v OF TYPE %v to a number!! (tried float32, float64, int32, int64, int)", input, reflect.TypeOf(input))
    } else {
      Log.Debug("Number -> %v", reflect.TypeOf(input))
    }
    break
  case String:
    _, ok = input.(string)
    break
  case Bool:
    _, ok = input.(bool)
    break
  case Time:
    _, ok = input.(time.Time)
    if !ok {
      _, ok = input.(*time.Time)
    }
    break
  }

  if _, isString := input.(string); !ok && isString && rule.Type != String { // try to convert a number/time/boolean string
    Log.Debug("Trying to convert string (%v) to %v", input, rule.Type)
    input = rule.convertString(input.(string))
    if input != nil {
      return rule.TypeOkFor(input)
    }
  }

  return input, ok
}

func (rule *Rule) convertString(input string) interface{} {
  Log.Debug("convertString <- %v", input)
  var converted interface{}
  var ok bool

  switch rule.Type {
  case Bool:
    converted, err := strconv.ParseBool(input)
    Log.Debug("Tried to convert bool '%v' to '%v': succeeded? %v", input, converted, err != nil)
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
  return converted
}

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
