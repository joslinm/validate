package validate

import (
	"errors"
	"log"
	"reflect"
	"regexp"
)

// Rule encompasses a single validation rule for a parameter
type Rule struct {
	// validations
	Key      string
	Type     int
	Required bool
	Regex    string
	Message  string
	Min      interface{}
	Max      interface{}

	// callbacks
	Customs  []CustomCallback
	Prepares []PrepareCallback
	Alters   []AlterCallback

	// input / output
	Input   map[string]interface{}
	Results []Result
}

/* Public */

// ProcessWith(...) set the rule's input and processes it
func (rule *Rule) ProcessWith(val map[string]interface{}) []Result {
	rule.Input = val
	return rule.Process()
}

// Process() examines its input & output and creates
// an array of `Result` structs. This array gets bubbled
// up to `validate.Data(input).With(rules), which means
// it DOES NOT PANIC. The library as a whole follows a convention
// to not panic unless a programmer error is recognized (e.g.,
// clearly setting the wrong type or setting disparate min/max
// values). This is meant to give the programmer the means
// to bundle up the errors himself.
func (rule *Rule) Process() []Result {
	input := rule.Input
	if input == nil {
		panic("Tried to process a rule without an input")
	}

	rule.Results = []Result{Result{}}
	return rule.Results

}

func (rule *Rule) checkAgainst(given map[string]interface{}) (interface{}, error) {
	log.Printf("\nChecking \"%v\"", rule.Key)
	log.Printf("\n--------")

	log.Printf("\nRequired...")
	val, ok := given[rule.Key]
	if !ok {
		if rule.Required {
			// Throw error indicating value is not in given input
			log.Printf("FAIL")
			return false, errors.New("Required key not found")
		} else {
			log.Printf("SKIP")
			return true, nil
		}
	} else {
		log.Printf("OK")
	}

	log.Printf("\nType... %v", reflect.TypeOf(val))

	log.Printf("\nCustom...")
	if len(rule.Customs) > 0 {
		for _, cb := range rule.Customs {
			if ok := cb(val); !ok {
				log.Printf("FAIL")
				return false, errors.New("Custom failed")
			} else {
				log.Printf("OK")
			}
		}
	} else {
		log.Printf("SKIP")
	}

	log.Printf("\nAlters...")
	if len(rule.Alters) > 0 {
		for _, cb := range rule.Alters {
			val = cb(val)
		}
	} else {
		log.Printf("SKIP")
	}

	log.Printf("\nRegex...")
	if len(rule.Regex) > 0 {
		re, err := regexp.Compile(rule.Regex)
		if err != nil {
			log.Printf("Invalid regex! %v", err)
			return false, errors.New("Invalid regex")
		}
		if re.MatchString(val.(string)) {
			log.Printf("OK")
		} else {
			log.Printf("FAIL")
			return false, nil
		}
	} else {
		log.Printf("SKIP")
	}

	log.Printf("\nType...")
	if rule.Type > 0 {
		ok := true
		if rule.Type == String {
			_, ok = val.(string)
		}
		if ok {
			log.Printf("OK")
		} else {
			log.Printf("FAIL")
		}

	} else {
		log.Printf("SKIP")
	}

	log.Printf("\n--------\n")

	return val, nil
}
