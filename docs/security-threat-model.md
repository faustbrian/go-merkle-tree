# Security threat model: go-merkle-tree

**Model version:** 1.0 (2026-09-30).

**Source baseline:** `bca41287d7c0ffb07d992b1e241aaab85dbc9539` on `main`.

**Scope:** the root module and its complete root, streaming builder, snapshot,
inclusion, consistency, multi-inclusion, and binary-encoding families. This is
a source assessment, not certification of a deployment or a public release.
The [security policy](../SECURITY.md) defines private reporting and application
responsibilities; [specification decisions](specification-decisions.md) define
the intentional differences from RFC 9162.

## Assets, attacker control, and consumers

Assets are commitment integrity, complete root identity, bounded CPU and
memory, immutable snapshot/proof ownership, and atomic builder state. Leaves
may contain confidential application data; hashing is not encryption.

An attacker can choose raw leaf bytes, proof and snapshot encodings, tree
sizes, counts, indexes, node references, byte-accounting metadata, and replayed
or stale roots. Trusted caller policy supplies limits, contexts, profiles,
authenticated root identities, and independently trusted resumption accounting.
The library does not perform authentication or authorization of a publisher.

Owned public consumers include the executable examples and external-package
tests in this repository. The only direct source dependency found in the
Golib sibling inventory is the blank import in
`go-library-tools/release/compatibility-consumer/consumer_test.go`, pinned to
v1.0.0. No runtime reverse consumer was found in that inventory. Applications
outside the inventory remain responsible for their own integration policies.

## Boundary and hostile-input matrix

| Boundary | Controls and named risks | Existing behavioral evidence |
| --- | --- | --- |
| Caller bytes to `RawLeaf` | `NewRawLeafWithLimit` rejects invalid or excessive size before copying. `NewRawLeaf` deliberately copies without imposing a limit, and `Bytes` copies already owned storage. Bound ingress before either allocating the input or using the convenience constructor. | `TestNewRawLeafWithLimitBoundsBeforeCopying`, `TestLeafAndDigestBytesNeverAliasCallerMemory` |
| Leaves to roots and streaming builders | Nonzero leaf-count, per-leaf, and total-byte policies are checked before hashing or derived allocation. Subtraction-based accounting avoids overflow. Root-only state is logarithmic. Builders publish batch state only after validation and hashing succeed. | `TestComputeRootRejectsInvalidConfigurationAndResourceClaims`, `TestRootBuilderResourceBoundsAndAtomicFailures`, `TestBuilderBatchAppendIsAtomic` |
| Leaves to retained snapshots | Construction limits and retained-node counts bound full-tree storage; integer representability is checked. Snapshot copies do not alias mutable builders. | `TestSnapshotRejectsRetainedNodeClaimsBeforeAllocation`, `TestSnapshotNodeCountChecksIntegerConversion`, `TestBuilderAppendSnapshotsMatchBatchConstruction` |
| Root/profile identity to cryptography | Only standard-library SHA-256 is supported. Leaf and branch prefixes are distinct. Profile, version, algorithm, tree size, and digest remain bound; identical digest bytes do not authorize profile conversion. | `TestProfileValidationRejectsPartiallyMatchingRFCIdentity`, `TestProfilesMakeRootConventionsExplicit`, `TestRFC9162PinnedReferenceFixture` |
| Inclusion and consistency proofs to verification | Count/depth/leaf limits precede work. Unique path length, compatible identities, and exact consumption reject missing/surplus paths. Equal-size consistency requires equal roots; zero-to-nonzero consistency is intentionally unsupported. Digest comparisons use `crypto/subtle`. | `TestVerifyInclusionRejectsMalformedAndNonVerifyingProofs`, `TestVerifyConsistencyRejectsMalformedAndNonVerifyingProofs`, `TestConsistencyPathLengthCoversUint64Boundary` |
| Multi-inclusion selection/frontier to verification | Selection count bounds copying/sorting; duplicate or out-of-range indexes fail. Leaf, aggregate-byte, frontier-count, and depth policies bound verification. Canonical sorted selection and exact frontier consumption prevent alternate representations. | `TestMultiInclusionProofCanonicalizesIndexesAndOwnsSlices`, `TestMultiProofExactLimits`, `TestVerifyMultiInclusionRejectsMalformedAndNonVerifyingProofs` |
| Binary objects to decoded values | Explicit object/profile/version/algorithm framing, exact sizes, overflow-checked vector arithmetic, encoded-byte and operation-specific count/depth limits precede derived allocation. Parsing a proof does not authenticate leaves or establish root trust. | `TestProofDecodersRejectTruncatedTrailingAndWrongOperation`, `TestProofEncodingRejectsResourceAndStructuralClaims`, `TestProofEncodingSizeArithmetic` |
| Persisted snapshots to restored state | Encoded bytes, leaves, total accounting, retained nodes, depth, node reads, and temporary node bytes are bounded. Exact canonical postorder topology, child sizes, branch hashes, and root validation reject cycles, reference bombs, corruption, and surplus nodes. `ResumeBuilder` requires independently trusted expected byte accounting. | `TestSnapshotPersistenceRejectsCorruptMetadataAndNodes`, `TestSnapshotPersistenceInputAndResourceLimits`, `TestResumeBuilderValidationCancellationAndLimits` |
| Cancellation and shared state | Context checks occur between leaves, merges, node copying/validation, and proof traversal. A bounded individual hash or sort is not interruptible. No locks or goroutines are created; callers synchronize mutable builders. Immutable values are safe for concurrent reads. | `TestBuilderAppendBatchAndSnapshotHonorEveryCancellationCheckpoint`, `TestRootBuilderHonorsEveryCancellationCheckpoint`, `TestSnapshotPersistenceCancellationDuringNodesAndValidation`, `TestMultiInclusionHonorsEveryCancellationCheckpoint` |
| Errors and diagnostics | Sentinel errors and `ResourceError` disclose categories and numeric resource metadata, not leaves, paths, or encoded payloads. No logging, panic recovery, diagnostic callback, or implicit environment/network/filesystem/process access exists in production APIs. | Source review of `errors.go` and production imports; ownership and resource-error assertions in `root_test.go` |

The tests above are an evidence inventory, not a claim that this document ran
them. Parser and verifier fuzz targets cover all object families, including
`FuzzParseSnapshot`; their input ceilings and canonical round-trip assertions
are separate from application admission controls.

Default policies are finite compatibility ceilings, not low-memory deployment
budgets: construction permits 16 MiB per leaf and 1 GiB of aggregate leaf
bytes; retained snapshots permit 2^19 leaves; decoders permit 32 MiB objects
or 64 MiB snapshots. Tighter budgets are needed for unauthenticated workloads.
Count/depth limits also bound structural validation and multi-proof sorting;
contexts do not promise interruption inside a single bounded primitive.

## Caller-owned boundaries and residual risks

These are conditional design limitations with explicit mitigation, not
acceptance of an unbounded attacker workload or an approval of a deployment.

| Residual risk | Owner and rationale | Mitigation and review condition |
| --- | --- | --- |
| Attacker data can consume memory before the library receives it, or through the unlimited copying constructor. | Application ingress owner; the library cannot undo a caller's already allocated input, and the convenience API supports trusted data. | Bound reads and batch accumulation upstream, use `NewRawLeafWithLimit`, and choose finite per-operation policies. Review on ingress, payload-size, constructor, or limit changes. |
| A valid proof under an attacker-selected or stale root proves no publisher authority, freshness, or semantic truth. | Application protocol/authentication owner; these require external trusted state. | Authenticate the complete requested identity and compare the proof's root before accepting it; enforce freshness and authorization separately. Review on publication, replay, or tenant-policy changes. |
| Snapshot byte accounting is not committed by the Merkle root; concurrent or partial writes can publish stale or unavailable state. | Persistence/publication owner; the library has no storage adapter or transaction owner. | Supply separately trusted accounting to `ResumeBuilder`, validate on read, write data durably before atomic publication, and use compare-and-swap for writers. Review on storage, recovery, accounting, or writer changes. |
| Caller-selected generous limits permit large bounded work; cancellation cannot interrupt individual hashes or sorting. | Application availability owner; finite bounds trade throughput and compatibility against latency and memory. | Measure and lower size/count limits, bound request concurrency and admission, and propagate deadlines. Review on traffic, capacity, or limit changes. |
| Mutable builder races can corrupt state; raw leaves and proofs can disclose application information through caller logs or dumps. | Application concurrency and observability owners; synchronization and telemetry are outside this package. | Serialize builder use, limit data lifetime, and redact payloads and generated failure corpora. Review on parallelism, tracing, or crash-reporting changes. |
| Compromised dependencies, fixtures, tools, actions, or release credentials can defeat source integrity. | Merkle maintainers and release owner; test-only peer dependencies and automation are supply-chain inputs even though production uses only the standard library. | Pin and review dependencies, fixtures, tooling, and actions; run the selected security gates and verify releases from clean public consumers. Review on dependency/tool/pin updates and every release. |

SSRF, path traversal, SQL injection, queue retries, decompression, credential
comparison, and callback/backend cancellation are not owned runtime boundaries:
the package accepts no URL, path, SQL, queue, compressed stream, credential,
`io.Reader`, plugin, or storage collaborator. Introducing any such adapter
requires a new threat-model boundary rather than inheriting this verdict.

## Supply chain, authority review, and release verdict

The owned workflow delegates to immutable go-library-tools v1.4.0 source and
fails its stable Required job when selected CI fails. There are no package-local
scanner suppressions in production source. The repository's shared secret-scan
configuration includes path-specific corpus exceptions, which must not be
expanded to cover application payloads or private credentials.

RFC 9162 text and the errata response were downloaded on 2026-09-30 and matched
both existing byte counts and SHA-256 pins in
[`specification/manifest.tsv`](../specification/manifest.tsv). Erratum 8670 is
Editorial, Held for Document Update: it removes the redundant `or fn is 0`
termination clause from section 2.1.3.2 step 5.b.ii and section 2.1.4.2 step
6.b.iii. Inclusion verification here uses recursive tree splitting rather than
that shift loop. Consistency verification rejects `sn == 0` before its merge
branch; entering that branch requires odd `fn` or `fn == sn`, so `fn` is nonzero
when its trailing-zero shift is computed. The clarification therefore requires
no behavior change for either inclusion or consistency decisions. The authority
response did not drift, and the 30-day review cadence is retained.

**Per-module verdict:** no confirmed runtime vulnerability was identified in
this bounded manual source audit. This documentation/authority-metadata batch
does not require a module release and does not establish completion of the
ecosystem security goal. Exact scanner results, hosted required-check state,
and any release/clean-consumer evidence belong in the coordinator's ledger;
neither existing test names nor a threat model substitute for those results.
