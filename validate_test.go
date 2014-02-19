package validate_test

import (
	. "github.com/franela/goblin"
	"github.com/joslinm/validate"
	"testing"
)

func expectError(g *G, ok bool, errors []error) {
	g.Assert(ok).IsFalse()
	g.Assert(len(errors) > 0).IsTrue()
}

func TestValidate(t *testing.T) {
	g := Goblin(t)
	g.Describe("Validation [Negative Tests]", func() {

		g.It("Should error if number is below min", func() {
			_, err := validate.Data(map[string]interface{}{
				"min": 4,
			}).With(validate.RuleBook{
				"min": validate.RB.Min(5),
			})
			g.Assert(err != nil).IsTrue()
		})

		g.It("Should error if number is above max", func() {
			_, err := validate.Data(map[string]interface{}{
				"max": 5,
			}).With(validate.RuleBook{
				"max": validate.RB.Max(4),
			})
			g.Assert(err != nil).IsTrue()
		})

	})

	g.Describe("Validation [Positive Tests]", func() {
	})
}
