package legifrance

type TextLink struct {
	Authority     string `xml:"autorite,attr"`
	Id            string `xml:"cid,attr"`
	PublishedDate string `xml:"date_publi,attr"`
	SignedDate    string `xml:"date_signature,attr"`
	Ministry      string `xml:"ministere,attr"`
	Nature        string `xml:"nature,attr"`
	Nor           string `xml:"nor,attr"`
	Number        string `xml:"num,attr"`

	Title *struct {
		Title      string `xml:",chardata"`
		ShortTitle string `xml:"c_titre_court,attr"`
		Id         string `xml:"id_txt,attr"`
		BeginDate  string `xml:"debut,attr"`
		EndDate    string `xml:"fin,attr"`
	} `xml:"TITRE_TXT"`

	TOC *TOC `xml:"TM"`
}
