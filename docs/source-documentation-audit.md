# Source documentation audit

[Documentation index](README.md)

## Scope and reproduction

Documentation-only audit of base `be75eb5515846dbc710cbd95f1675aa8e6395a49`.
One root module (`github.com/faustbrian/go-merkle-tree`), one public package
(`merkletree`), external tests (`merkletree_test`), and the build-excluded
fixture producer (`main` in `testdata`) are in scope. No interfaces or nested
modules are declared. Reproduce the inventory with `git ls-files '*.go'`,
`git ls-files '*go.mod'`, and `golib source-docs inventory --json`; inspect
`modules.json`, `packages.json`, `.golib.yaml`, and `verification/package.mk`
for catalog and native command ownership. Inventory includes tests and the
ignored-by-normal-build generator, not only Go's default package selection.

## Initial gaps and selected work

Initial structural inventory reported 26 undocumented exported fields across
resource policies and ResourceError. Existing constructors, proof identities,
copying accessors, error sentinels, and profile descriptions already have
substantial documentation. Semantic gaps: verifiers use proof-embedded roots
rather than separately supplied trusted roots; parse success is not authority
or proof verification; limits describe particular accounting measures, not a
process-wide heap cap; cancellation is cooperative, not preemptive; mutable
builder copies alias storage. Fuzz oracles and benchmark workload/timing
boundaries need adjacent rationale. Repair those contracts, not declaration
counts or private fixture fields with no consumer contract.

Owner: this repository's source documentation. Direct consumers are root/proof
construction, hostile-input decoders, immutable snapshot users and externally
synchronized builders. Preserve all executable source and specification
history. Tier A checks: native formatting, parsed structure/token/directive
equivalence to base, rendered root documentation, native docs/examples target,
tracked-file documentation inventory, local links, changed-content secrets,
and complete final diff. No unchanged runtime, fuzz, race, conformance,
benchmark or release campaign is required. The API baseline is `apidiff`,
`api/baseline.txt`; no native prose snapshot generation exists. Signatures and
that baseline remain unchanged.

## Investigated rationale and boundaries

- Both version-1 profiles use leaf domain 0x00, branch domain 0x01 and the
  largest-lower-power-of-two split, but retain distinct complete identities.
  Root digest equality alone does not authorize profile conversion. The
  package-owned MTRE frame is not Certificate Transparency protocol framing.
- Verification checks the proof's bound identity and relationship. Callers
  authenticate the complete requested root (both roots for consistency), plus
  the requested leaf index/set. It establishes no freshness, truth,
  authorization, uniqueness, availability or confidentiality.
- RawLeaf constructors copy; Bytes and proof vector accessors copy. Immutable
  values share only immutable internal storage. Builder and RootBuilder need
  external synchronization, must not be copied for independent mutation, and
  launch no goroutines. Batch state is published after local computation;
  an observed error leaves it unchanged. Cancellation may arrive after the
  last checkpoint; hashing, sorting, allocation and bulk copies are not
  preempted by context cancellation.
- Checked framing arithmetic uses multiplication/carry/native-int guards;
  raw-byte accumulation uses subtraction preflight and saturating diagnostics.
  Exact bounds are accepted. Retained snapshots have 2*n-1 nodes when nonempty;
  construction validates native-int capacity. Encoded node footprint accounting
  is not a measurement of total live process memory.
- Snapshot parsing checks canonical postorder graph/shape, branch digests and
  root. It cannot recover raw leaves or authenticate their cumulative lengths.
  Resume requires separately trusted byte accounting AND caller-authenticated
  root identity; old per-leaf bounds cannot be rechecked without old raw bytes.
- Inclusion paths are leaf-to-root; consistency nodes follow RFC SUBPROOF;
  multi-inclusion uses package-defined sorted unique indexes and left-to-right
  unselected-subtree frontier. Generation sorts a copy; verification leaves
  must follow returned index order. Empty-to-nonempty consistency is rejected;
  equal-sized proofs have zero nodes and matching digests.
- Proof parsers check framing/count/shape, not all hash relationships. Root
  parsing accepts an opaque digest for nonempty trees. Marshal operations
  return independent bytes but have no context or caller size budget; limits
  should be enforced on ingestion/construction first.

Evidence: source and deterministic assertions in the tracked files below;
`SECURITY.md`, `docs/security-threat-model.md`, `docs/encoding.md`, and resolved
MERKLETREE-DEC-001 through MERKLETREE-DEC-010 in the decision register. Pinned
RFC 9162 source identity and fixture checksum live in
`specification/manifest.tsv`. The fixture producer uses
transparency-dev/merkle v0.0.2 at revision
`036047b5d2f7faf3b1ee643d391e60fe5b1defcf`; fixture_test checks its identity and
checksum. Peers and generated vectors are evidence, not replacement authority.
No specification policy or historical record is changed.

## Variants, generation and assurance sources

The only build constraint is `//go:build ignore` on the fixture producer;
there are no platform suffix variants, generated Go outputs or go:generate
commands. The generator's documented explicit go run operation writes JSON
to stdout; testdata/README.md owns fixture regeneration/checksum procedure.
Existing security suppression rationales are retained verbatim. No unfinished
implementation markers or saved active-fuzz crash corpora are tracked.

Nine fuzz targets: bounded determinism/profile parity in FuzzComputeRoot;
mutated inclusion identity/path/leaf with classified-error oracle; consistency
and multi-proof mutation targets with panic-safety (not acceptance) oracles;
five byte-stream decoders with canonical re-encoding on successful parse.
Fourteen local benchmark entry points plus one peer-comparison entry measure
construction, growing mutable append state, proof generation/verification,
encoding and snapshot restore. Read timing boundaries and peer conventions;
no universal throughput/security equivalence is implied. Seven executable
examples cover roots, retained/streaming builders, three proof operations and
snapshot restore. Custom checkpoint contexts simulate cancellation at each
Err call deterministically, not asynchronous scheduling or deadline latency.
Private mutation helpers clone proof vectors to isolate corrupt fixtures;
independent recursive oracles and exhaustive small-prefix/subset tests
separate expected paths from the production traversal.

## Tracked Go coverage

Every listed file was audited for declarations, contracts, private algorithm
and fixture rationale, assertions, ownership, limits and applicable variants.
The native structural check also parses every file, including the generator.

| File | Role / audit boundary |
| --- | --- |
| `benchmark_test.go` | Native workload, growing append state, setup and timed regions |
| `builder.go` | Retained mutable frontier; atomic publication and copied snapshots |
| `builder_internal_test.go` | Internal corruption, limits and cancellation assertions |
| `builder_test.go` | Public contract, boundaries, independent oracle and ownership assertions |
| `comparison_benchmark_test.go` | Peer conventions and ownership differ; native construction tracks |
| `consistency.go` | Prefix identity, empty/equal cases, RFC path and verification |
| `consistency_fuzz_test.go` | Mutated identity/path, panic-safety oracle |
| `consistency_internal_test.go` | Internal corruption, limits and cancellation assertions |
| `consistency_test.go` | Public contract, boundaries, independent oracle and ownership assertions |
| `differential_test.go` | Pinned peer roots and paths; independent verification |
| `doc.go` | Canonical package, ownership, cancellation and trust contract |
| `encoding.go` | Strict root framing; structure versus authority |
| `encoding_fuzz_test.go` | Five decoder byte-stream/canonical re-encoding oracles |
| `encoding_internal_test.go` | Internal corruption, limits and cancellation assertions |
| `encoding_test.go` | Public contract, boundaries, independent oracle and ownership assertions |
| `errors.go` | Sentinels and bounded metadata-only resource errors |
| `example_test.go` | Seven executable public workflows |
| `fixture_test.go` | Checksummed reference corpus identity and vectors |
| `fuzz_test.go` | Bounded root determinism and profile digest parity |
| `inclusion.go` | Snapshot ownership and inclusion generation/verification |
| `inclusion_fuzz_test.go` | Mutated identity/path/leaf, classified errors and panic safety |
| `inclusion_internal_test.go` | Internal corruption, limits and cancellation assertions |
| `inclusion_test.go` | Public contract, boundaries, independent oracle and ownership assertions |
| `multi.go` | Sorted copied indexes, frontier minimality, leaf-order verification |
| `multi_fuzz_test.go` | Mutated indexes/identity/path/leaf, panic-safety oracle |
| `multi_internal_test.go` | Internal corruption, limits and cancellation assertions |
| `multi_test.go` | Public contract, boundaries, independent oracle and ownership assertions |
| `profile.go` | Versioned profiles, algorithm registry distinction |
| `proof_encoding.go` | Count/framing arithmetic, operation identity, decode versus verify |
| `proof_encoding_internal_test.go` | Internal corruption, limits and cancellation assertions |
| `proof_encoding_test.go` | Public contract, boundaries, independent oracle and ownership assertions |
| `root.go` | Root computation, leaf limits, checked sums and split/frontier |
| `root_internal_test.go` | Internal corruption, limits and cancellation assertions |
| `root_test.go` | Public contract, boundaries, independent oracle and ownership assertions |
| `snapshot_encoding.go` | Postorder graph integrity, footprint limits, trusted resume accounting |
| `snapshot_encoding_internal_test.go` | Internal corruption, limits and cancellation assertions |
| `snapshot_encoding_test.go` | Public contract, boundaries, independent oracle and ownership assertions |
| `stream.go` | Logarithmic mutable frontier, no proofs or raw-byte retention |
| `stream_internal_test.go` | Internal corruption, limits and cancellation assertions |
| `stream_test.go` | Public contract, boundaries, independent oracle and ownership assertions |
| `testdata/generate_vectors.go` | Build-excluded pinned reference JSON producer |
| `types.go` | RawLeaf/Digest/Root zero values, copied bytes and complete identity |

## Disposition

Initial inventory: 42 tracked Go files (13 production, 28 test/benchmark/fuzz/
example files, one fixture producer). Every exported declaration, method and
public struct field is reviewed; there are no exported interface methods.
Native parsed inventory covers 153 public declarations/members, including
all 26 public resource-policy/error fields. All 42 files preserve their AST,
non-comment token stream and compiler/security directives against the base.
Go 1.27.2 formatting and full root Go documentation passed; the native
`make -f verification/package.mk docs` target ran all seven examples and
rendered package documentation. Tracked-file source documentation checking
has zero findings. Local Markdown links, pinned fixture identity/checksum,
changed-content secret scanning and full-diff inspection passed. No API prose
snapshot is configured; the apidiff baseline and executable APIs are untouched.
All task-owned Go caches and disposable check artifacts are cleaned.
No unresolved source-policy question or newly demonstrated runtime defect was
identified. Application root authority, old raw-leaf accounting, resource
budget choice and storage durability intentionally remain caller-owned; they
are not missing package guarantees. Broader runtime and hosted CI boundaries
are not proven by this documentation-only pass.

## Objective Enforcement

The qualified public Tools v2.4.0 binary and immutable reusable CI source
provide the strict tracked-source check through `make source-docs-check`,
the existing local check/CI routes and Required aggregation. Release and
selected security/assurance routes retain their existing controls. This
structural status supplements the rendered semantic audit; it does not
certify prose meaning, generator execution or runtime conformance.

Continue with [Adoption and migration guide](adoption.md).
