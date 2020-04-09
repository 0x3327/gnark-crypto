// +build !amd64

// Copyright 2020 ConsenSys AG
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

// Code generated by goff (v0.2.1) DO NOT EDIT

// Package fp contains field arithmetic operations
package fp

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular,
// there is no security guarantees such as constant time implementation
// or side-channel attack resistance
// /!\ WARNING /!\

import "math/bits"

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *Element) Square(x *Element) *Element {

	var p [4]uint64

	var u, v uint64
	{
		// round 0
		u, p[0] = bits.Mul64(x[0], x[0])
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		var t uint64
		t, u, v = madd1sb(x[0], x[1], u)
		C, p[0] = madd2(m, 10917124144477883021, v, C)
		t, u, v = madd1s(x[0], x[2], t, u)
		C, p[1] = madd2(m, 13281191951274694749, v, C)
		_, u, v = madd1s(x[0], x[3], t, u)
		p[3], p[2] = madd3(m, 3486998266802970665, v, C, u)
	}
	{
		// round 1
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		u, v = madd1(x[1], x[1], p[1])
		C, p[0] = madd2(m, 10917124144477883021, v, C)
		var t uint64
		t, u, v = madd2sb(x[1], x[2], p[2], u)
		C, p[1] = madd2(m, 13281191951274694749, v, C)
		_, u, v = madd2s(x[1], x[3], p[3], t, u)
		p[3], p[2] = madd3(m, 3486998266802970665, v, C, u)
	}
	{
		// round 2
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		C, p[0] = madd2(m, 10917124144477883021, p[1], C)
		u, v = madd1(x[2], x[2], p[2])
		C, p[1] = madd2(m, 13281191951274694749, v, C)
		_, u, v = madd2sb(x[2], x[3], p[3], u)
		p[3], p[2] = madd3(m, 3486998266802970665, v, C, u)
	}
	{
		// round 3
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		C, z[0] = madd2(m, 10917124144477883021, p[1], C)
		C, z[1] = madd2(m, 13281191951274694749, p[2], C)
		u, v = madd1(x[3], x[3], p[3])
		z[3], z[2] = madd3(m, 3486998266802970665, v, C, u)
	}

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 3486998266802970665 || (z[3] == 3486998266802970665 && (z[2] < 13281191951274694749 || (z[2] == 13281191951274694749 && (z[1] < 10917124144477883021 || (z[1] == 10917124144477883021 && (z[0] < 4332616871279656263))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 4332616871279656263, 0)
		z[1], b = bits.Sub64(z[1], 10917124144477883021, b)
		z[2], b = bits.Sub64(z[2], 13281191951274694749, b)
		z[3], _ = bits.Sub64(z[3], 3486998266802970665, b)
	}
	return z

}
