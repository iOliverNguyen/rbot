package genapi

import (
	"go/types"
	"text/template"

	generator "github.com/olvrng/ggen"
	"github.com/olvrng/rbot/be/tools/genapi/defs"
)

var currentPrinter generator.Printer
var tpl = template.Must(template.New("tpl").Funcs(funcs).Parse(tplText))

var funcs = map[string]interface{}{
	"type": renderType,
	"new":  renderNew,
}

type Opts struct {
	BasePath string
}

func generateServices(printer generator.Printer, opts Opts, services []*defs.Service) error {
	currentPrinter = printer
	printer.Import("context", "context")
	printer.Import("fmt", "fmt")
	printer.Import("http", "net/http")
	printer.Import("httprpc", "github.com/olvrng/rbot/be/pkg/httprpc")
	vars := map[string]interface{}{
		"Services": services,
		"Opts":     opts,
	}
	return tpl.Execute(printer, vars)
}

func renderType(typ types.Type) string {
	return currentPrinter.TypeString(typ)
}

func renderNew(typ types.Type) string {
	named := typ.(*types.Pointer).Elem().(*types.Named)
	return "&" + currentPrinter.TypeString(named) + "{}"
}
