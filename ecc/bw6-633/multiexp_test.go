// Copyright 2020 Consensys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package bw6633

import (
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/0x3327/gnark-crypto/ecc"
	"github.com/0x3327/gnark-crypto/ecc/bw6-633/fr"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

func TestMultiExpG1(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = 3
	} else {
		parameters.MinSuccessfulTests = nbFuzzShort * 2
	}

	properties := gopter.NewProperties(parameters)

	genScalar := GenFr()

	// size of the multiExps
	const nbSamples = 73

	// multi exp points
	var samplePoints [nbSamples]G1Affine
	var g G1Jac
	g.Set(&g1Gen)
	for i := 1; i <= nbSamples; i++ {
		samplePoints[i-1].FromJacobian(&g)
		g.AddAssign(&g1Gen)
	}

	// sprinkle some points at infinity
	rand.Seed(time.Now().UnixNano())
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here

	// final scalar to use in double and add method (without mixer factor)
	// n(n+1)(2n+1)/6  (sum of the squares from 1 to n)
	var scalar big.Int
	scalar.SetInt64(nbSamples)
	scalar.Mul(&scalar, new(big.Int).SetInt64(nbSamples+1))
	scalar.Mul(&scalar, new(big.Int).SetInt64(2*nbSamples+1))
	scalar.Div(&scalar, new(big.Int).SetInt64(6))

	// ensure a multiexp that's splitted has the same result as a non-splitted one..
	properties.Property("[G1] Multi exponentiation (cmax) should be consistent with splitted multiexp", prop.ForAll(
		func(mixer fr.Element) bool {
			var samplePointsLarge [nbSamples * 13]G1Affine
			for i := 0; i < 13; i++ {
				copy(samplePointsLarge[i*nbSamples:], samplePoints[:])
			}

			var rmax, splitted1, splitted2 G1Jac

			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples * 13]fr.Element

			for i := 1; i <= nbSamples; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer)
			}

			rmax.MultiExp(samplePointsLarge[:], sampleScalars[:], ecc.MultiExpConfig{})
			splitted1.MultiExp(samplePointsLarge[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: 128})
			splitted2.MultiExp(samplePointsLarge[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: 51})
			return rmax.Equal(&splitted1) && rmax.Equal(&splitted2)
		},
		genScalar,
	))

	// cRange is generated from template and contains the available parameters for the multiexp window size
	cRange := []uint64{4, 5, 6, 8, 12, 16}
	if testing.Short() {
		// test only "odd" and "even" (ie windows size divide word size vs not)
		cRange = []uint64{5, 14}
	}

	properties.Property(fmt.Sprintf("[G1] Multi exponentiation (c in %v) should be consistent with sum of square", cRange), prop.ForAll(
		func(mixer fr.Element) bool {

			var expected G1Jac

			// compute expected result with double and add
			var finalScalar, mixerBigInt big.Int
			finalScalar.Mul(&scalar, mixer.BigInt(&mixerBigInt))
			expected.ScalarMultiplication(&g1Gen, &finalScalar)

			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples]fr.Element

			for i := 1; i <= nbSamples; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer)
			}

			results := make([]G1Jac, len(cRange))
			for i, c := range cRange {
				_innerMsmG1(&results[i], c, samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: runtime.NumCPU()})
			}
			for i := 1; i < len(results); i++ {
				if !results[i].Equal(&results[i-1]) {
					t.Logf("result for c=%d != c=%d", cRange[i-1], cRange[i])
					return false
				}
			}
			return true
		},
		genScalar,
	))

	properties.Property(fmt.Sprintf("[G1] Multi exponentiation (c in %v) of points at infinity should output a point at infinity", cRange), prop.ForAll(
		func(mixer fr.Element) bool {

			var samplePointsZero [nbSamples]G1Affine

			var expected G1Jac

			// compute expected result with double and add
			var finalScalar, mixerBigInt big.Int
			finalScalar.Mul(&scalar, mixer.BigInt(&mixerBigInt))
			expected.ScalarMultiplication(&g1Gen, &finalScalar)

			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples]fr.Element

			for i := 1; i <= nbSamples; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer)
				samplePointsZero[i-1].setInfinity()
			}

			results := make([]G1Jac, len(cRange))
			for i, c := range cRange {
				_innerMsmG1(&results[i], c, samplePointsZero[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: runtime.NumCPU()})
			}
			for i := 0; i < len(results); i++ {
				if !results[i].Z.IsZero() {
					t.Logf("result for c=%d is not infinity", cRange[i])
					return false
				}
			}
			return true
		},
		genScalar,
	))

	properties.Property(fmt.Sprintf("[G1] Multi exponentiation (c in %v) with a vector of 0s as input should output a point at infinity", cRange), prop.ForAll(
		func(mixer fr.Element) bool {
			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples]fr.Element

			results := make([]G1Jac, len(cRange))
			for i, c := range cRange {
				_innerMsmG1(&results[i], c, samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: runtime.NumCPU()})
			}
			for i := 0; i < len(results); i++ {
				if !results[i].Z.IsZero() {
					t.Logf("result for c=%d is not infinity", cRange[i])
					return false
				}
			}
			return true
		},
		genScalar,
	))

	// note : this test is here as we expect to have a different multiExp than the above bucket method
	// for small number of points
	properties.Property("[G1] Multi exponentiation (<50points) should be consistent with sum of square", prop.ForAll(
		func(mixer fr.Element) bool {

			var g G1Jac
			g.Set(&g1Gen)

			// mixer ensures that all the words of a fpElement are set
			samplePoints := make([]G1Affine, 30)
			sampleScalars := make([]fr.Element, 30)

			for i := 1; i <= 30; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer)
				samplePoints[i-1].FromJacobian(&g)
				g.AddAssign(&g1Gen)
			}

			var op1MultiExp G1Affine
			op1MultiExp.MultiExp(samplePoints, sampleScalars, ecc.MultiExpConfig{})

			var finalBigScalar fr.Element
			var finalBigScalarBi big.Int
			var op1ScalarMul G1Affine
			finalBigScalar.SetUint64(9455).Mul(&finalBigScalar, &mixer)
			finalBigScalar.BigInt(&finalBigScalarBi)
			op1ScalarMul.ScalarMultiplication(&g1GenAff, &finalBigScalarBi)

			return op1ScalarMul.Equal(&op1MultiExp)
		},
		genScalar,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestCrossMultiExpG1(t *testing.T) {
	const nbSamples = 1 << 14
	// multi exp points
	var samplePoints [nbSamples]G1Affine
	var g G1Jac
	g.Set(&g1Gen)
	for i := 1; i <= nbSamples; i++ {
		samplePoints[i-1].FromJacobian(&g)
		g.AddAssign(&g1Gen)
	}

	// sprinkle some points at infinity
	rand.Seed(time.Now().UnixNano())
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here

	var sampleScalars [nbSamples]fr.Element
	fillBenchScalars(sampleScalars[:])

	// sprinkle some doublings
	for i := 10; i < 100; i++ {
		samplePoints[i] = samplePoints[0]
		sampleScalars[i] = sampleScalars[0]
	}

	// cRange is generated from template and contains the available parameters for the multiexp window size
	cRange := []uint64{4, 5, 6, 8, 12, 16}
	if testing.Short() {
		// test only "odd" and "even" (ie windows size divide word size vs not)
		cRange = []uint64{5, 14}
	}

	results := make([]G1Jac, len(cRange))
	for i, c := range cRange {
		_innerMsmG1(&results[i], c, samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: runtime.NumCPU()})
	}

	var r G1Jac
	_innerMsmG1Reference(&r, samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: runtime.NumCPU()})

	var expected, got G1Affine
	expected.FromJacobian(&r)

	for i := 0; i < len(results); i++ {
		got.FromJacobian(&results[i])
		if !expected.Equal(&got) {
			t.Fatalf("cross msm failed with c=%d", cRange[i])
		}
	}

}

// _innerMsmG1Reference always do ext jacobian with c == 16
func _innerMsmG1Reference(p *G1Jac, points []G1Affine, scalars []fr.Element, config ecc.MultiExpConfig) *G1Jac {
	// partition the scalars
	digits, _ := partitionScalars(scalars, 16, config.NbTasks)

	nbChunks := computeNbChunks(16)

	// for each chunk, spawn one go routine that'll loop through all the scalars in the
	// corresponding bit-window
	// note that buckets is an array allocated on the stack and this is critical for performance

	// each go routine sends its result in chChunks[i] channel
	chChunks := make([]chan g1JacExtended, nbChunks)
	for i := 0; i < len(chChunks); i++ {
		chChunks[i] = make(chan g1JacExtended, 1)
	}

	// the last chunk may be processed with a different method than the rest, as it could be smaller.
	n := len(points)
	for j := int(nbChunks - 1); j >= 0; j-- {
		processChunk := processChunkG1Jacobian[bucketg1JacExtendedC16]
		go processChunk(uint64(j), chChunks[j], 16, points, digits[j*n:(j+1)*n], nil)
	}

	return msmReduceChunkG1Affine(p, int(16), chChunks[:])
}

func BenchmarkMultiExpG1(b *testing.B) {

	const (
		pow       = (bits.UintSize / 2) - (bits.UintSize / 8) // 24 on 64 bits arch, 12 on 32 bits
		nbSamples = 1 << pow
	)

	var (
		samplePoints             [nbSamples]G1Affine
		sampleScalars            [nbSamples]fr.Element
		sampleScalarsSmallValues [nbSamples]fr.Element
		sampleScalarsRedundant   [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	copy(sampleScalarsSmallValues[:], sampleScalars[:])
	copy(sampleScalarsRedundant[:], sampleScalars[:])

	// this means first chunk is going to have more work to do and should be split into several go routines
	for i := 0; i < len(sampleScalarsSmallValues); i++ {
		if i%5 == 0 {
			sampleScalarsSmallValues[i].SetZero()
			sampleScalarsSmallValues[i][0] = 1
		}
	}

	// bad case for batch affine because scalar distribution might look uniform
	// but over batchSize windows, we may hit a lot of conflicts and force the msm-affine
	// to process small batches of additions to flush its queue of conflicted points.
	for i := 0; i < len(sampleScalarsRedundant); i += 100 {
		for j := i + 1; j < i+100 && j < len(sampleScalarsRedundant); j++ {
			sampleScalarsRedundant[j] = sampleScalarsRedundant[i]
		}
	}

	fillBenchBasesG1(samplePoints[:])

	var testPoint G1Affine

	for i := 5; i <= pow; i++ {
		using := 1 << i

		b.Run(fmt.Sprintf("%d points", using), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				testPoint.MultiExp(samplePoints[:using], sampleScalars[:using], ecc.MultiExpConfig{})
			}
		})

		b.Run(fmt.Sprintf("%d points-smallvalues", using), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				testPoint.MultiExp(samplePoints[:using], sampleScalarsSmallValues[:using], ecc.MultiExpConfig{})
			}
		})

		b.Run(fmt.Sprintf("%d points-redundancy", using), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				testPoint.MultiExp(samplePoints[:using], sampleScalarsRedundant[:using], ecc.MultiExpConfig{})
			}
		})
	}
}

func BenchmarkMultiExpG1Reference(b *testing.B) {
	const nbSamples = 1 << 20

	var (
		samplePoints  [nbSamples]G1Affine
		sampleScalars [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	fillBenchBasesG1(samplePoints[:])

	var testPoint G1Affine

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		testPoint.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
	}
}

func BenchmarkManyMultiExpG1Reference(b *testing.B) {
	const nbSamples = 1 << 20

	var (
		samplePoints  [nbSamples]G1Affine
		sampleScalars [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	fillBenchBasesG1(samplePoints[:])

	var t1, t2, t3 G1Affine
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			t1.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
			wg.Done()
		}()
		go func() {
			t2.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
			wg.Done()
		}()
		go func() {
			t3.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
			wg.Done()
		}()
		wg.Wait()
	}
}

// WARNING: this return points that are NOT on the curve and is meant to be use for benchmarking
// purposes only. We don't check that the result is valid but just measure "computational complexity".
//
// Rationale for generating points that are not on the curve is that for large benchmarks, generating
// a vector of different points can take minutes. Using the same point or subset will bias the benchmark result
// since bucket additions in extended jacobian coordinates will hit doubling algorithm instead of add.
func fillBenchBasesG1(samplePoints []G1Affine) {
	var r big.Int
	r.SetString("340444420969191673093399857471996460938405", 10)
	samplePoints[0].ScalarMultiplication(&samplePoints[0], &r)

	one := samplePoints[0].X
	one.SetOne()

	for i := 1; i < len(samplePoints); i++ {
		samplePoints[i].X.Add(&samplePoints[i-1].X, &one)
		samplePoints[i].Y.Sub(&samplePoints[i-1].Y, &one)
	}
}

func TestMultiExpG2(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = 3
	} else {
		parameters.MinSuccessfulTests = nbFuzzShort * 2
	}

	properties := gopter.NewProperties(parameters)

	genScalar := GenFr()

	// size of the multiExps
	const nbSamples = 73

	// multi exp points
	var samplePoints [nbSamples]G2Affine
	var g G2Jac
	g.Set(&g2Gen)
	for i := 1; i <= nbSamples; i++ {
		samplePoints[i-1].FromJacobian(&g)
		g.AddAssign(&g2Gen)
	}

	// sprinkle some points at infinity
	rand.Seed(time.Now().UnixNano())
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here

	// final scalar to use in double and add method (without mixer factor)
	// n(n+1)(2n+1)/6  (sum of the squares from 1 to n)
	var scalar big.Int
	scalar.SetInt64(nbSamples)
	scalar.Mul(&scalar, new(big.Int).SetInt64(nbSamples+1))
	scalar.Mul(&scalar, new(big.Int).SetInt64(2*nbSamples+1))
	scalar.Div(&scalar, new(big.Int).SetInt64(6))

	// ensure a multiexp that's splitted has the same result as a non-splitted one..
	properties.Property("[G2] Multi exponentiation (cmax) should be consistent with splitted multiexp", prop.ForAll(
		func(mixer fr.Element) bool {
			var samplePointsLarge [nbSamples * 13]G2Affine
			for i := 0; i < 13; i++ {
				copy(samplePointsLarge[i*nbSamples:], samplePoints[:])
			}

			var rmax, splitted1, splitted2 G2Jac

			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples * 13]fr.Element

			for i := 1; i <= nbSamples; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer)
			}

			rmax.MultiExp(samplePointsLarge[:], sampleScalars[:], ecc.MultiExpConfig{})
			splitted1.MultiExp(samplePointsLarge[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: 128})
			splitted2.MultiExp(samplePointsLarge[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: 51})
			return rmax.Equal(&splitted1) && rmax.Equal(&splitted2)
		},
		genScalar,
	))

	// cRange is generated from template and contains the available parameters for the multiexp window size
	// for g2, CI suffers with large c size since it needs to allocate a lot of memory for the buckets.
	// test only "odd" and "even" (ie windows size divide word size vs not)
	cRange := []uint64{5, 14}

	properties.Property(fmt.Sprintf("[G2] Multi exponentiation (c in %v) should be consistent with sum of square", cRange), prop.ForAll(
		func(mixer fr.Element) bool {

			var expected G2Jac

			// compute expected result with double and add
			var finalScalar, mixerBigInt big.Int
			finalScalar.Mul(&scalar, mixer.BigInt(&mixerBigInt))
			expected.ScalarMultiplication(&g2Gen, &finalScalar)

			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples]fr.Element

			for i := 1; i <= nbSamples; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer)
			}

			results := make([]G2Jac, len(cRange))
			for i, c := range cRange {
				_innerMsmG2(&results[i], c, samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: runtime.NumCPU()})
			}
			for i := 1; i < len(results); i++ {
				if !results[i].Equal(&results[i-1]) {
					t.Logf("result for c=%d != c=%d", cRange[i-1], cRange[i])
					return false
				}
			}
			return true
		},
		genScalar,
	))

	properties.Property(fmt.Sprintf("[G2] Multi exponentiation (c in %v) of points at infinity should output a point at infinity", cRange), prop.ForAll(
		func(mixer fr.Element) bool {

			var samplePointsZero [nbSamples]G2Affine

			var expected G2Jac

			// compute expected result with double and add
			var finalScalar, mixerBigInt big.Int
			finalScalar.Mul(&scalar, mixer.BigInt(&mixerBigInt))
			expected.ScalarMultiplication(&g2Gen, &finalScalar)

			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples]fr.Element

			for i := 1; i <= nbSamples; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer)
				samplePointsZero[i-1].setInfinity()
			}

			results := make([]G2Jac, len(cRange))
			for i, c := range cRange {
				_innerMsmG2(&results[i], c, samplePointsZero[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: runtime.NumCPU()})
			}
			for i := 0; i < len(results); i++ {
				if !results[i].Z.IsZero() {
					t.Logf("result for c=%d is not infinity", cRange[i])
					return false
				}
			}
			return true
		},
		genScalar,
	))

	properties.Property(fmt.Sprintf("[G2] Multi exponentiation (c in %v) with a vector of 0s as input should output a point at infinity", cRange), prop.ForAll(
		func(mixer fr.Element) bool {
			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples]fr.Element

			results := make([]G2Jac, len(cRange))
			for i, c := range cRange {
				_innerMsmG2(&results[i], c, samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: runtime.NumCPU()})
			}
			for i := 0; i < len(results); i++ {
				if !results[i].Z.IsZero() {
					t.Logf("result for c=%d is not infinity", cRange[i])
					return false
				}
			}
			return true
		},
		genScalar,
	))

	// note : this test is here as we expect to have a different multiExp than the above bucket method
	// for small number of points
	properties.Property("[G2] Multi exponentiation (<50points) should be consistent with sum of square", prop.ForAll(
		func(mixer fr.Element) bool {

			var g G2Jac
			g.Set(&g2Gen)

			// mixer ensures that all the words of a fpElement are set
			samplePoints := make([]G2Affine, 30)
			sampleScalars := make([]fr.Element, 30)

			for i := 1; i <= 30; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer)
				samplePoints[i-1].FromJacobian(&g)
				g.AddAssign(&g2Gen)
			}

			var op1MultiExp G2Affine
			op1MultiExp.MultiExp(samplePoints, sampleScalars, ecc.MultiExpConfig{})

			var finalBigScalar fr.Element
			var finalBigScalarBi big.Int
			var op1ScalarMul G2Affine
			finalBigScalar.SetUint64(9455).Mul(&finalBigScalar, &mixer)
			finalBigScalar.BigInt(&finalBigScalarBi)
			op1ScalarMul.ScalarMultiplication(&g2GenAff, &finalBigScalarBi)

			return op1ScalarMul.Equal(&op1MultiExp)
		},
		genScalar,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestCrossMultiExpG2(t *testing.T) {
	const nbSamples = 1 << 14
	// multi exp points
	var samplePoints [nbSamples]G2Affine
	var g G2Jac
	g.Set(&g2Gen)
	for i := 1; i <= nbSamples; i++ {
		samplePoints[i-1].FromJacobian(&g)
		g.AddAssign(&g2Gen)
	}

	// sprinkle some points at infinity
	rand.Seed(time.Now().UnixNano())
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here
	samplePoints[rand.Intn(nbSamples)].setInfinity() //#nosec G404 weak rng is fine here

	var sampleScalars [nbSamples]fr.Element
	fillBenchScalars(sampleScalars[:])

	// sprinkle some doublings
	for i := 10; i < 100; i++ {
		samplePoints[i] = samplePoints[0]
		sampleScalars[i] = sampleScalars[0]
	}

	// cRange is generated from template and contains the available parameters for the multiexp window size
	// for g2, CI suffers with large c size since it needs to allocate a lot of memory for the buckets.
	// test only "odd" and "even" (ie windows size divide word size vs not)
	cRange := []uint64{5, 14}

	results := make([]G2Jac, len(cRange))
	for i, c := range cRange {
		_innerMsmG2(&results[i], c, samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: runtime.NumCPU()})
	}

	var r G2Jac
	_innerMsmG2Reference(&r, samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: runtime.NumCPU()})

	var expected, got G2Affine
	expected.FromJacobian(&r)

	for i := 0; i < len(results); i++ {
		got.FromJacobian(&results[i])
		if !expected.Equal(&got) {
			t.Fatalf("cross msm failed with c=%d", cRange[i])
		}
	}

}

// _innerMsmG2Reference always do ext jacobian with c == 16
func _innerMsmG2Reference(p *G2Jac, points []G2Affine, scalars []fr.Element, config ecc.MultiExpConfig) *G2Jac {
	// partition the scalars
	digits, _ := partitionScalars(scalars, 16, config.NbTasks)

	nbChunks := computeNbChunks(16)

	// for each chunk, spawn one go routine that'll loop through all the scalars in the
	// corresponding bit-window
	// note that buckets is an array allocated on the stack and this is critical for performance

	// each go routine sends its result in chChunks[i] channel
	chChunks := make([]chan g2JacExtended, nbChunks)
	for i := 0; i < len(chChunks); i++ {
		chChunks[i] = make(chan g2JacExtended, 1)
	}

	// the last chunk may be processed with a different method than the rest, as it could be smaller.
	n := len(points)
	for j := int(nbChunks - 1); j >= 0; j-- {
		processChunk := processChunkG2Jacobian[bucketg2JacExtendedC16]
		go processChunk(uint64(j), chChunks[j], 16, points, digits[j*n:(j+1)*n], nil)
	}

	return msmReduceChunkG2Affine(p, int(16), chChunks[:])
}

func BenchmarkMultiExpG2(b *testing.B) {

	const (
		pow       = (bits.UintSize / 2) - (bits.UintSize / 8) // 24 on 64 bits arch, 12 on 32 bits
		nbSamples = 1 << pow
	)

	var (
		samplePoints             [nbSamples]G2Affine
		sampleScalars            [nbSamples]fr.Element
		sampleScalarsSmallValues [nbSamples]fr.Element
		sampleScalarsRedundant   [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	copy(sampleScalarsSmallValues[:], sampleScalars[:])
	copy(sampleScalarsRedundant[:], sampleScalars[:])

	// this means first chunk is going to have more work to do and should be split into several go routines
	for i := 0; i < len(sampleScalarsSmallValues); i++ {
		if i%5 == 0 {
			sampleScalarsSmallValues[i].SetZero()
			sampleScalarsSmallValues[i][0] = 1
		}
	}

	// bad case for batch affine because scalar distribution might look uniform
	// but over batchSize windows, we may hit a lot of conflicts and force the msm-affine
	// to process small batches of additions to flush its queue of conflicted points.
	for i := 0; i < len(sampleScalarsRedundant); i += 100 {
		for j := i + 1; j < i+100 && j < len(sampleScalarsRedundant); j++ {
			sampleScalarsRedundant[j] = sampleScalarsRedundant[i]
		}
	}

	fillBenchBasesG2(samplePoints[:])

	var testPoint G2Affine

	for i := 5; i <= pow; i++ {
		using := 1 << i

		b.Run(fmt.Sprintf("%d points", using), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				testPoint.MultiExp(samplePoints[:using], sampleScalars[:using], ecc.MultiExpConfig{})
			}
		})

		b.Run(fmt.Sprintf("%d points-smallvalues", using), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				testPoint.MultiExp(samplePoints[:using], sampleScalarsSmallValues[:using], ecc.MultiExpConfig{})
			}
		})

		b.Run(fmt.Sprintf("%d points-redundancy", using), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				testPoint.MultiExp(samplePoints[:using], sampleScalarsRedundant[:using], ecc.MultiExpConfig{})
			}
		})
	}
}

func BenchmarkMultiExpG2Reference(b *testing.B) {
	const nbSamples = 1 << 20

	var (
		samplePoints  [nbSamples]G2Affine
		sampleScalars [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	fillBenchBasesG2(samplePoints[:])

	var testPoint G2Affine

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		testPoint.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
	}
}

func BenchmarkManyMultiExpG2Reference(b *testing.B) {
	const nbSamples = 1 << 20

	var (
		samplePoints  [nbSamples]G2Affine
		sampleScalars [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	fillBenchBasesG2(samplePoints[:])

	var t1, t2, t3 G2Affine
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			t1.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
			wg.Done()
		}()
		go func() {
			t2.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
			wg.Done()
		}()
		go func() {
			t3.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
			wg.Done()
		}()
		wg.Wait()
	}
}

// WARNING: this return points that are NOT on the curve and is meant to be use for benchmarking
// purposes only. We don't check that the result is valid but just measure "computational complexity".
//
// Rationale for generating points that are not on the curve is that for large benchmarks, generating
// a vector of different points can take minutes. Using the same point or subset will bias the benchmark result
// since bucket additions in extended jacobian coordinates will hit doubling algorithm instead of add.
func fillBenchBasesG2(samplePoints []G2Affine) {
	var r big.Int
	r.SetString("340444420969191673093399857471996460938405", 10)
	samplePoints[0].ScalarMultiplication(&samplePoints[0], &r)

	one := samplePoints[0].X
	one.SetOne()

	for i := 1; i < len(samplePoints); i++ {
		samplePoints[i].X.Add(&samplePoints[i-1].X, &one)
		samplePoints[i].Y.Sub(&samplePoints[i-1].Y, &one)
	}
}

func fillBenchScalars(sampleScalars []fr.Element) {
	// ensure every words of the scalars are filled
	for i := 0; i < len(sampleScalars); i++ {
		sampleScalars[i].SetRandom()
	}
}
