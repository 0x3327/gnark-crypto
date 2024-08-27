package fri

import (
	"path/filepath"

	"github.com/0x3327/gnark-crypto/internal/generator/config"
	"github.com/consensys/bavard"
)

func Generate(conf config.Curve, baseDir string, bgen *bavard.BatchGenerator) error {

	// fri commitment scheme
	conf.Package = "fri"
	entries := []bavard.Entry{
		{File: filepath.Join(baseDir, "doc.go"), Templates: []string{"doc.go.tmpl"}},
		{File: filepath.Join(baseDir, "fri.go"), Templates: []string{"fri.go.tmpl"}},
		{File: filepath.Join(baseDir, "fri_test.go"), Templates: []string{"fri.test.go.tmpl"}},
	}
	return bgen.Generate(conf, conf.Package, "./fri/template/", entries...)

}
