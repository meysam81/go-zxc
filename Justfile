vendor:
  #!/bin/bash

  find vendor/ -type f \( -name '*.h' -or -name '*.c' \) ! -path '*/cli/*' ! -path '*/tests/*' ! -path '*/build/*' -exec cp {} internal/czxc/ \;

  for file in internal/czxc/*.c internal/czxc/*.h; do
    sed -i 's|#include "../../include/|#include "|g' "$file"
  done

build:
  go build ./...


test:
  go test -v ./...

bench:
  go test -bench=. -benchmem ./...
