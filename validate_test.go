package validate_test

import (
	. "github.com/franela/goblin"
	"github.com/joslinm/validate"
	"testing"
	"time"
)

func Test(t *testing.T) {
	g := Goblin(t)
	g.Describe("Rules", func() {
		g.It("Should be capable of being constructed by builder", func() {
			msg := "This is a test"
			rule := validate.RuleBuilder.Required().Regex(".*").Min(5).Max(10).Message(msg).Build()
			g.Assert(rule.Required).Equal(true)
			g.Assert(rule.Min).Equal(5)
			g.Assert(rule.Max).Equal(10)
			g.Assert(rule.Message).Equal(msg)
			g.Assert(rule.Regex).Equal(".*")
		})

		g.It(`Should update Type property to number when Min(...) / Max(...) / Between(...)
          called with number`, func() {
			rule := validate.RuleBuilder.Min(5).Build()
			g.Assert(rule.Type).Equal(validate.Number)
			rule = validate.RuleBuilder.Max(5).Build()
			g.Assert(rule.Type).Equal(validate.Number)
			rule = validate.RuleBuilder.Between(5, 10).Build()
			g.Assert(rule.Type).Equal(validate.Number)
		})

		g.It(`Should update Type property to Time when Min(...) / Max(...) / Between(...)
          called with a time`, func() {
			rule := validate.RuleBuilder.Min(time.Now()).Build()
			g.Assert(rule.Type).Equal(validate.Time)
			rule = validate.RuleBuilder.Max(time.Now()).Build()
			g.Assert(rule.Type).Equal(validate.Time)
			rule = validate.RuleBuilder.Between(time.Now(), time.Now()).Build()
			g.Assert(rule.Type).Equal(validate.Time)
		})

		g.It("Should not allow disparate values passed into Between(...)", func(done Done) {
			defer func() {
				// catch panic..
				g.Assert(recover() != nil).IsTrue()
				done()
			}()
			validate.RuleBuilder.Between(5, time.Now()).Build()
		})

		g.It("Should not allow disparate values passed into Min & Max", func(done Done) {
			defer func() {
				// catch panic..
				g.Assert(recover() != nil).IsTrue()
				done()
			}()
			validate.RuleBuilder.Min(5).Max(time.Now()).Build()
		})

		g.It("Should be able to process a non-map input if key isn't set", func() {
		})

		g.It("Should reject processing an input that does not correspond with its type", func() {
		})

		g.It("Should process its output given an input", func() {
			rule := validate.RB.Min(5).Max(10).Key("Hi").Build()
			input := map[string]interface{}{"Hi": 6}
			rule.Input = input
			rule.Process()

			g.Assert(rule.Results != nil).IsTrue()
		})

	})

	g.Describe("Validation", func() {
		g.It("Should call 'custom' callbacks", func() {
			called_custom_cb := false
			rule := validate.RuleBuilder.Custom(func(value interface{}) bool {
				called_custom_cb = true
				return true
			})
			validate.Data(map[string]interface{}{
				"key": "val",
			}).With(validate.RuleBook{
				"key": rule,
			})

			g.Assert(called_custom_cb).Equal(true)
		})
	})
}
