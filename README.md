*Under Construction*

Rules
------
There are rules for your input. A rule is a pretty simple data structure:
```go

type AlterCallback func(value interface{}) interface{}
type PrepareCallback func(value interface{}) interface{}
type CustomCallback func(value interface{}) bool

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
```

If a property isn't set, it's not evaluated. 

When the data is validated and the `Is` member is set, the value will be coerced into the expecting type.
For example, a string value will be coerced to an integer and returned in the dictionary from `dict, err := validate.Data(input).With(rules)`

Rules are best built using the [builder](https://github.com/lann/builder).

```go
  // simple
  stringRule := validate.RuleBuilder.String()

  // complex
  stringRule = validate.RuleBuilder.Required().String().Regex("^\\w+$").Custom(func (val interface{}) bool {
    if val.(string) == "bad" { // this custom rule ensures no "bad" values get in
      return false
    }
    return true
  })
```

But you're probably going to be defining a lot of rules, so it's better you reuse builder(s).

```go
  required := validate.RuleBuilder.Required()
  optional := validate.RuleBuilder

  // then define rules on top of them..
  // just demonstrations; usually you'll make these calls in the RuleBook
  date := required.Date()
  minRule := required.Min(1) // type will automatically be set to Int
  rangeRule :=  required.Between(5.5, 7.5) // type will be set to Float 
  emailRule = optional.Email() // helper builder functions like this pre-set values. in this case regex becomes an email regex
```

You don't really need to create a bunch of rule variables though. You can just do something like this:

```go
  required := validate.RuleBuilder.Required()
  optional := validate.RuleBuilder

  params, err := validate.Data(data).With(validate.RuleBook {
    "first_name": required.String(),
    "last_name": required.String().Regex("\\w+"),
    "email": optional.Email(),
    "age": optional.Between(0, 150), // range of ints
    "born": optional.Between(time.Parse("1900-Jan-01"), time.Now()), // range of dates
    "iq": optional.Max(100) // keep out the smart people
  })
  
  // use your params...
```

RuleBook template from struct
-----
Perhaps you already have a struct and want a RuleBook right quick. Just pass an empty struct and you'll get a RuleBook with rules for all recognized types. After you get your RuleBook  back, you can modify any rules just as you would above.

```go
  type UserParams struct {
    FirstName string
    LastName string
    Email string
    Age int
  }
  
  rules = validate.RuleBookFor(&UserParams{}, true) // mark them all required
  rules["email"].(validate.RuleBuilder).Email()
  params, err := validate.Data(input).With(rules)
```

Recognized Types
------
* `int` | `float` -> `validate.NUMBER`
* `int`           -> `validate.INT`
* `float`         -> `validate.FLOAT`
* `string`        -> `validate.STRING`
* `bool`          -> `validate.BOOL`
* `time.Time`     -> `validate.TIME`

Nesting
------
You might want the ability to nest data structures. This is easily accomplished.

```go
params, err := validate.Data(map[string]interface{}{
    "date": map[string]time.Time {
      "start": time.Time.Now(),
      "end": time.Time.Now(),
    }
  }).With(validate.RuleBook{
      "date": validate.RuleBook{
        "start": RuleBuilder.Required().Date()
        "end": RuleBuilder.Required().Date()
      }
  })
```

Pre/Post Processing
--------
Before or after validation rules (which includes custom callbacks), you might want to transform the data. 

* pre processing  -> `Prepare()`
* post processing -> `Alter()`

```go
lc_then_uc := validate.RuleBuilder.Required().Prepare(func(val interface{}) interface{} {
  // Lowercase code
  return val
}).Regex("[a-z]+") // Validate that it's lower case
.Alter(validate.Alter { // Now upper case it
  // Upper case code
  return val
})
```

