VERSION := $(shell git describe --tags)

LDFLAGS := -X 'github.com/flowswiss/cli/v2/internal/commands.Version=${VERSION}'
BUILD := go build -v -ldflags "${LDFLAGS}"

OUTDIR := ./bin

all: outdir ${OUTDIR}/flow ${OUTDIR}/cloudbit

.PHONY: flow
${OUTDIR}/flow:
	${BUILD} -o $@ ./cmd/flow

.PHONY: cloudbit
${OUTDIR}/cloudbit:
	${BUILD} -o $@ -tags cloudbit ./cmd/flow

.PHONY: outdir
outdir:
	@mkdir -p $@