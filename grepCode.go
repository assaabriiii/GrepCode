package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	searchPattern  string
	searchDir      string
	caseSensitive  bool
	useRegex       bool
	fileExtensions string
	showContext    int
	showStats      bool
	excludeDirs    string
)

func initFlags() {
	flag.StringVar(&searchPattern, "pattern", "", "Search pattern (required)")
	flag.StringVar(&searchDir, "dir", ".", "Directory to search")
	flag.BoolVar(&caseSensitive, "case", false, "Case-sensitive search")
	flag.BoolVar(&useRegex, "regex", false, "Use regular expressions")
	flag.StringVar(&fileExtensions, "ext", "", "Comma-separated file extensions (e.g., go,js,ts)")
	flag.IntVar(&showContext, "context", 0, "Show N lines of context around match")
	flag.BoolVar(&showStats, "stats", false, "Show search statistics")
	flag.StringVar(&excludeDirs, "exclude-dir", "", "Comma-separated directories to exclude")
	flag.Parse()
}

func main() {
	initFlags()
	if searchPattern == "" {
		color.Red("Error: search pattern is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	startTime := time.Now()
	searcher := NewCodeSearcher()
	results, stats, err := searcher.Search()
	if err != nil {
		color.Red("Search error: %v", err)
		os.Exit(1)
	}

	for _, result := range results {
		result.Print()
	}

	if showStats {
		printStats(stats, time.Since(startTime))
	}
}

type SearchResult struct {
	File    string
	Line    int
	Content string
	Context []string
	IsMatch bool
}

type SearchStats struct {
	FilesScanned int
	MatchesFound int
	TotalLines   int
}

type CodeSearcher struct {
	fileFilter    map[string]struct{}
	excludeFilter map[string]struct{}
	pattern       *regexp.Regexp
}

func NewCodeSearcher() *CodeSearcher {
	cs := &CodeSearcher{}
	cs.initFilters()
	cs.compilePattern()
	return cs
}

func (cs *CodeSearcher) initFilters() {
	cs.fileFilter = make(map[string]struct{})
	if fileExtensions != "" {
		for _, ext := range strings.Split(fileExtensions, ",") {
			cs.fileFilter["."+strings.TrimSpace(ext)] = struct{}{}
		}
	}

	cs.excludeFilter = make(map[string]struct{})
	if excludeDirs != "" {
		for _, dir := range strings.Split(excludeDirs, ",") {
			cs.excludeFilter[strings.TrimSpace(dir)] = struct{}{}
		}
	}
}

func (cs *CodeSearcher) compilePattern() error {
	pattern := searchPattern
	if !useRegex {
		pattern = regexp.QuoteMeta(pattern)
	}
	if !caseSensitive {
		pattern = "(?i)" + pattern
	}

	var err error
	cs.pattern, err = regexp.Compile(pattern)
	return err
}

func (cs *CodeSearcher) Search() ([]SearchResult, SearchStats, error) {
	var results []SearchResult
	stats := SearchStats{}

	err := filepath.WalkDir(searchDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			if _, excluded := cs.excludeFilter[d.Name()]; excluded {
				return filepath.SkipDir
			}
			return nil
		}

		if !cs.isValidFile(path) {
			return nil
		}

		fileResults, fileStats, err := cs.searchFile(path)
		if err != nil {
			return nil
		}

		results = append(results, fileResults...)
		stats.FilesScanned++
		stats.MatchesFound += len(fileResults)
		stats.TotalLines += fileStats.TotalLines
		return nil
	})

	return results, stats, err
}

func (cs *CodeSearcher) isValidFile(path string) bool {
	if len(cs.fileFilter) == 0 {
		return true
	}
	ext := filepath.Ext(path)
	_, ok := cs.fileFilter[ext]
	return ok
}

func (cs *CodeSearcher) searchFile(path string) ([]SearchResult, SearchStats, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, SearchStats{}, err
	}
	defer file.Close()

	var results []SearchResult
	stats := SearchStats{}
	scanner := bufio.NewScanner(file)

	var contextBuffer []string
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		stats.TotalLines++
		line := scanner.Text()

		if cs.pattern.MatchString(line) {
			result := cs.createResult(path, line, lineNumber, contextBuffer)
			results = append(results, result)
			contextBuffer = nil
		} else if showContext > 0 {
			if len(contextBuffer) >= showContext {
				contextBuffer = contextBuffer[1:]
			}
			contextBuffer = append(contextBuffer, line)
		}
	}

	return results, stats, scanner.Err()
}

func (cs *CodeSearcher) createResult(path, line string, lineNum int, context []string) SearchResult {
	result := SearchResult{
		File:    path,
		Line:    lineNum,
		Content: line,
		IsMatch: true,
	}

	if showContext > 0 {
		var contextLines []string
		for _, ctxLine := range context {
			contextLines = append(contextLines, SearchResult{
				Content: ctxLine,
				IsMatch: false,
			}.formatLine())
		}
		contextLines = append(contextLines, result.formatLine())
		result.Context = contextLines
	}

	return result
}

func (sr SearchResult) Print() {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	if len(sr.Context) > 0 {
		fmt.Printf("\n%s\n", cyan(sr.File))
		for _, ctx := range sr.Context {
			fmt.Println(ctx)
		}
	} else {
		fmt.Printf("\n%s:%s\n%s\n",
			cyan(sr.File),
			yellow(sr.Line),
			red(sr.Content),
		)
	}
}

func (sr SearchResult) formatLine() string {
	lineNumColor := color.New(color.FgYellow).SprintFunc()
	matchColor := color.New(color.FgRed).SprintFunc()
	normalColor := color.New(color.FgWhite).SprintFunc()

	var formatted string
	if sr.IsMatch {
		formatted = matchColor(sr.Content)
	} else {
		formatted = normalColor(sr.Content)
	}

	return fmt.Sprintf("%s %s", lineNumColor(sr.Line), formatted)
}

func printStats(stats SearchStats, duration time.Duration) {
	fmt.Println("\n--- Search Statistics ---")
	fmt.Printf("Execution time:    %s\n", duration.Round(time.Millisecond))
	fmt.Printf("Files scanned:     %d\n", stats.FilesScanned)
	fmt.Printf("Total lines:       %d\n", stats.TotalLines)
	fmt.Printf("Matches found:     %d\n", stats.MatchesFound)
	fmt.Printf("Go routines used:  %d\n", runtime.NumGoroutine())
}
