package langs

import (
	"bytes"
	"strings"
)

type Dictionary struct {
	Name       string
	entries    map[string]string
	titleWords map[string]string
}

func ParseDictionary(name string, data []byte) (*Dictionary, error) {
	dict := &Dictionary{
		Name:       name,
		entries:    make(map[string]string, 64),
		titleWords: make(map[string]string, 16),
	}

	offset := 0
	sectionBodyBegin := 0
	sectionBodyEnd := 0
	sectionKey := ""
	nextSectionBodyBegin := 0
	nextSectionKey := ""
	for {
		lineEnd := bytes.IndexByte(data[offset:], '\n')
		flush := false
		stop := false
		if lineEnd == -1 || offset >= len(data) {
			stop = true
			offset = len(data)
			flush = true
			sectionBodyEnd = len(data)
		} else {
			line := data[offset : offset+lineEnd]
			if bytes.HasPrefix(line, []byte("##")) {
				colonPos := bytes.IndexByte(line, ':')
				flush = true
				sectionBodyEnd = offset
				if colonPos != -1 {
					nextSectionBodyBegin = offset + colonPos + 1
					nextSectionKey = string(bytes.TrimSpace(line[len("##"):colonPos]))
				} else {
					nextSectionBodyBegin = offset + lineEnd + 1
					nextSectionKey = string(line[len("##"):])
				}
			}
			offset += lineEnd + 1
		}
		if flush {
			if sectionKey != "" {
				s := strings.TrimSpace(string(data[sectionBodyBegin:sectionBodyEnd]))
				dict.entries[sectionKey] = s
			}
			sectionKey = nextSectionKey
			sectionBodyBegin = nextSectionBodyBegin
		}
		if stop {
			break
		}
	}

	return dict, nil
}

func (d *Dictionary) GetTitleCase(key string) string {
	s, ok := d.titleWords[key]
	if !ok {
		s2, ok := d.entries[key]
		if !ok {
			return "{{" + key + "}}"
		}
		s = strings.Title(s2)
		d.titleWords[key] = s
	}
	return s
}

func (d *Dictionary) Get(key string) string {
	s, ok := d.entries[key]
	if !ok {
		return "{{" + key + "}}"
	}
	return s
}
