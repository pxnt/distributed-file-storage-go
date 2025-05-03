package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

const DEFAULT_ROOT_FOLDER = "dfs_test_output"

type Location struct {
	Path     string
	Filename string
}

func (l Location) FullPath() string {
	return path.Join(l.Path, l.Filename)
}

func (l Location) FirstPathName() string {
	tokens := strings.Split(l.Path, "/")

	if len(tokens) == 0 {
		return ""
	}

	return tokens[0]
}

type PathTransformFunc func(string) Location

func DefaultPathTransform(key string) Location {
	return Location{
		Path:     key,
		Filename: key,
	}
}

func CASPathTransformFunc(key string) Location {
	pathHash := sha1.Sum([]byte(key))
	pathHashString := hex.EncodeToString(pathHash[:])

	PATH_TOKEN_SIZE := 5
	TOTAL_PATH_TOKENS := len(pathHashString) / PATH_TOKEN_SIZE

	pathTokens := make([]string, TOTAL_PATH_TOKENS)

	for i := range TOTAL_PATH_TOKENS {
		from, to := i*PATH_TOKEN_SIZE, (i+1)*PATH_TOKEN_SIZE

		pathTokens[i] = pathHashString[from:to]
	}

	return Location{
		Path:     strings.Join(pathTokens, "/"),
		Filename: pathHashString,
	}
}

type StoreOpts struct {
	Root              string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	location := s.PathTransformFunc(key)
	dirWithRoot := path.Join(s.Root, location.Path)

	if err := os.MkdirAll(dirWithRoot, os.ModePerm); err != nil {
		return err
	}

	fullPath := location.FullPath()
	fullPathWithRoot := path.Join(s.Root, fullPath)

	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("Wrote %d bytes to %s\n", n, fullPathWithRoot)

	return nil
}

func (s *Store) ReadStream(key string) (io.Reader, error) {
	location := s.PathTransformFunc(key)
	fullPathWithRoot := path.Join(s.Root, location.FullPath())

	f, err := os.Open(fullPathWithRoot)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) Delete(key string) error {
	location := s.PathTransformFunc(key)

	defer func() {
		log.Printf("Deleted %s", key)
	}()

	return os.RemoveAll(location.FirstPathName())
}

func (s *Store) Clear() error {
	defer func() {
		log.Printf("Deleted Root %s", s.Root)
	}()

	return os.RemoveAll(s.Root)
}

func (s *Store) Has(key string) bool {
	location := s.PathTransformFunc(key)
	fullPathWithRoot := path.Join(s.Root, location.FullPath())

	_, err := os.Stat(fullPathWithRoot)

	return errors.Is(err, os.ErrNotExist)
}
