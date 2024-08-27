package sumcheck

import (
	"path/filepath"

	"github.com/0x3327/gnark-crypto/internal/generator/config"
	"github.com/consensys/bavard"
)

func Generate(conf config.FieldDependency, baseDir string, bgen *bavard.BatchGenerator) error {
	entries := []bavard.Entry{
		{File: filepath.Join(baseDir, "sumcheck.go"), Templates: []string{"sumcheck.go.tmpl"}},
		{File: filepath.Join(baseDir, "sumcheck_test.go"), Templates: []string{"sumcheck.test.go.tmpl"}},
	}
	return bgen.Generate(conf, "sumcheck", "./sumcheck/template/", entries...)
}
