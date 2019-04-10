package shell

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/cheekybits/is"
	"github.com/godcong/go-ipfs-restapi/options"
)

const (
	examplesHash = "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
	shellURL     = "localhost:5001"
)

// TestAdd ...
func TestAdd(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	mhash, err := s.Add(bytes.NewBufferString("Hello IPFS Shell tests"))
	is.Nil(err)
	is.Equal(mhash, "QmUfZ9rAdhV5ioBzXKdUTh2ZNsz9bzbkaLVyQ8uc8pj21F")
}

// TestRedirect ...
func TestRedirect(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	err := s.
		Request("/version").
		Exec(context.Background(), nil)
	is.NotNil(err)
	is.True(strings.Contains(err.Error(), "unexpected redirect"))
}

// TestAddWithCat ...
func TestAddWithCat(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)
	s.SetTimeout(1 * time.Second)

	rand := randString(32)

	mhash, err := s.Add(bytes.NewBufferString(rand))
	is.Nil(err)

	reader, err := s.Cat(mhash)
	is.Nil(err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	catRand := buf.String()

	is.Equal(rand, catRand)
}

// TestAddOnlyHash ...
func TestAddOnlyHash(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)
	s.SetTimeout(1 * time.Second)

	rand := randString(32)

	mhash, err := s.Add(bytes.NewBufferString(rand), OnlyHash(true))
	is.Nil(err)

	_, err = s.Cat(mhash)
	is.Err(err) // we expect an http timeout error because `cat` won't find the `rand` string
}

// TestAddNoPin ...
func TestAddNoPin(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	h, err := s.Add(bytes.NewBufferString(randString(32)), Pin(false))
	is.Nil(err)

	pins, err := s.Pins()
	is.Nil(err)

	_, ok := pins[h]
	is.False(ok)
}

// TestAddNoPinDeprecated ...
func TestAddNoPinDeprecated(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	h, err := s.AddNoPin(bytes.NewBufferString(randString(32)))
	is.Nil(err)

	pins, err := s.Pins()
	is.Nil(err)

	_, ok := pins[h]
	is.False(ok)
}

// TestAddDir ...
func TestAddDir(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	cid, err := s.AddDir("./testdata")
	is.Nil(err)
	is.Equal(cid, "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv")
}

// TestLocalShell ...
func TestLocalShell(t *testing.T) {
	is := is.New(t)
	s := NewLocalShell()
	is.NotNil(s)

	mhash, err := s.Add(bytes.NewBufferString("Hello IPFS Shell tests"))
	is.Nil(err)
	is.Equal(mhash, "QmUfZ9rAdhV5ioBzXKdUTh2ZNsz9bzbkaLVyQ8uc8pj21F")
}

// TestCat ...
func TestCat(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	rc, err := s.Cat(fmt.Sprintf("/ipfs/%s/readme", examplesHash))
	is.Nil(err)

	md5 := md5.New()
	_, err = io.Copy(md5, rc)
	is.Nil(err)
	is.Equal(fmt.Sprintf("%x", md5.Sum(nil)), "3fdcaad186e79983a6920b4c7eeda949")
}

// TestList ...
func TestList(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	list, err := s.List(fmt.Sprintf("/ipfs/%s", examplesHash))
	is.Nil(err)

	is.Equal(len(list), 7)

	// TODO: document difference in size between 'ipfs ls' and 'ipfs file ls -v'. additional object encoding in data block?
	expected := map[string]LsLink{
		"about":          {Type: TFile, Hash: "QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V", Name: "about", Size: 1677},
		"contact":        {Type: TFile, Hash: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", Name: "contact", Size: 189},
		"help":           {Type: TFile, Hash: "QmY5heUM5qgRubMDD1og9fhCPA6QdkMp3QCwd4s7gJsyE7", Name: "help", Size: 311},
		"ping":           {Type: TFile, Hash: "QmejvEPop4D7YUadeGqYWmZxHhLc4JBUCzJJHWMzdcMe2y", Name: "ping", Size: 4},
		"quick-start":    {Type: TFile, Hash: "QmXgqKTbzdh83pQtKFb19SpMCpDDcKR2ujqk3pKph9aCNF", Name: "quick-start", Size: 1681},
		"readme":         {Type: TFile, Hash: "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB", Name: "readme", Size: 1091},
		"security-notes": {Type: TFile, Hash: "QmQ5vhrL7uv6tuoN9KeVBwd4PwfQkXdVVmDLUZuTNxqgvm", Name: "security-notes", Size: 1162},
	}
	for _, l := range list {
		el, ok := expected[l.Name]
		is.True(ok)
		is.NotNil(el)
		is.Equal(*l, el)
	}
}

// TestFileList ...
func TestFileList(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	list, err := s.FileList(fmt.Sprintf("/ipfs/%s", examplesHash))
	is.Nil(err)

	is.Equal(list.Type, "Directory")
	is.Equal(list.Size, 0)
	is.Equal(len(list.Links), 7)

	// TODO: document difference in sice betwen 'ipfs ls' and 'ipfs file ls -v'. additional object encoding in data block?
	expected := map[string]UnixLsLink{
		"about":          {Type: "File", Hash: "QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V", Name: "about", Size: 1677},
		"contact":        {Type: "File", Hash: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", Name: "contact", Size: 189},
		"help":           {Type: "File", Hash: "QmY5heUM5qgRubMDD1og9fhCPA6QdkMp3QCwd4s7gJsyE7", Name: "help", Size: 311},
		"ping":           {Type: "File", Hash: "QmejvEPop4D7YUadeGqYWmZxHhLc4JBUCzJJHWMzdcMe2y", Name: "ping", Size: 4},
		"quick-start":    {Type: "File", Hash: "QmXgqKTbzdh83pQtKFb19SpMCpDDcKR2ujqk3pKph9aCNF", Name: "quick-start", Size: 1681},
		"readme":         {Type: "File", Hash: "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB", Name: "readme", Size: 1091},
		"security-notes": {Type: "File", Hash: "QmQ5vhrL7uv6tuoN9KeVBwd4PwfQkXdVVmDLUZuTNxqgvm", Name: "security-notes", Size: 1162},
	}
	for _, l := range list.Links {
		el, ok := expected[l.Name]
		is.True(ok)
		is.NotNil(el)
		is.Equal(*l, el)
	}
}

// TestPins ...
func TestPins(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	// Add a thing, which pins it by default
	h, err := s.Add(bytes.NewBufferString("go-ipfs-api pins test 9F3D1F30-D12A-4024-9477-8F0C8E4B3A63"))
	is.Nil(err)

	pins, err := s.Pins()
	is.Nil(err)

	_, ok := pins[h]
	is.True(ok)

	err = s.Unpin(h)
	is.Nil(err)

	pins, err = s.Pins()
	is.Nil(err)

	_, ok = pins[h]
	is.False(ok)

	err = s.Pin(h)
	is.Nil(err)

	pins, err = s.Pins()
	is.Nil(err)

	info, ok := pins[h]
	is.True(ok)
	is.Equal(info.Type, RecursivePin)
}

// TestPatch_rmLink ...
func TestPatch_rmLink(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)
	newRoot, err := s.Patch(examplesHash, "rm-link", "about")
	is.Nil(err)
	is.Equal(newRoot, "QmPmCJpciopaZnKcwymfQyRAEjXReR6UL2rdSfEscZfzcp")
}

// TestPatchLink ...
func TestPatchLink(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)
	newRoot, err := s.PatchLink(examplesHash, "about", "QmUXTtySmd7LD4p6RG6rZW6RuUuPZXTtNMmRQ6DSQo3aMw", true)
	is.Nil(err)
	is.Equal(newRoot, "QmVfe7gesXf4t9JzWePqqib8QSifC1ypRBGeJHitSnF7fA")
	newRoot, err = s.PatchLink(examplesHash, "about", "QmUXTtySmd7LD4p6RG6rZW6RuUuPZXTtNMmRQ6DSQo3aMw", false)
	is.Nil(err)
	is.Equal(newRoot, "QmVfe7gesXf4t9JzWePqqib8QSifC1ypRBGeJHitSnF7fA")
	newHash, err := s.NewObject("unixfs-dir")
	is.Nil(err)
	_, err = s.PatchLink(newHash, "a/b/c", newHash, false)
	is.NotNil(err)
	newHash, err = s.PatchLink(newHash, "a/b/c", newHash, true)
	is.Nil(err)
	is.Equal(newHash, "QmQ5D3xbMWFQRC9BKqbvnSnHri31GqvtWG1G6rE8xAZf1J")
}

// TestResolvePath ...
func TestResolvePath(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	childHash, err := s.ResolvePath(fmt.Sprintf("/ipfs/%s/about", examplesHash))
	is.Nil(err)
	is.Equal(childHash, "QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V")
}

// TestPubSub ...
func TestPubSub(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	var (
		topic = "test"

		sub *PubSubSubscription
		err error
	)

	t.Log("subscribing...")
	sub, err = s.PubSubSubscribe(topic)
	is.Nil(err)
	is.NotNil(sub)
	t.Log("sub: done")

	time.Sleep(10 * time.Millisecond)

	t.Log("publishing...")
	is.Nil(s.PubSubPublish(topic, "Hello World!"))
	t.Log("pub: done")

	t.Log("next()...")
	r, err := sub.Next()
	t.Log("next: done. ")

	is.Nil(err)
	is.NotNil(r)
	is.Equal(r.Data, "Hello World!")

	sub2, err := s.PubSubSubscribe(topic)
	is.Nil(err)
	is.NotNil(sub2)

	is.Nil(s.PubSubPublish(topic, "Hallo Welt!"))

	r, err = sub2.Next()
	is.Nil(err)
	is.NotNil(r)
	is.Equal(r.Data, "Hallo Welt!")

	r, err = sub.Next()
	is.NotNil(r)
	is.Nil(err)
	is.Equal(r.Data, "Hallo Welt!")

	is.Nil(sub.Cancel())
}

// TestObjectStat ...
func TestObjectStat(t *testing.T) {
	obj := "QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V"
	is := is.New(t)
	s := NewShell(shellURL)
	stat, err := s.ObjectStat("QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V")
	is.Nil(err)
	is.Equal(stat.Hash, obj)
	is.Equal(stat.LinksSize, 3)
	is.Equal(stat.CumulativeSize, 1688)
}

// TestDagPut ...
func TestDagPut(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	c, err := s.DagPut(`{"x": "abc","y":"def"}`, "json", "cbor")
	is.Nil(err)
	is.Equal(c, "zdpuAt47YjE9XTgSxUBkiYCbmnktKajQNheQBGASHj3FfYf8M")
}

// TestDagPutWithOpts ...
func TestDagPutWithOpts(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	c, err := s.DagPutWithOpts(`{"x": "abc","y":"def"}`, options.Dag.Pin("true"))
	is.Nil(err)
	is.Equal(c, "zdpuAt47YjE9XTgSxUBkiYCbmnktKajQNheQBGASHj3FfYf8M")
}

// TestStatsBW ...
func TestStatsBW(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)
	_, err := s.StatsBW(context.Background())
	is.Nil(err)
}

// TestSwarmPeers ...
func TestSwarmPeers(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)
	_, err := s.SwarmPeers(context.Background())
	is.Nil(err)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// TestRefs ...
func TestRefs(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellURL)

	cid, err := s.AddDir("./testdata")
	is.Nil(err)
	is.Equal(cid, "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv")
	refs, err := s.Refs(cid, false)
	is.Nil(err)
	expected := []string{
		"QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V",
		"QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		"QmY5heUM5qgRubMDD1og9fhCPA6QdkMp3QCwd4s7gJsyE7",
		"QmejvEPop4D7YUadeGqYWmZxHhLc4JBUCzJJHWMzdcMe2y",
		"QmXgqKTbzdh83pQtKFb19SpMCpDDcKR2ujqk3pKph9aCNF",
		"QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB",
		"QmQ5vhrL7uv6tuoN9KeVBwd4PwfQkXdVVmDLUZuTNxqgvm",
	}
	var actual []string
	for r := range refs {
		actual = append(actual, r)
	}

	sort.Strings(expected)
	sort.Strings(actual)
	is.Equal(expected, actual)
}
