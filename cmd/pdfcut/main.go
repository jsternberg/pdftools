package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jsternberg/pdftools/pdf"
	flag "github.com/spf13/pflag"
	"github.com/unidoc/unidoc/pdf/creator"
)

func realMain() int {
	pages := flag.StringP("pages", "p", "", "The page range or single page to cut")
	output := flag.StringP("output", "o", "<stdout>", "The output pdf name")
	flag.Parse()

	if *pages == "" {
		fmt.Fprint(os.Stderr, "Error: A page or page range must be specified.\n")
		return 1
	}

	var start, end int64
	if strings.Contains(*pages, ":") {
		parts := strings.SplitN(*pages, ":", 2)

		var err error
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not parse start page number: %s.\n", err)
			return 1
		}

		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not parse end page number: %s.\n", err)
			return 1
		}
	} else {
		n, err := strconv.ParseInt(*pages, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not parse page number: %s.\n", err)
			return 1
		}
		start, end = n, n
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Must specify exactly one file as an argument.\n")
		return 1
	}
	fpath := args[0]

	f, err := pdf.Open(fpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not open file: %s.\n", err)
		return 1
	}

	newpdf := creator.New()
	for i := start; i <= end; i++ {
		page, err := f.GetPage(int(i))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not retrieve page %d: %s.\n", i, err)
			return 1
		}
		newpdf.AddPage(page)
	}

	switch *output {
	case "<stdout>":
		if err := newpdf.Write(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not write PDF document to stdout.\n")
			return 1
		}
	case "<stderr>":
		if err := newpdf.Write(os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not write PDF document to stderr.\n")
			return 1
		}
	default:
		if err := newpdf.WriteToFile(*output); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not write output file %s: %s.\n", *output, err)
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
