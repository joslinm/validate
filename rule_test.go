package validate_test

import (
	"fmt"
	. "github.com/franela/goblin"
	"github.com/joslinm/validate"
	"testing"
)

func TestRules(t *testing.T) {

	fmt.Println("-- Testing Rules --")

	g := Goblin(t)
	g.Describe("Construction", func() {
		g.It("should be capable of being constructed by builder", func() {
			msg := "This is a test"
			rule := validate.RuleBuilder.Required().Regex(".*").Min(5).Max(10).Message(msg).Build()
			g.Assert(rule.Required).Equal(true)
			g.Assert(rule.Min).Equal(5)
			g.Assert(rule.Max).Equal(10)
			g.Assert(rule.Message).Equal(msg)
			g.Assert(rule.Regex).Equal(".*")
		})

	})

	/* Types
	/*************
	/ - Number
	/ - String
	/ - Bool
	/ - Time
	/**************/
	g.Describe("Type", func() {
		/* Number */
		g.Describe("Number", func() {
			g.It("should error if input < min", func() {
				rule := validate.RB.Min(5).Build()
				ok, errors := rule.Process(4)
				g.Assert(ok).IsFalse()
				g.Assert(len(errors) > 0).IsTrue()
			})
			g.It("should error if input > max", func() {
				rule := validate.RB.Max(5).Build()
				ok, errors := rule.Process(6)
				g.Assert(ok).IsFalse()
				g.Assert(len(errors) > 0).IsTrue()
			})
			g.It("Should accept string input and perform validation", func() {
				rule := validate.RB.Min(4).Build()
				ok, _ := rule.Process("5")
				g.Assert(ok).IsTrue()
			})
			g.It("Should error if string input cannot be converted", func() {
				rule := validate.RB.Min(4).Build()
				ok, _ := rule.Process("z")
				g.Assert(ok).IsFalse()
			})
		})

		/* String */
		g.Describe("String", func() {
			// :]
			g.It("should succeed with valid regex", func() {
				rule := validate.RB.Regex("hi").Build()
				ok, _ := rule.Process("hi")
				g.Assert(ok).IsTrue()
			})
			g.It("should succeed if value within enum", func() {
				rule := validate.RB.In([]string{"hi", "ho", "he"}).Build()
				ok, _ := rule.Process("hi")
				g.Assert(ok).IsTrue()
			})
			// :[
			g.It("should error with invalid regex", func() {
				rule := validate.RB.Regex("r.*").Build()
				ok, _ := rule.Process("hi")
				g.Assert(ok).IsFalse()
			})
			g.It("should error if value not in enum", func() {
				rule := validate.RB.In([]string{"ho", "he"}).Build()
				ok, errors := rule.Process("hi")
				g.Assert(ok).IsFalse()
				g.Assert(len(errors) > 0).IsTrue()
			})
		})
	})
}
