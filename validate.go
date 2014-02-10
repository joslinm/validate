package validate

import (
	_ "net/http"
	"reflect"
)

type RuleBook map[string]ruleBuilder

type ValidationData struct {
	data map[string]interface{}
}

func Data(data map[string]interface{}) *ValidationData {
	return &ValidationData{data: data}
}

func (v *ValidationData) With(rules RuleBook) (interface{}, map[string][]error) {
	return Map(v.data, rules)
}

func sameType(vals ...interface{}) bool {
	expectedType := reflect.TypeOf(vals[0]).Kind()
	for _, val := range vals {
		if expectedType != reflect.TypeOf(val).Kind() {
			return false
		}
	}

	return true
}

func Map(given map[string]interface{}, expected RuleBook) (map[string]interface{}, map[string][]error) {
	params := make(map[string]interface{})
	paramErrors := make(map[string][]error)

	for k, v := range expected {
		rule := v.Build()
		if ok, errors := rule.Process(given[k]); !ok {
			paramErrors[k] = errors
		}
		params[k] = v
	}

	return params, paramErrors
}
