#!/bin/bash

rm *.json *.xml

PMIDS="21193628 32387127 29282247 26247036 25537714 16444918"

for PMID in ${PMIDS}

do
    echo "Fetching XML for PMID: $PMID"
    # Use the PubMed API to fetch the XML data
    curl -s "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=$PMID&retmode=xml" \
        -o "$PMID.xml"
done


PMIDS="32897388,39270656"

echo "Fetching XML for PMIDs: $PMIDS"

# Use the PubMed API to fetch the XML data
curl -s "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=$PMIDS&retmode=xml" \
    -o "pubmed_${PMIDS//,/}.xml"