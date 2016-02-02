package validate

import (
	"errors"
	"strings"
)

// Friendly user-facing errors
var (
	errBadEmail = errors.New("User doesn't exist.")
)

// Email does basic validation on the intended receiving address. It returns a
// friendly error message that can be served directly to connecting clients if
// validation fails.
func Email(to string) error {
	host := to[strings.IndexRune(to, '@')+1:]
	// TODO: use given configurable host (don't hardcode)
	if host != "writ.es" {
		return errBadEmail
	}

	toName := to[:strings.IndexRune(to, '@')]
	if toName == "anyone" {
		return nil
	}

	if len(toName) != 32 {
		return errBadEmail
	}

	return nil
}
