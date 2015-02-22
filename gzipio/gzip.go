package util

import (
	"compress/gzip"
	"io"
)

type gzipReadCloser struct {
	src io.ReadCloser
	gr  *gzip.Reader
}

func (grc *gzipReadCloser) Read(p []byte) (n int, err error) {
	return grc.gr.Read(p)
}

func (grc *gzipReadCloser) Close() (err error) {
	defer func() {
		err2 := grc.src.Close()
		if err == nil {
			err = err2
		}
	}()
	err = grc.gr.Close()
	return
}

// ReadCloser embeds a gzip.Reader, changing the Close method so that it also
// closes the underlying io.ReadCloser.
type ReadCloser struct {
	*gzip.Reader
	src io.ReadCloser
}

func (p *ReadCloser) Close() error {
	err := p.Reader.Close()
	err2 := p.src.Close()
	if err == nil {
		err = err2
	}
	return err
}

// NewReadCloser creates a new ReadCloser (embeds a gzip.Reader). Returns an
// error if gzip.NewReader fails (e.g. if src does not start with a valid
// gzip header).
func NewReadCloser(src io.ReadCloser) (*ReadCloser, error) {
	r, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}
	rc := &ReadCloser{
		Reader: r,
		src: src,
	}
	return rc, nil
}
