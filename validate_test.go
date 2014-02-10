package validate_test

import (
	"fmt"
	. "github.com/franela/goblin"
	"github.com/joslinm/validate"
	"testing"
)

func TestValidate(t *testing.T) {
	fmt.Println("-- Testing Validation --")

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
