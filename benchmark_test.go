package merkletree_test

import (
	"context"
	"fmt"
	"testing"

	merkletree "github.com/faustbrian/go-merkle-tree"
)

// BenchmarkComputeRoot measures one-shot roots for 1/256/16384 pre-owned 32-byte
// leaves; input construction is outside each timed sub-benchmark.
func BenchmarkComputeRoot(b *testing.B) {
	for _, size := range []int{1, 256, 16_384} {
		leaves := make([]merkletree.RawLeaf, size)
		for index := range leaves {
			value := make([]byte, 32)
			value[0] = byte(index)
			value[1] = byte(index >> 8)
			leaves[index] = merkletree.NewRawLeaf(value)
		}

		b.Run(fmt.Sprintf("leaves_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(size * 32))

			for range b.N {
				if _, err := merkletree.ComputeRoot(
					context.Background(),
					merkletree.CanonicalProfile(),
					leaves,
					merkletree.DefaultLimits(),
				); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkInclusionProof times path generation at index 8192 of 16384 leaves.
// Snapshot/input construction precedes ResetTimer.
func BenchmarkInclusionProof(b *testing.B) {
	leaves := benchmarkLeaves(16_384)
	snapshot, err := merkletree.NewSnapshot(
		context.Background(),
		merkletree.CanonicalProfile(),
		leaves,
		merkletree.DefaultSnapshotLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := snapshot.InclusionProof(
			context.Background(),
			8_192,
			merkletree.DefaultProofLimits(),
		); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkVerifyInclusion times verification with a prebuilt 16384-leaf proof and
// raw leaf at index 8192; proof construction is excluded.
func BenchmarkVerifyInclusion(b *testing.B) {
	leaves := benchmarkLeaves(16_384)
	snapshot, err := merkletree.NewSnapshot(
		context.Background(),
		merkletree.CanonicalProfile(),
		leaves,
		merkletree.DefaultSnapshotLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}
	proof, err := snapshot.InclusionProof(
		context.Background(),
		8_192,
		merkletree.DefaultProofLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := merkletree.VerifyInclusion(
			context.Background(),
			proof,
			leaves[8_192],
			merkletree.DefaultProofLimits(),
		); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuilderAppend times one 32-byte append into the same growing builder.
// Setup is excluded; node/leaf capacities depend on b.N, while the default
// total-byte ceiling remains in force. This is not fixed-size append latency.
func BenchmarkBuilderAppend(b *testing.B) {
	leaf := merkletree.NewRawLeaf(make([]byte, 32))
	limits := merkletree.DefaultSnapshotLimits()
	limits.Construction.MaxLeaves = uint64(b.N)
	limits.MaxRetainedNodes = 2*uint64(b.N) + 1
	builder, err := merkletree.NewBuilder(
		merkletree.CanonicalProfile(),
		limits,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(32)
	b.ResetTimer()
	for range b.N {
		if err := builder.Append(context.Background(), leaf); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuilderAppendBatch times repeated 256-leaf batches of 32-byte leaves
// into one growing builder. Setup is excluded and capacities depend on b.N;
// the default total-byte ceiling still bounds calibration.
func BenchmarkBuilderAppendBatch(b *testing.B) {
	leaves := benchmarkLeaves(256)
	limits := merkletree.DefaultSnapshotLimits()
	limits.Construction.MaxLeaves = uint64(b.N * len(leaves))
	limits.MaxRetainedNodes = 2*limits.Construction.MaxLeaves + 1
	builder, err := merkletree.NewBuilder(
		merkletree.CanonicalProfile(),
		limits,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(leaves) * 32))
	b.ResetTimer()
	for range b.N {
		if err := builder.AppendBatch(
			context.Background(),
			leaves,
		); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRootBuilderAppend times a 32-byte append into one growing streaming
// frontier, excluding setup; b.N determines leaf capacity, not the byte limit.
func BenchmarkRootBuilderAppend(b *testing.B) {
	leaf := merkletree.NewRawLeaf(make([]byte, 32))
	limits := merkletree.DefaultLimits()
	limits.MaxLeaves = uint64(b.N)
	builder, err := merkletree.NewRootBuilder(
		merkletree.CanonicalProfile(),
		limits,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(32)
	b.ResetTimer()
	for range b.N {
		if err := builder.Append(context.Background(), leaf); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRootBuilderAppendBatch times 256-leaf/32-byte batches into one growing
// streaming frontier. Setup is excluded; the default total-byte ceiling remains.
func BenchmarkRootBuilderAppendBatch(b *testing.B) {
	leaves := benchmarkLeaves(256)
	limits := merkletree.DefaultLimits()
	limits.MaxLeaves = uint64(b.N) * uint64(len(leaves))
	builder, err := merkletree.NewRootBuilder(
		merkletree.CanonicalProfile(),
		limits,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(leaves) * 32))
	b.ResetTimer()
	for range b.N {
		if err := builder.AppendBatch(
			context.Background(),
			leaves,
		); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRootBuilderRoot times folding a prebuilt 16383-leaf frontier.
// Leaf ingestion is excluded; the non-power-of-two size exercises final merges.
func BenchmarkRootBuilderRoot(b *testing.B) {
	leaves := benchmarkLeaves(16_383)
	builder, err := merkletree.NewRootBuilder(
		merkletree.CanonicalProfile(),
		merkletree.DefaultLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}
	if err := builder.AppendBatch(
		context.Background(),
		leaves,
	); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := builder.Root(context.Background()); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConsistencyProof times path generation for an 8191-to-16384 prefix.
// Both snapshots and the older root are constructed before ResetTimer.
func BenchmarkConsistencyProof(b *testing.B) {
	leaves := benchmarkLeaves(16_384)
	older, newer := benchmarkConsistencySnapshots(b, leaves, 8_191)
	olderRoot, err := older.Root()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := newer.ConsistencyProof(
			context.Background(),
			olderRoot,
			merkletree.DefaultConsistencyProofLimits(),
		); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkVerifyConsistency times verification of a prebuilt 8191-to-16384 proof,
// excluding snapshot and proof construction.
func BenchmarkVerifyConsistency(b *testing.B) {
	leaves := benchmarkLeaves(16_384)
	older, newer := benchmarkConsistencySnapshots(b, leaves, 8_191)
	olderRoot, err := older.Root()
	if err != nil {
		b.Fatal(err)
	}
	proof, err := newer.ConsistencyProof(
		context.Background(),
		olderRoot,
		merkletree.DefaultConsistencyProofLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := merkletree.VerifyConsistency(
			context.Background(),
			proof,
			merkletree.DefaultConsistencyProofLimits(),
		); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMultiInclusionProof times frontier generation for 256 evenly spaced
// indexes (stride 64) in 16384 leaves; input and snapshot setup are excluded.
func BenchmarkMultiInclusionProof(b *testing.B) {
	leaves := benchmarkLeaves(16_384)
	snapshot, err := merkletree.NewSnapshot(
		context.Background(),
		merkletree.CanonicalProfile(),
		leaves,
		merkletree.DefaultSnapshotLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}
	indexes := benchmarkMultiIndexes(256, 64)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := snapshot.MultiInclusionProof(
			context.Background(),
			indexes,
			merkletree.DefaultMultiProofLimits(),
		); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkVerifyMultiInclusion times a prebuilt 256-selected-leaf proof in a
// 16384-leaf tree; leaf selection and proof generation precede ResetTimer.
func BenchmarkVerifyMultiInclusion(b *testing.B) {
	leaves := benchmarkLeaves(16_384)
	snapshot, err := merkletree.NewSnapshot(
		context.Background(),
		merkletree.CanonicalProfile(),
		leaves,
		merkletree.DefaultSnapshotLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}
	indexes := benchmarkMultiIndexes(256, 64)
	selected := make([]merkletree.RawLeaf, len(indexes))
	for index, leafIndex := range indexes {
		selected[index] = leaves[leafIndex]
	}
	proof, err := snapshot.MultiInclusionProof(
		context.Background(),
		indexes,
		merkletree.DefaultMultiProofLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := merkletree.VerifyMultiInclusion(
			context.Background(),
			proof,
			selected,
			merkletree.DefaultMultiProofLimits(),
		); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkInclusionProofEncoding separates marshal, parse and truncated-frame
// rejection of a prebuilt 16384-leaf/index-8192 proof. Rates use encoded bytes;
// parsing does not include cryptographic verification.
func BenchmarkInclusionProofEncoding(b *testing.B) {
	leaves := benchmarkLeaves(16_384)
	snapshot, err := merkletree.NewSnapshot(
		context.Background(),
		merkletree.CanonicalProfile(),
		leaves,
		merkletree.DefaultSnapshotLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}
	proof, err := snapshot.InclusionProof(
		context.Background(),
		8_192,
		merkletree.DefaultProofLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}
	encoded, err := proof.MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}

	b.Run("marshal", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(encoded)))
		for range b.N {
			if _, err := proof.MarshalBinary(); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("parse", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(encoded)))
		for range b.N {
			if _, err := merkletree.ParseInclusionProof(
				context.Background(),
				encoded,
				merkletree.DefaultEncodingLimits(),
				merkletree.DefaultProofLimits(),
			); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("reject_malformed", func(b *testing.B) {
		truncated := encoded[:len(encoded)-1]
		b.ReportAllocs()
		b.SetBytes(int64(len(truncated)))
		for range b.N {
			if _, err := merkletree.ParseInclusionProof(
				context.Background(),
				truncated,
				merkletree.DefaultEncodingLimits(),
				merkletree.DefaultProofLimits(),
			); err == nil {
				b.Fatal("malformed proof accepted")
			}
		}
	})
}

// BenchmarkSnapshotPersistence separates full marshal, parse and builder resume
// of a prebuilt 16384-leaf snapshot. Setup is outside timed sub-benchmarks;
// marshal/parse rates use encoded bytes, not raw-leaf throughput.
func BenchmarkSnapshotPersistence(b *testing.B) {
	leaves := benchmarkLeaves(16_384)
	snapshot, err := merkletree.NewSnapshot(
		context.Background(),
		merkletree.CanonicalProfile(),
		leaves,
		merkletree.DefaultSnapshotLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}
	encoded, err := snapshot.MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}
	persistenceLimits := merkletree.DefaultSnapshotPersistenceLimits()

	b.Run("marshal", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(encoded)))
		for range b.N {
			if _, err := snapshot.MarshalBinary(); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("parse", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(encoded)))
		for range b.N {
			if _, err := merkletree.ParseSnapshot(
				context.Background(),
				encoded,
				persistenceLimits,
			); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("resume_builder", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			if _, err := merkletree.ResumeBuilder(
				context.Background(),
				snapshot,
				uint64(len(leaves)*32),
				merkletree.DefaultSnapshotLimits(),
			); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkMultiIndexes(count int, stride uint64) []uint64 {
	indexes := make([]uint64, count)
	for index := range indexes {
		indexes[index] = uint64(index) * stride
	}

	return indexes
}

func benchmarkConsistencySnapshots(
	b *testing.B,
	leaves []merkletree.RawLeaf,
	olderSize int,
) (merkletree.Snapshot, merkletree.Snapshot) {
	b.Helper()

	older, err := merkletree.NewSnapshot(
		context.Background(),
		merkletree.CanonicalProfile(),
		leaves[:olderSize],
		merkletree.DefaultSnapshotLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}
	newer, err := merkletree.NewSnapshot(
		context.Background(),
		merkletree.CanonicalProfile(),
		leaves,
		merkletree.DefaultSnapshotLimits(),
	)
	if err != nil {
		b.Fatal(err)
	}

	return older, newer
}

func benchmarkLeaves(size int) []merkletree.RawLeaf {
	leaves := make([]merkletree.RawLeaf, size)
	for index := range leaves {
		value := make([]byte, 32)
		value[0] = byte(index)
		value[1] = byte(index >> 8)
		leaves[index] = merkletree.NewRawLeaf(value)
	}

	return leaves
}
