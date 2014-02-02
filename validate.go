package validate

import (
	"errors"
	"fmt"
	"github.com/lann/builder"
	"log"
	_ "net/http"
	"reflect"
	"regexp"
	"time"
)

// Types validate supports
const (
	Unknown = iota
	String
	Number
	Bool
	Time
)

// Type of callbacks to be used in a Rule
type AlterCallback func(value interface{}) interface{}
type PrepareCallback func(value interface{}) interface{}
type CustomCallback func(value interface{}) bool

// Rule encompasses a single validation rule for a parameter
type Rule struct {
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
func (rb ruleBuilder) Type(is string) ruleBuilder {
	return builder.Set(rb, "Type", is).(ruleBuilder)
}
func (rb ruleBuilder) String() ruleBuilder {
	return builder.Set(rb, "Type", String).(ruleBuilder)
}
func (rb ruleBuilder) Number() ruleBuilder {
	return builder.Set(rb, "Type", Number).(ruleBuilder)
}
func (rb ruleBuilder) Bool() ruleBuilder {
	return builder.Set(rb, "Type", Bool).(ruleBuilder)
}
func (rb ruleBuilder) Time() ruleBuilder {
	return builder.Set(rb, "Type", Time).(ruleBuilder)
}

// message
func (rb ruleBuilder) Message(msg string) ruleBuilder {
	return builder.Set(rb, "Message", msg).(ruleBuilder)
}

// regex
func (rb ruleBuilder) Regex(regex string) ruleBuilder {
	return builder.Set(rb, "Regex", regex).(ruleBuilder)
}

func (rb ruleBuilder) updateTypeAccordingTo(val interface{}) ruleBuilder {
	// get type
	t := reflect.TypeOf(val)
	log.Println("Dynamically updating type from", t.Kind().String())

	switch t.Kind() {

	// primitives
	case reflect.Bool:
		log.Println("Updating type to number")
		rb = rb.Bool()
		break
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		log.Println("Updating type to number")
		rb = rb.Number()
		break
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Slice:
	case reflect.String:
	case reflect.Struct:
		if _, ok := val.(time.Time); ok {
			log.Println("Updating type to time")
			rb = rb.Time()
		}
		break

	case reflect.UnsafePointer:
	default:
		panic("Do not understand this type: " + reflect.TypeOf(val).Kind().String())
	}

	return rb
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

func (rb ruleBuilder) ensureMinAndMaxEqualTypes() {
	rule := rb.Build()
	if rule.Max != nil && rule.Min != nil { // if max/min set
		if !sameType(rule.Max, rule.Min) {
			panic("You set min/max rules to different types")
		}
	}
}

// min, max, between
func (rb ruleBuilder) Min(min interface{}) ruleBuilder {
	rb = builder.Set(rb, "Min", min).(ruleBuilder)
	rb.ensureMinAndMaxEqualTypes()
	rb = rb.updateTypeAccordingTo(min)
	return rb
}

func (rb ruleBuilder) Max(max interface{}) ruleBuilder {
	rb = builder.Set(rb, "Max", max).(ruleBuilder)
	rb.ensureMinAndMaxEqualTypes()
	rb = rb.updateTypeAccordingTo(max)
	return rb
}

func (rb ruleBuilder) Between(min interface{}, max interface{}) ruleBuilder {
	if reflect.TypeOf(min).Kind() != reflect.TypeOf(max).Kind() {
		panic("Disparate values passed into Between(...) \n" +
			"\nMin: " + reflect.TypeOf(min).Kind().String() +
			"\nMax: " + reflect.TypeOf(max).Kind().String())
	}
	rb = rb.updateTypeAccordingTo(min)
	builder.Set(rb, "Min", min)
	builder.Set(rb, "Max", max)

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
	if rule.Type > 0 {
		ok := true
		if rule.Type == String {
			_, ok = val.(string)
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
