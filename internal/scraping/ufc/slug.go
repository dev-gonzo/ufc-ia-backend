package ufc

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func AthleteSlugFromName(name string) (string, error) {
	s := strings.TrimSpace(name)
	if s == "" {
		return "", ErrInvalidAthleteSlug
	}

	s = strings.ToLower(s)
	s = norm.NFD.String(s)
	var out []rune
	prevDash := false

	for _, r := range s {
		if unicode.Is(unicode.Mn, r) {
			continue
		}

		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			out = append(out, r)
			prevDash = false
		case r == '-' || unicode.IsSpace(r) || r == '_' || r == '\'' || r == '’' || r == '.' || r == ',' || r == '/' || r == '\\' || r == '(' || r == ')' || r == '[' || r == ']' || r == '&':
			if !prevDash && len(out) > 0 {
				out = append(out, '-')
				prevDash = true
			}
		default:
			if !prevDash && len(out) > 0 {
				out = append(out, '-')
				prevDash = true
			}
		}
	}

	slug := strings.Trim(strings.TrimSpace(string(out)), "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	if slug == "" {
		return "", ErrInvalidAthleteSlug
	}
	return slug, nil
}
