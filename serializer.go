package main

import "fmt"
import "encoding/json"
import "path/filepath"

type serializer func([]match) string

func getSerializer(format string) (serializer, error) {
	switch format {
	case "text":
		return SerializeText, nil
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
		snip := m.Results[0].Snippet()
		str = str + fmt.Sprintf("%02d %s: %s\n", m.Score, m.Path, snip)
	}

	return str
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
		if m.Results != nil {
			snip = m.Results[0].Snippet()
		} else {
			snip = "Create new file"
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
