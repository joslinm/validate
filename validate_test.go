package validate_test

import (
	"fmt"
	. "github.com/franela/goblin"
	"github.com/joslinm/validate"
	"testing"
)

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
				g.Assert(len(paramErrors["x"])).Equal(1)
			})
			g.Describe("http.Request", func() {
			})
		})
	})
}
