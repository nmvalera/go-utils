package tag

func MapTags(m map[string]any) []*Tag {
	tags := make([]*Tag, 0, len(m))

	for k, v := range m {
		switch v := v.(type) {
		case string:
			tags = append(tags, Key(k).String(v))
		case int64:
			tags = append(tags, Key(k).Int64(v))
		case int:
			tags = append(tags, Key(k).Int64(int64(v)))
		case float64:
			tags = append(tags, Key(k).Float64(v))
		case bool:
			tags = append(tags, Key(k).Bool(v))
		case map[string]any:
			tags = append(tags, Key(k).Map(MapTags(v)...))
		default:
			tags = append(tags, Key(k).Object(v))
		}
	}
	return tags
}
