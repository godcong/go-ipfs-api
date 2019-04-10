package shell

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/ipfs/go-ipfs-files"
)

// Object ...
type Object struct {
	Hash string
	Name string
	Size uint64
}

// UnmarshalJSON redraw json unmarshal
func (r *Object) UnmarshalJSON(b []byte) error {
	out := struct {
		Hash string
		Name string
		Size string
	}{}
	e := json.Unmarshal(b, &out)
	if e != nil {
		return e
	}
	r.Size, _ = strconv.ParseUint(out.Size, 10, 64)
	r.Hash = out.Hash
	_, r.Name = path.Split(out.Name)
	return nil
}

// AddOpts ...
type AddOpts = func(*RequestBuilder) error

// OnlyHash ...
func OnlyHash(enabled bool) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("only-hash", enabled)
		return nil
	}
}

// Pin ...
func Pin(enabled bool) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("pin", enabled)
		return nil
	}
}

// Progress ...
func Progress(enabled bool) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("progress", enabled)
		return nil
	}
}

// RawLeaves ...
func RawLeaves(enabled bool) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("raw-leaves", enabled)
		return nil
	}
}

// Add ...
func (s *Shell) Add(r io.Reader, options ...AddOpts) (*Object, error) {
	fr := files.NewReaderFile(r)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	fileReader := files.NewMultiFileReader(slf, true)

	var out Object
	rb := s.Request("add")
	for _, option := range options {
		option(rb)
	}
	err := rb.Body(fileReader).Exec(context.Background(), &out)
	if err != nil {
		return &Object{}, err
	}
	return &out, nil
}

// AddNoPin adds a file to ipfs without pinning it
// Deprecated: Use Add() with option functions instead
func (s *Shell) AddNoPin(r io.Reader) (*Object, error) {
	return s.Add(r, Pin(false))
}

// AddWithOpts adds a file to ipfs with some additional options
// Deprecated: Use Add() with option functions instead
func (s *Shell) AddWithOpts(r io.Reader, pin bool, rawLeaves bool) (*Object, error) {
	return s.Add(r, Pin(pin), RawLeaves(rawLeaves))
}

// AddLink ...
func (s *Shell) AddLink(target string) (*Object, error) {
	link := files.NewLinkFile(target, nil)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", link)})
	reader := files.NewMultiFileReader(slf, true)

	var out Object
	err := s.Request("add").Body(reader).Exec(context.Background(), &out)
	if err != nil {
		return &Object{}, err
	}
	return &out, nil
}

// AddDir adds a directory recursively with all of the files under it
func (s *Shell) AddDir(dir string) ([]*Object, error) {
	stat, err := os.Lstat(dir)
	if err != nil {
		return nil, err
	}

	sf, err := files.NewSerialFile(dir, false, stat)
	if err != nil {
		return nil, err
	}
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry(path.Base(dir), sf)})
	reader := files.NewMultiFileReader(slf, true)

	resp, err := s.Request("add").
		Option("recursive", true).
		Body(reader).
		Send(context.Background())
	if err != nil {
		return nil, nil
	}

	defer resp.Close()

	if resp.Error != nil {
		return nil, resp.Error
	}

	dec := json.NewDecoder(resp.Output)
	var final []*Object
	for {
		var out Object
		err = dec.Decode(&out)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		final = append(final, &out)
	}

	if final == nil {
		return nil, errors.New("no results received")
	}

	return final, nil
}
