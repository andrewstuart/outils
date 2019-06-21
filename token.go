package outils

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
)

type Wrapper interface {
	Wrap(ts oauth2.TokenSource) oauth2.TokenSource
}

type WrapperFunc func(oauth2.TokenSource) oauth2.TokenSource

func (w WrapperFunc) Wrap(ts oauth2.TokenSource) oauth2.TokenSource {
	return w(ts)
}

type SourceFunc func() (*oauth2.Token, error)

func (s SourceFunc) Token() (*oauth2.Token, error) {
	return s()
}

var _ oauth2.TokenSource = SourceFunc(nil)

type DiskCache struct {
	Filename    string
	TokenSource oauth2.TokenSource
}

func (dc *DiskCache) Token() (t *oauth2.Token, err error) {
	f, err := os.OpenFile(dc.Filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return dc.TokenSource.Token()
	}
	defer func() {
		f.Truncate(0)
		f.Seek(0, 0)
		json.NewEncoder(f).Encode(t)
		f.Close()
	}()

	t = &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	if err == nil && t.Valid() {
		return t, nil
	}

	return dc.TokenSource.Token()
}

var _ oauth2.TokenSource = &DiskCache{}