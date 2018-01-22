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
	force := flag.BoolP("force", "f", false, "Force the pdf files to merge")
	remove := flag.BoolP("rm", "", false, "Remove the original files")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Error: Must have exactly two arguments.\n")
		return 1
	}

	sideA, err := pdf.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot open PDF document %s: %s.\n", args[0], err)
		return 1
	}
	defer sideA.Close()
	sideALen, _ := sideA.GetNumPages()

	sideB, err := pdf.Open(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot open PDF document %s: %s.\n", args[1], err)
		return 1
	}
	defer sideB.Close()
	sideBLen, _ := sideB.GetNumPages()

	if !*force && sideALen != sideBLen {
		fmt.Fprintf(os.Stderr, "Error: PDF documents do not have the same number of pages. This will potentially cause the document to merge incorrectly. Use --force to perform the merge anyway.\n")
		return 1
	}

	i, j := 1, sideBLen
	concat := creator.New()
	for i <= sideALen || j > 0 {
		if i <= sideALen {
			page, err := sideA.GetPage(i)
			if err == nil {
				concat.AddPage(page)
			} else {
				fmt.Fprintf(os.Stderr, "Skipping page %d for side A: %s", i, err)
			}
			i++
		}

		if j > 0 {
			page, err := sideB.GetPage(j)
			if err == nil {
				concat.AddPage(page)
			} else {
				fmt.Fprintf(os.Stderr, "Skipping page %d for side B: %s", j, err)
			}
			j--
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

	if *remove {
		if err := os.Remove(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not remove file %s: %s.\n", args[0], err)
		}
		if err := os.Remove(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not remove file %s: %s.\n", args[1], err)
		}
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
