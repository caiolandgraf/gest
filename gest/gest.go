package gest

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

// ── ANSI ──────────────────────────────────────────────────────────────────────

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[38;2;220;60;60m"
	green   = "\033[38;2;40;210;90m"
	yellow  = "\033[38;2;230;200;40m"
	bgGreen = "\033[48;2;40;180;80m\033[38;2;255;255;255m"
	bgRed   = "\033[48;2;195;55;55m\033[38;2;255;255;255m"
)

// ── failure ───────────────────────────────────────────────────────────────────

type failure struct {
	matcher  string
	expected interface{}
	received interface{}
	file     string
	line     int
}

// ── snippet reader ────────────────────────────────────────────────────────────

func readSnippet(file string, targetLine, context int) []string {
	f, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer func() { _ = f.Close() }()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	start := targetLine - context - 1
	end := targetLine + context - 1

	if start < 0 {
		start = 0
	}
	if end >= len(lines) {
		end = len(lines) - 1
	}

	var out []string
	for i := start; i <= end; i++ {
		lineNum := i + 1
		marker := "  "
		color := dim
		if lineNum == targetLine {
			marker = "> "
			color = reset
		}
		out = append(
			out,
			fmt.Sprintf(
				"    %s%s%3d │ %s%s\n",
				color,
				marker,
				lineNum,
				lines[i],
				reset,
			),
		)
	}
	return out
}

func shortPath(full string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return full
	}
	rel, err := filepath.Rel(cwd, full)
	if err != nil {
		return full
	}
	if strings.HasPrefix(rel, "..") {
		return full
	}
	return rel
}

// ── T ─────────────────────────────────────────────────────────────────────────

// T is the test context passed to each It() callback.
// Use t.Expect(value) to start an assertion chain.
type T struct {
	gt       *testing.T
	failed   bool
	failures []failure
}

func (t *T) record(f failure) {
	t.failed = true
	t.failures = append(t.failures, f)

	msg := fmt.Sprintf(
		"\n    %sexpect(received).%s(expected)%s\n\n    %sExpected:%s %s%v%s\n    %sReceived:%s %s%v%s",
		dim,
		f.matcher,
		reset,
		green+bold,
		reset,
		green,
		f.expected,
		reset,
		red+bold,
		reset,
		red,
		f.received,
		reset,
	)

	snippet := readSnippet(f.file, f.line, 2)
	if len(snippet) > 0 {
		msg += "\n\n"
		for _, sl := range snippet {
			msg += sl
		}
	}

	msg += fmt.Sprintf(
		"\n    %sat %s:%d%s\n",
		dim,
		shortPath(f.file),
		f.line,
		reset,
	)

	t.gt.Helper()
	t.gt.Errorf("%s", msg)
}

// Expect starts an assertion chain for the given value.
func (t *T) Expect(actual interface{}) *Expectation {
	return &Expectation{t: t, actual: actual}
}

// ── Expectation ───────────────────────────────────────────────────────────────

// Expectation holds an actual value and lets you chain matchers against it.
type Expectation struct {
	t      *T
	actual interface{}
	inv    bool
}

// Not negates the next matcher in the chain.
func (e *Expectation) Not() *Expectation { e.inv = !e.inv; return e }

func (e *Expectation) pass(ok bool) bool {
	if e.inv {
		return !ok
	}
	return ok
}

func (e *Expectation) caller() (string, int) {
	// walk up the stack past gest internals to find the user's assertion line
	for skip := 2; skip < 10; skip++ {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		if !strings.Contains(file, "/gest/gest.go") {
			return file, line
		}
	}
	_, file, line, _ := runtime.Caller(2)
	return file, line
}

func mname(inv bool, name string) string {
	if inv {
		return "not." + name
	}
	return name
}

// ── Matchers ──────────────────────────────────────────────────────────────────

// ToBe asserts strict equality (==).
func (e *Expectation) ToBe(expected interface{}) {
	ok := e.actual == expected
	if !e.pass(ok) {
		file, line := e.caller()
		exp, rec := expected, e.actual
		if e.inv {
			exp, rec = e.actual, expected
		}
		e.t.record(failure{mname(e.inv, "toBe"), exp, rec, file, line})
	}
}

// ToEqual asserts deep equality via reflect.DeepEqual.
func (e *Expectation) ToEqual(expected interface{}) {
	ok := reflect.DeepEqual(e.actual, expected)
	if !e.pass(ok) {
		file, line := e.caller()
		exp, rec := expected, e.actual
		if e.inv {
			exp, rec = e.actual, expected
		}
		e.t.record(failure{mname(e.inv, "toEqual"), exp, rec, file, line})
	}
}

// ToBeNil asserts the value is nil.
func (e *Expectation) ToBeNil() {
	var ok bool
	if e.actual == nil {
		ok = true
	} else {
		v := reflect.ValueOf(e.actual)
		switch v.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Slice,
			reflect.Map, reflect.Chan, reflect.Func:
			ok = v.IsNil()
		}
	}
	if !e.pass(ok) {
		file, line := e.caller()
		if e.inv {
			e.t.record(
				failure{
					mname(e.inv, "toBeNil"),
					"not nil",
					e.actual,
					file,
					line,
				},
			)
		} else {
			e.t.record(
				failure{mname(e.inv, "toBeNil"), nil, e.actual, file, line},
			)
		}
	}
}

// ToBeTrue asserts the value is exactly true.
func (e *Expectation) ToBeTrue() {
	ok := e.actual == true
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toBeTrue"), true, e.actual, file, line},
		)
	}
}

// ToBeFalse asserts the value is exactly false.
func (e *Expectation) ToBeFalse() {
	ok := e.actual == false
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toBeFalse"), false, e.actual, file, line},
		)
	}
}

// ToContain asserts a string contains the given substring.
func (e *Expectation) ToContain(sub string) {
	s, ok2 := e.actual.(string)
	if !ok2 {
		file, line := e.caller()
		e.t.record(
			failure{
				"toContain",
				sub,
				fmt.Sprintf("%T (not a string)", e.actual),
				file,
				line,
			},
		)
		return
	}
	ok := strings.Contains(s, sub)
	if !e.pass(ok) {
		file, line := e.caller()
		if e.inv {
			e.t.record(
				failure{
					mname(e.inv, "toContain"),
					fmt.Sprintf("not containing %q", sub),
					s,
					file,
					line,
				},
			)
		} else {
			e.t.record(failure{mname(e.inv, "toContain"), sub, s, file, line})
		}
	}
}

// ToHaveLength asserts the length of a string, slice, array, or map.
func (e *Expectation) ToHaveLength(n int) {
	v := reflect.ValueOf(e.actual)
	var length int
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		length = v.Len()
	default:
		file, line := e.caller()
		e.t.record(
			failure{
				"toHaveLength",
				n,
				fmt.Sprintf("%T (unsupported)", e.actual),
				file,
				line,
			},
		)
		return
	}
	ok := length == n
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(failure{mname(e.inv, "toHaveLength"), n, length, file, line})
	}
}

// ToBeGreaterThan asserts a number is strictly greater than n.
func (e *Expectation) ToBeGreaterThan(n interface{}) {
	ok := toFloat(e.actual) > toFloat(n)
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toBeGreaterThan"), n, e.actual, file, line},
		)
	}
}

// ToBeLessThan asserts a number is strictly less than n.
func (e *Expectation) ToBeLessThan(n interface{}) {
	ok := toFloat(e.actual) < toFloat(n)
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toBeLessThan"), n, e.actual, file, line},
		)
	}
}

// ToBeCloseTo asserts a float is approximately equal to n within delta (default ±0.001).
func (e *Expectation) ToBeCloseTo(n interface{}, delta ...float64) {
	d := 0.001
	if len(delta) > 0 {
		d = delta[0]
	}
	ok := math.Abs(toFloat(e.actual)-toFloat(n)) <= d
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toBeCloseTo"), n, e.actual, file, line},
		)
	}
}

// ── Suite ─────────────────────────────────────────────────────────────────────

type testCase struct {
	name string
	fn   func(*T)
}

// Suite represents a group of related tests.
type Suite struct {
	name  string
	tests []testCase
}

// Describe creates a new test suite with the given name.
func Describe(name string) *Suite { return &Suite{name: name} }

// It adds a test case to the suite. Returns *Suite for chaining.
func (s *Suite) It(name string, fn func(*T)) *Suite {
	s.tests = append(s.tests, testCase{name: name, fn: fn})
	return s
}

// Run executes all It() cases as native go test subtests via t.Run().
// Call this at the end of any TestXxx function.
//
//	func TestCalculator(t *testing.T) {
//	    gest.Describe("calculator").
//	        It("adds correctly", func(t *gest.T) {
//	            t.Expect(calc.Add(2, 2)).ToBe(float64(4))
//	        }).
//	        Run(t)
//	}
func (s *Suite) Run(gt *testing.T) {
	gt.Helper()
	for _, tc := range s.tests {
		tc := tc
		gt.Run(tc.name, func(gt *testing.T) {
			gt.Helper()
			t := &T{gt: gt}
			func() {
				defer func() {
					if r := recover(); r != nil {
						gt.Errorf("%spanicked: %v%s", red, r, reset)
					}
				}()
				tc.fn(t)
			}()
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func fmtDur(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.3fs", d.Seconds())
}

// fmtDur is used by the CLI — exported via the package for potential reuse.
var _ = fmtDur

func toFloat(v interface{}) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case float32:
		return float64(n)
	case float64:
		return n
	}
	return 0
}

// color helpers exported for use by the CLI
var (
	ColorReset   = reset
	ColorBold    = bold
	ColorDim     = dim
	ColorRed     = red
	ColorGreen   = green
	ColorYellow  = yellow
	ColorBgGreen = bgGreen
	ColorBgRed   = bgRed
)
