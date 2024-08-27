package pairing

import (
	"path/filepath"
	"strings"

	"github.com/0x3327/gnark-crypto/internal/generator/config"
	"github.com/consensys/bavard"
)

func Generate(conf config.Curve, baseDir string, bgen *bavard.BatchGenerator) error {

	packageName := strings.ReplaceAll(conf.Name, "-", "")
	return bgen.Generate(conf, packageName, "./pairing/template", bavard.Entry{
		File: filepath.Join(baseDir, "pairing_test.go"), Templates: []string{"tests/pairing.go.tmpl"},
	})

}
