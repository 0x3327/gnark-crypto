package eddsa

import (
	"path/filepath"

	"github.com/0x3327/gnark-crypto/internal/generator/config"
	"github.com/consensys/bavard"
)

func Generate(conf config.TwistedEdwardsCurve, baseDir string, bgen *bavard.BatchGenerator) error {
	// eddsa
	conf.Package = "eddsa"
	baseDir = filepath.Join(baseDir, conf.Package)

	entries := []bavard.Entry{
		{File: filepath.Join(baseDir, "doc.go"), Templates: []string{"doc.go.tmpl"}},
		{File: filepath.Join(baseDir, "eddsa.go"), Templates: []string{"eddsa.go.tmpl"}},
		{File: filepath.Join(baseDir, "eddsa_test.go"), Templates: []string{"eddsa.test.go.tmpl"}},
		{File: filepath.Join(baseDir, "marshal.go"), Templates: []string{"marshal.go.tmpl"}},
	}
	return bgen.Generate(conf, conf.Package, "./edwards/eddsa/template", entries...)

}
