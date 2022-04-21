#!/usr/bin/bash

shopt -s globstar

for f in articles/**/*.md ; do
	echo "Reading $f"
	mdparser/mdparser $f mdparser/templates/ > $(dirname $f)/index.html
done
