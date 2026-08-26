# merkle-tree

[![CI](https://github.com/faustbrian/go-merkle-tree/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-merkle-tree/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-merkle-tree/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-merkle-tree.svg)](https://pkg.go.dev/github.com/faustbrian/go-merkle-tree)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-merkle-tree?sort=semver)](https://github.com/faustbrian/go-merkle-tree/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`merkle-tree` is a storage-independent library for explicitly profiled,
ordered Merkle trees. It computes and streams bounded roots, incrementally
appends leaves, creates immutable snapshots, and generates and verifies
inclusion, multi-inclusion, and consistency proofs.

The stable v1 surface supports the package canonical binary profile and the
RFC 9162 Certificate Transparency profile.

## Installation

```sh
go get github.com/faustbrian/go-merkle-tree
```

## Quick start

```go
profile := merkletree.CanonicalProfile()
leaves := []merkletree.RawLeaf{
	merkletree.NewRawLeaf([]byte("first")),
	merkletree.NewRawLeaf([]byte("second")),
}

root, err := merkletree.ComputeRoot(
	ctx,
	profile,
	leaves,
	merkletree.DefaultLimits(),
)
if err != nil {
	return err
}
```

Use `NewSnapshot` when proof generation is required and `NewRootBuilder` for
streaming root computation with logarithmic retained state.

## Guarantees and limits

- Leaves, digests, roots, snapshots, and proofs own their retained bytes.
- Hashing and proof behavior is profile-bound and domain-separated.
- Batch append is atomic and cancellation-safe.
- Proofs bind profile, version, algorithm, root, tree size, indexes, and paths.
- Input, retained state, temporary work, and proof sizes are explicitly
  bounded.
- A valid proof authenticates against its supplied root; it does not establish
  that the root is trusted, recent, or authorized.

## Documentation

Use the [documentation index](docs/README.md) for profiles, proof semantics,
security, compatibility, and operations. The [detailed reference](docs/reference.md)
contains the complete builder, proof, specification, and v1 boundary notes.

## Development

Run `make check` before changing profile, hashing, proof, or serialization
behavior.

## License

MIT. See [LICENSE](LICENSE).
