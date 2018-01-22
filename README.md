# PDF Tools

These are a few simple command line tools for manipulating pdf documents.

## pdfcat

`pdfcat` takes a list of pdf documents as command line arguments and outputs the concatenation of all of them into a single pdf document.

## pdfmerge

`pdfmerge` merges two pdfs together. The pdfs should be the exact same length and the second pdf should be in reverse order. The primary function of this tool is to merge two pdfs that come from the same set of pages. My scanner can only use the paper feed for one side at a time and the backside of the pages gets read backwards so the merge tool will merge a pdf document as the following page order.

    1 -> 3 ---|
              |----> 1 -> 2 -> 3 -> 4
    4 -> 2 ---|
