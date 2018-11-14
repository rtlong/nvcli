package main

import "fmt"
import "encoding/json"
import "path/filepath"

type serializer func([]match) string

const newFileDescription = "NEW"

func getSerializer(format string) (serializer, error) {
	switch format {
	case "text":
		return SerializeText, nil
	case "grep":
		return SerializeGrep, nil
	case "alfred-json":
		return SerializeAlfredJSON, nil
	default:
		return nil, fmt.Errorf("Format %s is not valid", format)
	}
}

// SerializeText formats results as simple newline-separated list of filenames
func SerializeText(matches []match) string {
	str := ""

	for _, m := range matches {
		var snip string
		if len(m.Results) == 0 {
			snip = newFileDescription
		} else {
			snip = m.Snippet()
		}
		str = str + fmt.Sprintf("%s: %s\n", m.Path, snip)
	}

	return str
}

// SerializeGrep formats results as a format like grep's output (but this isn't exactly appropriate)
func SerializeGrep(matches []match) string {
	str := ""

	for _, m := range matches {
		if len(m.Results) == 0 {
			str = str + fmt.Sprintf("%s:%d:%f:%s\n", m.Path, 0, m.Score, newFileDescription)
		} else {
			for _, r := range m.Results {
				str = str + fmt.Sprintf("%s:%d:%f:%s\n", m.Path, r.LineNo(), m.Score, r.Snippet())
			}
		}
	}

	return str
}

// SerializeAlfredJSON formats results as JSON for Alfred.app workflows on OSX
/* example of output we want:
{"items": [
    {
        "type": "file",
        "title": "relative path from notes dir",
        "subtitle": "snippet",
        "arg": "full path",
    }
]} */
func SerializeAlfredJSON(matches []match) string {
	alfred := alfredResult{
		Items: []alfredItem{},
	}

	for _, m := range matches {
		var snip string
		if len(m.Results) == 0 {
			snip = newFileDescription
		} else {
			snip = m.Snippet()
		}

		relPath, err := filepath.Rel(*notesPath, m.Path)
		if err != nil {
			relPath = m.Path
		}
		alfred.Items = append(alfred.Items, alfredItem{
			Type:     "file",
			Title:    relPath,
			Subtitle: snip,
			Arg:      m.Path,
		})
	}

	bytes, err := json.Marshal(alfred)
	handleError(err)
	return string(bytes)
}

type alfredResult struct {
	Items []alfredItem `json:"items"`
}

type alfredItem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg"`
}
