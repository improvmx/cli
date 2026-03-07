package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

var JSONOutput bool

func PrintJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func NewTable(headers ...string) *Table {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	return &Table{writer: w, headers: headers}
}

type Table struct {
	writer  *tabwriter.Writer
	headers []string
	printed bool
}

func (t *Table) AddRow(cols ...string) {
	if !t.printed {
		fmt.Fprintln(t.writer, strings.Join(t.headers, "\t"))
		t.printed = true
	}
	fmt.Fprintln(t.writer, strings.Join(cols, "\t"))
}

func (t *Table) Render() {
	t.writer.Flush()
}

func Success(msg string) {
	fmt.Printf("✓ %s\n", msg)
}

func Error(msg string) {
	fmt.Fprintf(os.Stderr, "✗ %s\n", msg)
}
