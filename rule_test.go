package validate_test

import (
	. "github.com/franela/goblin"
	. "github.com/joslinm/validate"
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
			// Positive Tests
			g.It("Should accept string input and perform validation", func() {
				rule := RB.Min(4).Build()
				ok, _ := rule.Process("5")
				g.Assert(ok).IsTrue()
			})
			g.It("Should error if string input cannot be converted", func() {
				rule := RB.Min(4).Build()
				ok, _ := rule.Process("z")
				g.Assert(ok).IsFalse()
			})

			// Negative Tests
			g.It("should error if input < min", func() {
				rule := RB.Min(5).Build()
				ok, errors := rule.Process(4)
				g.Assert(ok).IsFalse()
				g.Assert(len(errors) > 0).IsTrue()
			})
			g.It("should error if input > max", func() {
				rule := RB.Max(5).Build()
				ok, errors := rule.Process(6)
				g.Assert(ok).IsFalse()
				g.Assert(len(errors) > 0).IsTrue()
			})
		})

		/* String */
		g.Describe("String", func() {
			// :]
			g.It("should succeed with valid regex", func() {
				rule := RB.Regex("hi").Build()
				ok, _ := rule.Process("hi")
				g.Assert(ok).IsTrue()
			})
			g.It("should succeed if value within enum", func() {
				rule := RB.In([]string{"hi", "ho", "he"}).Build()
				ok, _ := rule.Process("hi")
				g.Assert(ok).IsTrue()
			})
			// :[
			g.It("should error with invalid regex", func() {
				rule := RB.Regex("r.*").Build()
				ok, _ := rule.Process("hi")
				g.Assert(ok).IsFalse()
			})
			g.It("should error if value not in enum", func() {
				rule := RB.In([]string{"ho", "he"}).Build()
				ok, errors := rule.Process("hi")
				g.Assert(ok).IsFalse()
				g.Assert(len(errors) > 0).IsTrue()
			})
		})

		/* Boolean */
		g.Describe("Boolean", func() {
			g.It("should succeed regardless of value given", func() {
				rule := RB.Bool().Build()
				ok, errors := rule.Process(true)
				g.Assert(ok).IsTrue()
				g.Assert(len(errors) == 0).IsTrue()
				ok, errors = rule.Process(false)
				g.Assert(ok).IsTrue()
				g.Assert(len(errors) == 0).IsTrue()
			})

			g.It("Should not error if value is string capable of converting to boolean", func() {
				rule := RB.Bool().Build()
				ok, errors := rule.Process("1")
				g.Assert(ok).IsTrue()
				ok, errors = rule.Process("true")
				g.Assert(ok).IsTrue()
				ok, errors = rule.Process("TRUE")
				g.Assert(ok).IsTrue()
				g.Assert(len(errors) == 0).IsTrue()
			})

			g.It("Should error if value is not a boolean", func() {
				rule := RB.Bool().Build()
				ok, errors := rule.Process("hi")
				expectError(g, ok, errors)
			})
		})

		/* Time */
		g.Describe("Time", func() {
			g.It("should succeed if given time is after set minimum", func() {
				tim, _ := time.Parse("2006-Jan-02", "2006-Jan-02")
				rule := RB.After(tim).Build()
				ok, _ := rule.Process(time.Now())
				g.Assert(ok).IsTrue()
			})
			g.It("should succeed if given time is before set maximum", func() {
				tim, _ := time.Parse("2006-Jan-02", "2020-Jan-02")
				rule := RB.Before(tim).Build()
				ok, _ := rule.Process(time.Now())
				g.Assert(ok).IsTrue()
			})

			g.It("should error if given time is before set minimum", func() {
				tim, _ := time.Parse("2006-Jan-02", "2020-Jan-02")
				rule := RB.After(tim).Build()
				ok, errors := rule.Process(time.Now())
				expectError(g, ok, errors)
			})
			g.It("should error if given time is after set maximum", func() {
				tim, _ := time.Parse("2006-Jan-02", "2006-Jan-02")
				rule := RB.Before(tim).Build()
				ok, errors := rule.Process(time.Now())
				expectError(g, ok, errors)
			})
		})

	})
}
