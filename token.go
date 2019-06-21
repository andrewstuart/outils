package outils

import (
	"encoding/json"
	"os"
	"path"

	"golang.org/x/oauth2"
)

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
	err = os.MkdirAll(path.Dir(dc.Filename), 0755)
	if err != nil {
		return dc.TokenSource.Token()
	}

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
