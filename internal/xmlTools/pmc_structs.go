package xmlTools

import "encoding/xml"

// PMCArticle represents the root element of a JATS XML article.
type PMCArticle struct {
	XMLName     xml.Name        `xml:"article"`
	ArticleType string          `xml:"article-type,attr"`
	Front       PMCFront        `xml:"front"`
	Body        *PMCBody        `xml:"body,omitempty"`
	Back        *PMCBack        `xml:"back,omitempty"`
	FloatsGroup *PMCFloatsGroup `xml:"floats-group,omitempty"`
}

// PMCFloatsGroup represents a group of floating objects such as figures and tables.
type PMCFloatsGroup struct {
	Figures []PMCFigure    `xml:"fig"`
	Tables  []PMCTableWrap `xml:"table-wrap"`
}

// PMCFront contains metadata about the article.
type PMCFront struct {
	JournalMeta PMCJournalMeta `xml:"journal-meta"`
	ArticleMeta PMCArticleMeta `xml:"article-meta"`
}

// PMCJournalMeta holds metadata about the publishing journal.
type PMCJournalMeta struct {
	JournalID    []PMCID      `xml:"journal-id"`
	JournalTitle string       `xml:"journal-title-group>journal-title"`
	ISSN         []PMCISSN    `xml:"issn"`
	Publisher    PMCPublisher `xml:"publisher"`
}

// PMCID is a journal identifier.
type PMCID struct {
	IDType string `xml:"journal-id-type,attr"`
	Value  string `xml:",chardata"`
}

type PMCISSN struct {
	PubType string `xml:"pub-type,attr"`
	Value   string `xml:",chardata"`
}

type PMCPublisher struct {
	PublisherName string `xml:"publisher-name"`
	PublisherLoc  string `xml:"publisher-loc"`
}

// PMCArticleMeta describes the metadata of the article.
type PMCArticleMeta struct {
	ArticleID         []PMCArticleID      `xml:"article-id"`
	ArticleCategories []PMCSubjectGroup   `xml:"article-categories>subj-group"`
	TitleGroup        PMCTitleGroup       `xml:"title-group"`
	ContribGroup      []PMCContribGroup   `xml:"contrib-group"`
	AuthorNotes       *PMCAuthorNotes     `xml:"author-notes"`
	PubDate           []PMCPubDate        `xml:"pub-date"`
	History           []PMCDate           `xml:"history>date"`
	Abstract          *PMCAbstract        `xml:"abstract"`
	Permissions       *PMCPermissions     `xml:"permissions"`
	SelfURI           *PMCSelfURI         `xml:"self-uri"`
	RelatedArticle    *PMCRelatedArticle  `xml:"related-article"`
	CustomMetaGroup   *PMCCustomMetaGroup `xml:"custom-meta-group"`
	Volume            string              `xml:"volume"`
	Issue             string              `xml:"issue"`
	FPage             string              `xml:"fpage"`
	LPage             string              `xml:"lpage"`
	AffList           []PMCAff            `xml:"aff"`
}

type PMCArticleID struct {
	IDType string `xml:"pub-id-type,attr"`
	Value  string `xml:",chardata"`
}

type PMCSubjectGroup struct {
	SubjectGroupType string   `xml:"subj-group-type,attr"`
	Subjects         []string `xml:"subject"`
}

type PMCTitleGroup struct {
	ArticleTitle string `xml:"article-title"`
}

type PMCContribGroup struct {
	Contrib []PMCContrib `xml:"contrib"`
}

type PMCContrib struct {
	ContribType string  `xml:"contrib-type,attr"`
	Name        PMCName `xml:"name"`
	Degrees     string  `xml:"degrees"`
	Aff         *PMCAff `xml:"aff,omitempty"`
	Corresp     string  `xml:"corresp,attr,omitempty"`
}

type PMCName struct {
	Surname    string `xml:"surname"`
	GivenNames string `xml:"given-names"`
}

type PMCAff struct {
	ID   string `xml:"id,attr,omitempty"`
	Text string `xml:",chardata"`
}

type PMCAuthorNotes struct {
	Corresp []PMCCorresp `xml:"corresp"`
}

type PMCCorresp struct {
	ID    string `xml:"id,attr"`
	Email string `xml:"email"`
	Text  string `xml:",chardata"`
}

type PMCPubDate struct {
	PubType string `xml:"pub-type,attr"`
	Year    string `xml:"year"`
	Month   string `xml:"month,omitempty"`
	Day     string `xml:"day,omitempty"`
}

type PMCDate struct {
	DateType string `xml:"date-type,attr"`
	Year     string `xml:"year"`
	Month    string `xml:"month,omitempty"`
	Day      string `xml:"day,omitempty"`
}

type PMCAbstract struct {
	Title      string           `xml:"title"`
	Paragraphs []string         `xml:"p"`
	Sec        []PMCAbstractSec `xml:"sec"`
}

type PMCAbstractSec struct {
	Title      string   `xml:"title"`
	Paragraphs []string `xml:"p"`
}

type PMCPermissions struct {
	CopyrightStatement string `xml:"copyright-statement"`
}

type PMCElementCitation struct {
	PublicationType string    `xml:"publication-type,attr"`
	ArticleTitle    string    `xml:"article-title"`
	Source          string    `xml:"source"`
	Year            string    `xml:"year"`
	Volume          string    `xml:"volume"`
	FPage           string    `xml:"fpage"`
	LPage           string    `xml:"lpage"`
	PubID           string    `xml:"pub-id"`
	Name            []PMCName `xml:"name"`
}

type PMCSelfURI struct {
	Href string `xml:"xlink:href,attr"`
}

type PMCRelatedArticle struct {
	Type string `xml:"related-article-type,attr"`
	ID   string `xml:"id,attr"`
	Href string `xml:"xlink:href,attr"`
}

type PMCCustomMetaGroup struct {
	CustomMeta []PMCCustomMeta `xml:"custom-meta"`
}

type PMCCustomMeta struct {
	Name  string `xml:"meta-name"`
	Value string `xml:"meta-value"`
}

type PMCBody struct {
	Sections []PMCSection `xml:"sec"`
}

type PMCSection struct {
	ID          string         `xml:"id,attr,omitempty"`
	SecType     string         `xml:"sec-type,attr,omitempty"`
	Title       string         `xml:"title"`
	Paragraphs  []string       `xml:"p"`
	SubSections []PMCSection   `xml:"sec"`
	Figures     []PMCFigure    `xml:"fig"`
	Tables      []PMCTableWrap `xml:"table-wrap"`
	XRefs       []PMCXRef      `xml:"xref"`
}

type PMCXRef struct {
	RefType string `xml:"ref-type,attr"`
	RID     string `xml:"rid,attr"`
	Text    string `xml:",chardata"`
}

type PMCTableWrap struct {
	ID      string     `xml:"id,attr"`
	Label   string     `xml:"label"`
	Caption PMCCaption `xml:"caption"`
	Graphic PMCGraphic `xml:"graphic"`
}

type PMCFigure struct {
	ID      string     `xml:"id,attr"`
	Label   string     `xml:"label"`
	Caption PMCCaption `xml:"caption"`
	Graphic PMCGraphic `xml:"graphic"`
}

type PMCCaption struct {
	Paragraphs []string `xml:"p"`
}

type PMCGraphic struct {
	Href string `xml:"xlink:href,attr"`
}

type PMCBack struct {
	Acknowledgments *PMCAcknowledgments `xml:"ack,omitempty"`
	References      *PMCReferences      `xml:"ref-list,omitempty"`
	FnGroup         *PMCFnGroup         `xml:"fn-group,omitempty"`
}

type PMCAcknowledgments struct {
	Paragraphs []string `xml:"p"`
}

type PMCReferences struct {
	References []PMCReference `xml:"ref"`
}

type PMCReference struct {
	ID              string              `xml:"id,attr"`
	ElementCitation *PMCElementCitation `xml:"element-citation"`
	MixedCitation   *PMCMixedCitation   `xml:"mixed-citation"`
}
type PMCMixedCitation struct {
	PublicationType string `xml:"publication-type,attr"`
	ArticleTitle    string `xml:"article-title"`
	Source          string `xml:"source"`
	Year            string `xml:"year"`
	Volume          string `xml:"volume"`
	FPage           string `xml:"fpage"`
	LPage           string `xml:"lpage"`
	PubID           string `xml:"pub-id"`
}

type PMCFnGroup struct {
	Footnotes []PMCFootnote `xml:"fn"`
}

type PMCFootnote struct {
	Type string   `xml:"fn-type,attr"`
	Text []string `xml:"p"`
}
