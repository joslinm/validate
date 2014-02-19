package validate

import (
	"github.com/lann/builder"
	"reflect"
	"time"
)

type ruleBuilder builder.Builder

func (rb ruleBuilder) Build() Rule {
	return builder.GetStruct(rb).(Rule)
}

func (rb ruleBuilder) updateTypeAccordingTo(val interface{}) ruleBuilder {
	Log.Debug("updateTypeAccordingTo(...) <- %v ", val)
	// get type
	t := reflect.TypeOf(val)

	switch t.Kind() {
	case reflect.Bool:
		Log.Debug("Type to boolean")
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
		Log.Debug("Type to number")
		rb = rb.Number()
		break
	case reflect.String:
		Log.Debug("Type to string")
		rb = rb.String()
		break
	case reflect.UnsafePointer:
		fallthrough
	case reflect.Ptr:
		if _, ok := val.(*time.Time); ok {
			Log.Debug("Type to time")
			rb = rb.Time()
		}
		break
	case reflect.Struct:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Slice:
	default:
		panic("Do not understand this type: " + reflect.TypeOf(val).Kind().String())
	}

	return rb
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
	rb = builder.Set(rb, "Regex", regex).(ruleBuilder)
	rb = rb.updateTypeAccordingTo(regex)
	return rb
}

// min, max, equals, between
func (rb ruleBuilder) Min(min float64) ruleBuilder {
	rb = builder.Set(rb, "Min", min).(ruleBuilder)
	rb = builder.Set(rb, "DidSetMin", true).(ruleBuilder)
	rb = rb.updateTypeAccordingTo(min)
	return rb
}

func (rb ruleBuilder) Max(max float64) ruleBuilder {
	rb = builder.Set(rb, "Max", max).(ruleBuilder)
	rb = builder.Set(rb, "DidSetMax", true).(ruleBuilder)
	rb = rb.updateTypeAccordingTo(max)
	return rb
}

func (rb ruleBuilder) In(val []string) ruleBuilder {
	rb = builder.Set(rb, "In", val).(ruleBuilder)
	rb = rb.updateTypeAccordingTo("a string")
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

// time
func (rb ruleBuilder) After(t time.Time) ruleBuilder {
	ptr := &t
	rb = builder.Set(rb, "After", ptr).(ruleBuilder)
	rb = rb.updateTypeAccordingTo(ptr)
	return rb
}
func (rb ruleBuilder) Before(t time.Time) ruleBuilder {
	ptr := &t
	rb = builder.Set(rb, "Before", ptr).(ruleBuilder)
	rb = rb.updateTypeAccordingTo(ptr)
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
var RB = RuleBuilder
