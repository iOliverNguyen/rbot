package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/olvrng/ggen"
	"github.com/olvrng/rbot/be/tools/genapi"
)

var flClean = flag.Bool("clean", false, "clean generated files without generating new files")
var flPlugins = flag.String("plugins", "", "comma separated list of plugins for generating (default to all plugins)")
var flNamespace = flag.String("namespace", "", "only parse and generate packages under this namespace (example: github.com/foo)")

func usage() {
	const text = `
Usage: ggen [OPTION] PATTERN ...

Options:
`
	fmt.Print(text[1:])
	flag.PrintDefaults()
}

func main() {
	Start(
		genapi.New(),
	)
}

func Start(plugins ...ggen.Plugin) {
	flag.Parse()
	patterns := flag.Args()
	if len(patterns) == 0 {
		usage()
		os.Exit(2)
	}

	enabledPlugins := allPluginNames(plugins)
	if *flPlugins != "" {
		enabledPlugins = strings.Split(*flPlugins, ",")
	}

	cfg := ggen.Config{
		CleanOnly:      *flClean,
		Namespace:      *flNamespace,
		EnabledPlugins: enabledPlugins,
		GoimportsArgs:  []string{}, // example: -local github.com/foo
	}

	must(ggen.RegisterPlugin(plugins...))
	must(ggen.Start(cfg, patterns...))
}

func must(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func allPluginNames(plugins []ggen.Plugin) []string {
	names := make([]string, len(plugins))
	for i, p := range plugins {
		names[i] = p.Name()
	}
	return names
}
