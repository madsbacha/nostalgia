package util

import "strings"

func RemoveEmptyTags(tags []string) []string {
	var result []string
	for _, tag := range tags {
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

func ParseTags(rawTags string) []string {
	if strings.TrimSpace(rawTags) == "" {
		return []string{}
	}
	tags := strings.Split(strings.ToLower(rawTags), ",")
	if tags == nil {
		return []string{}
	}
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return tags
}

func JoinTags(tags []string) string {
	if tags == nil {
		return ""
	}
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return strings.Join(tags, ",")
}
