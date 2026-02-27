package main

import "github.com/caiolandgraf/gest/gest"

func init() {
	calc := Calculator{}
	s := gest.Describe("calculator")

	s.It("adding 2 + 2 should return 4", func(t *gest.T) {
		t.Expect(calc.Add(2, 2)).ToBe(float64(5))
	})

	s.It("adding 2 + 2 should return 4", func(t *gest.T) {
		t.Expect(calc.Add(2, 2)).ToBe(float64(4))
	})

	s.It("subtracting 10 - 3 should return 7", func(t *gest.T) {
		t.Expect(calc.Subtract(10, 3)).ToBe(float64(7))
	})

	s.It("multiplying 3 * 4 should return 12", func(t *gest.T) {
		t.Expect(calc.Multiply(3, 4)).ToBe(float64(12))
	})

	s.It("dividing 10 / 2 should return 5", func(t *gest.T) {
		result, err := calc.Divide(10, 2)
		t.Expect(err).ToBeNil()
		t.Expect(result).ToBe(float64(5))
	})

	s.It("dividing by zero should return an error", func(t *gest.T) {
		_, err := calc.Divide(10, 0)
		t.Expect(err).Not().ToBeNil()
	})

	gest.Register(s)
}
