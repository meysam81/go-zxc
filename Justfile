vendor:
  #!/bin/bash

  find vendor/ -type f \( -name '*.h' -or -name '*.c' \) ! -path '*/cli/*' ! -path '*/tests/*' ! -path '*/build/*' -exec cp {} internal/czxc/ \;

build:
  go build ./...


test:
  go test -v ./...

bench:
  go test -bench=. -benchmem ./...
