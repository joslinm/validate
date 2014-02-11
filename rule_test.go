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

	g.Describe("Number validation", func() {
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
}
