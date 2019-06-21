package outils

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/oauth2"
)

type FirstValid []oauth2.TokenSource

func (f FirstValid) Token() (*oauth2.Token, error) {
	var errs error
	for _, ts := range f {
		tok, err := ts.Token()
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		if !tok.Valid() {
			errs = multierror.Append(errs, fmt.Errorf("token was invalid"))
			continue
		}

		return tok, nil
	}
	return nil, errs
}
