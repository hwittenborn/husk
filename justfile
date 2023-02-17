fmt:
    #!/usr/bin/env bash
    cd "$(git rev-parse --show-toplevel)"
    cargo fmt
    cd src/go && go fmt
