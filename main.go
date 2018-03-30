package main

import "os"
import "fmt"

import "gopkg.in/alecthomas/kingpin.v2"

// 0. if stdout is TTY, we're in Interactive mode, else Batch mode
// 1. accept filter query as args for "batch mode" (or interactively)
// 2. find matches by path
// 3. find matches by content (probably just use ag)
// 4. if no matches, include a result to create a new file using the filter query with the extension appended
// 5. output list of results in requested format (flags)
// 6. if interactive, ENTER will open selected file using $EDITOR

var (
	app         = kingpin.New("nvcli", "find notes and edit selected, interactively")
	interactive = app.Command("interactive", "Run interactively").Default()
	batch       = app.Command("batch", "Run a single query and output matches in FORMAT specified")
	format      = batch.Flag("format", "Set the output format").Required().String()
	query       = batch.Arg("query", "Filter string for finding notes").Required().String()
	notesPath   = app.Flag("notes-path", "Directory your notes live in").Required().ExistingDir()
)

func main() {
	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	switch command {
	case interactive.FullCommand():
		// interactive
		fmt.Println("interactive not yet implemented")
	case batch.FullCommand():
		// batch
		q := parseQuery(*query)
		results := searchNotes(q)
		serialize, err := getSerializer(*format)
		handleError(err)
		fmt.Println(serialize(results))
	}
}

func handleError(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}
