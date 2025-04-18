#!/bin/bash

###############################################################################
# generate_test_data.sh
#
# This script prepares test input files for the `pubparse` tool by downloading:
#   - PubMed article XMLs using the NCBI E-utilities API
#   - PMC open-access full-text tarballs from the NCBI FTP server
#
# It extracts all XMLs into a consistent flat structure and cleans old files.
#
# Usage:
#   bash generate_test_data.sh
#
# Output directory structure:
#   test/data/
#   ├── test_pubmed/
#   │   └── xml/             # Individual and grouped PubMed XMLs
#   └── test_pmc/
#       ├── zipped_xml/      # Downloaded PMC .tar.gz files
#       └── xml/             # Extracted PMC XMLs (flattened, no subdirs)
###############################################################################

set -e  # Exit immediately if any command fails

echo ">>> Cleaning previous test data..."
rm -rf test/data/*

# Create necessary directory structure
echo ">>> Creating test directories..."
mkdir -p test/data/test_pubmed/xml
mkdir -p test/data/test_pmc/zipped_xml
mkdir -p test/data/test_pmc/xml

# -------------------------------
# Step 1: Fetch individual PubMed articles
# -------------------------------

echo ">>> Fetching individual PubMed XML files..."

# List of single PubMed PMIDs to download individually
PMIDS="21193628 32387127 29282247 26247036 25537714 16444918 30725926"

# Loop through and fetch each article as a separate XML file
for PMID in ${PMIDS}; do
    echo "Fetching PubMed XML for PMID: $PMID"
    curl -s "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=$PMID&retmode=xml" \
        -o "test/data/test_pubmed/xml/$PMID.xml"
done

# -------------------------------
# Step 2: Fetch bulk PubMed articles in one request
# -------------------------------

PMIDS="32897388,39270656"
echo ">>> Fetching bulk PubMed XML for PMIDs: $PMIDS"

# Download multiple articles into one file
curl -s "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=$PMIDS&retmode=xml" \
    -o "test/data/test_pubmed/xml/pubmed_${PMIDS//,/}.xml"

# -------------------------------
# Step 3: Download and extract PMC full-text XML tarballs
# -------------------------------

echo ">>> Downloading PMC archives..."

# Base URL for PMC open-access archives
FTP_BASE="https://ftp.ncbi.nlm.nih.gov/pub/pmc/oa_bulk"

# List of full tarball paths (include `/xml/` subfolder)
PMC_TARBALLS=(
    "oa_comm/xml/oa_comm_xml.PMC000xxxxxx.baseline.2024-12-18.tar.gz"
    "oa_other/xml/oa_other_xml.PMC000xxxxxx.baseline.2024-12-18.tar.gz"
)

# Directories for downloaded archives and extracted XMLs
DOWNLOAD_DIR="test/data/test_pmc/zipped_xml"
EXTRACT_DIR="test/data/test_pmc/xml"

# Loop through and download + extract each tarball
for TARBALL in "${PMC_TARBALLS[@]}"; do
    BASENAME=$(basename "$TARBALL")
    TARBALL_URL="${FTP_BASE}/${TARBALL}"
    DEST="${DOWNLOAD_DIR}/${BASENAME}"

    echo "Downloading: $BASENAME"
    if ! curl -s -f "$TARBALL_URL" -o "$DEST"; then
        echo "Failed to download $TARBALL_URL — skipping"
        continue
    fi

    echo "Extracting: $BASENAME"
    if tar -xzf "$DEST" -C "$EXTRACT_DIR" --strip-components=1; then
        echo "Extracted: $BASENAME"
        rm "$DEST"
    else
        echo "Failed to extract $BASENAME (invalid archive?)"
        break
    fi
done

echo "Test data generation complete."
