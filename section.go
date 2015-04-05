package legifrance

import (
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Section struct {
	Code *Code `xml:"-"`

	Id           string     `xml:"ID"`
	Title        string     `xml:"TITRE_TA"`
	TOC          []*TOCLink `xml:"STRUCTURE_TA>LIEN_SECTION_TA"`
	ArticleLinks []*TOCLink `xml:"STRUCTURE_TA>LIEN_ART"`

	Text TextLink `xml:"CONTEXTE>TEXTE"`
}

func NewSectionFromReader(r io.Reader) (*Section, error) {
	s := &Section{}
	if err := xml.NewDecoder(r).Decode(&s); err != nil {
		return nil, err
	}
	s.Title = strings.TrimSpace(s.Title)
	return s, nil
}

func NewSectionFromFile(path string) (*Section, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewSectionFromReader(f)
}

func (s *Section) Article(articleId string) (*Article, error) {
	path := filepath.Join(s.Code.l.Path, "global", "code_et_TNC_en_vigueur", "code_en_vigueur", legiIdToSlug(s.Code.Id), s.Code.Id, "article", legiIdToSlug(articleId), articleId) + ".xml"
	article, err := NewArticleFromFile(path)
	if err != nil {
		return nil, err
	}
	article.Code = s.Code
	article.Section = s
	return article, nil
}

func (s *Section) Articles() <-chan *Article {
	ch := make(chan *Article)

	go func() {
		for _, articleLink := range s.ArticleLinks {
			article, err := s.Article(articleLink.Id)
			if err == nil {
				ch <- article
			}
		}
		close(ch)
	}()

	return ch
}

func (s *Section) SubSections() <-chan *Section {
	ch := make(chan *Section)

	go func() {
		for _, sectionLink := range s.TOC {
			section, err := s.Code.Section(sectionLink.Id)
			if err == nil {
				ch <- section
			}
		}
		close(ch)
	}()

	return ch
}

func walkSectionTree(s *Section, ch chan *Article) {
	for article := range s.Articles() {
		ch <- article
	}
	for section := range s.SubSections() {
		walkSectionTree(section, ch)
	}
}

func (s *Section) ArticlesWithChildren() <-chan *Article {
	ch := make(chan *Article)

	go func() {
		walkSectionTree(s, ch)
		close(ch)
	}()

	return ch
}
