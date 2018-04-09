package zipfile

import (
	"archive/zip"
	"context"
)

type ReaderCloser interface {
	Reader
	Close()
}

type Reader interface {
	FileCollection
}

type reader struct {
	file *zip.Reader
}

type readerCloser struct {
	file *zip.ReadCloser
}

func OpenReader(file string) (ReaderCloser, error) {
	z, err := zip.OpenReader(file)
	if err != nil {
		return nil, err
	}

	r := &readerCloser{
		file: z,
	}

	return r, nil
}

func NewReader(file *zip.Reader) Reader {
	return &reader{
		file: file,
	}
}

func (c *reader) Find(filters ...FileFilter) FileCollection {
	return find(c, filters...)
}

func (c *reader) Slice() []*zip.File {
	return slice(c)
}

func (c *reader) Iter(ctx context.Context) <-chan *zip.File {
	ch := make(chan *zip.File)

	go func(c *reader, ch chan<- *zip.File) {
		defer close(ch)
		for _, file := range c.file.File {
			select {
			case <-ctx.Done():
				return
			default:
			}
			ch <- file
		}
	}(c, ch)

	return ch
}

func (c *readerCloser) Find(filters ...FileFilter) FileCollection {
	return find(c, filters...)
}

func (c *readerCloser) Slice() []*zip.File {
	return slice(c)
}

func (c *readerCloser) Iter(ctx context.Context) <-chan *zip.File {
	ch := make(chan *zip.File)

	go func(c *readerCloser, ch chan<- *zip.File) {
		defer close(ch)
		for _, file := range c.file.File {
			select {
			case <-ctx.Done():
				return
			default:
			}
			ch <- file
		}
	}(c, ch)

	return ch
}

func (c *readerCloser) Close() {
	c.file.Close()
}
