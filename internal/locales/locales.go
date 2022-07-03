package locales

import (
	"regexp"
	"strings"
)

func InferLanguages() []string {
	langs := inferLanguages()
	normalizeLangs(langs)
	return langs
}

func normalizeLangs(langs []string) {
	for i, l := range langs {
		langs[i] = normalizeLang(l)
	}
}

func normalizeLang(l string) string {
	l = strings.ToLower(l)
	if linuxLangRE.MatchString(l) {
		l = l[:strings.IndexByte(l, '.')]
		l = strings.Replace(l, "_", "-", 1)
	}
	if simpleLangRE.MatchString(l) {
		return l[:strings.IndexByte(l, '-')]
	}
	return l
}

var (
	simpleLangRE = regexp.MustCompile(`\w+-\w+`)
	linuxLangRE  = regexp.MustCompile(`[a-zA-Z]_\w+.[\-\w]`)
)
