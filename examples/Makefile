NO_COLOR=\033[0m
OK_COLOR=\033[32;01m

ALL_GO_FILES = $(wildcard *.go)
ALL_FILES = $(patsubst %.go,%,$(ALL_GO_FILES))

all: format vet lint
	@$(foreach file,$(ALL_GO_FILES), echo "$(OK_COLOR)==> Building $(file) $(NO_COLOR)" && go build $(file);)
clean:
	@$(foreach file,$(ALL_FILES),rm -f $(file);)

format:
	@echo "$(OK_COLOR)==> Formatting the code $(NO_COLOR)"
	@gofmt -s -w *.go
	@goimports -w *.go

vet:
	@echo "$(OK_COLOR)==> Running go vet $(NO_COLOR)"
	@`which go` vet .

lint:
	@echo "$(OK_COLOR)==> Running golint $(NO_COLOR)"
	@`which golint` .

.PHONY: all clean format vet lint
