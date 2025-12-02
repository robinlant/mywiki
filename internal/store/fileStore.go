package store

import (
	"os"
	"path"
)

type FileStore struct {
	Dir string
}

func NewFileStore(dir string) *FileStore {
	return &FileStore{
		Dir: dir,
	}
}

func (f *FileStore) SavePage(p *Page) error {
	filename := p.Title + ".txt"
	return os.WriteFile(path.Join(f.Dir, filename), p.Body, 0600)
}

func (f *FileStore) LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(path.Join(f.Dir, filename))
	if err != nil {
		return nil, err
	}
	return &Page{
		Title: title,
		Body:  body,
	}, nil
}
