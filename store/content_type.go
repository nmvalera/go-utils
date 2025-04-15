package store

import "fmt"

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
