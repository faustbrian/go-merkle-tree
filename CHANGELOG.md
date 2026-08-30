# Changelog

## Unreleased

### Changed

- Replace copied repository-local verification tooling with the released,
  checksum-pinned `go-library-tools` v1.0.14 CLI and the immutable
  `0bbff11f25d74203018019fa5b26ae1443310fe7` workflow while preserving the
  package's source-owned conformance, benchmark, and mutation evidence.

### Added

- Add machine-enforced RFC 9162 source and errata monitoring, decision and
  conformance records, preserved decision history, change control, and
  maintained-peer differential classifications. See the
  [specification decision register](docs/specification-decisions.md).

### Specification Decisions

- MERKLETREE-DEC-001 sha256:fa3421258576aa3899d8782dd7e304111f5df49a99991077bf56607d145987b5
- MERKLETREE-DEC-002 sha256:5f24b34120614a05f9dad87961f65941c6317fc1bfe48d06fd9b854e0b7fcb5d
- MERKLETREE-DEC-003 sha256:4488606b3445e521712b3b60287623c05640f1cfdee6deac982138eaa1b14783
- MERKLETREE-DEC-004 sha256:7a474513e24c5867ac652498f3c058e4604341757afeb96f2b171f88a3afb0f6
- MERKLETREE-DEC-005 sha256:21d7a3ea1e852de9a13d45af0fbbc2e2e49e58b5db8399de7ab6d3e016f5d3d6
- MERKLETREE-DEC-006 sha256:861a972c86744474044c13a11d0fceb16019d9e687b2b6d65e380410bc9590f5
- MERKLETREE-DEC-007 sha256:ab62348a7244890d5a98b2d52286fa531e722799c35af86918fa99c97c33ed90
- MERKLETREE-DEC-008 sha256:858f53d0d73d548035a0464028cea2d192aefa560df297a7c28e066c934d8895
- MERKLETREE-DEC-009 sha256:43717310bf90946ab65f19048c1d0c534d0e034cd3d3af1a456006e6c0eae4eb
- MERKLETREE-DEC-010 sha256:fe09d5afe9f074364a5c999dcf5b778756de8d22c4094bfa5e2e01ebbd39e4c6

### Documentation

- Replace archived monorepo links and completed execution artifacts with a
  standalone, human-oriented documentation structure.

## 1.0.0 - 2026-08-25

### Changed

- Upgrade the cryptographic dependency set to its current secure releases.

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Link the package README to package-owned documentation.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-merkle-tree` identity while preserving its documented API and behavior.
- Link compatibility guidance directly to the canonical specification decision
  register.
- Pin the RFC 9162 reference corpus and record the security, resource,
  compatibility, and wire consequences of every profile and format decision.

### Added

- Explicit versioned identities for the canonical binary and RFC 9162
  profiles, with SHA-256 selected by stable algorithm identity.
- Bounded, cancellation-aware root construction from ordered raw leaves using
  RFC 9162 domain separation and non-power-of-two tree shape.
- Pre-copy raw-leaf byte limits for untrusted leaf ingestion.
- Immutable root identities that bind the profile, version, algorithm, digest,
  and exact tree size without aliasing caller memory.
- Immutable snapshots that retain node digests for deterministic logarithmic
  RFC 9162 inclusion-path generation without retaining raw leaf bytes.
- Independently verifiable inclusion proofs binding the complete operation
  identity, with typed malformed, unsupported, resource, and authentication
  failures.
- Caller-owned incremental builders with atomic append and batch-append
  operations whose immutable snapshots remain stable after later mutations.
- RFC 9162 consistency proof generation and independent verification binding
  both complete root identities, with bounded hostile-input traversal.
- Deterministic multi-inclusion proofs with canonical index ordering, minimal
  frontier nodes, explicit resource limits, and independent verification.
- Atomic streaming root construction that retains only a logarithmic digest
  frontier and never retains raw leaves or the full node tree.
- Versioned canonical binary encodings for roots and all proof operations,
  with strict structural decoding, operation-specific resource limits, and
  independently owned decoded state.
- Canonical persisted snapshots with complete node-integrity validation,
  cumulative byte accounting, explicit hostile-input limits, and independent
  mutable builder resumption.
- Pinned RFC 9162 reference fixtures and differential root, inclusion, and
  consistency evidence against `transparency-dev/merkle`.
- Reproducible native benchmark tracks for required ecosystem packages, with
  explicit semantic and ownership boundaries that prevent unlike speed claims.
- Adoption, migration, error recovery, hostile-input, structure-selection, and
  security guidance for the complete pre-v1 API.
