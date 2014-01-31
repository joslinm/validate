package validate

import (
  "errors"
  "fmt"
  "github.com/lann/builder"
  _ "net/http"
  "reflect"
  "regexp"
  "time"
)

// Types validate supports
const (
  STRING = "STRING"
  NUMBER = "NUMBER"
  INT    = "INT"
  FLOAT  = "FLOAT"
  BOOL   = "BOOL"
  TIME   = "TIME"
)

// Type of callbacks to be used in a Rule
type AlterCallback func(value interface{}) interface{}
type PrepareCallback func(value interface{}) interface{}
type CustomCallback func(value interface{}) bool

// Rule encompasses a single validation rule for a parameter
type Rule struct {
  Key      string
  Is       string
  Required bool
  Regex    string
  Min      interface{}
  Max      interface{}

  // callbacks
  Customs  []CustomCallback
  Prepares []PrepareCallback
  Alters   []AlterCallback
}

type ruleBuilder builder.Builder
type RuleBook map[string]ruleBuilder

type ValidationData struct {
  data interface{}
}

func Data(data interface{}) *ValidationData {
  fmt.Println("Validation data: ", data)
  return &ValidationData{data: data}
}

func (v *ValidationData) With(rules RuleBook) (interface{}, error) {
  fmt.Printf("Validation data type: %v", reflect.TypeOf(v.data))
  if _, ok := v.data.(map[string]interface{}); ok {
    return Map(v.data.(map[string]interface{}), rules)
  }
  return nil, errors.New("Unrecognized input data")
}

func (rb ruleBuilder) updateTypeTo(t string) ruleBuilder {
  return builder.Set(rb, "Is", t).(ruleBuilder)
}

func (rb ruleBuilder) Build() Rule {
  s := builder.GetStruct(rb).(Rule)
  fmt.Println("\nGot struct: ", s)
  return builder.GetStruct(rb).(Rule)
}

// required / optional
func (rb ruleBuilder) Required() ruleBuilder {
  return builder.Set(rb, "Required", true).(ruleBuilder)
}

// key
func (rb ruleBuilder) Key(key string) ruleBuilder {
  return builder.Set(rb, "Key", key).(ruleBuilder)
}

// type
func (rb ruleBuilder) Is(is string) ruleBuilder {
  return builder.Set(rb, "Is", is).(ruleBuilder)
}
func (rb ruleBuilder) String() ruleBuilder {
  return builder.Set(rb, "Is", STRING).(ruleBuilder)
}
func (rb ruleBuilder) Number() ruleBuilder {
  return builder.Set(rb, "Is", NUMBER).(ruleBuilder)
}
func (rb ruleBuilder) Bool() ruleBuilder {
  return builder.Set(rb, "Is", BOOL).(ruleBuilder)
}

// regex
func (rb ruleBuilder) Regex(regex string) ruleBuilder {
  b := rb.updateTypeTo(STRING)
  return builder.Set(b, "Regex", regex).(ruleBuilder)
}

// min, max, between
func (rb ruleBuilder) Min(min int) ruleBuilder {
  rb = rb.updateTypeTo(NUMBER)
  return builder.Set(rb, "Min", min).(ruleBuilder)
}

func (rb ruleBuilder) Max(max int) ruleBuilder {
  rb = rb.updateTypeTo(NUMBER)
  return builder.Set(rb, "Max", max).(ruleBuilder)
}

func (rb ruleBuilder) Between(min int, max int) ruleBuilder {
  rb = rb.updateTypeTo(NUMBER)
  builder.Set(rb, "Min", min)
  builder.Set(rb, "Max", max)

  return rb
}

// date
func (rb ruleBuilder) Before(time time.Time) ruleBuilder {
  rb.updateTypeTo(TIME)
  return builder.Set(rb, "Before", time).(ruleBuilder)
}
func (rb ruleBuilder) After(time time.Time) ruleBuilder {
  rb.updateTypeTo(TIME)
  return builder.Set(rb, "After", time).(ruleBuilder)
}
func (rb ruleBuilder) BetweenTimes(min time.Time, max time.Time) ruleBuilder {
  rb.updateTypeTo(NUMBER)
  builder.Set(rb, "Before", max)
  builder.Set(rb, "After", min)

  return rb
}

// callback
func (rb ruleBuilder) Custom(cb CustomCallback) ruleBuilder {
  return builder.Append(rb, "Customs", cb).(ruleBuilder)
}
func (rb ruleBuilder) Alter(cb AlterCallback) ruleBuilder {
  return builder.Append(rb, "Alters", cb).(ruleBuilder)
}
func (rb ruleBuilder) Prepare(cb PrepareCallback) ruleBuilder {
  return builder.Append(rb, "Prepares", cb).(ruleBuilder)
}

// custom
func (rb ruleBuilder) Email() ruleBuilder {
  return builder.Set(rb, "Regex", `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`).(ruleBuilder)
}

var RuleBuilder = builder.Register(ruleBuilder{}, Rule{}).(ruleBuilder)

func (rule *Rule) checkAgainst(given map[string]interface{}) (interface{}, error) {
  fmt.Printf("\nChecking \"%v\"", rule.Key)
  fmt.Printf("\n--------")

  fmt.Printf("\nRequired...")
  val, ok := given[rule.Key]
  if !ok {
    if rule.Required {
      // Throw error indicating value is not in given input
      fmt.Printf("FAIL")
      return false, errors.New("Required key not found")
    } else {
      fmt.Printf("SKIP")
      return true, nil
    }
  } else {
    fmt.Printf("OK")
  }

  fmt.Printf("\nType... %v", reflect.TypeOf(val))

  fmt.Printf("\nCustom...")
  if len(rule.Customs) > 0 {
    for _, cb := range rule.Customs {
      if ok := cb(val); !ok {
        fmt.Printf("FAIL")
        return false, errors.New("Custom failed")
      } else {
        fmt.Printf("OK")
      }
    }
  } else {
    fmt.Printf("SKIP")
  }

  fmt.Printf("\nAlters...")
  if len(rule.Alters) > 0 {
    for _, cb := range rule.Alters {
      val = cb(val)
    }
  } else {
    fmt.Printf("SKIP")
  }

  fmt.Printf("\nRegex...")
  if len(rule.Regex) > 0 {
    re, err := regexp.Compile(rule.Regex)
    if err != nil {
      fmt.Printf("Invalid regex! %v", err)
      return false, errors.New("Invalid regex")
    }
    if re.MatchString(val.(string)) {
      fmt.Printf("OK")
    } else {
      fmt.Printf("FAIL")
      return false, nil
    }
  } else {
    fmt.Printf("SKIP")
  }

  fmt.Printf("\nType...")
  if len(rule.Is) > 0 {
    ok := true
    switch rule.Is {
    case STRING:
      _, ok = val.(string)
      break
    }
    if ok {
      fmt.Printf("OK")
    } else {
      fmt.Printf("FAIL")
    }

  } else {
    fmt.Printf("SKIP")
  }

  fmt.Printf("\n--------\n")

  return val, nil
}

func RuleBookFor(obj interface{}, required bool) RuleBook {
  fmt.Println("\nGot obj ", obj)

  typ := reflect.TypeOf(obj)
  // if a pointer to a struct is passed, get the type of the dereferenced object
  if typ.Kind() == reflect.Ptr {
    typ = typ.Elem()
  }
  fmt.Println("\nGot obj ", typ)
  ruleBook := make(RuleBook)

  for i := 0; i < typ.NumField(); i++ {
    p := typ.Field(i)
    var rb ruleBuilder

    switch p.Type.Kind() {
    case reflect.Float64:
    case reflect.Float32:
    case reflect.Int:
      rb = RuleBuilder.Number()
      break
    case reflect.String:
      rb.String()
      break
    case reflect.Bool:
      rb.Bool()
      break
    default:
      break
    }
    if required {
      ruleBook[p.Name] = rb.Required()
    } else {
      ruleBook[p.Name] = rb
    }
  }

  return ruleBook
}

func Map(given map[string]interface{}, expected RuleBook) (map[string]interface{}, error) {
  retVals := make(map[string]interface{})
  for k, v := range expected {
    fmt.Printf("\n%v = %v", k, reflect.TypeOf(v))
    rule := v.Key(k).Build()
    val, err := rule.checkAgainst(given)
    if err != nil {
      return nil, err
    }
    retVals[k] = val
  }

  return retVals, nil
}
