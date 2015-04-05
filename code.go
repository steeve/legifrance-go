package legifrance

import (
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
)

type Code struct {
	*CommonMeta

	l *LegiFrance

	Id                   string   `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>CID"`
	Number               string   `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>NUM"`
	SequenceNumber       string   `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>NUM_SEQUENCE"`
	Nor                  string   `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>NOR"`
	PublishedDate        string   `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>DATE_PUBLI"`
	TextDate             string   `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>DATE_TEXTE"`
	LastModification     string   `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>DERNIERE_MODIFICATION"`
	VersionsToCome       []string `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>VERSIONS_A_VENIR>VERSION_A_VENIR"`
	PublicationOrigin    string   `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>ORIGINE_PUBLI"`
	PublicationBeginPage string   `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>PAGE_DEB_PUBLI"`
	PublicationEndPage   string   `xml:"META>META_SPEC>META_TEXTE_CHRONICLE>PAGE_FIN_PUBLI"`

	TOC []*TOCLink `xml:"STRUCT>LIEN_SECTION_TA"`
}

func NewCodeFromReader(r io.Reader) (*Code, error) {
	c := &Code{}
	if err := xml.NewDecoder(r).Decode(&c); err != nil {
		return nil, err
	}
	return c, nil
}

func NewCodeFromFile(path string) (*Code, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewCodeFromReader(f)
}

func (c *Code) Section(sectionId string) (*Section, error) {
	path := filepath.Join(c.l.Path, "global", "code_et_TNC_en_vigueur", "code_en_vigueur", legiIdToSlug(c.Id), c.Id, "section_ta", legiIdToSlug(sectionId), sectionId) + ".xml"
	section, err := NewSectionFromFile(path)
	if err != nil {
		return nil, err
	}
	section.Code = c
	return section, nil
}

func (c *Code) Article(articleId string) (*Article, error) {
	path := filepath.Join(c.l.Path, "global", "code_et_TNC_en_vigueur", "code_en_vigueur", legiIdToSlug(c.Id), c.Id, "article", legiIdToSlug(articleId), articleId) + ".xml"
	article, err := NewArticleFromFile(path)
	if err != nil {
		return nil, err
	}
	article.Code = c
	return article, nil
}

func (c *Code) Sections() <-chan *Section {
	ch := make(chan *Section)

	go func() {
		for _, sectionLink := range c.TOC {
			section, err := c.Section(sectionLink.Id)
			if err == nil {
				ch <- section
			}
		}
		close(ch)
	}()

	return ch
}

func (c *Code) Articles() <-chan *Article {
	ch := make(chan *Article)

	go func() {
		for section := range c.Sections() {
			for article := range section.ArticlesWithChildren() {
				ch <- article
			}
		}
		close(ch)
	}()

	return ch
}
