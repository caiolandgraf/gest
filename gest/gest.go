package gest

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ── ANSI ──────────────────────────────────────────────────────────────────────
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[38;2;220;60;60m"                       // vermelho
	green   = "\033[38;2;40;210;90m"                       // verde
	yellow  = "\033[38;2;230;200;40m"                      // amarelo
	bgGreen = "\033[48;2;40;180;80m\033[38;2;255;255;255m" // bg verde + branco
	bgRed   = "\033[48;2;195;55;55m\033[38;2;255;255;255m" // bg vermelho + branco
)

// ── Registry ──────────────────────────────────────────────────────────────────

var registry []*Suite

// Register adds a suite to the global registry.
// Call it inside init() in each _spec.go file.
func Register(s *Suite) {
	registry = append(registry, s)
}

// RunRegistered runs all suites registered via Register().
// This is all your main.go needs to call.
//
// If --watch (or -w) is present in os.Args, RunRegistered enters watch mode:
// it runs tests immediately and then re-runs them whenever a .go file changes.
// The flag is stripped before the args are forwarded to the subprocess.
func RunRegistered() bool {
	var watchMode bool
	var remaining []string

	for _, arg := range os.Args[1:] {
		if arg == "--watch" || arg == "-w" {
			watchMode = true
		} else {
			remaining = append(remaining, arg)
		}
	}

	if watchMode {
		WatchTests(remaining)
		return true
	}

	return RunAll(registry...)
}

// ── Watch mode ────────────────────────────────────────────────────────────────

// WatchTests compiles the project into a temporary OS directory, runs it
// immediately, then re-runs it automatically on every .go file change.
// args are forwarded to the binary on each invocation.
// The temp directory is removed when the process exits.
func WatchTests(args []string) {
	tmpDir, err := os.MkdirTemp("", "gest-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "gest: failed to create temp dir: %v\n", err)
		os.Exit(1)
	}

	// Remove the temp dir when the user presses Ctrl+C or the process is terminated.
	watchSetupCleanup(tmpDir)

	fmt.Println("\ngest: running tests…")
	watchBuildAndRun(tmpDir, args)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gest: failed to create watcher: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = w.Close() }()

	if err := watchAddDirs(w); err != nil {
		fmt.Fprintf(os.Stderr, "gest: failed to watch directories: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\ngest: watching for changes… (Ctrl+C to stop)")

	var (
		mu    sync.Mutex
		timer *time.Timer
	)

	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			if !strings.HasSuffix(event.Name, ".go") {
				continue
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
				continue
			}

			// Debounce: reset the timer on every burst event so rapid saves
			// (e.g. auto-format on save) collapse into a single re-run.
			mu.Lock()
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(200*time.Millisecond, func() {
				fmt.Print("\033[2J\033[3J\033[H")
				watchBuildAndRun(tmpDir, args)
			})
			mu.Unlock()

		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "gest: watcher error: %v\n", err)
		}
	}
}

// watchSetupCleanup removes the temp directory on Ctrl+C or SIGTERM.
func watchSetupCleanup(tmpDir string) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		_ = os.RemoveAll(tmpDir)
		os.Exit(0)
	}()
}

// watchBuildAndRun compiles the project into a temp binary and executes it.
// If compilation fails the errors are printed and the run is skipped.
func watchBuildAndRun(tmpDir string, args []string) {
	bin := tmpDir + "/runner"
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	build := exec.Command("go", "build", "-o", bin, ".")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		return
	}

	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

// watchAddDirs walks the cwd and registers every sub-directory with the
// watcher, skipping vendor and hidden directories (e.g. .git).
func watchAddDirs(w *fsnotify.Watcher) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return filepath.Walk(
		cwd,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				return nil
			}
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" {
				return filepath.SkipDir
			}
			return w.Add(path)
		},
	)
}

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
	defer f.Close()

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

	matchers := []string{
		".ToBe(", ".ToEqual(", ".ToContain(", ".ToBeNil(",
		".ToBeCloseTo(", ".ToHaveLength(", ".ToBeGreaterThan(",
		".ToBeLessThan(", ".ToBeTrue(", ".ToBeFalse(",
	}

	var result []string
	for i := start; i <= end; i++ {
		lineNum := i + 1
		content := lines[i]

		if lineNum == targetLine {
			result = append(
				result,
				fmt.Sprintf(
					"    %s> %3d | %s%s\n",
					red+bold,
					lineNum,
					content,
					reset,
				),
			)
			col := -1
			for _, m := range matchers {
				if idx := strings.Index(content, m); idx >= 0 {
					parenIdx := strings.Index(content[idx:], "(")
					col = idx + parenIdx + 1
					break
				}
			}
			if col >= 0 {
				padding := strings.Repeat(" ", col)
				result = append(
					result,
					fmt.Sprintf(
						"    %s  %3s | %s^%s\n",
						red,
						"",
						padding,
						reset,
					),
				)
			}
		} else {
			result = append(
				result,
				fmt.Sprintf(
					"    %s  %3d | %s%s\n",
					dim,
					lineNum,
					content,
					reset,
				),
			)
		}
	}
	return result
}

func shortPath(full string) string {
	for _, marker := range []string{"example/", "src/", "cmd/", "internal/"} {
		if idx := strings.LastIndex(full, marker); idx >= 0 {
			return full[idx:]
		}
	}
	parts := strings.Split(full, "/")
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return full
}

// ── T ─────────────────────────────────────────────────────────────────────────

// T is the test context passed to each test function.
type T struct {
	failed   bool
	failures []failure
}

func (t *T) record(f failure) {
	t.failed = true
	t.failures = append(t.failures, f)
}

// Expect starts an assertion chain on the given value.
func (t *T) Expect(v interface{}) *Expectation {
	return &Expectation{t: t, actual: v}
}

// ── Expectation ───────────────────────────────────────────────────────────────

// Expectation represents a chainable assertion.
type Expectation struct {
	t      *T
	actual interface{}
	inv    bool
}

// Not negates the next matcher.
func (e *Expectation) Not() *Expectation { e.inv = !e.inv; return e }

func (e *Expectation) pass(b bool) bool {
	if e.inv {
		return !b
	}
	return b
}

func (e *Expectation) caller() (string, int) {
	_, file, line, _ := runtime.Caller(2)
	return file, line
}

func mname(inv bool, name string) string {
	if inv {
		return "not." + name
	}
	return name
}

// ToBe asserts strict equality (==).
func (e *Expectation) ToBe(expected interface{}) {
	file, line := e.caller()
	if !e.pass(e.actual == expected) {
		e.t.record(
			failure{mname(e.inv, "toBe"), expected, e.actual, file, line},
		)
	}
}

// ToEqual asserts deep equality (reflect.DeepEqual).
func (e *Expectation) ToEqual(expected interface{}) {
	file, line := e.caller()
	if !e.pass(reflect.DeepEqual(e.actual, expected)) {
		e.t.record(
			failure{mname(e.inv, "toEqual"), expected, e.actual, file, line},
		)
	}
}

// ToBeNil asserts the value is nil.
func (e *Expectation) ToBeNil() {
	file, line := e.caller()
	isNil := e.actual == nil
	if !isNil {
		v := reflect.ValueOf(e.actual)
		if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			isNil = v.IsNil()
		}
	}
	if !e.pass(isNil) {
		exp := interface{}(nil)
		rec := e.actual
		if e.inv {
			exp = "not nil"
			rec = nil
		}
		e.t.record(failure{mname(e.inv, "toBeNil"), exp, rec, file, line})
	}
}

// ToBeTrue asserts the value is true.
func (e *Expectation) ToBeTrue() {
	file, line := e.caller()
	if !e.pass(e.actual == true) {
		e.t.record(
			failure{mname(e.inv, "toBeTrue"), true, e.actual, file, line},
		)
	}
}

// ToBeFalse asserts the value is false.
func (e *Expectation) ToBeFalse() {
	file, line := e.caller()
	if !e.pass(e.actual == false) {
		e.t.record(
			failure{mname(e.inv, "toBeFalse"), false, e.actual, file, line},
		)
	}
}

// ToContain asserts a string contains the given substring.
func (e *Expectation) ToContain(sub string) {
	file, line := e.caller()
	s, ok := e.actual.(string)
	if !ok {
		e.t.record(
			failure{
				"toContain",
				"string",
				fmt.Sprintf("%T", e.actual),
				file,
				line,
			},
		)
		return
	}
	if !e.pass(strings.Contains(s, sub)) {
		e.t.record(failure{mname(e.inv, "toContain"), sub, s, file, line})
	}
}

// ToHaveLength asserts the length of a string, slice, array or map.
func (e *Expectation) ToHaveLength(n int) {
	file, line := e.caller()
	v := reflect.ValueOf(e.actual)
	var l int
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		l = v.Len()
	default:
		e.t.record(
			failure{
				"toHaveLength",
				n,
				fmt.Sprintf("unsupported type %T", e.actual),
				file,
				line,
			},
		)
		return
	}
	if !e.pass(l == n) {
		e.t.record(failure{mname(e.inv, "toHaveLength"), n, l, file, line})
	}
}

// ToBeGreaterThan asserts the number is greater than n.
func (e *Expectation) ToBeGreaterThan(n float64) {
	file, line := e.caller()
	if !e.pass(toFloat(e.actual) > n) {
		e.t.record(
			failure{
				mname(e.inv, "toBeGreaterThan"),
				fmt.Sprintf("> %v", n),
				e.actual,
				file,
				line,
			},
		)
	}
}

// ToBeLessThan asserts the number is less than n.
func (e *Expectation) ToBeLessThan(n float64) {
	file, line := e.caller()
	if !e.pass(toFloat(e.actual) < n) {
		e.t.record(
			failure{
				mname(e.inv, "toBeLessThan"),
				fmt.Sprintf("< %v", n),
				e.actual,
				file,
				line,
			},
		)
	}
}

// ToBeCloseTo asserts two floats are approximately equal (default delta: 0.001).
func (e *Expectation) ToBeCloseTo(expected float64, delta ...float64) {
	file, line := e.caller()
	d := 0.001
	if len(delta) > 0 {
		d = delta[0]
	}
	if !e.pass(math.Abs(toFloat(e.actual)-expected) <= d) {
		e.t.record(failure{
			mname(e.inv, "toBeCloseTo"),
			fmt.Sprintf("%v (±%v)", expected, d),
			e.actual, file, line,
		})
	}
}

// ── Suite ─────────────────────────────────────────────────────────────────────

type testCase struct {
	name     string
	fn       func(*T)
	passed   bool
	elapsed  time.Duration
	failures []failure
	panicMsg string
}

// Suite represents a group of related tests.
type Suite struct {
	name    string
	tests   []testCase
	passed  int
	failed  int
	elapsed time.Duration
}

// Describe creates a new test suite with the given name.
func Describe(name string) *Suite { return &Suite{name: name} }

// It adds a test case to the suite.
func (s *Suite) It(name string, fn func(*T)) {
	s.tests = append(s.tests, testCase{name: name, fn: fn})
}

func (s *Suite) run() {
	start := time.Now()
	for i := range s.tests {
		tc := &s.tests[i]
		t := &T{}
		ts := time.Now()
		func() {
			defer func() {
				if r := recover(); r != nil {
					tc.panicMsg = fmt.Sprintf("panic: %v", r)
				}
			}()
			tc.fn(t)
		}()
		tc.elapsed = time.Since(ts)
		if t.failed || tc.panicMsg != "" {
			tc.failures = t.failures
			s.failed++
		} else {
			tc.passed = true
			s.passed++
		}
	}
	s.elapsed = time.Since(start)
}

func (s *Suite) printHeader() {
	if s.failed == 0 {
		fmt.Printf(" %s PASS %s  %s%s%s\n", bgGreen, reset, bold, s.name, reset)
	} else {
		fmt.Printf(" %s FAIL %s  %s%s%s\n", bgRed, reset, bold, s.name, reset)
	}
	for _, tc := range s.tests {
		ms := fmtDur(tc.elapsed)
		if tc.passed {
			fmt.Printf(
				"  %s✓%s %s%s%s %s(%s)%s\n",
				green,
				reset,
				dim,
				tc.name,
				reset,
				dim,
				ms,
				reset,
			)
		} else {
			fmt.Printf(
				"  %s✕%s %s%s%s %s(%s)%s\n",
				red,
				reset,
				red+bold,
				tc.name,
				reset,
				dim,
				ms,
				reset,
			)
		}
	}
}

func (s *Suite) printFailures() {
	for _, tc := range s.tests {
		if tc.passed {
			continue
		}

		fmt.Printf("  %s● %s%s\n\n", red+bold, tc.name, reset)

		if tc.panicMsg != "" {
			fmt.Printf("    %s%s%s\n\n", red, tc.panicMsg, reset)
			continue
		}

		for _, f := range tc.failures {
			fmt.Printf(
				"    %sexpect(received).%s(expected)%s\n\n",
				dim,
				f.matcher,
				reset,
			)
			fmt.Printf(
				"    %sExpected:%s %s%v%s\n",
				green+bold,
				reset,
				green,
				f.expected,
				reset,
			)
			fmt.Printf(
				"    %sReceived:%s %s%v%s\n",
				red+bold,
				reset,
				red,
				f.received,
				reset,
			)

			snippet := readSnippet(f.file, f.line, 2)
			if len(snippet) > 0 {
				fmt.Println()
				for _, sl := range snippet {
					fmt.Print(sl)
				}
			}

			fmt.Printf(
				"\n    %sat %s:%d%s\n\n",
				dim,
				shortPath(f.file),
				f.line,
				reset,
			)
		}
	}
}

// ── Coverage table ────────────────────────────────────────────────────────────

func pctToColor(pct float64) string {
	switch {
	case pct >= 80:
		return green
	case pct >= 50:
		return yellow
	default:
		return red
	}
}

func coverageBar(pct float64, width int) string {
	filled := int(math.Round(pct / 100 * float64(width)))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	color := pctToColor(pct)
	empty := width - filled

	// rounded pip-style bar: filled ╺━━━━╸  empty ╺──────╴
	var bar string
	if filled == 0 {
		bar = dim + "╺" + strings.Repeat("─", width-2) + "╴" + reset
	} else if filled == width {
		bar = color + "╺" + strings.Repeat("━", width-2) + "╸" + reset
	} else {
		filledPart := color + "╺" + strings.Repeat("━", filled-1)
		emptyPart := dim + strings.Repeat("─", empty-1) + "╴" + reset
		bar = filledPart + emptyPart
	}
	return bar
}

func printCoverageTable(suites []*Suite) {
	const barWidth = 28

	maxName := len("All suites")
	for _, s := range suites {
		if len(s.name) > maxName {
			maxName = len(s.name)
		}
	}

	tableWidth := 3 + maxName + 2 + barWidth + 2 + 7 + 3 + 12
	sep := dim + strings.Repeat("─", tableWidth) + reset

	fmt.Println()
	fmt.Println(sep)
	fmt.Printf("   %s%-*s  %-*s  %-7s  %s%s\n",
		bold, maxName, "Suite", barWidth, "", "Coverage", "Passed/Total", reset)
	fmt.Println(sep)

	totalPassed, totalTests := 0, 0

	for _, s := range suites {
		total := s.passed + s.failed
		totalPassed += s.passed
		totalTests += total

		var pct float64
		if total > 0 {
			pct = float64(s.passed) / float64(total) * 100
		}

		bar := coverageBar(pct, barWidth)
		pctColor := pctToColor(pct)
		icon := green + "✓" + reset
		if s.failed > 0 {
			icon = red + "✕" + reset
		}

		fmt.Printf(" %s  %s%-*s%s  %s  %s%5.1f%%%s   %s%d/%d%s\n",
			icon,
			bold, maxName, s.name, reset,
			bar,
			pctColor, pct, reset,
			dim, s.passed, total, reset,
		)
	}

	fmt.Println(sep)
	var totalPct float64
	if totalTests > 0 {
		totalPct = float64(totalPassed) / float64(totalTests) * 100
	}
	bar := coverageBar(totalPct, barWidth)
	pctColor := pctToColor(totalPct)

	allIcon := green + "✓" + reset
	if totalPassed < totalTests {
		allIcon = red + "✕" + reset
	}

	fmt.Printf(" %s  %s%-*s%s  %s  %s%5.1f%%%s   %s%d/%d%s\n",
		allIcon,
		bold, maxName, "All suites", reset,
		bar,
		pctColor+bold, totalPct, reset,
		dim, totalPassed, totalTests, reset,
	)
	fmt.Println(sep)
	fmt.Println()
}

// ── RunAll ────────────────────────────────────────────────────────────────────

// RunAll runs all given suites and prints the full Jest-style output.
// Pass -c or --coverage to display the coverage table.
func RunAll(suites ...*Suite) bool {
	showCoverage := false
	for _, arg := range os.Args[1:] {
		if arg == "-c" || arg == "--coverage" {
			showCoverage = true
		}
	}

	for _, s := range suites {
		s.run()
	}

	fmt.Println()
	for _, s := range suites {
		s.printHeader()
	}

	hasFailed := false
	for _, s := range suites {
		if s.failed > 0 {
			hasFailed = true
			break
		}
	}

	if hasFailed {
		fmt.Println()
		for _, s := range suites {
			if s.failed > 0 {
				s.printFailures()
			}
		}
	}

	totalSuitesFailed := 0
	totalPassed, totalFailed := 0, 0
	var totalTime time.Duration
	for _, s := range suites {
		if s.failed > 0 {
			totalSuitesFailed++
		}
		totalPassed += s.passed
		totalFailed += s.failed
		totalTime += s.elapsed
	}
	suitesPassed := len(suites) - totalSuitesFailed

	fmt.Println()
	if totalSuitesFailed > 0 {
		fmt.Printf(
			"%sTest Suites:%s %s%d failed%s, %d passed, %d total\n",
			bold,
			reset,
			red+bold,
			totalSuitesFailed,
			reset,
			suitesPassed,
			len(suites),
		)
	} else {
		fmt.Printf("%sTest Suites:%s %s%d passed%s, %d total\n",
			bold, reset, green+bold, suitesPassed, reset, len(suites))
	}
	if totalFailed > 0 {
		fmt.Printf(
			"%sTests:%s       %s%d failed%s, %d passed, %d total\n",
			bold,
			reset,
			red+bold,
			totalFailed,
			reset,
			totalPassed,
			totalPassed+totalFailed,
		)
	} else {
		fmt.Printf("%sTests:%s       %s%d passed%s, %d total\n",
			bold, reset, green+bold, totalPassed, reset, totalPassed)
	}
	fmt.Printf("%sSnapshots:%s   0 total\n", bold, reset)
	fmt.Printf("%sTime:%s        %s\n", bold, reset, fmtDur(totalTime))
	fmt.Printf("%sRan all test suites.%s\n", dim, reset)

	if showCoverage {
		printCoverageTable(suites)
	}

	return totalFailed == 0
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
