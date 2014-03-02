package validate

import (
	"github.com/op/go-logging"
	"net/http"
	"reflect"
)

type RuleBook map[string]interface{}

var Log = logging.MustGetLogger("validate")

type ValidationData struct {
	data interface{}
}

func Validate(data map[string]interface{}) *ValidationData {
	return &ValidationData{data: data}
}

func (v *ValidationData) With(rules RuleBook) (map[string]interface{}, map[string][]error) {
	if _, ok := v.data.(*http.Request); ok {
		return Request(v.data.(*http.Request), rules)
	} else {
		return Map(v.data.(map[string]interface{}), rules)
	}
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

func Request(given *http.Request, expected RuleBook) (map[string]interface{}, map[string][]error) {
	// TODO
	return nil, nil
}

func Map(given map[string]interface{}, expected RuleBook) (map[string]interface{}, map[string][]error) {
	params := make(map[string]interface{})
	paramErrors := make(map[string][]error)

	for k, v := range expected {
		builder, ok := v.(ruleBuilder)
		if ok {
			rule := builder.Build()
			input, errors := rule.Process(given[k])
			if errors != nil {
				paramErrors[k] = errors
			} else {
				params[k] = input
			}
		}
	}

	return params, paramErrors
}

func SetLoggingLevel(level logging.Level) {
	logging.SetLevel(level, "validate")
}
