# RFC 9162 conformance map

The [specification decision register](../docs/specification-decisions.md) is
the human contract; `decisions.json`, `conformance.json`,
`decision-history.json`, and `monitoring.json` are its machine-enforced source,
evidence, history, and change-monitoring records.

`manifest.tsv` pins the December 2021 RFC text, its current errata response,
and `testdata/rfc9162-sha256-v1.json`. RFC 9162 publishes no official
machine-readable Merkle corpus, so the vector record pins the maintained
`transparency-dev/merkle` version, revision, module checksum, license, and
regeneration procedure documented in [`testdata/README.md`](../testdata/README.md).

| Decision | Source boundary | Conformance evidence | Maintained-peer differential status |
| --- | --- | --- | --- |
| MERKLETREE-DEC-001 | RFC 9162 section 2.1.1 | Root, profile, ownership, fixture, and fuzz checks | `transparency-dev/merkle` agrees on leaf and branch domains; package metadata is a deliberate policy difference. |
| MERKLETREE-DEC-002 | RFC 9162 section 2.1.1 | Empty root, encoding, persistence, fixture, and fuzz checks | `transparency-dev/merkle` agrees on SHA-256 of the empty string; package metadata is a deliberate policy difference. |
| MERKLETREE-DEC-003 | RFC 9162 section 2.1.1 | Shape, builder-prefix, fixture, and fuzz checks | Live roots agree for every differential corpus size; package profile identity is a deliberate policy difference. |
| MERKLETREE-DEC-004 | RFC 9162 section 2.1 omission | Profile and proof identity plus root-parser fuzz checks | Not assessed because complete package root identity has no peer wire contract. |
| MERKLETREE-DEC-005 | RFC 9162 sections 2.1.3.1-2.1.3.2 | Generation, verification, hostile input, fixture, and fuzz checks | Live audit paths and peer verification agree; transport representation is a deliberate policy difference. |
| MERKLETREE-DEC-006 | RFC 9162 sections 2.1.4.1-2.1.4.2 omission | RFC examples, edge cases, hostile input, fixture, and fuzz checks | Live strict-prefix nodes and peer verification agree; equal-size and zero-size handling is a deliberate policy difference. |
| MERKLETREE-DEC-007 | RFC 9162 section 2.1.3 omission | Exhaustive canonical subsets, hostile input, and fuzz checks | Not assessed because RFC 9162 has no multi-proof representation. |
| MERKLETREE-DEC-008 | Outside RFC 9162 | Atomicity, ownership, prefix equivalence, and cancellation checks | Not assessed because Go builder ownership is package policy. |
| MERKLETREE-DEC-009 | RFC 9162 section 2.1 omission | Canonical fixtures, malformed input, and parser fuzz checks | Not assessed because RFC 9162 has no standalone binary envelope. |
| MERKLETREE-DEC-010 | Outside RFC 9162 | Round-trip, corruption, recovery, limits, and parser fuzz checks | Not assessed because RFC 9162 has no snapshot persistence contract. |

The `conformance` operation runs the RFC root, inclusion, consistency, pinned
fixture, and live peer differential checks. Scheduled and pull-request CI also
runs the pinned tooling's offline register validator; scheduled source review
additionally downloads both monitored authorities and fails on byte changes.
