package legifrance

const (
	StateActive   = "VIGUEUR"
	StateModified = "MODIFIE"
	StateRevoked  = "ABROGE"
)

const (
	LinkCodifies     = "CODIFIE"
	LinkCodification = "CODIFICATION"
	LinkCreates      = "CREE"
	LinkCreation     = "CREATION"
	LinkModifies     = "MODIFIE"
	LinkModification = "MODIFICATION"
	LinkCitation     = "CITATION"
	LinkRetirement   = "ABROGATION"

	LinkSource = "source"
	LinkTarget = "cible"
)

const (
	NatureCode       = "CODE"
	NatureDecree     = "DECRET"
	NatureDecreeLaw  = "DECRET_LOI"
	NatureLaw        = "LOI"
	NatureSentence   = "ARRETE"
	NatureOrganicLaw = "LOI_ORGANIQUE"
)

type CommonMeta struct {
	Id     string `xml:"META>META_COMMUN>ID"`
	OldId  string `xml:"META>META_COMMUN>ANCIEN_ID"`
	Origin string `xml:"META>META_COMMUN>ORIGINE"`
	URL    string `xml:"META>META_COMMUN>URL"`
	Nature string `xml:"META>META_COMMUN>NATURE"`
}

type Version struct {
	State string `xml:"etat,attr"`
	Link  *struct {
		Id        string `xml:"id,attr"`
		BeginDate string `xml:"debut,attr"`
		EndDate   string `xml:"fin,attr"`
		Number    string `xml:"num,attr"`
		Origin    string `xml:"origine,attr"`
	} `xml:"LIEN_ART"`
}

type TOCLink struct {
	Text        string `xml:",chardata"`
	Id          string `xml:"id,attr"`
	ChronicleId string `xml:"cid,attr"`
	Level       int    `xml:"niv,attr"`
	URL         string `xml:"url,attr"`
	BeginDate   string `xml:"debut,attr"`
	EndDate     string `xml:"fin,attr"`
	State       string `xml:"etat,attr"`
}

type TOC struct {
	Title *struct {
		Id        string `xml:"id,attr"`
		Name      string `xml:",chardata"`
		BeginDate string `xml:"debut,attr"`
		EndDate   string `xml:"fin,attr"`
	} `xml:"TITRE_TM"`
	Next *TOC `xml:"TM"`
}
