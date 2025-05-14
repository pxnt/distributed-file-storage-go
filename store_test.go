package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	TEST_PATH := "test_folder_1"
	location := CASPathTransformFunc(TEST_PATH)

	expectedPath := "753ba/205f2/11917/3e905/f67d5/4ea62/37418/3ca04"
	expectedFilename := "753ba205f2119173e905f67d54ea62374183ca04"

	if location.Path != expectedPath {
		t.Errorf("Incorrect Path: have %s want %s", location.Path, expectedPath)
	}

	if location.Filename != expectedFilename {
		t.Errorf("Incorrect Filename: have %s want %s", location.Filename, expectedFilename)
	}
}

func TestStore(t *testing.T) {
	storeOpts := StoreOpts{
		Root:              "meoww",
		PathTransformFunc: CASPathTransformFunc,
	}

	TEST_FOLDER := "folder_1"
	TEST_FILE_CONTENT := "Hello, World!"

	store := NewStore(storeOpts)
	defer teardown(t, store)

	type A struct {
		Q int
		S string
		E []int
	}

	data := bytes.NewReader([]byte(TEST_FILE_CONTENT))

	err := store.WriteStream(TEST_FOLDER, data)

	if err != nil {
		t.Fatalf("Failed to write stream: %v", err)
	}

	r, err := store.ReadStream(TEST_FOLDER)

	if err != nil {
		t.Fatalf("Failed to read stream: %v", err)
	}

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, r)

	if err != nil {
		t.Fatalf("Failed to copy stream: %v", err)
	}

	if buf.String() != TEST_FILE_CONTENT {
		t.Fatalf("Incorrect data: have %s want %s", buf.String(), TEST_FILE_CONTENT)
	}
}

func teardown(t *testing.T, store *Store) {
	err := store.Clear()

	if err != nil {
		t.Fatalf("Failed to clear store: %v", err)
	}
}
