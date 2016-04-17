package main

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

type store struct {
	mut      sync.RWMutex
	fileName string
}

func (s store) withRWFile(f func(*os.File) error) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	file, err := os.OpenFile(s.fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return f(file)
}

func (s store) withRoFile(f func(*os.File) error) error {
	s.mut.RLock()
	defer s.mut.RUnlock()

	file, err := os.OpenFile(s.fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return f(file)
}

func loadPosts(f *os.File) ([]post, error) {
	var posts []post
	err := json.NewDecoder(f).Decode(&posts)
	switch err {
	case nil:
		return posts, nil
	case io.EOF:
		return make([]post, 0), nil
	default:
		return nil, err
	}
}

func putPosts(f *os.File, posts []post) error {
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(&posts)
}

func newStore(dbFile string) store {
	return store{
		fileName: dbFile,
	}
}
