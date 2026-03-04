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
	display string // friendly name shown in output (strips "Test" prefix)
	passed  bool
	failed  bool
	skipped bool
	elapsed time.Duration
	output  []string // output lines captured for this top-level test
}

type suite struct {
	pkg         string
	name        string
	tests       []*testCase
	byName      map[string]*testCase
	passed      int
	failed      int
	skipped     int
	elapsed     time.Duration
	buildFailed bool
	buildOutput []string
}

func newSuite(pkg string) *suite {
	name := pkg
	if idx := strings.LastIndex(pkg, "/"); idx >= 0 {
		name = pkg[idx+1:]
	}
	return &suite{pkg: pkg, name: name, byName: make(map[string]*testCase)}
}

// getOrCreate returns the top-level testCase for the given test name.
// Subtests (names containing "/") are grouped under their parent.
func (s *suite) getOrCreate(testName string) *testCase {
	top := testName
	if idx := strings.Index(testName, "/"); idx >= 0 {
		top = testName[:idx]
	}
	if tc, ok := s.byName[top]; ok {
		return tc
	}
	// build a friendly display name: strip leading "Test" prefix and
	// replace underscores with spaces so "TestCalculator" → "Calculator"
	display := top
	if strings.HasPrefix(display, "Test") {
		display = display[4:]
	}
	display = strings.ReplaceAll(display, "_", " ")
	tc := &testCase{name: top, display: display}
	s.tests = append(s.tests, tc)
	s.byName[top] = tc
	return tc
}

// ── stream parser ─────────────────────────────────────────────────────────────

func parseStream(r io.Reader) ([]*suite, bool) {
	var suites []*suite
	byPkg := map[string]*suite{}
	allPassed := true

	dec := json.NewDecoder(bufio.NewReader(r))
	for {
		var ev testEvent
		if err := dec.Decode(&ev); err != nil {
			break
		}
		if ev.Package == "" {
			continue
		}

		s, exists := byPkg[ev.Package]
		if !exists {
			s = newSuite(ev.Package)
			byPkg[ev.Package] = s
			suites = append(suites, s)
		}

		isSubtest := strings.Contains(ev.Test, "/")

		switch ev.Action {
		case "output":
			if ev.Test == "" {
				// package-level output (build errors, etc.)
				s.buildOutput = append(s.buildOutput, ev.Output)
			} else {
				// accumulate subtest output on the parent so we can print it on failure
				tc := s.getOrCreate(ev.Test)
				if isSubtest {
					tc.output = append(tc.output, ev.Output)
				}
			}

		case "run":
			if ev.Test != "" {
				s.getOrCreate(ev.Test)
			}

		case "pass":
			if ev.Test == "" {
				s.elapsed = time.Duration(ev.Elapsed * float64(time.Second))
			} else {
				tc := s.getOrCreate(ev.Test)
				elapsed := time.Duration(ev.Elapsed * float64(time.Second))
				if !isSubtest {
					tc.passed = true
					tc.elapsed = elapsed
					s.passed++
				} else if elapsed > tc.elapsed {
					// accumulate subtest elapsed onto parent so timing is meaningful
					tc.elapsed += elapsed
				}
			}

		case "fail":
			if ev.Test == "" {
				s.elapsed = time.Duration(ev.Elapsed * float64(time.Second))
				if len(s.tests) == 0 {
					s.buildFailed = true
				}
				allPassed = false
			} else {
				tc := s.getOrCreate(ev.Test)
				elapsed := time.Duration(ev.Elapsed * float64(time.Second))
				if !isSubtest {
					tc.failed = true
					tc.elapsed = elapsed
					s.failed++
					allPassed = false
				} else if elapsed > tc.elapsed {
					tc.elapsed += elapsed
				}
			}

		case "skip":
			if ev.Test != "" && !isSubtest {
				tc := s.getOrCreate(ev.Test)
				tc.skipped = true
				tc.elapsed = time.Duration(ev.Elapsed * float64(time.Second))
				s.skipped++
			}
		}
	}

	return suites, allPassed
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
				strings.HasPrefix(stripped, "--- SKIP") {
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
