package store

import (
	"context"
	"fmt"
	"io"
	"strings"
)

// Store is an interface for storing and loading objects.
//
//go:generate mockgen -destination=./mock/store.go -package=mock github.com/kkrt-labs/go-utils/store Store
type Store interface {
	// Store stores an object in the store.
	//
	// The key is the identifier for the object.
	// The reader is the object to store.
	// The headers are optional metadata about the object.
	Store(ctx context.Context, key string, reader io.Reader, headers *Headers) error

	// Load loads an object from the store.
	//
	// The key is the identifier for the object.
	// The headers are optional metadata about the object.
	// It is the responsibility of the caller to close the returned reader
	Load(ctx context.Context, key string) (io.ReadCloser, *Headers, error)
}

// Headers are optional metadata about an object to store/load
type Headers struct {
	// ContentType is the type of the object
	ContentType ContentType

	// ContentEncoding is the compression algorithm used to store the object.
	ContentEncoding ContentEncoding

	// KeyValue is a map of key-value pairs to store/load with the object.
	KeyValue map[string]string
}

func (h *Headers) GetContentType() (string, error) {
	return strings.TrimPrefix(h.ContentType.String(), "application/"), nil
}
func (h *Headers) GetContentEncoding() (ContentEncoding, error) {
	switch h.ContentEncoding {
	case ContentEncodingGzip:
		return ContentEncodingGzip, nil
	case ContentEncodingZlib:
		return ContentEncodingZlib, nil
	case ContentEncodingFlate:
		return ContentEncodingFlate, nil
	case ContentEncodingPlain:
		return ContentEncodingPlain, nil
	}
	return -1, fmt.Errorf("invalid compression: %s", h.ContentEncoding)
}

var unknown = "unknown"

type ContentType int

const (
	ContentTypeUnknown ContentType = iota
	ContentTypeText
	ContentTypeJSON
	ContentTypeProtobuf
)

var contentTypeStrings = [...]string{
	"",
	"text/plain",
	"application/json",
	"application/protobuf",
}

func (ct ContentType) String() string {
	if ct < 0 || int(ct) >= len(contentTypeStrings) {
		return unknown
	}
	return contentTypeStrings[ct]
}

var contentTypeFileExtensions = map[ContentType]string{
	ContentTypeText:     "",
	ContentTypeJSON:     "json",
	ContentTypeProtobuf: "protobuf",
}

func (ct ContentType) FileExtension() string {
	ext, ok := contentTypeFileExtensions[ct]
	if !ok {
		return ""
	}
	return ext
}

func (ct ContentType) FilePath(key string) string {
	ext := ct.FileExtension()
	if ext == "" {
		return key
	}
	return fmt.Sprintf("%s.%s", key, ext)
}

var contentTypes = map[string]ContentType{
	contentTypeStrings[ContentTypeJSON]:     ContentTypeJSON,
	contentTypeStrings[ContentTypeProtobuf]: ContentTypeProtobuf,
}

func ParseContentType(contentType string) (ContentType, error) {
	if ct, ok := contentTypes[contentType]; ok {
		return ct, nil
	}
	return -1, fmt.Errorf("invalid content type: %s", contentType)
}

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
