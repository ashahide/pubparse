package xml

type PubmedBookArticleSet struct {
	PubmedBookArticle string `xml:"PubmedBookArticle"`
}

type PubmedArticle struct {
	MedlineCitation string `xml:"MedlineCitation"`
	PubmedData      string `xml:"PubmedData"`
}

type PubmedBookArticle struct {
	BookDocument   string `xml:"BookDocument"`
	PubmedBookData string `xml:"PubmedBookData"`
}

type MedlineCitation struct {
	PMID                    string `xml:"PMID"`
	DateCompleted           string `xml:"DateCompleted"`
	DateRevised             string `xml:"DateRevised"`
	Article                 string `xml:"Article"`
	MedlineJournalInfo      string `xml:"MedlineJournalInfo"`
	ChemicalList            string `xml:"ChemicalList"`
	SupplMeshList           string `xml:"SupplMeshList"`
	CitationSubset          string `xml:"CitationSubset"`
	CommentsCorrectionsList string `xml:"CommentsCorrectionsList"`
	GeneSymbolList          string `xml:"GeneSymbolList"`
	MeshHeadingList         string `xml:"MeshHeadingList"`
	NumberOfReferences      string `xml:"NumberOfReferences"`
	PersonalNameSubjectList string `xml:"PersonalNameSubjectList"`
	OtherID                 string `xml:"OtherID"`
	OtherAbstract           string `xml:"OtherAbstract"`
	KeywordList             string `xml:"KeywordList"`
	CoiStatement            string `xml:"CoiStatement"`
	SpaceFlightMission      string `xml:"SpaceFlightMission"`
	InvestigatorList        string `xml:"InvestigatorList"`
	GeneralNote             string `xml:"GeneralNote"`
	Owner                   string `xml:"Owner,attr"`
	Status                  string `xml:"Status,attr"`
	Medline                 string `xml:"MEDLINE,attr"`
	VersionID               string `xml:"VersionID,attr"`
	VersionDate             string `xml:"VersionDate,attr"`
	IndexingMethod          string `xml:"IndexingMethod,attr"`
}

type PubmedBookData struct {
	History           string `xml:"History"`
	PublicationStatus string `xml:"PublicationStatus"`
	ArticleIdList     string `xml:"ArticleIdList"`
	ObjectList        string `xml:"ObjectList"`
}

type Abstract struct {
	AbstractText         string `xml:"AbstractText"`
	CopyrightInformation string `xml:"CopyrightInformation"`
}

type AffiliationInfo struct {
	Affiliation string `xml:"Affiliation"`
	Identifier  string `xml:"Identifier"`
}

type ArticleIdList struct {
	ArticleId string `xml:"ArticleId"`
}

type Day struct{}

type GeneSymbolList struct {
	GeneSymbol string `xml:"GeneSymbol"`
}

type Grant struct {
	GrantID string `xml:"GrantID"`
	Acronym string `xml:"Acronym"`
	Agency  string `xml:"Agency"`
	Country string `xml:"Country"`
}

type GrantList struct {
	Grant      string `xml:"Grant"`
	CompleteYN string `xml:"CompleteYN,attr"`
}

type Item struct{}

type ItemList struct {
	Item     string `xml:"Item"`
	ListType string `xml:"ListType,attr"`
}

type Journal struct {
	ISSN            string `xml:"ISSN"`
	JournalIssue    string `xml:"JournalIssue"`
	Title           string `xml:"Title"`
	ISOAbbreviation string `xml:"ISOAbbreviation"`
}

type LocationLabel struct {
	Type string `xml:"Type,attr"`
}

type MeshHeading struct {
	DescriptorName string `xml:"DescriptorName"`
	QualifierName  string `xml:"QualifierName"`
}

type MeshHeadingList struct {
	MeshHeading string `xml:"MeshHeading"`
}

type Object struct {
	Param string `xml:"Param"`
	Type  string `xml:"Type,attr"`
}

type PubMedPubDate struct {
	Year       string `xml:"Year"`
	Month      string `xml:"Month"`
	Day        string `xml:"Day"`
	Hour       string `xml:"Hour"`
	Minute     string `xml:"Minute"`
	Second     string `xml:"Second"`
	PubStatus  string `xml:"PubStatus,attr"`
	PPublish   string `xml:"ppublish,attr"`
	Retracted  string `xml:"retracted,attr"`
	PreMedline string `xml:"premedline,attr"`
}

type SupplMeshList struct {
	SupplMeshName string `xml:"SupplMeshName"`
}
