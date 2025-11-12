#!/bin/bash
set -e

for file in "$(go env GOROOT)"/src/os/stat_{linux,netbsd,openbsd,freebsd,darwin}.go; do
    cp "$file" .
    name="$(basename "$file")"
    sed -i "./${name}" -E -e 's|package os|package mmap|' \
                          -e 's|internal/filepathlite|path|' \
                          -e 's|filepathlite|path|g' \
                          -e 's| FileMode| os.FileMode|g' \
                          -e 's| Mode| os.Mode|g' \
                          -e 's| FileInfo| os.FileInfo|g'
    gofmt -w "./${name}"; goimports -w "./${name}"
done
