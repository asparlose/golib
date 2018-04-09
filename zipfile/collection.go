package zipfile

import (
	"archive/zip"
	"context"
)

type FileCollection interface {
	Find(...FileFilter) FileCollection
	Iter(context.Context) <-chan *zip.File
	Slice() []*zip.File
}

type findResult struct {
	parent  FileCollection
	filters []FileFilter
}

func find(collection FileCollection, filters ...FileFilter) FileCollection {
	return &findResult{
		parent:  collection,
		filters: filters,
	}
}

func slice(collection FileCollection) []*zip.File {
	s := make([]*zip.File, 0)
	for file := range collection.Iter(context.Background()) {
		s = append(s, file)
	}
	return s
}

func (c *findResult) Find(filters ...FileFilter) FileCollection {
	return find(c, filters...)
}

func (c *findResult) Slice() []*zip.File {
	return slice(c)
}

func (c *findResult) Iter(ctx context.Context) <-chan *zip.File {
	ch := make(chan *zip.File)

	go func(c *findResult, ch chan<- *zip.File) {
		defer close(ch)
		for file := range c.parent.Iter(ctx) {
			select {
			case <-ctx.Done():
				return
			default:
			}

			for _, filter := range c.filters {
				if filter(file) {
					ch <- file
					break
				}
			}
		}
	}(c, ch)

	return ch
}
