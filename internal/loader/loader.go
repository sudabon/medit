package loader

import (
	"errors"
	"fmt"
	"os"
)

// Document holds the loaded Markdown file state.
type Document struct {
	Path    string
	Content string
	IsNew   bool
}

// Load reads path as UTF-8. A nonexistent path opens as a new empty file
// without creating it on disk. Directories and unreadable paths return an error.
func Load(path string) (*Document, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Document{Path: path, Content: "", IsNew: true}, nil
		}
		return nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}

	return &Document{Path: path, Content: string(data), IsNew: false}, nil
}

// Reload re-reads the document from disk. On failure the previous Content is
// left unchanged and an error is returned. If a previously new path now exists,
// IsNew is cleared.
func (d *Document) Reload() error {
	info, err := os.Stat(d.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if d.IsNew {
				d.Content = ""
				return nil
			}
			return fmt.Errorf("file not found: %s", d.Path)
		}
		return fmt.Errorf("cannot open %s: %w", d.Path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory", d.Path)
	}

	data, err := os.ReadFile(d.Path)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", d.Path, err)
	}

	d.Content = string(data)
	d.IsNew = false
	return nil
}
