package main

import (
	"fmt"
	"os"

	"github.com/jsternberg/pdftools/pdf"
	flag "github.com/spf13/pflag"
	"github.com/unidoc/unidoc/pdf/creator"
)

func realMain() int {
	output := flag.StringP("output", "o", "<stdout>", "The output pdf name")
	flag.Parse()

	args := flag.Args()
	concat := creator.New()
	for _, fpath := range args {
		if err := func() error {
			f, err := pdf.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()

			pages, err := f.GetNumPages()
			if err != nil {
				return err
			}

			for i := 1; i <= pages; i++ {
				page, err := f.GetPage(i)
				if err != nil {
					return err
				}
				concat.AddPage(page)
			}
			return nil
		}(); err != nil {
			fmt.Fprintf(os.Stderr, "Error while processing %s: %s.\n", fpath, err)
			return 1
		}
	}

	switch *output {
	case "<stdout>":
		if err := concat.Write(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not write PDF document to stdout.\n")
			return 1
		}
	case "<stderr>":
		if err := concat.Write(os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not write PDF document to stderr.\n")
			return 1
		}
	default:
		if err := concat.WriteToFile(*output); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not write output file %s: %s.\n", *output, err)
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
