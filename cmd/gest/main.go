package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
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
	red     = "\033[38;2;220;60;60m"
	green   = "\033[38;2;40;210;90m"
	yellow  = "\033[38;2;230;200;40m"
	bgGreen = "\033[48;2;40;180;80m\033[38;2;255;255;255m"
	bgRed   = "\033[48;2;195;55;55m\033[38;2;255;255;255m"
)

// ── go test -json event ───────────────────────────────────────────────────────

type testEvent struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
}

// ── internal model ────────────────────────────────────────────────────────────

type testCase struct {
	name    string
	display string // friendly name shown in output (the It() description)
	passed  bool
	failed  bool
	skipped bool
	elapsed time.Duration
	output  []string // output lines captured for failures
}

type suite struct {
	pkg         string
	name        string // the Describe("name") value, or TestXxx fallback
	topTest     string // the TestXxx function name
	tests       []*testCase
	byName      map[string]*testCase
	passed      int
	failed      int
	skipped     int
	elapsed     time.Duration
	buildFailed bool
	buildOutput []string
}

func newSuite(pkg, topTest string) *suite {
	// default display name: strip "Test" prefix, underscores → spaces
	name := topTest
	if strings.HasPrefix(name, "Test") {
		name = name[4:]
	}
	name = strings.ReplaceAll(name, "_", " ")
	return &suite{
		pkg:     pkg,
		topTest: topTest,
		name:    name,
		byName:  make(map[string]*testCase),
	}
}

// getOrCreateItem returns the testCase for a subtest name (the It() description).
// It un-escapes Go's test name encoding (underscores back to spaces, etc.).
func (s *suite) getOrCreateItem(subName string) *testCase {
	if tc, ok := s.byName[subName]; ok {
		return tc
	}
	display := strings.ReplaceAll(subName, "_", " ")
	tc := &testCase{name: subName, display: display}
	s.tests = append(s.tests, tc)
	s.byName[subName] = tc
	return tc
}

// ── stream parser ─────────────────────────────────────────────────────────────

// parseStream reads go test -json events and builds one suite per top-level
// TestXxx function. Each It() subtest becomes a testCase inside that suite.
// The Describe("name") value is captured from the "gest:describe:<name>" log
// line emitted by Suite.Run and used as the suite display name.
func parseStream(r io.Reader) ([]*suite, bool) {
	// key: "pkg::TestXxx"
	byKey := map[string]*suite{}
	var suites []*suite
	allPassed := true

	// pkgBuildOutput holds package-level output lines (build errors, etc.)
	pkgBuildOutput := map[string][]string{}
	pkgBuildFailed := map[string]bool{}

	key := func(pkg, top string) string { return pkg + "::" + top }

	getSuite := func(pkg, top string) *suite {
		k := key(pkg, top)
		if s, ok := byKey[k]; ok {
			return s
		}
		s := newSuite(pkg, top)
		byKey[k] = s
		suites = append(suites, s)
		return s
	}

	// topOf returns the top-level TestXxx name from any test name.
	topOf := func(testName string) string {
		if idx := strings.Index(testName, "/"); idx >= 0 {
			return testName[:idx]
		}
		return testName
	}

	dec := json.NewDecoder(bufio.NewReader(r))
	for {
		var ev testEvent
		if err := dec.Decode(&ev); err != nil {
			break
		}
		if ev.Package == "" {
			continue
		}

		isSubtest := strings.Contains(ev.Test, "/")

		switch ev.Action {
		case "output":
			if ev.Test == "" {
				pkgBuildOutput[ev.Package] = append(
					pkgBuildOutput[ev.Package],
					ev.Output,
				)
				continue
			}
			top := topOf(ev.Test)
			s := getSuite(ev.Package, top)

			// Check for the gest:describe sentinel emitted by Suite.Run.
			// go test formats t.Log output as "    file:line: message\n"
			// It can arrive on the top-level test or a subtest output line.
			const descPrefix = "gest:describe:"
			if trimmed := strings.TrimSpace(
				ev.Output,
			); strings.Contains(
				trimmed,
				descPrefix,
			) {
				idx := strings.Index(trimmed, descPrefix)
				s.name = trimmed[idx+len(descPrefix):]
				continue
			}

			if isSubtest {
				subName := stripDupSuffix(
					ev.Test[strings.Index(ev.Test, "/")+1:],
				)
				tc := s.getOrCreateItem(subName)
				tc.output = append(tc.output, ev.Output)
			}

		case "run":
			if ev.Test != "" && !isSubtest {
				getSuite(ev.Package, ev.Test)
			}

		case "pass":
			if ev.Test == "" {
				// package done — attach any build output
				for _, s := range suites {
					if s.pkg == ev.Package {
						s.elapsed = time.Duration(
							ev.Elapsed * float64(time.Second),
						)
					}
				}
				continue
			}
			top := topOf(ev.Test)
			s := getSuite(ev.Package, top)
			elapsed := time.Duration(ev.Elapsed * float64(time.Second))

			if !isSubtest {
				// top-level TestXxx finished — elapsed already set by subtest accumulation
				if s.elapsed == 0 {
					s.elapsed = time.Duration(ev.Elapsed * float64(time.Second))
				}
			} else {
				subName := stripDupSuffix(
					ev.Test[strings.Index(ev.Test, "/")+1:],
				)
				tc := s.getOrCreateItem(subName)
				tc.passed = true
				tc.elapsed = elapsed
				s.elapsed += elapsed
				s.passed++
			}

		case "fail":
			if ev.Test == "" {
				allPassed = false
				// mark build failure for suites in this package that have no tests
				if _, hasSuites := func() (*suite, bool) {
					for _, s := range suites {
						if s.pkg == ev.Package {
							return s, true
						}
					}
					return nil, false
				}(); !hasSuites {
					// create a sentinel build-failed suite for the package
					pkg := ev.Package
					pkgName := pkg
					if idx := strings.LastIndex(pkg, "/"); idx >= 0 {
						pkgName = pkg[idx+1:]
					}
					s := &suite{
						pkg:         pkg,
						name:        pkgName,
						buildFailed: true,
						buildOutput: pkgBuildOutput[pkg],
						byName:      make(map[string]*testCase),
					}
					suites = append(suites, s)
					pkgBuildFailed[pkg] = true
				} else {
					for _, s := range suites {
						if s.pkg == ev.Package && len(s.tests) == 0 {
							s.buildFailed = true
							s.buildOutput = pkgBuildOutput[ev.Package]
						}
					}
				}
				continue
			}
			top := topOf(ev.Test)
			s := getSuite(ev.Package, top)
			elapsed := time.Duration(ev.Elapsed * float64(time.Second))

			if !isSubtest {
				if s.elapsed == 0 {
					s.elapsed = time.Duration(ev.Elapsed * float64(time.Second))
				}
				// if no subtests were recorded, this is a build/compile failure
				if len(s.tests) == 0 {
					s.buildFailed = true
					s.buildOutput = pkgBuildOutput[ev.Package]
				}
				allPassed = false
			} else {
				subName := stripDupSuffix(
					ev.Test[strings.Index(ev.Test, "/")+1:],
				)
				tc := s.getOrCreateItem(subName)
				tc.failed = true
				tc.elapsed = elapsed
				s.elapsed += elapsed
				s.failed++
				allPassed = false
			}

		case "skip":
			if ev.Test == "" || !isSubtest {
				continue
			}
			top := topOf(ev.Test)
			s := getSuite(ev.Package, top)
			subName := stripDupSuffix(ev.Test[strings.Index(ev.Test, "/")+1:])
			tc := s.getOrCreateItem(subName)
			tc.skipped = true
			elapsed := time.Duration(ev.Elapsed * float64(time.Second))
			tc.elapsed = elapsed
			s.elapsed += elapsed
			s.skipped++
		}
	}

	return suites, allPassed
}

// stripDupSuffix removes Go's automatic "#NN" deduplication suffix from
// subtest names that arise when two It() cases share the same description.
// e.g. "adding_2_+_2_should_return_4#01" → "adding_2_+_2_should_return_4"
func stripDupSuffix(name string) string {
	if idx := strings.LastIndex(name, "#"); idx >= 0 {
		suffix := name[idx+1:]
		allDigits := len(suffix) > 0
		for _, c := range suffix {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			return name[:idx]
		}
	}
	return name
}

// ── renderer ──────────────────────────────────────────────────────────────────

func printSuiteHeader(s *suite) {
	if s.buildFailed || s.failed > 0 {
		fmt.Printf(" %s FAIL %s  %s%s%s\n", bgRed, reset, bold, s.name, reset)
	} else {
		fmt.Printf(" %s PASS %s  %s%s%s\n", bgGreen, reset, bold, s.name, reset)
	}

	if s.buildFailed {
		for _, line := range s.buildOutput {
			fmt.Printf("  %s%s%s\n", red, strings.TrimRight(line, "\n"), reset)
		}
		return
	}

	for _, tc := range s.tests {
		ms := fmtDur(tc.elapsed)
		switch {
		case tc.skipped:
			fmt.Printf("  %s-%s %s%s%s %s(%s)%s\n",
				yellow, reset, dim, tc.display, reset, dim, ms, reset)
		case tc.passed:
			fmt.Printf("  %s✓%s %s%s%s %s(%s)%s\n",
				green, reset, dim, tc.display, reset, dim, ms, reset)
		default:
			fmt.Printf("  %s✕%s %s%s%s %s(%s)%s\n",
				red, reset, red+bold, tc.display, reset, dim, ms, reset)
		}
	}
}

func printSuiteFailures(s *suite) {
	for _, tc := range s.tests {
		if !tc.failed {
			continue
		}
		fmt.Printf("  %s● %s%s\n\n", red+bold, tc.display, reset)
		for _, line := range tc.output {
			trimmed := strings.TrimRight(line, "\n")
			if trimmed == "" {
				continue
			}
			// skip go test runner noise lines
			stripped := strings.TrimSpace(trimmed)
			if strings.HasPrefix(stripped, "=== RUN") ||
				strings.HasPrefix(stripped, "=== PAUSE") ||
				strings.HasPrefix(stripped, "=== CONT") ||
				strings.HasPrefix(stripped, "--- PASS") ||
				strings.HasPrefix(stripped, "--- FAIL") ||
				strings.HasPrefix(stripped, "--- SKIP") ||
				strings.Contains(stripped, "gest:describe:") {
				continue
			}
			// lines already formatted by the gest lib (contain ANSI ESC) pass through as-is
			if strings.Contains(trimmed, "\033[") {
				fmt.Println(trimmed)
			} else {
				fmt.Printf("    %s%s%s\n", dim, trimmed, reset)
			}
		}
		fmt.Println()
	}
}

// ── coverage table ────────────────────────────────────────────────────────────

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

	var bar string
	switch {
	case filled == 0:
		bar = dim + "╺" + strings.Repeat("─", width-2) + "╴" + reset
	case filled == width:
		bar = color + "╺" + strings.Repeat("━", width-2) + "╸" + reset
	default:
		bar = color + "╺" + strings.Repeat("━", filled-1) +
			dim + strings.Repeat("─", empty-1) + "╴" + reset
	}
	return bar
}

func printCoverageTable(suites []*suite) {
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
		total := s.passed + s.failed + s.skipped
		totalPassed += s.passed
		totalTests += total

		var pct float64
		if total > 0 {
			pct = float64(s.passed) / float64(total) * 100
		}
		bar := coverageBar(pct, barWidth)
		pctColor := pctToColor(pct)
		icon := green + "✓" + reset
		if s.failed > 0 || s.buildFailed {
			icon = red + "✕" + reset
		}
		fmt.Printf(" %s  %s%-*s%s  %s  %s%5.1f%%%s   %s%d/%d%s\n",
			icon, bold, maxName, s.name, reset, bar,
			pctColor, pct, reset, dim, s.passed, total, reset)
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
		allIcon, bold, maxName, "All suites", reset, bar,
		pctColor+bold, totalPct, reset, dim, totalPassed, totalTests, reset)
	fmt.Println(sep)
	fmt.Println()
}

// ── run ───────────────────────────────────────────────────────────────────────

func runAll(args []string, showCoverage bool) bool {
	goArgs := append([]string{"test", "-v", "-json"}, args...)
	cmd := exec.Command("go", goArgs...)
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gest: %v\n", err)
		return false
	}
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "gest: %v\n", err)
		return false
	}

	suites, allPassed := parseStream(stdout)
	_ = cmd.Wait()

	fmt.Println()
	for _, s := range suites {
		printSuiteHeader(s)
	}

	hasFailed := false
	for _, s := range suites {
		if s.failed > 0 || s.buildFailed {
			hasFailed = true
			break
		}
	}
	if hasFailed {
		fmt.Println()
		for _, s := range suites {
			if s.failed > 0 {
				printSuiteFailures(s)
			}
		}
	}

	totalSuitesFailed := 0
	totalPassed, totalFailed, totalSkipped := 0, 0, 0
	var totalTime time.Duration
	for _, s := range suites {
		if s.failed > 0 || s.buildFailed {
			totalSuitesFailed++
		}
		totalPassed += s.passed
		totalFailed += s.failed
		totalSkipped += s.skipped
		totalTime += s.elapsed
	}
	suitesPassed := len(suites) - totalSuitesFailed
	total := totalPassed + totalFailed + totalSkipped

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
		fmt.Printf("%sTests:%s       %s%d failed%s, %d passed, %d total\n",
			bold, reset, red+bold, totalFailed, reset, totalPassed, total)
	} else {
		fmt.Printf("%sTests:%s       %s%d passed%s, %d total\n",
			bold, reset, green+bold, totalPassed, reset, total)
	}
	fmt.Printf("%sSnapshots:%s   0 total\n", bold, reset)
	fmt.Printf("%sTime:%s        %s\n", bold, reset, fmtDur(totalTime))
	fmt.Printf("%sRan all test suites.%s\n", dim, reset)

	if showCoverage {
		printCoverageTable(suites)
	}

	return allPassed
}

// ── watch mode ────────────────────────────────────────────────────────────────

func watchTests(args []string, showCoverage bool) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("\ngest: stopping.")
		os.Exit(0)
	}()

	fmt.Println("\ngest: running tests…")
	runAll(args, showCoverage)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gest: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = w.Close() }()

	if err := watchAddDirs(w); err != nil {
		fmt.Fprintf(os.Stderr, "gest: %v\n", err)
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
			mu.Lock()
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(30*time.Millisecond, func() {
				fmt.Print("\033[2J\033[3J\033[H")
				runAll(args, showCoverage)
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

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	var watchMode bool
	var showCoverage bool
	var passThrough []string

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--watch", "-w":
			watchMode = true
		case "--coverage", "-c":
			showCoverage = true
		default:
			passThrough = append(passThrough, arg)
		}
	}

	if len(passThrough) == 0 {
		passThrough = []string{"./..."}
	}

	if watchMode {
		watchTests(passThrough, showCoverage)
		return
	}

	if !runAll(passThrough, showCoverage) {
		os.Exit(1)
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
