package fileio

import (
	"compress/gzip"
	"io"
	"os"
	"strings"

	"github.com/golang/glog"

	"github.com/jamessynge/go_io/gzipio"
	"github.com/jamessynge/go_io/pumpio"
)

// OpenReadFile opens the file filePath for reading, optionally decompressing
// the file it its name ends with .gz (i.e. if it is apparently a gzip
// compressed file). Returns an io.ReadCloser providing the file's contents
// if successful, and an error otherwise (e.g. permissions failure or file
// name ends with .gz but the file doesn't start with a valid gzip header).
// TODO(jamessynge): Add options for controlling decompression (e.g. for gzip
// decompression, or choose another decompressor).
func OpenReadFile(filePath string) (io.ReadCloser, error) {
	// Open the file for reading.
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(filePath, ".gz") {
		glog.V(1).Infof("Opened file for reading: %s", filePath)
		return f, nil
	}
	grc, err := gzipio.NewReadCloser(f)
	if err != nil {
		glog.Warningf("Unable to create gzip Reader for %q\nError: %v", filePath, err)
		f.Close()
		return nil, err
	} else {
		glog.V(1).Infof("Opened compressed file for reading: %s", filePath)
		return grc, nil
	}
}

// OpenReadFileAndPump is like OpenReadFile, but uses pumpio to start
// go routines for reading the file from the OS, and decompressing it if
// the name ends with .gz, prior to the caller reading from the returned
// io.ReadCloser.
// TODO(jamessynge): Add options for blockSize and blockCount.
func OpenReadFileAndPump(filePath string) (io.ReadCloser, error) {
	// Open the file for reading.
	var rc io.ReadCloser
	rc, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	glog.V(1).Infof("Opened file for reading: %s", filePath)
	blockSize, blockCount := 4096, 16
	if strings.HasSuffix(filePath, ".gz") {
		glog.V(1).Infof("Opened compressed file for reading: %s", filePath)
		// Use half the blocks for the file pump, half for the decompressed pump
		// (not sure what the "right" ratio is, if such a thing exists).
		blockCount /= 2
		rc = pumpio.NewReadCloserPump(rc, blockSize, blockCount)
		gr, err2 := gzip.NewReader(rc)
		if err2 != nil {
			rc.Close()
			return nil, err2
		}
		rc, err = gzipio.NewReadCloser(rc)
		if err != nil {
			return nil, err
		}
	} else {
		glog.V(1).Infof("Opened file for reading: %s", filePath)
	}
	rc = pumpio.NewReadCloserPump(rc, blockSize, blockCount)
	return rc, nil
}

// Generalization of ioutil.WriteFile to handle data that is split into
// fragments (i.e. is a slice of byte slices).
func WriteFile(fp string, fragments [][]byte, perm os.FileMode) (err error) {
	var f *os.File
	if f, err = os.OpenFile(fp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm); err != nil {
		return
	}
	for _, fragment := range fragments {
		if len(fragment) == 0 {
			continue
		}
		var n int
		n, err = f.Write(fragment)
		if err != nil {
			break
		}
		if n != len(fragment) {
			err = io.ErrShortWrite
			break
		}
	}
	err2 := f.Close()
	if err == nil {
		err = err2
	}
	return
}
