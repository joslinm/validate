package validate_test

import (
	. "github.com/franela/goblin"
	. "github.com/joslinm/validate"
	_ "reflect"
	"testing"
	"time"
)

func TestRules(t *testing.T) {
	//SetLoggingLevel(0) // critical

	g := Goblin(t)

	// Sanity Tests
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

	// Unit Tests
	g.Describe("[Rule] Unit Tests", func() {
		g.Describe("TypeOkFor(input) (val, bool)", func() {
			g.It("Should convert a boolean string", func() {
				rule := RB.Bool().Build()
				for _, val := range []string{"1", "true", "TRUE"} {
					output, ok := rule.TypeOkFor(val)
					g.Assert(output).Equal(true)
					g.Assert(ok).IsTrue()
				}
			})
			g.It("Should convert a time string", func() {
				rule := RB.Time().Build()
				time_, _ := time.Parse("2006-Jan-02", "2006-Jan-02")
				output, ok := rule.TypeOkFor("2006-Jan-02")
				g.Assert(output).Equal(time_)
				g.Assert(ok).IsTrue()
			})
			g.It("Should convert a numeric string", func() {
				rule := RB.Number().Build()
				for _, val := range []string{"5", "5.55", "5.39495", "5.0", "-1", "-1.03"} {
					output, ok := rule.TypeOkFor(val)
					_, isFloat := output.(float64)
					g.Assert(isFloat).IsTrue()
					g.Assert(ok).IsTrue()
				}
			})
			g.It("Should reject a bad boolean input", func() {
				rule := RB.Bool().Build()
				_, ok := rule.TypeOkFor("l")
				g.Assert(ok).IsFalse()
			})
			g.It("Should reject a bad time input", func() {
				rule := RB.Time().Build()
				input, ok := rule.TypeOkFor("lkjasdf")
				Log.Debug("GOT INPUT FOR BAD TIME INPUT --> %v", input)
				g.Assert(ok).IsFalse()
			})
			g.It("Should reject a bad numeric input", func() {
				rule := RB.Number().Build()
				_, ok := rule.TypeOkFor("lkjasdf")
				g.Assert(ok).IsFalse()
			})
		})
		g.Describe("Process(input) (val, []error)", func() {
			g.Describe("Types", func() {
				/* Number */
				g.Describe("Number", func() {
					g.It("Should accept string input and perform validation", func() {
						rule := RB.Min(4).Build()
						input, _ := rule.Process("5")
						g.Assert(input != nil).IsTrue()
					})

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

		})

		// Integration Tests
		/* Types
		/*************
		/ - Number
		/ - String
		/ - Bool
		/ - Time
		/**************/

	})
}
