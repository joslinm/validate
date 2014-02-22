package validate_test

import (
  "fmt"
  . "github.com/franela/goblin"
  "github.com/joslinm/validate"
  "testing"
)

func expectError(g *G, errors []error) {
  g.Assert(len(errors) > 0).IsTrue()
}

func TestValidate(t *testing.T) {
  g := Goblin(t)
  g.Describe("Validate", func() {
    g.Describe("Data", func() {
      g.Describe("Map", func() {
        params, paramErrors := validate.Validate(map[string]interface{}{
          "x": 4,
        }).With(validate.RuleBook{
          "x": validate.RB.Min(5),
        })
        fmt.Println("Params\n", params)
        fmt.Println("Param Errors\n", paramErrors)
      })
      g.Describe("http.Request", func() {
      })
    })
  })
}
