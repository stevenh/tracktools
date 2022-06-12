package gopro

import (
	"fmt"
	"sort"
)

// File represents a GoPro video file.
type File struct {
	// Name is the raw filename.
	Name string

	// Chapter is the chapter of the video typically a two character number.
	// Depending on generation the first chapter may be 00 or 01.
	Chapter string

	// Index is video index of the video typically a 4 character number.
	Index string
}

// FileSlice attaches the methods of sort.Interface to []File,
// sorting in increasing chapter order.
type FileSlice []File

// Len implements sort.Interface.
func (s FileSlice) Len() int {
	return len(s)
}

// Less implements sort.Interface by Chapter.
func (s FileSlice) Less(i, j int) bool {
	return s[i].Chapter < s[j].Chapter
}

// Swap implements sort.Interface.
func (s FileSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Validate validates the FileSlice.
func (s FileSlice) Validate() error {
	if len(s) == 0 {
		return errNoChapters
	}

	var inc int
	switch s[0].Chapter {
	case "00":
	case "01":
		inc = 1
	default:
		return chapterError{chapter: s[0].Chapter}
	}

	for i, c := range s {
		if c.Chapter != fmt.Sprintf("%02d", i+inc) {
			return chapterError{chapter: c.Chapter, idx: i}
		}
	}

	return nil
}

// FileSet represents a set of files for a single GoPro video.
type FileSet struct {
	Number   string
	Chapters FileSlice
}

// Chapter adds file as a chapter ensuring correct order.
func (s *FileSet) Chapter(f *File) {
	s.Chapters = append(s.Chapters, *f)
	sort.Sort(s.Chapters)
}
