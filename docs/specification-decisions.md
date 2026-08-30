# Merkle-tree specification decisions

This register is the human-readable view of
[`specification/decisions.json`](../specification/decisions.json). RFC 9162
is authoritative only for the identified Merkle tree operations; package
formats and lifecycle policies are explicitly classified as package policy.

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in BCP 14
[RFC 2119](https://www.rfc-editor.org/rfc/rfc2119) and
[RFC 8174](https://www.rfc-editor.org/rfc/rfc8174) when, and only when, they
appear in all capitals, as shown here.

Statuses are `resolved`, `unresolved`, or `superseded`. Changes preserve old
digests in `specification/decision-history.json`, require compatibility and
changelog review, and never let peer behavior override the authoritative source.

## MERKLETREE-DEC-001: Leaf and branch domain separation

| Field | Decision |
| --- | --- |
| Status | resolved |
| Owner | merkle-tree maintainers |
| Classification | interoperability policy |
| Decision scope | normative |
| Specification | RFC 9162 |
| Exact version | RFC 9162 (December 2021) |
| Source authority | rfc9162-source |
| Authoritative URL | https://www.rfc-editor.org/rfc/rfc9162.txt |
| Section | 2.1.1 |
| Requirement strength | not specified |
| Issue | Raw leaves, leaf digests, and branch digests can be confused unless their domains and input types are explicit. |
| Credible interpretations | Hash raw leaves and branches without prefixes.<br>Accept caller-provided leaf digests.<br>Use the RFC 9162 leaf and branch prefixes over distinct input types. |
| Known peer behavior | transparency-dev/merkle v0.0.2 uses the RFC 9162 0x00 leaf and 0x01 branch domains. |
| Selected behavior | Both version-1 profiles hash RawLeaf values with the 0x00 prefix and branches with the 0x01 prefix over left and right digests; RawLeaf and Digest remain distinct types. |
| Normative rationale | This preserves the RFC algorithm, second-preimage domain separation, and an API boundary that prevents accidental double hashing. |
| Security consequences | Leaf and branch substitution and accidental digest reuse fail at the commitment or type boundary. |
| Resource consequences | Per-leaf and cumulative byte limits apply before hashing; domain separation adds one prefix byte per hash invocation. |
| Compatibility consequences | Unprefixed or pre-hashed trees are incompatible and must be rebuilt. |
| Wire consequences | Identical ordered raw leaves have RFC-compatible SHA-256 digests; package envelopes still carry package-owned identity metadata. |
| Executable evidence | TestComputeRootMatchesRFC9162TreeHash<br>TestProfilesMakeRootConventionsExplicit<br>TestLeafAndDigestBytesNeverAliasCallerMemory<br>TestRFC9162MatchesTransparencyDevMerkle |
| Official or pinned fixtures | testdata/rfc9162-sha256-v1.json |
| Fuzz evidence | FuzzComputeRoot |
| Interoperability evidence | differential_test.go |
| Affected public APIs | RawLeaf<br>Digest<br>ComputeRoot<br>Snapshot<br>Builder<br>RootBuilder |
| Affected documentation | docs/structures.md<br>docs/reference.md |
| Upstream status | RFC 9162 defines both prefixes directly; no applicable erratum is currently recorded. |
| Reconsider when | RFC 9162 errata or a separately versioned profile changes the commitment scheme. |

## MERKLETREE-DEC-002: Empty-tree commitment

| Field | Decision |
| --- | --- |
| Status | resolved |
| Owner | merkle-tree maintainers |
| Classification | interoperability policy |
| Decision scope | normative |
| Specification | RFC 9162 |
| Exact version | RFC 9162 (December 2021) |
| Source authority | rfc9162-source |
| Authoritative URL | https://www.rfc-editor.org/rfc/rfc9162.txt |
| Section | 2.1.1 |
| Requirement strength | not specified |
| Issue | Empty trees can be rejected, represented by a zero digest, or committed as the hash of an empty string. |
| Credible interpretations | Reject empty trees.<br>Use an all-zero digest.<br>Use SHA-256 of the empty string as RFC 9162 defines. |
| Known peer behavior | transparency-dev/merkle v0.0.2 commits the empty tree as SHA-256 of the empty string. |
| Selected behavior | Both version-1 profiles support an empty tree committed as SHA-256 of the empty string; decoders reject any other zero-size digest. |
| Normative rationale | One RFC-compatible empty identity avoids zero-value and forged-empty ambiguity. |
| Security consequences | The invalid Go zero value cannot be confused with the canonical empty commitment. |
| Resource consequences | Empty construction retains no leaves or nodes. |
| Compatibility consequences | Consumers using a zero digest or absent root must migrate to the RFC commitment. |
| Wire consequences | The root digest interoperates with RFC 9162 while the package binary envelope remains package-owned. |
| Executable evidence | TestRFC9162PinnedReferenceFixture<br>TestRootCanonicalBinaryEncodingFixture<br>TestSnapshotPersistenceEmptyAndRFCProfiles<br>TestRFC9162MatchesTransparencyDevMerkle |
| Official or pinned fixtures | testdata/rfc9162-sha256-v1.json |
| Fuzz evidence | FuzzComputeRoot<br>FuzzParseRoot |
| Interoperability evidence | differential_test.go |
| Affected public APIs | ComputeRoot<br>Root<br>Snapshot<br>ParseRoot<br>ParseSnapshot |
| Affected documentation | docs/structures.md<br>docs/encoding.md |
| Upstream status | RFC 9162 defines the empty hash directly; no applicable erratum is currently recorded. |
| Reconsider when | RFC 9162 errata or a separately identified profile defines another empty commitment. |

## MERKLETREE-DEC-003: Ordered recursive tree shape

| Field | Decision |
| --- | --- |
| Status | resolved |
| Owner | merkle-tree maintainers |
| Classification | interoperability policy |
| Decision scope | normative |
| Specification | RFC 9162 |
| Exact version | RFC 9162 (December 2021) |
| Source authority | rfc9162-source |
| Authoritative URL | https://www.rfc-editor.org/rfc/rfc9162.txt |
| Section | 2.1.1 |
| Requirement strength | not specified |
| Issue | Odd-width trees can duplicate, promote, pad, sort, or recursively split leaves and therefore produce incompatible roots. |
| Credible interpretations | Duplicate the final node.<br>Promote the final node unchanged.<br>Pad to a power of two or sort leaves or pairs.<br>Split at the largest power of two smaller than the tree size. |
| Known peer behavior | transparency-dev/merkle v0.0.2 preserves order and uses the RFC recursive largest-lower-power-of-two split. |
| Selected behavior | Preserve caller order and recursively split at the largest power of two smaller than the subtree size without sorting, duplication, promotion, or padding. |
| Normative rationale | RFC 9162 uniquely defines this shape and authenticates order. |
| Security consequences | Caller order cannot be normalized away without changing the commitment. |
| Resource consequences | Depth and node work are bounded before recursion. |
| Compatibility consequences | Duplicate-last, sorted-pair, promotion, and padding conventions are incompatible profiles. |
| Wire consequences | RFC-profile roots and proof paths match RFC 9162 for supported sizes. |
| Executable evidence | TestCanonicalProfileUsesDocumentedRFC9162Shape<br>TestRFC9162MatchesTransparencyDevMerkle<br>TestBuilderMatchesBatchConstructionForEverySmallPrefixAndProfile |
| Official or pinned fixtures | testdata/rfc9162-sha256-v1.json |
| Fuzz evidence | FuzzComputeRoot |
| Interoperability evidence | differential_test.go |
| Affected public APIs | ComputeRoot<br>Snapshot<br>Builder<br>RootBuilder<br>InclusionProof<br>ConsistencyProof |
| Affected documentation | docs/structures.md<br>docs/reference.md |
| Upstream status | RFC 9162 defines the recursive split algorithm; no applicable erratum is currently recorded. |
| Reconsider when | RFC 9162 errata or a new explicit profile defines another tree shape. |

## MERKLETREE-DEC-004: Complete root identity

| Field | Decision |
| --- | --- |
| Status | resolved |
| Owner | merkle-tree maintainers |
| Classification | omission |
| Decision scope | defensive |
| Specification | RFC 9162 |
| Exact version | RFC 9162 (December 2021) |
| Source authority | rfc9162-source |
| Authoritative URL | https://www.rfc-editor.org/rfc/rfc9162.txt |
| Section | 2.1 |
| Requirement strength | not specified |
| Issue | RFC 9162 tree algorithms do not define the complete package value identity needed to prevent cross-profile, cross-size, or cross-algorithm proof use. |
| Credible interpretations | Expose only digest bytes.<br>Infer metadata from each verifier.<br>Bind profile, version, algorithm, size, and digest in an immutable value. |
| Known peer behavior | Certificate Transparency interfaces commonly carry tree size beside the root hash but do not define this package envelope. |
| Selected behavior | Every Root binds profile ID, profile version, hash algorithm, tree size, and digest; verification and decoding reject partial or mismatched identities. |
| Normative rationale | Interpretation metadata is part of the authentication claim and must not be implicit. |
| Security consequences | Complete binding prevents cross-profile, cross-size, and cross-algorithm proof substitution. |
| Resource consequences | Fixed-size metadata is validated before path allocation or hashing. |
| Compatibility consequences | Consumers must preserve the complete Root rather than only digest bytes. |
| Wire consequences | Package binary formats encode identity metadata explicitly and are not Certificate Transparency wire formats. |
| Executable evidence | TestProfilesMakeRootConventionsExplicit<br>TestProfileValidationRejectsPartiallyMatchingRFCIdentity<br>TestInclusionProofBindsOperationIdentityAndOwnsReturnedSlices |
| Official or pinned fixtures | None. |
| Fuzz evidence | FuzzParseRoot |
| Interoperability evidence | None. |
| Affected public APIs | Profile<br>Root<br>InclusionProof<br>ConsistencyProof<br>MultiInclusionProof |
| Affected documentation | docs/structures.md<br>docs/encoding.md<br>docs/reference.md |
| Upstream status | The complete package identity is outside RFC 9162 and has no upstream erratum. |
| Reconsider when | A standardized envelope binds equivalent metadata with a compatible registry. |

## MERKLETREE-DEC-005: Inclusion proof path and verification

| Field | Decision |
| --- | --- |
| Status | resolved |
| Owner | merkle-tree maintainers |
| Classification | interoperability policy |
| Decision scope | normative |
| Specification | RFC 9162 |
| Exact version | RFC 9162 (December 2021) |
| Source authority | rfc9162-source |
| Authoritative URL | https://www.rfc-editor.org/rfc/rfc9162.txt |
| Section | 2.1.3.1 and 2.1.3.2 |
| Requirement strength | not specified |
| Issue | Audit paths can be ordered differently or accepted with missing, surplus, or structurally impossible siblings. |
| Credible interpretations | Store root-to-leaf nodes.<br>Store leaf-to-root nodes and accept any reconstructing path.<br>Require the unique RFC path for the bound index and size. |
| Known peer behavior | transparency-dev/merkle v0.0.2 generates and verifies RFC 9162 leaf-to-root audit paths. |
| Selected behavior | Store siblings leaf-to-root and require the unique RFC audit-path length; reject missing, surplus, malformed, wrong-leaf, wrong-root, wrong-profile, and over-limit proofs. |
| Normative rationale | A unique path and complete identity prevent malleable or context-free proof acceptance. |
| Security consequences | Verification binds every claim component and rejects surplus or structurally impossible input. |
| Resource consequences | Path count, bytes, depth, hashing work, and cancellation are checked before and during verification. |
| Compatibility consequences | RFC paths interoperate after adapting transport representation. |
| Wire consequences | Package binary proof encoding is not a Certificate Transparency wire format. |
| Executable evidence | TestRFC9162InclusionProofMatchesIndependentAuditPaths<br>TestInclusionProofRejectsWrongLeafAndInvalidRequests<br>TestVerifyInclusionRejectsMalformedAndNonVerifyingProofs |
| Official or pinned fixtures | testdata/rfc9162-sha256-v1.json |
| Fuzz evidence | FuzzVerifyInclusion<br>FuzzParseInclusionProof |
| Interoperability evidence | differential_test.go |
| Affected public APIs | Snapshot.InclusionProof<br>VerifyInclusion<br>InclusionProof |
| Affected documentation | docs/structures.md<br>docs/encoding.md<br>docs/reference.md |
| Upstream status | RFC 9162 defines generation and verification algorithms; no applicable erratum is currently recorded. |
| Reconsider when | RFC 9162 errata changes audit-path generation or verification semantics. |

## MERKLETREE-DEC-006: Consistency proof edge cases

| Field | Decision |
| --- | --- |
| Status | resolved |
| Owner | merkle-tree maintainers |
| Classification | omission |
| Decision scope | defensive |
| Specification | RFC 9162 |
| Exact version | RFC 9162 (December 2021) |
| Source authority | rfc9162-source |
| Authoritative URL | https://www.rfc-editor.org/rfc/rfc9162.txt |
| Section | 2.1.4.1 and 2.1.4.2 |
| Requirement strength | not specified |
| Issue | RFC 9162 defines nonempty strict-prefix proofs but omits zero-to-nonzero and equal-size API behavior. |
| Credible interpretations | Invent zero-to-nonzero proof semantics.<br>Accept any empty proof for equal sizes.<br>Require identical equal-size roots and reject the undefined zero-to-nonzero transition. |
| Known peer behavior | transparency-dev/merkle v0.0.2 matches nonempty strict-prefix SUBPROOF nodes; the peer API does not define this package's equal-size identity policy. |
| Selected behavior | Equal sizes require identical complete roots and an empty proof; zero-to-nonzero generation, encoding, and verification are rejected; other proofs use RFC SUBPROOF order. |
| Normative rationale | Undefined transitions must not acquire ad hoc trust semantics, while equal identity needs no path. |
| Security consequences | An equal-size digest mismatch cannot be hidden by an empty path. |
| Resource consequences | Proof node, byte, depth, and hash limits remain enforced. |
| Compatibility consequences | Nonempty strict-prefix behavior matches RFC 9162; callers needing zero-to-nonzero policy must establish it separately. |
| Wire consequences | Undefined transitions have no package proof encoding; equal-size proofs encode an empty path bound to both roots. |
| Executable evidence | TestConsistencyProofMatchesRFC9162Examples<br>TestConsistencyProofEqualSizeRequiresIdenticalRoots<br>TestConsistencyProofSupportsEqualEmptyTrees<br>TestConsistencyProofRejectsInvalidRequestsAndResourceClaims |
| Official or pinned fixtures | testdata/rfc9162-sha256-v1.json |
| Fuzz evidence | FuzzVerifyConsistency<br>FuzzParseConsistencyProof |
| Interoperability evidence | differential_test.go |
| Affected public APIs | Snapshot.ConsistencyProof<br>VerifyConsistency<br>ConsistencyProof |
| Affected documentation | docs/structures.md<br>docs/encoding.md<br>docs/reference.md |
| Upstream status | RFC 9162 omits zero-to-nonzero and equal-size proof operations; no accepted erratum currently defines them. |
| Reconsider when | An accepted RFC erratum or successor specification defines either transition. |

## MERKLETREE-DEC-007: Multi-inclusion proof authority

| Field | Decision |
| --- | --- |
| Status | resolved |
| Owner | merkle-tree maintainers |
| Classification | omission |
| Decision scope | application-policy |
| Specification | RFC 9162 |
| Exact version | RFC 9162 (December 2021) |
| Source authority | rfc9162-source |
| Authoritative URL | https://www.rfc-editor.org/rfc/rfc9162.txt |
| Section | 2.1.3 |
| Requirement strength | not specified |
| Issue | RFC 9162 defines single-leaf paths but no multi-inclusion proof representation, ordering, or duplicate policy. |
| Credible interpretations | Concatenate single proofs.<br>Retain caller index order.<br>Canonicalize indexes and emit one minimal deterministic frontier. |
| Known peer behavior | No maintained RFC 9162 peer defines a normative multi-inclusion proof representation. |
| Selected behavior | Copy, sort, and require unique in-range indexes, then emit the minimal left-to-right depth-first frontier; verification requires canonical leaf order and exact frontier consumption. |
| Normative rationale | One deterministic minimal form prevents duplicate and ordering ambiguity without claiming RFC interoperability. |
| Security consequences | Duplicate, missing, reordered, and surplus claims fail closed. |
| Resource consequences | Selected leaves, indexes, frontier nodes, bytes, depth, work, and temporary memory are bounded. |
| Compatibility consequences | A future standard multi-proof needs a separately versioned profile. |
| Wire consequences | The frontier and binary encoding are package-owned even when the root uses the RFC profile. |
| Executable evidence | TestMultiInclusionProofCanonicalizesIndexesAndOwnsSlices<br>TestMultiInclusionProofExhaustivelyAuthenticatesSmallSubsets<br>TestVerifyMultiInclusionRejectsMalformedAndNonVerifyingProofs |
| Official or pinned fixtures | None. |
| Fuzz evidence | FuzzVerifyMultiInclusion<br>FuzzParseMultiInclusionProof |
| Interoperability evidence | None. |
| Affected public APIs | Snapshot.MultiInclusionProof<br>VerifyMultiInclusion<br>MultiInclusionProof |
| Affected documentation | docs/structures.md<br>docs/encoding.md<br>docs/reference.md |
| Upstream status | RFC 9162 defines no multi-inclusion proof format. |
| Reconsider when | A suitable standard multi-proof format is adopted as a separately versioned profile. |

## MERKLETREE-DEC-008: Mutable builders and snapshot ownership

| Field | Decision |
| --- | --- |
| Status | resolved |
| Owner | merkle-tree maintainers |
| Classification | omission |
| Decision scope | application-policy |
| Specification | RFC 9162 |
| Exact version | RFC 9162 (December 2021) |
| Source authority | rfc9162-source |
| Authoritative URL | https://www.rfc-editor.org/rfc/rfc9162.txt |
| Section | Not specified by RFC 9162 |
| Requirement strength | not specified |
| Issue | RFC 9162 does not define mutable Go builder concurrency, failure atomicity, buffer ownership, or immutable snapshot behavior. |
| Credible interpretations | Make every builder method internally concurrent.<br>Permit partial batch append.<br>Require caller synchronization while making operations atomic and snapshots independent. |
| Known peer behavior | Maintained Merkle peers expose different ownership models, none of which is an RFC 9162 interoperability requirement. |
| Selected behavior | Builders are caller-synchronized; batches validate and compute before publication; inputs are copied or hashed immediately; returned snapshots own immutable state. |
| Normative rationale | Explicit ownership and failure atomicity prevent partial commitments without hidden synchronization costs. |
| Security consequences | Failed work cannot publish partial commitments or retain hostile caller buffers. |
| Resource consequences | Preflight limits and cancellation bound temporary state; callers must prevent data races on mutable builders. |
| Compatibility consequences | Successful incremental roots equal one-shot construction for every prefix. |
| Wire consequences | Concurrency and ownership do not alter Merkle commitments. |
| Executable evidence | TestBuilderBatchAppendIsAtomic<br>TestBuilderAppendSnapshotsMatchBatchConstruction<br>TestRootBuilderBatchAppendIsAtomicAndBounded<br>TestLeafAndDigestBytesNeverAliasCallerMemory |
| Official or pinned fixtures | None. |
| Fuzz evidence | None. |
| Interoperability evidence | None. |
| Affected public APIs | Builder<br>RootBuilder<br>Snapshot<br>RawLeaf |
| Affected documentation | docs/structures.md<br>docs/reference.md |
| Upstream status | Builder ownership is package policy outside RFC 9162. |
| Reconsider when | A separate concurrency-safe builder can preserve atomicity, ownership, and existing commitments. |

## MERKLETREE-DEC-009: Canonical binary encoding

| Field | Decision |
| --- | --- |
| Status | resolved |
| Owner | merkle-tree maintainers |
| Classification | omission |
| Decision scope | application-policy |
| Specification | RFC 9162 |
| Exact version | RFC 9162 (December 2021) |
| Source authority | rfc9162-source |
| Authoritative URL | https://www.rfc-editor.org/rfc/rfc9162.txt |
| Section | 2.1 |
| Requirement strength | not specified |
| Issue | RFC 9162 defines algorithms but no standalone package binary envelope for roots, proofs, or snapshots. |
| Credible interpretations | Serialize only digests.<br>Use a generic self-describing format.<br>Define a strict versioned envelope binding operation and complete identity. |
| Known peer behavior | Certificate Transparency protocol structures are intentionally not reused or claimed by this package. |
| Selected behavior | Version 1 uses MTRE, explicit object type and identity, fixed-width big-endian integers, exact SHA-256 digests, and operation-specific canonical structure; decoders reject all non-canonical input. |
| Normative rationale | A versioned typed envelope makes identity, allocation bounds, and compatibility explicit. |
| Security consequences | Type and identity binding prevent cross-operation substitution. |
| Resource consequences | Counts and size arithmetic are validated before allocation; parsing enforces byte, element, depth, work, memory, and cancellation limits. |
| Compatibility consequences | Any incompatible layout requires a new encoding version. |
| Wire consequences | Encoded bytes are stable package wire contracts but are not RFC 9162 or Certificate Transparency wire artifacts. |
| Executable evidence | TestRootCanonicalBinaryEncodingFixture<br>TestProofCanonicalBinaryEncodingFixtures<br>TestProofDecodersRejectTruncatedTrailingAndWrongOperation<br>TestSnapshotPersistenceRejectsMalformedAndBoundedInput |
| Official or pinned fixtures | None. |
| Fuzz evidence | FuzzParseRoot<br>FuzzParseInclusionProof<br>FuzzParseConsistencyProof<br>FuzzParseMultiInclusionProof<br>FuzzParseSnapshot |
| Interoperability evidence | None. |
| Affected public APIs | Root.MarshalBinary<br>InclusionProof.MarshalBinary<br>ConsistencyProof.MarshalBinary<br>MultiInclusionProof.MarshalBinary<br>Snapshot.MarshalBinary<br>ParseRoot<br>ParseInclusionProof<br>ParseConsistencyProof<br>ParseMultiInclusionProof<br>ParseSnapshot |
| Affected documentation | docs/encoding.md<br>docs/compatibility.md<br>docs/reference.md |
| Upstream status | RFC 9162 defines no equivalent standalone envelope. |
| Reconsider when | A standardized envelope meets the same identity, canonicality, and fail-closed requirements. |

## MERKLETREE-DEC-010: Persisted snapshots and resume trust

| Field | Decision |
| --- | --- |
| Status | resolved |
| Owner | merkle-tree maintainers |
| Classification | omission |
| Decision scope | defensive |
| Specification | RFC 9162 |
| Exact version | RFC 9162 (December 2021) |
| Source authority | rfc9162-source |
| Authoritative URL | https://www.rfc-editor.org/rfc/rfc9162.txt |
| Section | Not specified by RFC 9162 |
| Requirement strength | not specified |
| Issue | RFC 9162 does not define persisted tree topology, retained data, integrity validation, or trusted accounting for resume. |
| Credible interpretations | Trust serialized topology.<br>Persist raw leaves.<br>Persist canonical digest nodes and independently validate every structural and cryptographic invariant. |
| Known peer behavior | RFC 9162 peers do not define this package's snapshot persistence or resume policy. |
| Selected behavior | Persist canonical postorder digest nodes and metadata but never raw leaves; validate topology, hashes, root identity, uniqueness, and limits; ResumeBuilder requires a trusted expected raw-byte count. |
| Normative rationale | Only authenticated and independently validated state may resume mutation, while unauthenticated accounting remains an explicit side input. |
| Security consequences | Corrupt, cyclic, shared, reordered, or oversized state fails closed before use. |
| Resource consequences | Raw application data is not retained and untrusted accounting metadata cannot silently control resumed limits. |
| Compatibility consequences | Canonical snapshots round-trip and resume to roots identical to live construction. |
| Wire consequences | The trusted byte-count side input is required for resume and is not part of the Merkle commitment. |
| Executable evidence | TestSnapshotPersistenceRoundTripAndResume<br>TestSnapshotPersistenceRejectsCorruptMetadataAndNodes<br>TestResumeBuilderValidationCancellationAndLimits<br>TestSnapshotPersistenceEmptyMetadataAndInternalCorruption |
| Official or pinned fixtures | None. |
| Fuzz evidence | FuzzParseSnapshot |
| Interoperability evidence | None. |
| Affected public APIs | Snapshot.MarshalBinary<br>ParseSnapshot<br>ResumeBuilder |
| Affected documentation | docs/encoding.md<br>docs/errors-and-recovery.md<br>docs/reference.md |
| Upstream status | Snapshot persistence is outside RFC 9162. |
| Reconsider when | A new commitment profile authenticates retained accounting metadata. |

## Unresolved decisions

No known material RFC 9162 interpretation or package wire decision is
unresolved at this revision. New ambiguities remain unresolved until assigned
a stable identifier, authority analysis, executable evidence, and maintainer
disposition in both machine and human registers.
