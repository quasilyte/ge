//go:build linux

package locales

import "os"

func inferLanguages() []string {
	var langs []string
	if lang, ok := os.LookupEnv("LANG"); ok {
		langs = append(langs, lang)
	}
	return langs
}
