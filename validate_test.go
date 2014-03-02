package validate_test

import (
	_ "fmt"
	. "github.com/franela/goblin"
	. "github.com/joslinm/validate"
	"testing"
)

func TestValidate(t *testing.T) {
	g := Goblin(t)
	g.Describe("Validate", func() {
		g.Describe("Output", func() {
			g.It("Should return a hash of params sent in if they pass their rules", func() {
				params, _ := Validate(map[string]interface{}{
					"x": "4",
					"y": "hi!",
				}).With(RuleBook{
					"x": RB.Min(1),
					"y": RB.Regex("hi.*"),
				})
				g.Assert(params["x"]).Equal(4)
				g.Assert(params["y"]).Equal("hi!")
			})

			g.It("Should return a hash of errors for all failed params", func() {
				_, errors := Validate(map[string]interface{}{
					"x": "4",
					"y": "hi!",
				}).With(RuleBook{
					"x": RB.Min(5),
					"y": RB.Regex("ddd.*"),
				})
				g.Assert(errors["x"] != nil).IsTrue()
				g.Assert(errors["y"] != nil).IsTrue()
			})
		})
	})
}
