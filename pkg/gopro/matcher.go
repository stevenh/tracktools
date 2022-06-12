package gopro

import (
	"fmt"
	"path/filepath"
	"regexp"
)

// Formats are determined from:
// https://community.gopro.com/s/article/GoPro-Camera-File-Naming-Convention?language=en_US
var (
	Hero5  = NewMatcherMust("hero5", `^(?i)GOPR([0-9]{4})\.mp4$`, `^(?i)GP([0-9]{2})([0-9]{4}).mp4$`)
	Hero10 = NewMatcherMust("hero10", `^(?i)G[HX]([0-9]{2})([0-9]{4})\.mp4$`)
)

// Matcher provides the ability to identify GoPro video files.
type Matcher struct {
	name    string
	regexps []*regexp.Regexp
}

// NewMatcher returns a fully initialised Matcher.
func NewMatcher(name, first string, others ...string) (*Matcher, error) {
	m := &Matcher{
		name:    name,
		regexps: make([]*regexp.Regexp, len(others)+1),
	}
	for i, p := range append([]string{first}, others...) {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("regexp %q: %w", p, err)
		}

		m.regexps[i] = re
	}

	return m, nil
}

// NewMatcherMust returns a fully initialised Matcher and panics if
// any error occurs.
func NewMatcherMust(name, first string, others ...string) *Matcher {
	m, err := NewMatcher(name, first, others...)
	if err != nil {
		panic(err)
	}

	return m
}

// Match checks file if matches and returns the decoded File.
// If no match is found ErrNoMatch is returned.
func (m *Matcher) Match(file string) (*File, error) {
	f := filepath.Base(file)
	for _, re := range m.regexps {
		ma := re.FindStringSubmatch(f)
		switch len(ma) {
		case 0:
			// No match ignore for now.
		case 2:
			return &File{Name: f, Index: ma[1], Chapter: "00"}, nil
		case 3:
			return &File{Name: f, Index: ma[2], Chapter: ma[1]}, nil
		default:
			return nil, matchError(len(ma))
		}
	}

	return nil, ErrNoMatch
}
