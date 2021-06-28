package genutil

import (
	"fmt"
	"strconv"
	"strings"
)

type StructTag struct {
	Name  string
	Value string
}

func (t StructTag) String() string {
	return fmt.Sprintf("%s:%q", t.Name, t.Value)
}

type StructTags []StructTag

func (tags StructTags) String() string {
	s := make([]string, 0, len(tags))
	for _, tag := range tags {
		s = append(s, tag.String())
	}
	return "`" + strings.Join(s, " ") + "`"
}

func (tags StructTags) Get(name string) string {
	for _, tag := range tags {
		if tag.Name == name {
			return tag.Value
		}
	}
	return ""
}

func MustParseStructTags(tag string) StructTags {
	tags, err := ParseStructTags(tag)
	if err != nil {
		panic(err)
	}
	return tags
}

// ParseStructTags returns the full set of fields in a struct tag in the order they appear in
// the struct tag.
func ParseStructTags(tag string) (StructTags, error) {
	tags := StructTags{}
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := tag[:i]
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := tag[:i+1]
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return nil, err
		}
		tags = append(tags, StructTag{Name: name, Value: value})
	}
	return tags, nil
}
