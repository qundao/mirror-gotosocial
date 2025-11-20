#!/bin/bash
set -e

for file in "$(go env GOROOT)"/src/os/stat_{linux,netbsd,openbsd,freebsd,darwin}.go; do
    out="./$(basename "$file")"
    modtime=$(sed "$file" -n -e '/fs.modTime = /{ s|.*fs.modTime = ||; s|fs.||; p; }')
    echo > "$out" # reset file
    echo "package disk"                                  >> "$out"
    echo ""                                              >> "$out"
    echo "import ("                                      >> "$out"
    echo "    \"syscall\""                               >> "$out"
    echo "    \"time\""                                  >> "$out"
    echo ")"                                             >> "$out"
    echo ""                                              >> "$out"
    echo "func modtime(sys *syscall.Stat_t) time.Time {" >> "$out"
    echo "    return ${modtime}"                         >> "$out"
    echo "}"                                             >> "$out"
    gofmt -w "$out"
done
