package store

import (
	"fmt"
	"strings"
)

func Extension(ct ContentType, ce ContentEncoding) string {
	var parts []string

	ctExt := ct.FileExtension()
	if ctExt != "" {
		parts = append(parts, ctExt)
	}

	ceExt := ce.FileExtension()
	if ceExt != "" {
		parts = append(parts, ceExt)
	}

	return strings.Join(parts, ".")
}

func FilePath(key string, headers *Headers) string {
	var (
		ct ContentType
		ce ContentEncoding
	)

	if headers != nil {
		ct = headers.ContentType
		ce = headers.ContentEncoding
	}
	return fmt.Sprintf("%s.%s", key, Extension(ct, ce))
}
