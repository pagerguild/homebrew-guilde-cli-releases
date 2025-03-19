@lint_and_test dir="./mkrelease":
    cd {{dir}} && go fmt ./...
    cd {{dir}} && golangci-lint run --fix ./...
    cd {{dir}} && go vet ./...
    cd {{dir}} && go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest --fix --test ./...
    cd {{dir}} && go test ./...

build version:
    rm public/* Formula/*
    cd mkrelease && go run . {{version}}  ../
