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

package fptower

import "github.com/0x3327/gnark-crypto/ecc/bn254/fp"

// Frobenius set z to Frobenius(x), return z
func (z *E12) Frobenius(x *E12) *E12 {
	// Algorithm 28 from https://eprint.iacr.org/2010/354.pdf
	var t [6]E2

	// Frobenius acts on fp2 by conjugation
	t[0].Conjugate(&x.C0.B0)
	t[1].Conjugate(&x.C0.B1)
	t[2].Conjugate(&x.C0.B2)
	t[3].Conjugate(&x.C1.B0)
	t[4].Conjugate(&x.C1.B1)
	t[5].Conjugate(&x.C1.B2)

	t[1].MulByNonResidue1Power2(&t[1])
	t[2].MulByNonResidue1Power4(&t[2])
	t[3].MulByNonResidue1Power1(&t[3])
	t[4].MulByNonResidue1Power3(&t[4])
	t[5].MulByNonResidue1Power5(&t[5])

	z.C0.B0 = t[0]
	z.C0.B1 = t[1]
	z.C0.B2 = t[2]
	z.C1.B0 = t[3]
	z.C1.B1 = t[4]
	z.C1.B2 = t[5]

	return z
}

// FrobeniusSquare set z to Frobenius^2(x), and return z
func (z *E12) FrobeniusSquare(x *E12) *E12 {
	// Algorithm 29 from https://eprint.iacr.org/2010/354.pdf
	z.C0.B0 = x.C0.B0
	z.C0.B1.MulByNonResidue2Power2(&x.C0.B1)
	z.C0.B2.MulByNonResidue2Power4(&x.C0.B2)
	z.C1.B0.MulByNonResidue2Power1(&x.C1.B0)
	z.C1.B1.MulByNonResidue2Power3(&x.C1.B1)
	z.C1.B2.MulByNonResidue2Power5(&x.C1.B2)

	return z
}

// FrobeniusCube set z to Frobenius^3(x), return z
func (z *E12) FrobeniusCube(x *E12) *E12 {
	// Algorithm 30 from https://eprint.iacr.org/2010/354.pdf
	var t [6]E2

	// Frobenius^3 acts on fp2 by conjugation
	t[0].Conjugate(&x.C0.B0)
	t[1].Conjugate(&x.C0.B1)
	t[2].Conjugate(&x.C0.B2)
	t[3].Conjugate(&x.C1.B0)
	t[4].Conjugate(&x.C1.B1)
	t[5].Conjugate(&x.C1.B2)

	t[1].MulByNonResidue3Power2(&t[1])
	t[2].MulByNonResidue3Power4(&t[2])
	t[3].MulByNonResidue3Power1(&t[3])
	t[4].MulByNonResidue3Power3(&t[4])
	t[5].MulByNonResidue3Power5(&t[5])

	z.C0.B0 = t[0]
	z.C0.B1 = t[1]
	z.C0.B2 = t[2]
	z.C1.B0 = t[3]
	z.C1.B1 = t[4]
	z.C1.B2 = t[5]

	return z
}

// declaring these here instead of in the functions allow to inline the calls
var nonRes1Pow1to5 [5]E2
var nonRes3Pow1To5 [5]E2

func init() {
	// (11697423496358154304825782922584725312912383441159505038794027105778954184319,303847389135065887422783454877609941456349188919719272345083954437860409601)
	nonRes3Pow1To5[0] = E2{
		A0: fp.Element{
			3914496794763385213,
			790120733010914719,
			7322192392869644725,
			581366264293887267,
		},
		A1: fp.Element{
			12817045492518885689,
			4440270538777280383,
			11178533038884588256,
			2767537931541304486,
		},
	}

	// (3772000881919853776433695186713858239009073593817195771773381919316419345261,2236595495967245188281701248203181795121068902605861227855261137820944008926)
	nonRes3Pow1To5[1] = E2{
		A0: fp.Element{
			14532872967180610477,
			12903226530429559474,
			1868623743233345524,
			2316889217940299650,
		},
		A1: fp.Element{
			12447993766991532972,
			4121872836076202828,
			7630813605053367399,
			740282956577754197,
		},
	}

	// (19066677689644738377698246183563772429336693972053703295610958340458742082029,18382399103927718843559375435273026243156067647398564021675359801612095278180)
	nonRes3Pow1To5[2] = E2{
		A0: fp.Element{
			6297350639395948318,
			15875321927225446337,
			9702569988553770230,
			805825149519570764,
		},
		A1: fp.Element{
			11117433864585119104,
			10363184613815941297,
			5420513773305887730,
			278429812070195549,
		},
	}

	// (5324479202449903542726783395506214481928257762400643279780343368557297135718,16208900380737693084919495127334387981393726419856888799917914180988844123039)
	nonRes3Pow1To5[3] = E2{
		A0: fp.Element{
			4938922280314430175,
			13823286637238282975,
			15589480384090068090,
			481952561930628184,
		},
		A1: fp.Element{
			3105754162722846417,
			11647802298615474591,
			13057042392041828081,
			1660844386505564338,
		},
	}

	// (8941241848238582420466759817324047081148088512956452953208002715982955420483,10338197737521362862238855242243140895517409139741313354160881284257516364953)
	nonRes3Pow1To5[4] = E2{
		A0: fp.Element{
			16193900971494954399,
			13995139551301264911,
			9239559758168096094,
			1571199014989505406,
		},
		A1: fp.Element{
			3254114329011132839,
			11171599147282597747,
			10965492220518093659,
			2657556514797346915,
		},
	}

	// (8376118865763821496583973867626364092589906065868298776909617916018768340080,16469823323077808223889137241176536799009286646108169935659301613961712198316)
	nonRes1Pow1to5[0] = E2{
		A0: fp.Element{
			12653890742059813127,
			14585784200204367754,
			1278438861261381767,
			212598772761311868,
		},
		A1: fp.Element{
			11683091849979440498,
			14992204589386555739,
			15866167890766973222,
			1200023580730561873,
		},
	}

	// (21575463638280843010398324269430826099269044274347216827212613867836435027261,10307601595873709700152284273816112264069230130616436755625194854815875713954)
	nonRes1Pow1to5[1] = E2{
		A0: fp.Element{
			13075984984163199792,
			3782902503040509012,
			8791150885551868305,
			1825854335138010348,
		},
		A1: fp.Element{
			7963664994991228759,
			12257807996192067905,
			13179524609921305146,
			2767831111890561987,
		},
	}

	// (2821565182194536844548159561693502659359617185244120367078079554186484126554,3505843767911556378687030309984248845540243509899259641013678093033130930403)
	nonRes1Pow1to5[2] = E2{
		A0: fp.Element{
			16482010305593259561,
			13488546290961988299,
			3578621962720924518,
			2681173117283399901,
		},
		A1: fp.Element{
			11661927080404088775,
			553939530661941723,
			7860678177968807019,
			3208568454732775116,
		},
	}

	// (2581911344467009335267311115468803099551665605076196740867805258568234346338,19937756971775647987995932169929341994314640652964949448313374472400716661030)
	nonRes1Pow1to5[3] = E2{
		A0: fp.Element{
			8314163329781907090,
			11942187022798819835,
			11282677263046157209,
			1576150870752482284,
		},
		A1: fp.Element{
			6763840483288992073,
			7118829427391486816,
			4016233444936635065,
			2630958277570195709,
		},
	}

	// (685108087231508774477564247770172212460312782337200605669322048753928464687,8447204650696766136447902020341177575205426561248465145919723016860428151883)
	nonRes1Pow1to5[4] = E2{
		A0: fp.Element{
			14515217250696892391,
			16303087968080972555,
			3656613296917993960,
			1345095164996126785,
		},
		A1: fp.Element{
			957117326806663081,
			367382125163301975,
			15253872307375509749,
			3396254757538665050,
		},
	}
}

// MulByNonResidue1Power1 set z=x*(9,1)^(1*(p^1-1)/6) and return z
func (z *E2) MulByNonResidue1Power1(x *E2) *E2 {
	// (8376118865763821496583973867626364092589906065868298776909617916018768340080,16469823323077808223889137241176536799009286646108169935659301613961712198316)
	z.Mul(x, &nonRes1Pow1to5[0])
	return z
}

// MulByNonResidue1Power2 set z=x*(9,1)^(2*(p^1-1)/6) and return z
func (z *E2) MulByNonResidue1Power2(x *E2) *E2 {
	// (21575463638280843010398324269430826099269044274347216827212613867836435027261,10307601595873709700152284273816112264069230130616436755625194854815875713954)
	z.Mul(x, &nonRes1Pow1to5[1])
	return z
}

// MulByNonResidue1Power3 set z=x*(9,1)^(3*(p^1-1)/6) and return z
func (z *E2) MulByNonResidue1Power3(x *E2) *E2 {
	// (2821565182194536844548159561693502659359617185244120367078079554186484126554,3505843767911556378687030309984248845540243509899259641013678093033130930403)
	z.Mul(x, &nonRes1Pow1to5[2])
	return z
}

// MulByNonResidue1Power4 set z=x*(9,1)^(4*(p^1-1)/6) and return z
func (z *E2) MulByNonResidue1Power4(x *E2) *E2 {
	// (2581911344467009335267311115468803099551665605076196740867805258568234346338,19937756971775647987995932169929341994314640652964949448313374472400716661030)
	z.Mul(x, &nonRes1Pow1to5[3])
	return z
}

// MulByNonResidue1Power5 set z=x*(9,1)^(5*(p^1-1)/6) and return z
func (z *E2) MulByNonResidue1Power5(x *E2) *E2 {
	// (685108087231508774477564247770172212460312782337200605669322048753928464687,8447204650696766136447902020341177575205426561248465145919723016860428151883)
	z.Mul(x, &nonRes1Pow1to5[4])
	return z
}

// MulByNonResidue2Power1 set z=x*(9,1)^(1*(p^2-1)/6) and return z
func (z *E2) MulByNonResidue2Power1(x *E2) *E2 {
	// 21888242871839275220042445260109153167277707414472061641714758635765020556617
	b := fp.Element{
		14595462726357228530,
		17349508522658994025,
		1017833795229664280,
		299787779797702374,
	}
	z.A0.Mul(&x.A0, &b)
	z.A1.Mul(&x.A1, &b)
	return z
}

// MulByNonResidue2Power2 set z=x*(9,1)^(2*(p^2-1)/6) and return z
func (z *E2) MulByNonResidue2Power2(x *E2) *E2 {
	// 21888242871839275220042445260109153167277707414472061641714758635765020556616
	b := fp.Element{
		3697675806616062876,
		9065277094688085689,
		6918009208039626314,
		2775033306905974752,
	}
	z.A0.Mul(&x.A0, &b)
	z.A1.Mul(&x.A1, &b)
	return z
}

// MulByNonResidue2Power3 set z=x*(9,1)^(3*(p^2-1)/6) and return z
func (z *E2) MulByNonResidue2Power3(x *E2) *E2 {
	// 21888242871839275222246405745257275088696311157297823662689037894645226208582
	b := fp.Element{
		7548957153968385962,
		10162512645738643279,
		5900175412809962033,
		2475245527108272378,
	}
	z.A0.Mul(&x.A0, &b)
	z.A1.Mul(&x.A1, &b)
	return z
}

// MulByNonResidue2Power4 set z=x*(9,1)^(4*(p^2-1)/6) and return z
func (z *E2) MulByNonResidue2Power4(x *E2) *E2 {
	// 2203960485148121921418603742825762020974279258880205651966
	b := fp.Element{
		8183898218631979349,
		12014359695528440611,
		12263358156045030468,
		3187210487005268291,
	}
	z.A0.Mul(&x.A0, &b)
	z.A1.Mul(&x.A1, &b)
	return z
}

// MulByNonResidue2Power5 set z=x*(9,1)^(5*(p^2-1)/6) and return z
func (z *E2) MulByNonResidue2Power5(x *E2) *E2 {
	// 2203960485148121921418603742825762020974279258880205651967
	b := fp.Element{
		634941064663593387,
		1851847049789797332,
		6363182743235068435,
		711964959896995913,
	}
	z.A0.Mul(&x.A0, &b)
	z.A1.Mul(&x.A1, &b)
	return z
}

// MulByNonResidue3Power1 set z=x*(9,1)^(1*(p^3-1)/6) and return z
func (z *E2) MulByNonResidue3Power1(x *E2) *E2 {
	// (11697423496358154304825782922584725312912383441159505038794027105778954184319,303847389135065887422783454877609941456349188919719272345083954437860409601)
	z.Mul(x, &nonRes3Pow1To5[0])
	return z
}

// MulByNonResidue3Power2 set z=x*(9,1)^(2*(p^3-1)/6) and return z
func (z *E2) MulByNonResidue3Power2(x *E2) *E2 {
	// (3772000881919853776433695186713858239009073593817195771773381919316419345261,2236595495967245188281701248203181795121068902605861227855261137820944008926)
	z.Mul(x, &nonRes3Pow1To5[1])
	return z
}

// MulByNonResidue3Power3 set z=x*(9,1)^(3*(p^3-1)/6) and return z
func (z *E2) MulByNonResidue3Power3(x *E2) *E2 {
	// (19066677689644738377698246183563772429336693972053703295610958340458742082029,18382399103927718843559375435273026243156067647398564021675359801612095278180)
	z.Mul(x, &nonRes3Pow1To5[2])
	return z
}

// MulByNonResidue3Power4 set z=x*(9,1)^(4*(p^3-1)/6) and return z
func (z *E2) MulByNonResidue3Power4(x *E2) *E2 {
	z.Mul(x, &nonRes3Pow1To5[3])
	return z
}

// MulByNonResidue3Power5 set z=x*(9,1)^(5*(p^3-1)/6) and return z
func (z *E2) MulByNonResidue3Power5(x *E2) *E2 {
	z.Mul(x, &nonRes3Pow1To5[4])
	return z
}
