package store

import "fmt"

type ContentEncoding int

const (
	ContentEncodingPlain ContentEncoding = iota
	ContentEncodingGzip
	ContentEncodingZlib
	ContentEncodingFlate
)

var contentEncodingStrings = [...]string{
	"plain",
	"gzip",
	"zlib",
	"flate",
}

var contentEncodings = map[string]ContentEncoding{
	contentEncodingStrings[ContentEncodingPlain]: ContentEncodingPlain,
	contentEncodingStrings[ContentEncodingGzip]:  ContentEncodingGzip,
	contentEncodingStrings[ContentEncodingZlib]:  ContentEncodingZlib,
	contentEncodingStrings[ContentEncodingFlate]: ContentEncodingFlate,
}

func (ce ContentEncoding) String() string {
	if ce < 0 || int(ce) >= len(contentEncodingStrings) {
		return unknown
	}
	return contentEncodingStrings[ce]
}

var contentEncodingFileExtensions = map[ContentEncoding]string{
	ContentEncodingPlain: "",
	ContentEncodingGzip:  "gz",
	ContentEncodingZlib:  "zlib",
	ContentEncodingFlate: "flate",
}

func (ce ContentEncoding) FileExtension() string {
	ext, ok := contentEncodingFileExtensions[ce]
	if !ok {
		return ""
	}
	return ext
}

func (ce ContentEncoding) FilePath(key string) string {
	ext := ce.FileExtension()
	if ext == "" {
		return key
	}
	return fmt.Sprintf("%s.%s", key, ext)
}

func ParseContentEncoding(compression string) (ContentEncoding, error) {
	if ce, ok := contentEncodings[compression]; ok {
		return ce, nil
	}
	return -1, fmt.Errorf("invalid compression: %s", compression)
}
