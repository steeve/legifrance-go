package legifrance

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	ReHtmlTag       = regexp.MustCompile(`<[^>]*>`)
	ReWhiteSpaces   = regexp.MustCompile(`\s+`)
	ReCR            = regexp.MustCompile(`\n{2,}`)
	ReSpaces        = regexp.MustCompile(` +`)
	ReSpacesAfterCR = regexp.MustCompile(`\n[ ]+`)
	ReLawNumber     = regexp.MustCompile(`[0-9]{2,4}-[0-9]{2,4}(-[0-9]{2})?`)
)

type Article struct {
	*CommonMeta

	Code    *Code    `xml:"-"`
	Section *Section `xml:"-"`

	VersionId string `xml:"-"`
	Number    string `xml:"META>META_SPEC>META_ARTICLE>NUM"`
	State     string `xml:"META>META_SPEC>META_ARTICLE>ETAT"`
	BeginDate string `xml:"META>META_SPEC>META_ARTICLE>DATE_DEBUT"`
	EndDate   string `xml:"META>META_SPEC>META_ARTICLE>DATE_FIN"`
	Type      string `xml:"META>META_SPEC>META_ARTICLE>TYPE"`
	Notice    string `xml:"NOTA>CONTENU>p"`
	Body      struct {
		Text string `xml:",innerxml"`
	} `xml:"BLOC_TEXTUEL>CONTENU"`

	VersionLinks []*Version `xml:"VERSIONS>VERSION"`

	Text *TextLink `xml:"CONTEXTE>TEXTE"`

	Links []*struct {
		Text         string `xml:",chardata"`
		TextId       string `xml:"cidTexte,attr"`
		Id           string `xml:"id,attr"`
		TextSignDate string `xml:"datesignatexte,attr"`
		TextNature   string `xml:"naturetexte,attr"`
		TextNor      string `xml:"nortexte,attr"`
		Number       string `xml:"num,attr"`
		TextNumber   string `xml:"numtexte,attr"`
		Direction    string `xml:"sens,attr"`
		LinkType     string `xml:"typelien,attr"`
	} `xml:"LIENS>LIEN"`
}

func NewArticleFromReader(r io.Reader) (*Article, error) {
	a := &Article{}
	if err := xml.NewDecoder(r).Decode(&a); err != nil {
		return nil, err
	}
	a.VersionId = fmt.Sprintf("%s/%s", a.Text.Id, a.Number)

	a.Body.Text = ReHtmlTag.ReplaceAllString(a.Body.Text, "\n")
	a.Body.Text = ReSpacesAfterCR.ReplaceAllString(a.Body.Text, "\n")
	a.Body.Text = ReSpaces.ReplaceAllString(a.Body.Text, " ")
	a.Body.Text = ReCR.ReplaceAllString(a.Body.Text, "\n\n")
	a.Body.Text = strings.TrimSpace(a.Body.Text)

	for _, link := range a.Links {
		link.Text = ReWhiteSpaces.ReplaceAllString(link.Text, " ")
		if link.TextNumber == "" {
			link.TextNumber = ReLawNumber.FindString(link.Text)
		}
		if link.TextSignDate == "" && link.TextNumber != "" {
			if _, err := time.Parse("2006-01-02", link.TextNumber); err == nil {
				link.TextSignDate = link.TextNumber
			}
		}
	}

	return a, nil
}

func NewArticleFromFile(path string) (*Article, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewArticleFromReader(f)
}

func (a *Article) Versions() chan *Article {
	ch := make(chan *Article)

	go func() {
		for _, versionLink := range a.VersionLinks {
			if versionLink.Link.Id == a.Id {
				ch <- a
			} else {
				article, err := a.Code.Article(versionLink.Link.Id)
				if err == nil {
					ch <- article
				}
			}
		}
		close(ch)
	}()

	return ch
}

func (a *Article) Path() []string {
	path := []string{
		a.Text.Title.ShortTitle,
	}
	current := a.Text.TOC
	for current != nil {
		path = append(path, current.Title.Name)
		current = current.Next
	}
	return path
}
