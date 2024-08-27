package mimc

import (
	"os"
	"path/filepath"

	"github.com/0x3327/gnark-crypto/internal/generator/config"
	"github.com/consensys/bavard"
)

func Generate(conf config.Curve, baseDir string, bgen *bavard.BatchGenerator) error {

	conf.Package = "mimc"
	entries := []bavard.Entry{
		{File: filepath.Join(baseDir, "doc.go"), Templates: []string{"doc.go.tmpl"}},
		{File: filepath.Join(baseDir, "mimc.go"), Templates: []string{"mimc.go.tmpl"}},
		{File: filepath.Join(baseDir, "options.go"), Templates: []string{"options.go.tmpl"}},
	}
	os.Remove(filepath.Join(baseDir, "utils.go"))
	os.Remove(filepath.Join(baseDir, "utils_test.go"))

	return bgen.Generate(conf, conf.Package, "./crypto/hash/mimc/template", entries...)

}
