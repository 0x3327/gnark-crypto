package sis

import (
	"math/bits"
	"path/filepath"

	"github.com/0x3327/gnark-crypto/internal/generator/config"
	"github.com/consensys/bavard"
)

func Generate(conf config.Curve, baseDir string, bgen *bavard.BatchGenerator) error {

	conf.Package = "sis"
	entries := []bavard.Entry{
		{File: filepath.Join(baseDir, "sis_fft.go"), Templates: []string{"fft.go.tmpl"}},
	}

	funcs := make(map[string]interface{})
	funcs["bitReverse"] = func(n, i int64) uint64 {
		nn := uint64(64 - bits.TrailingZeros64(uint64(n)))
		r := make([]uint64, n)
		for i := 0; i < len(r); i++ {
			r[i] = uint64(i)
		}
		for i := 0; i < len(r); i++ {
			irev := bits.Reverse64(r[i]) >> nn
			if irev > uint64(i) {
				r[i], r[irev] = r[irev], r[i]
			}
		}
		return r[i]
	}

	bavardOpts := []func(*bavard.Bavard) error{bavard.Funcs(funcs)}

	return bgen.GenerateWithOptions(conf, conf.Package, "./sis/template/", bavardOpts, entries...)
}
