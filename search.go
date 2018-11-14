package main

import "sync"
import "regexp"
import "fmt"
import "os/exec"
import "os"
import "bufio"
import "strings"
import "strconv"
import p "path"
import "path/filepath"

type searchQuery struct {
	original string
	regex    *regexp.Regexp
}

func parseQuery(input string) *searchQuery {
	q := &searchQuery{
		original: input,
	}
	var err error
	q.regex, err = regexp.Compile(input)
	if err != nil {
		fmt.Println(err)
	}
	return q
}

func searchNotes(q *searchQuery) []match {
	var wg sync.WaitGroup
	matchset := newMatchset()
	matchset.Add(filepath.Join(*notesPath, q.original+*extension))

	results := make(chan result)
	wg.Add(2)
	go searchByPath(q, results, &wg)
	go searchByContent(q, results, &wg)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		matchset.AddResult(r)
	}

	return matchset.SortedMatches()
}

func searchByPath(query *searchQuery, results chan result, wg *sync.WaitGroup) {
	paths := make(chan string)
	go walkTreeForFiles(*notesPath, paths)
	for path := range paths {
		if pathMatchesQuery(query, p.Base(path)) {
			results <- pathResult{path: path, basenameMatch: true}
			continue
		}
		if pathMatchesQuery(query, path) {
			results <- pathResult{path: path}
		}
	}
	wg.Done()
}

func pathMatchesQuery(query *searchQuery, path string) bool {
	if query.regex != nil {
		return query.regex.MatchString(path)
	}
	return strings.Contains(path, query.original)
}

func walkTreeForFiles(dir string, results chan string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "failure accessing path %q: %v\n", path, err)
			return err
		}
		if info.Mode().IsRegular() {
			results <- path
		}
		return nil
	})
	handleError(err)
	close(results)
}

func searchByContent(query *searchQuery, results chan result, wg *sync.WaitGroup) {
	cmd := exec.Command("ag", "--nogroup", query.original, *notesPath)
	stdout, err := cmd.StdoutPipe()
	handleError(err)
	err = cmd.Start()
	handleError(err)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 3)
		lineno, _ := strconv.Atoi(parts[1])
		// fmt.Printf("%#v\n", parts)
		results <- contentResult{
			path:    parts[0],
			lineno:  lineno,
			snippet: parts[2],
		}
	}

	err = cmd.Wait()
	// handleError(err)
	wg.Done()
}
