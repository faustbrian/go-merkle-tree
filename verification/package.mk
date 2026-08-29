GO ?= go
BENCH_TIME ?= 100ms
FUZZ_TIME ?= 10000x

.PHONY: benchmark conformance docs fuzz

benchmark:
	$(GO) test -run '^$$' -bench '^Benchmark' \
		-benchmem -benchtime="$(BENCH_TIME)"

conformance:
	$(GO) test -run \
		'^(TestComputeRootMatchesRFC9162TreeHash|TestRFC9162InclusionProofMatchesIndependentAuditPaths|TestConsistencyProofMatchesRFC9162Examples|TestRFC9162PinnedReferenceFixture|TestRFC9162MatchesTransparencyDevMerkle)$$' \
		-count=1 .

docs:
	$(GO) test -run '^Example' -count=1 ./...
	$(GO) list -f '{{if .GoFiles}}{{.ImportPath}}{{end}}' ./... | \
		xargs -n 1 $(GO) doc >/dev/null

fuzz:
	$(GO) test -run '^$$' -fuzz '^FuzzComputeRoot$$' \
		-fuzztime="$(FUZZ_TIME)"
	$(GO) test -run '^$$' -fuzz '^FuzzVerifyInclusion$$' \
		-fuzztime="$(FUZZ_TIME)"
	$(GO) test -run '^$$' -fuzz '^FuzzVerifyConsistency$$' \
		-fuzztime="$(FUZZ_TIME)"
	$(GO) test -run '^$$' -fuzz '^FuzzVerifyMultiInclusion$$' \
		-fuzztime="$(FUZZ_TIME)"
	$(GO) test -run '^$$' -fuzz '^FuzzParseRoot$$' \
		-fuzztime="$(FUZZ_TIME)"
	$(GO) test -run '^$$' -fuzz '^FuzzParseInclusionProof$$' \
		-fuzztime="$(FUZZ_TIME)"
	$(GO) test -run '^$$' -fuzz '^FuzzParseConsistencyProof$$' \
		-fuzztime="$(FUZZ_TIME)"
	$(GO) test -run '^$$' -fuzz '^FuzzParseMultiInclusionProof$$' \
		-fuzztime="$(FUZZ_TIME)"
	$(GO) test -run '^$$' -fuzz '^FuzzParseSnapshot$$' \
		-fuzztime="$(FUZZ_TIME)"
