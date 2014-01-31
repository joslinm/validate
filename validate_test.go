package validate_test

import (
  "fmt"
  "github.com/joslinm/validate"
  "testing"
)

// request

func TestValidateDataChecksCustoms(t *testing.T) {
  did_check_callback := false
  required := validate.RuleBuilder.Required()
  expecting := validate.RuleBook{
    "string": required.Regex(`\w+`).Custom(
      func(value interface{}) bool {
        did_check_callback = true
        return true
      },
    ),
    "int": required.Number(),
  }

  _, err := validate.Data(map[string]interface{}{"int": 123}).With(expecting)
  if !did_check_callback {
    fmt.Printf("\nDid not check callback")
    t.Fail()
  }
  if err != nil {
    fmt.Printf("\nError: %v", err)
    t.Fail()
  }

}

func TestRuleBuilder(t *testing.T) {
  rule := validate.RuleBuilder.Required().Number().Build()
  fmt.Printf("\nrule = %v", rule)
  if rule.Required != true {
    fmt.Printf("\nrule.required = %v", rule.Required)
    fmt.Printf("\nexpected = %v\n", true)
    t.Fail()
  }
  if rule.Is != validate.NUMBER {
    t.Fail()
  }
}

func TestRuleBuilderDifferentRules(t *testing.T) {
  rule := validate.RuleBuilder

  rules := validate.RuleBook{
    "lc": rule.Regex("[a-z]+"),
    "uc": rule.Regex("[A-Z]+"),
  }
  params, err := validate.Data(map[string]interface{}{
    "lc":  "mark",
    "uc":  "JOSLIN",
    "age": 25,
  }).With(rules)

  if err != nil {
    fmt.Printf("\nError: %v", err)
    t.Fail()
  }

  fmt.Printf("\nParams:\n%v", params)
}

func TestRuleBookForStruct(t *testing.T) {
  type UserParams struct {
    // you must export struct declarations
    first_name string
    last_name  string
    age        int
  }

  rules := validate.RuleBookFor(&UserParams{}, true)
  rules["first_name"].Alter(func(value interface{}) interface{} {
    fmt.Printf("Value: %v", value)
    return value
  })
  params, _ := validate.Data(map[string]interface{}{
    "first_name": "Mark",
    "last_name":  "Joslin",
    "age":        25,
  }).With(rules)
  fmt.Printf("Got params.. %v", params)
}

func TestNestedRuleBooks(t *testing.T) {
}
