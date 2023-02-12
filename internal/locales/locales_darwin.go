//go:build darwin

package locales

import "os"

func inferLanguages() []string {
	// Is it a correct way for osx?
	// It's the same as in Linux right now.
	var langs []string
	if lang, ok := os.LookupEnv("LANG"); ok {
		langs = append(langs, lang)
	}
	return langs
}
