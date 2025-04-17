package xmlTools

import "encoding/xml"

// PubmedArticleSet is the root of a regular PubMed XML file.
// It contains a list of PubmedArticle elements.
type PubmedArticleSet struct {
	PubmedArticles []PubmedArticle `xml:"PubmedArticle"`
}

// PubmedArticle represents one article in the PubMed XML.
// It includes citation details and PubMed-specific metadata.
// Unknown captures unmapped tags.
type PubmedArticle struct {
	MedlineCitation MedlineCitation  `xml:"MedlineCitation"`
	PubmedData      PubmedData       `xml:"PubmedData"`
	Unknown         []UnknownElement `xml:",any"`
}

// PubmedBookArticleSet is the root for PubMed Book XML files.
type PubmedBookArticleSet struct {
	PubmedBookArticles []PubmedBookArticle `xml:"PubmedBookArticle"`
}

// PubmedBookArticle represents one book chapter or article.
type PubmedBookArticle struct {
	BookDocument   BookDocument   `xml:"BookDocument"`
	PubmedBookData PubmedBookData `xml:"PubmedBookData"`
}

// BookDocument holds metadata about a book section or article.
type BookDocument struct {
	PMID             string        `xml:"PMID"`
	ArticleIdList    ArticleIdList `xml:"ArticleIdList"`
	Book             string        `xml:"Book"`
	ArticleTitle     string        `xml:"ArticleTitle"`
	Abstract         Abstract      `xml:"Abstract"`
	KeywordList      []Keyword     `xml:"KeywordList>Keyword" json:"KeywordList"`
	GrantList        GrantList     `xml:"GrantList"`
	ReferenceList    []Reference   `xml:"ReferenceList>Reference" json:"ReferenceList"`
	PublicationType  string        `xml:"PublicationType"`
	InvestigatorList string        `xml:"InvestigatorList"`
	ContributionDate string        `xml:"ContributionDate"`
	DateRevised      string        `xml:"DateRevised"`
	ItemList         ItemList      `xml:"ItemList"`
	LocationLabel    string        `xml:"LocationLabel"`
}

// MedlineCitation holds the main bibliographic content.
type MedlineCitation struct {
	PMID                    string          `xml:"PMID"`
	DateCompleted           PubMedPubDate   `xml:"DateCompleted"`
	DateRevised             PubMedPubDate   `xml:"DateRevised"`
	Article                 Article         `xml:"Article"`
	MedlineJournalInfo      Journal         `xml:"MedlineJournalInfo"`
	ChemicalList            []Chemical      `xml:"ChemicalList>Chemical"`
	SupplMeshList           SupplMeshList   `xml:"SupplMeshList"`
	CitationSubset          string          `xml:"CitationSubset"`
	CommentsCorrectionsList string          `xml:"CommentsCorrectionsList"`
	GeneSymbolList          GeneSymbolList  `xml:"GeneSymbolList"`
	MeshHeadingList         MeshHeadingList `xml:"MeshHeadingList"`
	NumberOfReferences      string          `xml:"NumberOfReferences"`
	PersonalNameSubjectList string          `xml:"PersonalNameSubjectList"`
	OtherID                 string          `xml:"OtherID"`
	OtherAbstract           Abstract        `xml:"OtherAbstract"`
	KeywordList             []string        `xml:"KeywordList>Keyword" json:"KeywordList"`
	CoiStatement            string          `xml:"CoiStatement"`
	SpaceFlightMission      string          `xml:"SpaceFlightMission"`
	InvestigatorList        string          `xml:"InvestigatorList"`
	GeneralNote             string          `xml:"GeneralNote"`
	Owner                   string          `xml:"Owner,attr"`
	Status                  string          `xml:"Status,attr"`
	Medline                 string          `xml:"MEDLINE,attr"`
	VersionID               string          `xml:"VersionID,attr"`
	VersionDate             string          `xml:"VersionDate,attr"`
	IndexingMethod          string          `xml:"IndexingMethod,attr"`
}

// PubmedBookData contains metadata for books like IDs and objects.
type PubmedBookData struct {
	History           string        `xml:"History"`
	PublicationStatus string        `xml:"PublicationStatus"`
	ArticleIdList     ArticleIdList `xml:"ArticleIdList"`
	ObjectList        []Object      `xml:"ObjectList>Object"`
}

// PubmedData contains reference lists and metadata.
type PubmedData struct {
	History           []PubMedPubDate `xml:"History>PubMedPubDate"`
	PublicationStatus string          `xml:"PublicationStatus"`
	ArticleIdList     ArticleIdList   `xml:"ArticleIdList"`
	ObjectList        []Object        `xml:"ObjectList>Object"`
	ReferenceList     []Reference     `xml:"ReferenceList>Reference" json:"ReferenceList"`
}

// Abstract represents the article’s abstract.
type Abstract struct {
	AbstractText         string `xml:"AbstractText"`
	CopyrightInformation string `xml:"CopyrightInformation"`
}

// AffiliationInfo holds author affiliation metadata.
type AffiliationInfo struct {
	Affiliation string `xml:"Affiliation"`
	Identifier  string `xml:"Identifier"`
}

// ArticleId represents one identifier (DOI, PMID, etc.).
type ArticleId struct {
	ID     string `xml:",chardata"`
	IdType string `xml:"IdType,attr"`
}

// ArticleIdList is a wrapper for multiple ArticleIds.
type ArticleIdList struct {
	ArticleIds []ArticleId `xml:"ArticleId"`
}

type Day struct{} // Placeholder if day-specific parsing is needed

// GeneSymbolList represents a list of gene symbols.
type GeneSymbolList struct {
	GeneSymbols []string `xml:"GeneSymbol"`
}

// Grant holds one funding entry.
type Grant struct {
	GrantID string `xml:"GrantID"`
	Acronym string `xml:"Acronym"`
	Agency  string `xml:"Agency"`
	Country string `xml:"Country"`
}

// GrantList contains multiple grants.
type GrantList struct {
	Grants     []Grant `xml:"Grant"`
	CompleteYN string  `xml:"CompleteYN,attr"`
}

// UnknownElement captures unparsed or unexpected XML.
type UnknownElement struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

// ItemList represents a named list of items.
type ItemList struct {
	Items    []string `xml:"Item"`
	ListType string   `xml:"ListType,attr"`
}

// Journal holds metadata about the journal.
type Journal struct {
	ISSN            string `xml:"ISSN"`
	JournalIssue    string `xml:"JournalIssue"`
	Title           string `xml:"Title"`
	ISOAbbreviation string `xml:"ISOAbbreviation"`
}

// Keyword is a keyword term.
type Keyword struct {
	Text string `xml:",chardata"`
}

// Reference is a cited publication.
type Reference struct {
	Citation      string        `xml:"Citation"`
	ArticleIdList ArticleIdList `xml:"ArticleIdList"`
}

// Chemical holds chemical tag information.
type Chemical struct {
	RegistryNumber  string `xml:"RegistryNumber"`
	NameOfSubstance string `xml:"NameOfSubstance"`
}

// LocationLabel provides structural location (e.g., "Chapter 2").
type LocationLabel struct {
	Type string `xml:"Type,attr"`
}

// QualifierName is a MeSH qualifier.
type QualifierName struct {
	Text         string `xml:",chardata"`
	UI           string `xml:"UI,attr"`
	MajorTopicYN string `xml:"MajorTopicYN,attr"`
}

// MeshHeading represents a subject term.
type MeshHeading struct {
	DescriptorName string          `xml:"DescriptorName"`
	Qualifiers     []QualifierName `xml:"QualifierName"`
}

// MeshHeadingList contains all MeSH headings.
type MeshHeadingList struct {
	MeshHeadings []MeshHeading `xml:"MeshHeading"`
}

// Object holds supplementary identifiers.
type Object struct {
	Param string `xml:"Param"`
	Type  string `xml:"Type,attr"`
}

// PubMedPubDate holds timestamped metadata (e.g., publication or update).
type PubMedPubDate struct {
	Year      string `xml:"Year"`
	Month     string `xml:"Month"`
	Day       string `xml:"Day"`
	Hour      string `xml:"Hour"`
	Minute    string `xml:"Minute"`
	PubStatus string `xml:"PubStatus,attr"`
}

// SupplMeshList contains supplemental MeSH terms.
type SupplMeshList struct {
	SupplMeshNames []string `xml:"SupplMeshName"`
}

// Article contains core article metadata.
type Article struct {
	Journal             Journal           `xml:"Journal"`
	ArticleTitle        string            `xml:"ArticleTitle"`
	Abstract            Abstract          `xml:"Abstract"`
	AuthorList          []Author          `xml:"AuthorList>Author"`
	Language            []string          `xml:"Language"`
	PublicationTypeList []PublicationType `xml:"PublicationTypeList>PublicationType"`
	ArticleDate         string            `xml:"ArticleDate"`
}

// Author contains contributor metadata.
type Author struct {
	LastName        string            `xml:"LastName"`
	ForeName        string            `xml:"ForeName"`
	Initials        string            `xml:"Initials"`
	Suffix          string            `xml:"Suffix"`
	CollectiveName  string            `xml:"CollectiveName"`
	AffiliationInfo []AffiliationInfo `xml:"AffiliationInfo"`
	ValidYN         string            `xml:"ValidYN,attr"`
}

// PublicationType describes the article type (e.g., "Review").
type PublicationType struct {
	Text string `xml:",chardata"`
	UI   string `xml:"UI,attr"`
}
