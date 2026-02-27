package main

import (
	"math"
	"strings"

	"github.com/caiolandgraf/gest/gest"
)

func init() {
	s := gest.Describe("strings")

	s.It("'Hello, World!' should contain 'World'", func(t *gest.T) {
		t.Expect("Hello, World!").ToContain("World")
	})

	s.It("'gest' should have length 4", func(t *gest.T) {
		t.Expect("gest").ToHaveLength(4)
	})

	s.It("ToUpper should work correctly", func(t *gest.T) {
		t.Expect(strings.ToUpper("gest")).ToBe("GEST")
	})

	s.It("3.14159 should be close to π", func(t *gest.T) {
		t.Expect(3.14159).ToBeCloseTo(math.Pi, 0.01)
	})

	gest.Register(s)
}
