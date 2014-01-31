
Usage:

So, you have to pass a dictionary containing some rules. We call this the RuleBook.
** Under Construction **

Here's an example of a rule book that ensures the input contains a title & genre key with an optional duration integer. 

Rules
------
There are rules for your input. Let's define some.


A rule is a pretty simple data structure:
```
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
`

If a property isn't set, it's not evaluated. 

When the data is validated and the `Is` member is set, the value will be coerced into the expecting type.
For example, a string value will be coerced to an integer if need be.

The dictionary returned back to you from `dict, err := validate.Data(input).With(rules)` will contain
values typed accordingly. 

Rules are best built using the [builder](https://github.com/lann/builder).

```
  // simple
  stringRule := validate.RuleBuilder.String()

  // complex
  stringRule = validate.RuleBuilder.Required().String().Regex("^\\w+$").Custom(func (val interface{}) bool {
    if val.(string) == "bad" { // this custom rule ensures no "bad" values get in
      return false
    }
    return true
  }
`

But you're probably going to be defining a lot, so it's better you reuse builders.

```
  required := validate.RuleBuilder.Required(true)
  optional := validate.RuleBuilder.Required(false) // or required.Required(false) if you want to be silly

  // then define rules on top of them..
  // just demonstrations; usually you'll make these calls in the RuleBook
  date := required.Date()
  minRule := required.Min(0) // type will automatically be set to Int
  rangeRule :=  required.Between(5.5, 7.5) // type will be set to Float 
  emailRule = optional.Email() // helper builder functions like this pre-set values. in this case regex becomes an email regex
`

You don't really need to create a bunch of rule variables though. You can just do something like this:

```go
  required := validate.RuleBuilder.Required()
  optional := validate.RuleBuilder

  // all this really is just a `map[string]interface{}` 
  rules := validate.RuleBook {
    "first_name": required.String(),
    "last_name": required.String().Regex("\\w+"),
    "email": optional.Email(),
    "age": optional.Between(0, 150), // range of ints
    "born": optional.Between(time.Parse("1900-Jan-01"), time.Now()), // range of dates
    "iq": optional.Max(100) // keep out the smart people
  }
  params, err := validate.Data(data).With(rules) 
`

```
// let's setup a required rule template and just keep re-using it
required := RuleBuilder.Required(true)
// our adult filter
adult_filter := func (genre string) bool {
    if genre == "Adult" { // keep the bad stuff out
      return false
    }
    return true
}
rules := RuleBook{
  "title": required.Regex("[A-Z]\\w+"),
  "genre": required.Custom(adult_filter),
  "adult": required.Boolean()
  "release": required.After(time.Time.Now()).Before(time.Time.Now() + 100)
  "duration": Rule.Min(5).Max(10) // or .Between(5, 10)
  },
} 
params, err = validate.Data(input).With(rules)
if (err != nil) {
  fmt.Println("Error! ", err)
  // Handle err..
}
`

There's also helper functions to make your life easier:

```
  expecting := RuleBook {
    "email": RuleBuilder.Email().Required() 
  } `


RuleBook template from struct
-----
Perhaps you already have a struct and want a RuleBook right quick. Just pass an empty struct and you'll get a RuleBook with rules for 
all recognized types. For example, below, the rule would be `RuleBuilder.Required().String()`. After you get your RuleBook ( `map[string]*Rule` ) back, you
can modify any rules just as you would above.

```go
  type UserParams struct {
    FirstName string
    LastName string
    Email string
    Age int
  }
  
  rules = validate.RuleBookFor(&UserParams{}, true) // mark them all required
  rules["email"].Email()
  params, err := validate.Data(input).With(rules)
`

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
`

Pre/Post Processing
--------
Before or after validation rules (which includes custom callbacks), you might want to transform the data. You don't have
to alter the data. You can log it, but you really shouldn't. The terminology `Prepare()` & `Alter()` describe their intended usage.
Logging can come before & after the validation step. 

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
`

