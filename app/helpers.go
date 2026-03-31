package app

import "github.com/nmvalera/go-utils/tag"

func ComponentTag(component string) *tag.Tag {
	return tag.Key("component").String(component)
}
