package validate

import (
	"github.com/op/go-logging"
	_ "net/http"
	"reflect"
)

type RuleBook map[string]ruleBuilder

var log = logging.MustGetLogger("validate")
var Log = log

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
	log.Debug("Input: %v", given)

	for k, v := range expected {
		rule := v.Build()
		if ok, errors := rule.Process(given[k]); !ok {
			log.Debug(": %v", given)
			paramErrors[k] = errors
		}
		params[k] = v
	}

	return params, paramErrors
}
