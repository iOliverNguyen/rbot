package genapi

import (
	"golang.org/x/tools/go/packages"

	"github.com/olvrng/ggen"
	"github.com/olvrng/rbot/be/pkg/l"
	"github.com/olvrng/rbot/be/tools/genapi/defs"
	"github.com/olvrng/rbot/be/tools/genapi/parse"
	"github.com/olvrng/rbot/be/tools/genutil"
)

var ll = l.New()
var ls = ll.Sugar()
var _ ggen.Plugin = &plugin{}

type plugin struct {
	ggen.Filterer
	ggen.Qualifier
	ng ggen.Engine
}

func New() ggen.Plugin {
	return &plugin{
		Filterer:  ggen.FilterByCommand("gen:api"),
		Qualifier: genutil.Qualifier{},
	}
}

func (p *plugin) Name() string { return "api" }

func (p *plugin) Generate(ng ggen.Engine) error {
	p.ng = ng
	return ng.GenerateEachPackage(p.generatePackage)
}

func (p *plugin) generatePackage(ng ggen.Engine, pkg *packages.Package, printer ggen.Printer) (_err error) {
	ls.Debugf("api: generating package %v", pkg.PkgPath)

	pkgDirectives := ng.GetDirectivesByPackage(pkg)
	basePath := pkgDirectives.GetArg("gen:api:base-path")
	if basePath == "" {
		basePath = "/api"
	}
	opts := Opts{
		BasePath: basePath,
	}

	services, err := parse.Services(ng, pkg, []defs.Kind{defs.KindService})
	if err != nil {
		return err
	}
	for _, service := range services {
		if service.APIPath == "" {
			return ggen.Errorf(nil, "no api path for %v", service.Name)
		}
	}
	if err2 := generateServices(printer, opts, services); err2 != nil {
		return err2
	}
	return nil
}
