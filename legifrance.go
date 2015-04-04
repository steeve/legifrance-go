package legifrance

import (
	"path/filepath"
	"regexp"
)

var (
	slugRe = regexp.MustCompile(`[A-Z]{4}|[0-9]{2}`)
)

type LegiFrance struct {
	Path string
}

func NewLegiFrance(dumpPath string) *LegiFrance {
	return &LegiFrance{
		Path: dumpPath,
	}
}

func legiIdToSlug(id string) string {
	slug := slugRe.ReplaceAllStringFunc(id, func(m string) string {
		return m + "/"
	})
	return slug[:len(slug)-3]
}

func (l *LegiFrance) Code(codeId string) (*Code, error) {
	path := filepath.Join(l.Path, legiIdToSlug(codeId), codeId, "texte", "struct", codeId) + ".xml"
	code, err := NewCodeFromFile(path)
	if err != nil {
		return nil, err
	}
	code.l = l
	return code, nil
}
