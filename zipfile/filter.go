package zipfile

import (
	"archive/zip"
	"path/filepath"
	"strings"
)

type FileFilter func(*zip.File) bool

func And(filters ...FileFilter) FileFilter {
	return func(file *zip.File) bool {
		for _, filter := range filters {
			if !filter(file) {
				return false
			}
		}
		return true
	}
}

func Or(filters ...FileFilter) FileFilter {
	return func(file *zip.File) bool {
		for _, filter := range filters {
			if filter(file) {
				return true
			}
		}
		return false
	}
}

func FullName(fullname string) FileFilter {
	return func(file *zip.File) bool {
		return file.Name == fullname
	}
}

func File() FileFilter {
	return func(file *zip.File) bool {
		return !strings.HasSuffix(file.Name, "/")
	}
}

func Directory() FileFilter {
	return func(file *zip.File) bool {
		return strings.HasSuffix(file.Name, "/")
	}
}

func Name(name string) FileFilter {
	return func(file *zip.File) bool {
		return filepath.Base(file.Name) == name
	}
}

func ChildOf(name string) FileFilter {
	return func(file *zip.File) bool {
		return filepath.Dir(file.Name) == name
	}
}

func split(path string) []string {
	p := path
	a := make([]string, 0)
	for {
		d, f := filepath.Split(p)
		a = append([]string{f}, a...)
		if d == "" {
			break
		}
		p = d[:len(d)-1]
	}
	return a
}

func DescendantsOf(name string) FileFilter {
	return func(file *zip.File) bool {
		sf := split(file.Name)
		sp := split(name)

		if len(sf) < len(sp) {
			return false
		}

		for i := 0; i < len(sp); i++ {
			if sf[i] != sp[i] {
				return false
			}
		}
		return true
	}
}

func Match(pattern string) FileFilter {
	return func(file *zip.File) bool {
		r, _ := filepath.Match(pattern, file.Name)
		return r
	}
}
