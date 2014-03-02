package validate_test

import (
	. "github.com/franela/goblin"
	. "github.com/joslinm/validate"
	_ "reflect"
	"testing"
	"time"
)

func TestRules(t *testing.T) {
	//validate.SetLoggingLevel(0) // critical

	g := Goblin(t)
	g.Describe("Construction", func() {
		g.It("should be capable of being constructed by builder", func() {
			msg := "This is a test"
			rule := RuleBuilder.Required().Regex(".*").Min(5).Max(10).Message(msg).Build()
			g.Assert(rule.Required).Equal(true)
			g.Assert(rule.Min).Equal(5)
			g.Assert(rule.Max).Equal(10)
			g.Assert(rule.Message).Equal(msg)
			g.Assert(rule.Regex).Equal(".*")
		})

	})

	g.Describe("Type Coercion", func() {
		g.Describe("Number", func() {
			g.It("Should accept a string number input and convert it", func() {
				rule := RB.Min(4).Build()
				input, _ := rule.Process("5")
				_, isString := input.(string)
				g.Assert(isString).IsFalse()
			})
			g.It("Should error if string input cannot be converted", func() {
				rule := RB.Min(4).Build()
				input, errors := rule.Process("z")
				g.Assert(input).Equal("z")
				g.Assert(len(errors)).Equal(1)
			})
		})
		g.Describe("Time", func() {
			g.It("Should accept a string input and attempt conversion", func() {
				tim, _ := time.Parse("2006-Jan-02", "2006-Jan-02")
				rule := RB.Before(tim).Build()
				input, errors := rule.Process("2006-Jan-03")
				_, isTime := input.(time.Time)
				g.Assert(len(errors)).Equal(0)
				g.Assert(isTime).IsTrue()
			})
		})
	})
	/* Types
	/*************
	/ - Number
	/ - String
	/ - Bool
	/ - Time
	/**************/
	g.Describe("Type Logic", func() {
		/* Number */
		g.Describe("Number", func() {
			// Positive Tests
			g.It("Should accept string input and perform validation", func() {
				rule := RB.Min(4).Build()
				input, _ := rule.Process("5")
				g.Assert(input != nil).IsTrue()
			})

			// Negative Tests
			g.It("should error if input < min", func() {
				rule := RB.Min(5).Build()
				input, errors := rule.Process(4)
				Log.Debug("GOT INPUT %v", input)
				g.Assert(input != nil).IsTrue()
				g.Assert(len(errors) > 0).IsTrue()
			})
			g.It("should error if input > max", func() {
				rule := RB.Max(5).Build()
				input, errors := rule.Process(6)
				g.Assert(input != nil).IsTrue()
				g.Assert(len(errors) > 0).IsTrue()
			})
		})

		/* String */
		g.Describe("String", func() {
			// :]
			g.It("should succeed with valid regex", func() {
				rule := RB.Regex("hi").Build()
				input, _ := rule.Process("hi")
				g.Assert(input != nil).IsTrue()
			})
			g.It("should succeed if value within enum", func() {
				rule := RB.In([]string{"hi", "ho", "he"}).Build()
				input, _ := rule.Process("hi")
				g.Assert(input != nil).IsTrue()
			})
			// :[
			g.It("should error with invalid regex", func() {
				rule := RB.Regex("r.*").Build()
				_, errors := rule.Process("hi")
				g.Assert(len(errors)).Equal(1)
			})
			g.It("should error if value not in enum", func() {
				rule := RB.In([]string{"ho", "he"}).Build()
				input, errors := rule.Process("hi")
				g.Assert(input != nil).IsTrue()
				g.Assert(len(errors) > 0).IsTrue()
			})
		})

		/* Boolean */
		g.Describe("Boolean", func() {
			// :]
			g.It("should succeed regardless of value given", func() {
				rule := RB.Bool().Build()
				input, errors := rule.Process(true)
				g.Assert(input != nil).IsTrue()
				g.Assert(len(errors) == 0).IsTrue()
				input, errors = rule.Process(false)
				g.Assert(input != nil).IsTrue()
				g.Assert(len(errors) == 0).IsTrue()
			})

			g.It("Should not error if value is string capable of converting to boolean", func() {
				rule := RB.Bool().Build()

				// 1
				input, errors := rule.Process("1")
				g.Assert(input).Equal(true)
				g.Assert(len(errors) == 0).IsTrue()

				// true | TRUE
				input, errors = rule.Process("true")
				g.Assert(input).Equal(true)
				g.Assert(len(errors) == 0).IsTrue()

				input, errors = rule.Process("TRUE")
				g.Assert(input != nil).IsTrue()
				g.Assert(len(errors) == 0).IsTrue()
			})

			// :[
			g.It("Should error if value is not a boolean", func() {
				rule := RB.Bool().Build()
				_, errors := rule.Process("hi")
				g.Assert(len(errors)).Equal(1)
			})
		})

		/* Time */
		g.Describe("Time", func() {
			// :]
			g.It("should succeed if given time is after set minimum", func() {
				tim, _ := time.Parse("2006-Jan-02", "2006-Jan-02")
				rule := RB.After(tim).Build()
				input, _ := rule.Process(time.Now())
				g.Assert(input != nil).IsTrue()
			})
			g.It("should succeed if given time is before set maximum", func() {
				tim, _ := time.Parse("2006-Jan-02", "2020-Jan-02")
				rule := RB.Before(tim).Build()
				input, _ := rule.Process(time.Now())
				g.Assert(input != nil).IsTrue()
			})

			// :[
			g.It("should error if given time is before set minimum", func() {
				tim, _ := time.Parse("2006-Jan-02", "2020-Jan-02")
				rule := RB.After(tim).Build()
				_, errors := rule.Process(time.Now())
				g.Assert(len(errors)).Equal(1)
			})
			g.It("should error if given time is after set maximum", func() {
				tim, _ := time.Parse("2006-Jan-02", "2006-Jan-02")
				rule := RB.Before(tim).Build()
				_, errors := rule.Process(time.Now())
				g.Assert(len(errors)).Equal(1)
			})
		})

	})
}
