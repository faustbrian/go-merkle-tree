// Package merkletree constructs explicitly profiled cryptographic Merkle-tree
// roots, streams roots through a logarithmic-memory builder, incrementally
// appends ordered leaves into proof-retaining builders, creates immutable
// snapshots, and generates and independently verifies inclusion,
// multi-inclusion, and append-only consistency proofs. Roots and proofs have
// bounded, versioned canonical binary encodings. Immutable snapshots can be
// persisted, structurally validated, and resumed without raw leaf bytes.
//
// The package does not define a universal "Merkle tree" convention. Every
// root carries a profile, profile version, hash algorithm, and tree size.
// Merkle proofs establish a relationship to a trusted root; they do not
// establish that a leaf is true, fresh, authorized, or semantically valid.
// Verify functions use the roots embedded in proofs. Authenticate their complete
// identities and check the requested indexes separately; parsing alone only
// validates structure. Canonical binary envelopes are package formats, not
// Certificate Transparency protocol messages.
//
// RawLeaf constructors and byte/vector accessors copy. Snapshots and proofs
// may share immutable storage and support concurrent read-only use. Builder
// and RootBuilder require caller synchronization and must not be copied for
// independent mutation. The package starts no background goroutines.
//
// Context cancellation is cooperative at explicit checkpoints. It does not
// interrupt an individual hash, sort, allocation or bulk copy, and an operation
// may succeed if cancellation arrives after its final checkpoint. Resource
// limits bound their named accounting measures, not total process memory.
package merkletree
