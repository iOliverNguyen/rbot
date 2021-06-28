package genutil

import (
	"go/types"
	"path/filepath"
	"strings"

	"github.com/olvrng/ggen"
)

var _ ggen.Qualifier = &Qualifier{}

type Qualifier struct{}

func (q Qualifier) Qualify(pkg *types.Package) string {
	alias := pkg.Name()
	if alias == "model" || alias == "types" || alias == "convert" {
		super := filepath.Base(filepath.Dir(pkg.Path()))
		alias = strings.ToLower(super) + alias
	}
	return alias
}

func HasPrefixCamel(s string, prefix string) bool {
	ln := len(prefix)
	return len(s) > ln &&
		s[:ln] == prefix &&
		!(s[ln] >= 'a' && s[ln] <= 'z')
}
