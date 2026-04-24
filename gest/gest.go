// gest/gest.go
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
	"sync"
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

// ── snippet cache ─────────────────────────────────────────────────────────────

// snippetCache avoids reading the same source file from disk multiple times
// across different failures in the same test run.
var snippetCache struct {
	sync.Mutex
	files map[string][]string
}

func init() {
	snippetCache.files = make(map[string][]string)
}

func cachedLines(file string) []string {
	snippetCache.Lock()
	defer snippetCache.Unlock()

	if lines, ok := snippetCache.files[file]; ok {
		return lines
	}

	f, err := os.Open(file)
	if err != nil {
		snippetCache.files[file] = nil
		return nil
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	snippetCache.files[file] = lines
	return lines
}

func readSnippet(file string, targetLine, context int) []string {
	lines := cachedLines(file)
	if lines == nil {
		return nil
	}

	start := targetLine - context - 1
	end := targetLine + context - 1

	if start < 0 {
		start = 0
	}
	if end >= len(lines) {
		end = len(lines) - 1
	}

	out := make([]string, 0, end-start+1)
	for i := start; i <= end; i++ {
		lineNum := i + 1
		marker := "  "
		color := dim
		if lineNum == targetLine {
			marker = "> "
			color = reset
		}
		out = append(out, fmt.Sprintf(
			"    %s%s%3d │ %s%s\n",
			color, marker, lineNum, lines[i], reset,
		))
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

// ── failure ───────────────────────────────────────────────────────────────────

type failure struct {
	matcher  string
	expected interface{}
	received interface{}
	file     string
	line     int
}

// ── T ─────────────────────────────────────────────────────────────────────────

// T is the test context passed to each It() callback.
type T struct {
	gt       *testing.T
	failed   bool
	failures []failure
}

func (t *T) record(f failure) {
	t.failed = true
	t.failures = append(t.failures, f)

	var b strings.Builder
	b.Grow(256)

	b.WriteString("\n    ")
	b.WriteString(dim)
	b.WriteString("expect(received).")
	b.WriteString(f.matcher)
	b.WriteString("(expected)")
	b.WriteString(reset)
	b.WriteString("\n\n    ")
	b.WriteString(green + bold)
	b.WriteString("Expected:")
	b.WriteString(reset)
	b.WriteString(" ")
	b.WriteString(green)
	fmt.Fprintf(&b, "%v", f.expected)
	b.WriteString(reset)
	b.WriteString("\n    ")
	b.WriteString(red + bold)
	b.WriteString("Received:")
	b.WriteString(reset)
	b.WriteString(" ")
	b.WriteString(red)
	fmt.Fprintf(&b, "%v", f.received)
	b.WriteString(reset)

	if snippet := readSnippet(f.file, f.line, 2); len(snippet) > 0 {
		b.WriteString("\n\n")
		for _, sl := range snippet {
			b.WriteString(sl)
		}
	}

	b.WriteString("\n    ")
	b.WriteString(dim)
	fmt.Fprintf(&b, "at %s:%d", shortPath(f.file), f.line)
	b.WriteString(reset)
	b.WriteString("\n")

	t.gt.Helper()
	t.gt.Errorf("%s", b.String())
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

// caller returns the file and line of the user's assertion.
// skip=3 resolves to: caller → matcher method → user code.
func (e *Expectation) caller() (string, int) {
	_, file, line, _ := runtime.Caller(3)
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

// ToContain asserts a string or slice contains the given element.
func (e *Expectation) ToContain(expected interface{}) {
	var ok bool
	switch v := e.actual.(type) {
	case string:
		if sub, isString := expected.(string); isString {
			ok = strings.Contains(v, sub)
		}
	case reflect.Value: // Handle reflect.Value if it comes from a slice/array
		for i := 0; i < v.Len(); i++ {
			if reflect.DeepEqual(v.Index(i).Interface(), expected) {
				ok = true
				break
			}
		}
	default: // Try to handle slices/arrays directly
		val := reflect.ValueOf(e.actual)
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			for i := 0; i < val.Len(); i++ {
				if reflect.DeepEqual(val.Index(i).Interface(), expected) {
					ok = true
					break
				}
			}
		}
	}

	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(failure{mname(e.inv, "toContain"), expected, e.actual, file, line})
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

// ToBeGreaterThanOrEqual asserts a number is greater than or equal to n.
func (e *Expectation) ToBeGreaterThanOrEqual(n interface{}) {
	ok := toFloat(e.actual) >= toFloat(n)
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toBeGreaterThanOrEqual"), n, e.actual, file, line},
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

// ToBeLessThanOrEqual asserts a number is less than or equal to n.
func (e *Expectation) ToBeLessThanOrEqual(n interface{}) {
	ok := toFloat(e.actual) <= toFloat(n)
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toBeLessThanOrEqual"), n, e.actual, file, line},
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

// ToBeInstanceOf asserts that the actual value is of the same type as expected.
func (e *Expectation) ToBeInstanceOf(expected interface{}) {
	ok := reflect.TypeOf(e.actual) == reflect.TypeOf(expected)
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(
			failure{
				mname(e.inv, "toBeInstanceOf"),
				reflect.TypeOf(expected),
				reflect.TypeOf(e.actual),
				file,
				line,
			},
		)
	}
}

// ToBeUndefined asserts that the value is nil or its zero value for its type.
func (e *Expectation) ToBeUndefined() {
	var ok bool
	if e.actual == nil {
		ok = true
	} else {
		v := reflect.ValueOf(e.actual)
		ok = !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toBeUndefined"), "undefined/zero value", e.actual, file, line},
		)
	}
}

// ToBeNaN asserts that the value is Not-a-Number (for float types).
func (e *Expectation) ToBeNaN() {
	var ok bool
	if f, isFloat := e.actual.(float64); isFloat {
		ok = math.IsNaN(f)
	} else if f, isFloat32 := e.actual.(float32); isFloat32 {
		ok = math.IsNaN(float64(f))
	}
	if !e.pass(ok) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toBeNaN"), "NaN", e.actual, file, line},
		)
	}
}

// ToMatchObject asserts that the actual map contains all key-value pairs from the expected map.
func (e *Expectation) ToMatchObject(expected map[string]interface{}) {
	actualMap, ok := e.actual.(map[string]interface{})
	if !ok {
		file, line := e.caller()
		e.t.record(
			failure{
				mname(e.inv, "toMatchObject"),
				"map[string]interface{}",
				fmt.Sprintf("%T", e.actual),
				file,
				line,
			},
		)
		return
	}

	match := true
	for k, v := range expected {
		if actualVal, exists := actualMap[k]; !exists || !reflect.DeepEqual(actualVal, v) {
			match = false
			break
		}
	}

	if !e.pass(match) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toMatchObject"), expected, e.actual, file, line},
		)
	}
}

// ToThrow asserts that the given function panics.
func (e *Expectation) ToThrow(fn func()) {
	var didPanic bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		fn()
	}()

	if !e.pass(didPanic) {
		file, line := e.caller()
		e.t.record(
			failure{mname(e.inv, "toThrow"), "a panic", "no panic", file, line},
		)
	}
}

// ── Suite ─────────────────────────────────────────────────────────────────────

type testCase struct {
	name string
	fn   func(*T)
}

// Suite represents a group of related tests, optionally with nested Describes.
type Suite struct {
	name       string
	tests      []testCase
	children   []*Suite // New field for nested describes
	beforeAll  []func(*T)
	afterAll   []func(*T)
	beforeEach []func(*T)
	afterEach  []func(*T)
}

// Describe creates a new test suite with the given name.
func Describe(name string) *Suite {
	return &Suite{name: name}
}

// It adds a test case to the suite. Returns *Suite for chaining.
func (s *Suite) It(name string, fn func(*T)) *Suite {
	s.tests = append(s.tests, testCase{name: name, fn: fn})
	return s
}

// Describe adds a nested suite (nested describe block).
// The fn receives a *Suite where you can call It(), BeforeEach(), etc.
func (s *Suite) Describe(name string, fn func(*Suite)) *Suite {
	child := &Suite{name: name}
	fn(child)
	s.children = append(s.children, child)
	return s
}

// BeforeAll registers a function to be run once before all tests in the current suite.
func (s *Suite) BeforeAll(fn func(*T)) *Suite {
	s.beforeAll = append(s.beforeAll, fn)
	return s
}

// AfterAll registers a function to be run once after all tests in the current suite.
func (s *Suite) AfterAll(fn func(*T)) *Suite {
	s.afterAll = append(s.afterAll, fn)
	return s
}

// BeforeEach registers a function to be run before each test in the current suite.
func (s *Suite) BeforeEach(fn func(*T)) *Suite {
	s.beforeEach = append(s.beforeEach, fn)
	return s
}

// AfterEach registers a function to be run after each test in the current suite.
func (s *Suite) AfterEach(fn func(*T)) *Suite {
	s.afterEach = append(s.afterEach, fn)
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
	fmt.Printf("gest:describe:%s\n", s.name)

	t := &T{gt: gt}

	// Run BeforeAll hooks
	for _, h := range s.beforeAll {
		h(t)
		if t.failed {
			gt.Errorf("BeforeAll failed in suite %q", s.name)
			return
		}
	}

	// Run AfterAll hooks after all tests and sub-suites
	defer func() {
		for _, h := range s.afterAll {
			h(t)
			if t.failed {
				gt.Errorf("AfterAll failed in suite %q", s.name)
			}
		}
	}()

	// Run individual test cases
	for _, tc := range s.tests {
		tc := tc
		gt.Run(tc.name, func(gt *testing.T) {
			gt.Helper()
			t := &T{gt: gt} // New T for each test case

			// Run BeforeEach hooks
			for _, h := range s.beforeEach {
				h(t)
				if t.failed {
					gt.Errorf("BeforeEach failed for test %q", tc.name)
					return
				}
			}

			// Run AfterEach hooks
			defer func() {
				for _, h := range s.afterEach {
					h(t)
					if t.failed {
						gt.Errorf("AfterEach failed for test %q", tc.name)
					}
				}
			}()

			// Execute the test function with panic recovery
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

	// Run nested suites
	for _, child := range s.children {
		child := child
		gt.Run(child.name, func(gt *testing.T) {
			gt.Helper()
			child.Run(gt)
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
	case int8:
		return float64(n)
	case int16:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case uint:
		return float64(n)
	case uint8:
		return float64(n)
	case uint16:
		return float64(n)
	case uint32:
		return float64(n)
	case uint64:
		return float64(n)
	case float32:
		return float64(n)
	case float64:
		return n
	}
	return 0
}
