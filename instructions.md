# Introduction

USFM is a markdown format for scriptures. Its defined here: https://docs.usfm.bible/usfm/3.1/index.html. There is a javascript parser here: https://github.com/Bridgeconn/usfm-grammar?tab=readme-ov-file.
I want to create a golang parser that can serve as a standalone quick binary parser and a Go library for applications wanting to parse such files.

# Specs

The application should have 2 aspects:
- A standalone commandline interface. Use cobra for this. Use sensible flags. I want to input either a directory with SFM files or a single file. The output should be either a *.txt file, a JSON file (USFM-JSON, USJ)), a pdf, or just a TSV showing the contents. There should be a quiet, verbose, and normal output mode.
- A Go library that can parse such files. Output should be as a JSON struct. Create a new struct type for this.

I favor standard library functions, readability, and documentation. Be generous on comments.
Don't make references to AI or Claude.
Create unit tests.

bsb_usfm includes a sample directory to use for testing and to post examples on the README.

# Ansulary elements

* Make a git and jujutsu repository (using isma@arenzana.org). Create a .gitignore to go with it.
* Create github workflow to build the binary and tests. Builds on master branch should post a new version with the binary. Use goreleaser for build assistance. Do security scans with grype on the workflow.* 
* Create a Makefile to build the software
* Create a README, MIT-license, CHANGELOG
* Use semantic versioning starting with 0.0.1
