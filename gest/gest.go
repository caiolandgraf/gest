package gest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
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
	display string
	passed  bool
	failed  bool
	skipped bool
	elapsed time.Duration
	output  []string
}

type suite struct {
	pkg         string
	name        string
	topTest     string
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
	name := strings.TrimPrefix(topTest, "Test")
	name = strings.ReplaceAll(name, "_", " ")
	return &suite{
		pkg:     pkg,
		topTest: topTest,
		name:    name,
		byName:  make(map[string]*testCase),
	}
}

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

type pkgResult struct {
	pkg     string
	elapsed time.Duration
}

func parseStream(r io.Reader) ([]*suite, []*pkgResult, bool) {
	byKey := make(map[string]*suite, 32)
	suitesByPkg := make(map[string][]*suite, 8)
	var suites []*suite
	allPassed := true

	var pkgResults []*pkgResult
	pkgBuildOutput := make(map[string][]string, 4)

	key := func(pkg, top string) string { return pkg + "::" + top }

	getSuite := func(pkg, top string) *suite {
		k := key(pkg, top)
		if s, ok := byKey[k]; ok {
			return s
		}
		s := newSuite(pkg, top)
		byKey[k] = s
		suites = append(suites, s)
		suitesByPkg[pkg] = append(suitesByPkg[pkg], s)
		return s
	}

	topOf := func(testName string) string {
		if idx := strings.IndexByte(testName, '/'); idx >= 0 {
			return testName[:idx]
		}
		return testName
	}

	dec := json.NewDecoder(bufio.NewReaderSize(r, 64*1024))

	for {
		var ev testEvent
		if err := dec.Decode(&ev); err != nil {
			break
		}
		if ev.Package == "" {
			continue
		}

		slashIdx := strings.IndexByte(ev.Test, '/')
		isSubtest := slashIdx >= 0

		switch ev.Action {

		case "output":
			if ev.Test == "" {
				pkgBuildOutput[ev.Package] = append(pkgBuildOutput[ev.Package], ev.Output)
				continue
			}
			top := topOf(ev.Test)
			s := getSuite(ev.Package, top)

			const descPrefix = "gest:describe:"
			if trimmed := strings.TrimSpace(ev.Output); strings.Contains(trimmed, descPrefix) {
				if idx := strings.Index(trimmed, descPrefix); idx >= 0 {
					s.name = trimmed[idx+len(descPrefix):]
				}
				continue
			}

			if isSubtest {
				subName := stripDupSuffix(ev.Test[slashIdx+1:])
				tc := s.getOrCreateItem(subName)
				tc.output = append(tc.output, ev.Output)
			}

		case "run":
			if ev.Test != "" && !isSubtest {
				getSuite(ev.Package, ev.Test)
			}

		case "pass":
			if ev.Test == "" {
				pkgResults = append(pkgResults, &pkgResult{
					pkg:     ev.Package,
					elapsed: time.Duration(ev.Elapsed * float64(time.Second)),
				})
				continue
			}
			top := topOf(ev.Test)
			s := getSuite(ev.Package, top)
			elapsed := time.Duration(ev.Elapsed * float64(time.Second))

			if !isSubtest {
				s.elapsed = elapsed
			} else {
				subName := stripDupSuffix(ev.Test[slashIdx+1:])
				tc := s.getOrCreateItem(subName)
				tc.passed = true
				tc.elapsed = elapsed
				s.elapsed += elapsed
				s.passed++
			}

		case "fail":
			if ev.Test == "" {
				allPassed = false
				if len(suitesByPkg[ev.Package]) == 0 {
					pkg := ev.Package
					pkgName := pkg
					if idx := strings.LastIndexByte(pkg, '/'); idx >= 0 {
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
					suitesByPkg[pkg] = append(suitesByPkg[pkg], s)
				} else {
					for _, s := range suitesByPkg[ev.Package] {
						if len(s.tests) == 0 {
							s.buildFailed = true
							s.buildOutput = pkgBuildOutput[ev.Package]
						}
					}
				}
				continue
			}
			top := topOf(ev.Test)
			s := getSuite(ev.Package, top)

			if !isSubtest {
				s.elapsed = time.Duration(ev.Elapsed * float64(time.Second))
				if len(s.tests) == 0 {
					s.buildFailed = true
					s.buildOutput = pkgBuildOutput[ev.Package]
				}
				allPassed = false
			} else {
				subName := stripDupSuffix(ev.Test[slashIdx+1:])
				tc := s.getOrCreateItem(subName)
				tc.failed = true
				tc.elapsed = time.Duration(ev.Elapsed * float64(time.Second))
				s.elapsed += tc.elapsed
				s.failed++
				allPassed = false
			}

		case "skip":
			if ev.Test == "" || !isSubtest {
				continue
			}
			top := topOf(ev.Test)
			s := getSuite(ev.Package, top)
			subName := stripDupSuffix(ev.Test[slashIdx+1:])
			tc := s.getOrCreateItem(subName)
			tc.skipped = true
			tc.elapsed = time.Duration(ev.Elapsed * float64(time.Second))
			s.elapsed += tc.elapsed
			s.skipped++
		}
	}

	return suites, pkgResults, allPassed
}

// stripDupSuffix removes Go's "#NN" deduplication suffix.
func stripDupSuffix(name string) string {
	if idx := strings.LastIndexByte(name, '#'); idx >= 0 {
		suffix := name[idx+1:]
		if len(suffix) > 0 {
			allDigits := true
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
			if strings.ContainsRune(trimmed, '\033') {
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

	switch {
	case filled == 0:
		return dim + "╺" + strings.Repeat("─", width-2) + "╴" + reset
	case filled == width:
		return color + "╺" + strings.Repeat("━", width-2) + "╸" + reset
	default:
		return color + "╺" + strings.Repeat("━", filled-1) +
			dim + strings.Repeat("─", empty-1) + "╴" + reset
	}
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
	allIcon := green + "✓" + reset
	if totalPassed < totalTests {
		allIcon = red + "✕" + reset
	}
	fmt.Printf(" %s  %s%-*s%s  %s  %s%5.1f%%%s   %s%d/%d%s\n",
		allIcon, bold, maxName, "All suites", reset, coverageBar(totalPct, barWidth),
		pctToColor(totalPct)+bold, totalPct, reset, dim, totalPassed, totalTests, reset)
	fmt.Println(sep)
	fmt.Println()
}

// ── run ───────────────────────────────────────────────────────────────────────

// runAll executa os testes para os pacotes especificados.
// Se packagePaths estiver vazio, executa para "./...".
func runAll(packagePaths []string, extraGoTestArgs []string, showCoverage bool, noCache bool) bool {
	goArgs := []string{"test", "-json"}

	if noCache {
		goArgs = append(goArgs, "-count=1")
	}

	if len(packagePaths) == 0 {
		goArgs = append(goArgs, "./...") // Padrão: todos os pacotes
	} else {
		goArgs = append(goArgs, packagePaths...) // Rodar apenas os pacotes especificados
	}

	goArgs = append(goArgs, extraGoTestArgs...) // Argumentos adicionais passados pelo usuário

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

	suites, pkgResults, allPassed := parseStream(stdout)
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
	for _, s := range suites {
		if s.failed > 0 || s.buildFailed {
			totalSuitesFailed++
		}
		totalPassed += s.passed
		totalFailed += s.failed
		totalSkipped += s.skipped
	}

	var totalTime time.Duration
	for _, pr := range pkgResults {
		totalTime += pr.elapsed
	}
	if totalTime == 0 {
		for _, s := range suites {
			totalTime += s.elapsed
		}
	}

	suitesPassed := len(suites) - totalSuitesFailed
	total := totalPassed + totalFailed + totalSkipped

	fmt.Println()
	if totalSuitesFailed > 0 {
		fmt.Printf("%sTest Suites:%s %s%d failed%s, %d passed, %d total\n",
			bold, reset, red+bold, totalSuitesFailed, reset, suitesPassed, len(suites))
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

// packageFileMap armazena o mapeamento de caminho de arquivo para o import path do pacote Go.
var packageFileMap map[string]string

// buildPackageFileMap preenche packageFileMap executando 'go list'.
func buildPackageFileMap() error {
	packageFileMap = make(map[string]string)
	cmd := exec.Command("go", "list", "-f", "{{.Dir}} {{.ImportPath}}", "./...")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run 'go list': %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		dir := parts[0]
		importPath := parts[1]

		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
				fullPath := filepath.Join(dir, file.Name())
				packageFileMap[fullPath] = importPath
			}
		}
	}
	return nil
}

func watchTests(extraGoTestArgs []string, showCoverage bool, debounce time.Duration, noCache bool) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("\ngest: stopping.")
		os.Exit(0)
	}()

	fmt.Println("\ngest: running initial tests…")
	// Na primeira execução, rodamos todos os testes.
	runAll([]string{"./..."}, extraGoTestArgs, showCoverage, noCache)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gest: %v\n", err)
		os.Exit(1)
	}
	defer w.Close()

	if err := watchAddDirs(w); err != nil {
		fmt.Fprintf(os.Stderr, "gest: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\ngest: watching for changes… (Ctrl+C to stop)")

	var (
		mu              sync.Mutex
		timer           *time.Timer
		changedPackages = make(map[string]struct{}) // Usamos um map para garantir pacotes únicos
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

			pkgPath, found := packageFileMap[event.Name]
			if !found {
				// Se o arquivo não foi mapeado (ex: novo arquivo em um novo pacote, ou arquivo temporário do editor)
				// Por simplicidade, vamos re-executar todos os testes para garantir que nada seja perdido.
				// Em um sistema mais robusto, poderíamos tentar re-mapear ou ignorar.
				mu.Lock()
				changedPackages["./..."] = struct{}{} // Força a rodar tudo
				mu.Unlock()
			} else {
				mu.Lock()
				changedPackages[pkgPath] = struct{}{} // Adiciona o pacote ao conjunto de pacotes modificados
				mu.Unlock()
			}

			mu.Lock()
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounce, func() {
				mu.Lock()
				defer mu.Unlock()

				if len(changedPackages) == 0 {
					return // Nada para rodar
				}

				packagesToRun := make([]string, 0, len(changedPackages))
				for pkg := range changedPackages {
					packagesToRun = append(packagesToRun, pkg)
				}
				// Limpa o map para a próxima rodada.
				// Usar um loop para compatibilidade com Go < 1.21
				for k := range changedPackages {
					delete(changedPackages, k)
				}

				fmt.Print("\033[2J\033[3J\033[H")
				runAll(packagesToRun, extraGoTestArgs, showCoverage, noCache)
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
	return filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
			return filepath.SkipDir
		}
		return w.Add(path)
	})
}

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	var watchMode bool
	var showCoverage bool
	var noCache bool
	var debounce = 80 * time.Millisecond
	var extraGoTestArgs []string // Renomeado para clareza

	for _, arg := range os.Args[1:] {
		switch {
		case arg == "--watch" || arg == "-w":
			watchMode = true
		case arg == "--coverage" || arg == "-c":
			showCoverage = true
		case arg == "--no-cache": // Nova flag para desativar o cache do go test
			noCache = true
		case strings.HasPrefix(arg, "--debounce="):
			if ms, err := time.ParseDuration(arg[len("--debounce="):]); err == nil {
				debounce = ms
			}
		default:
			extraGoTestArgs = append(extraGoTestArgs, arg)
		}
	}

	if len(extraGoTestArgs) == 0 {
		extraGoTestArgs = []string{"./..."} // Default para go test se nenhum argumento for passado
	}

	if watchMode {
		if err := buildPackageFileMap(); err != nil {
			fmt.Fprintf(os.Stderr, "gest: error building package map: %v\n", err)
			os.Exit(1)
		}
		watchTests(extraGoTestArgs, showCoverage, debounce, noCache)
		return
	}

	// Modo de execução normal (não watch)
	if !runAll(
		extraGoTestArgs,
		nil,
		showCoverage,
		noCache,
	) { // No modo normal, extraGoTestArgs já contém os pacotes
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
