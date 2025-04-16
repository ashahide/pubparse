package xmlTools

import "encoding/xml"

// PubmedArticleSet is the root of a regular PubMed XML file.
// It contains a list of articles.
type PubmedArticleSet struct {
	PubmedArticles []PubmedArticle `xml:"PubmedArticle"`
}

// PubmedArticle represents one article in the PubMed XML.
// It includes citation details and PubMed-specific metadata.
// UnknownElement captures any tags that don't map to a known field.
type PubmedArticle struct {
	MedlineCitation MedlineCitation  `xml:"MedlineCitation"` // main article metadata
	PubmedData      PubmedData       `xml:"PubmedData"`      // publication + reference info
	Unknown         []UnknownElement `xml:",any"`            // fallback for unmapped XML
}

// PubmedBookArticleSet is the root of a PubMed Books XML file.
// Similar to PubmedArticleSet, but for book chapters or books.
type PubmedBookArticleSet struct {
	PubmedBookArticles []PubmedBookArticle `xml:"PubmedBookArticle"`
}

// PubmedBookArticle represents one book chapter/article.
type PubmedBookArticle struct {
	BookDocument   BookDocument   `xml:"BookDocument"`
	PubmedBookData PubmedBookData `xml:"PubmedBookData"`
}

// BookDocument holds the metadata about the book section or article.
type BookDocument struct {
	PMID             string        `xml:"PMID"`
	ArticleIdList    ArticleIdList `xml:"ArticleIdList"` // multiple IDs (PubMed, DOI, etc.)
	Book             string        `xml:"Book"`          // name of the book
	ArticleTitle     string        `xml:"ArticleTitle"`
	Abstract         Abstract      `xml:"Abstract"`
	KeywordList      []Keyword     `xml:"KeywordList>Keyword"`
	GrantList        GrantList     `xml:"GrantList"`
	ReferenceList    []Reference   `xml:"ReferenceList>Reference"`
	PublicationType  string        `xml:"PublicationType"`
	InvestigatorList string        `xml:"InvestigatorList"`
	ContributionDate string        `xml:"ContributionDate"`
	DateRevised      string        `xml:"DateRevised"`
	ItemList         ItemList      `xml:"ItemList"`
	LocationLabel    string        `xml:"LocationLabel"`
}

// MedlineCitation contains bibliographic metadata and tagging.
type MedlineCitation struct {
	PMID                    string          `xml:"PMID"`
	DateCompleted           PubMedPubDate   `xml:"DateCompleted"`
	DateRevised             PubMedPubDate   `xml:"DateRevised"`
	Article                 Article         `xml:"Article"` // actual article contents
	MedlineJournalInfo      Journal         `xml:"MedlineJournalInfo"`
	ChemicalList            []Chemical      `xml:"ChemicalList>Chemical"`
	SupplMeshList           SupplMeshList   `xml:"SupplMeshList"`
	CitationSubset          string          `xml:"CitationSubset"`
	CommentsCorrectionsList string          `xml:"CommentsCorrectionsList"`
	GeneSymbolList          GeneSymbolList  `xml:"GeneSymbolList"`
	MeshHeadingList         MeshHeadingList `xml:"MeshHeadingList"` // subject tags
	NumberOfReferences      string          `xml:"NumberOfReferences"`
	PersonalNameSubjectList string          `xml:"PersonalNameSubjectList"`
	OtherID                 string          `xml:"OtherID"`
	OtherAbstract           Abstract        `xml:"OtherAbstract"`
	KeywordList             []Keyword       `xml:"KeywordList>Keyword"`
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

// PubmedBookData holds metadata for books, like article IDs or objects.
type PubmedBookData struct {
	History           string        `xml:"History"`
	PublicationStatus string        `xml:"PublicationStatus"`
	ArticleIdList     ArticleIdList `xml:"ArticleIdList"`
	ObjectList        []Object      `xml:"ObjectList>Object"`
}

// PubmedData contains reference lists and publication metadata.
type PubmedData struct {
	History           []PubMedPubDate `xml:"History>PubMedPubDate"` // multiple status timestamps
	PublicationStatus string          `xml:"PublicationStatus"`
	ArticleIdList     ArticleIdList   `xml:"ArticleIdList"`
	ObjectList        []Object        `xml:"ObjectList>Object"`
	ReferenceList     []Reference     `xml:"ReferenceList>Reference"`
}

// Abstract represents the abstract of the article.
type Abstract struct {
	AbstractText         string `xml:"AbstractText"`
	CopyrightInformation string `xml:"CopyrightInformation"`
}

// AffiliationInfo describes the author's institutional affiliation.
type AffiliationInfo struct {
	Affiliation string `xml:"Affiliation"`
	Identifier  string `xml:"Identifier"`
}

// ArticleId represents one ID (e.g., DOI, PubMed ID, PII).
type ArticleId struct {
	ID     string `xml:",chardata"`
	IdType string `xml:"IdType,attr"`
}

// ArticleIdList is a list of ArticleIds.
type ArticleIdList struct {
	ArticleIds []ArticleId `xml:"ArticleId"`
}

type Day struct{} // Not used, but placeholder if needed

// GeneSymbolList contains genetic symbols associated with the paper.
type GeneSymbolList struct {
	GeneSymbols []string `xml:"GeneSymbol"`
}

// Grant contains funding information.
type Grant struct {
	GrantID string `xml:"GrantID"`
	Acronym string `xml:"Acronym"`
	Agency  string `xml:"Agency"`
	Country string `xml:"Country"`
}

// GrantList holds multiple Grant entries.
type GrantList struct {
	Grants     []Grant `xml:"Grant"`
	CompleteYN string  `xml:"CompleteYN,attr"`
}

// UnknownElement catches any unparsed XML.
type UnknownElement struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

// ItemList is a named list of generic items.
type ItemList struct {
	Items    []string `xml:"Item"`
	ListType string   `xml:"ListType,attr"`
}

// Journal holds journal-level metadata.
type Journal struct {
	ISSN            string `xml:"ISSN"`
	JournalIssue    string `xml:"JournalIssue"`
	Title           string `xml:"Title"`
	ISOAbbreviation string `xml:"ISOAbbreviation"`
}

// Keyword represents a tag/term assigned to the article.
type Keyword struct {
	Text string `xml:",chardata"`
}

// Reference holds the citation of a referenced article.
type Reference struct {
	Citation      string        `xml:"Citation"`
	ArticleIdList ArticleIdList `xml:"ArticleIdList"`
}

// Chemical is used for tagging chemicals related to the article.
type Chemical struct {
	RegistryNumber  string `xml:"RegistryNumber"`
	NameOfSubstance string `xml:"NameOfSubstance"`
}

// LocationLabel is used for book/article structural info like "chapter".
type LocationLabel struct {
	Type string `xml:"Type,attr"`
}

// QualifierName is used in MeSH tagging to add context to a descriptor.
type QualifierName struct {
	Text         string `xml:",chardata"`
	UI           string `xml:"UI,attr"`
	MajorTopicYN string `xml:"MajorTopicYN,attr"`
}

// MeshHeading links descriptors with qualifiers to define subjects.
type MeshHeading struct {
	DescriptorName string          `xml:"DescriptorName"`
	Qualifiers     []QualifierName `xml:"QualifierName"`
}

// MeshHeadingList holds all MeSH headings for the article.
type MeshHeadingList struct {
	MeshHeadings []MeshHeading `xml:"MeshHeading"`
}

// Object is used for extra identifiers (datasets, repositories).
type Object struct {
	Param string `xml:"Param"`
	Type  string `xml:"Type,attr"`
}

// PubMedPubDate stores a date + time + status like "published", "entrez", etc.
type PubMedPubDate struct {
	Year      string `xml:"Year"`
	Month     string `xml:"Month"`
	Day       string `xml:"Day"`
	Hour      string `xml:"Hour"`
	Minute    string `xml:"Minute"`
	PubStatus string `xml:"PubStatus,attr"`
}

// SupplMeshList is a list of supplementary MeSH terms.
type SupplMeshList struct {
	SupplMeshNames []string `xml:"SupplMeshName"`
}

// Article holds core metadata of the article itself.
type Article struct {
	Journal             Journal           `xml:"Journal"`
	ArticleTitle        string            `xml:"ArticleTitle"`
	Abstract            Abstract          `xml:"Abstract"`
	AuthorList          []Author          `xml:"AuthorList>Author"`
	Language            []string          `xml:"Language"`
	PublicationTypeList []PublicationType `xml:"PublicationTypeList>PublicationType"`
	ArticleDate         string            `xml:"ArticleDate"`
}

// Author contains information about each author listed.
type Author struct {
	LastName        string            `xml:"LastName"`
	ForeName        string            `xml:"ForeName"`
	Initials        string            `xml:"Initials"`
	Suffix          string            `xml:"Suffix"`
	CollectiveName  string            `xml:"CollectiveName"`
	AffiliationInfo []AffiliationInfo `xml:"AffiliationInfo"`
	ValidYN         string            `xml:"ValidYN,attr"`
}

// PublicationType describes how the article was classified (e.g., Review, Guideline).
type PublicationType struct {
	Text string `xml:",chardata"`
	UI   string `xml:"UI,attr"`
}
