package vfsutil

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
)

func init() {
	gob.Register(&GobFS{})
}

type GobFS struct {
	vfs.Filesystem
}

// NewGob returns a vfs.Filesystem that can be encoded and decoded using encoding/gob.
// It assumes that fs consists only of directories and regular files.
func NewGob(fs vfs.Filesystem) *GobFS {
	return &GobFS{fs}
}

func (fs *GobFS) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	err := Walk(fs, "/", func(path string, info os.FileInfo, err error) error {
		hdr := &tar.Header{
			Name:    path,
			ModTime: info.ModTime(),
			Mode:    int64(info.Mode().Perm()),
		}

		if info.IsDir() {
			hdr.Typeflag = tar.TypeDir
		} else {
			hdr.Typeflag = tar.TypeReg
			hdr.Size = info.Size()
		}

		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := vfs.Open(fs, path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})

	if err != nil {
		return nil, err
	}

	err = tw.Close()
	if err != nil {
		return nil, err
	}

	err = gw.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (fs *GobFS) GobDecode(b []byte) error {
	gr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return err
	}
	tr := tar.NewReader(gr)

	fs.Filesystem = memfs.Create()

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			err = vfs.MkdirAll(fs, hdr.Name, 0700)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := fs.OpenFile(hdr.Name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, hdr.FileInfo().Mode()|0600)
			if err != nil {
				return err
			}
			_, err = io.Copy(f, tr)
			f.Close()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown header type %d", hdr.Typeflag)
		}
	}
}
