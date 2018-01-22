package pdf

import (
	"io"
	"os"

	"github.com/unidoc/unidoc/pdf/model"
)

type Reader struct {
	*model.PdfReader
	io.Closer
}

func Open(fpath string) (*Reader, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	pdf, err := model.NewPdfReader(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	return &Reader{
		PdfReader: pdf,
		Closer:    f,
	}, nil
}
