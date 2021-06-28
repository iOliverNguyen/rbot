package parse

import (
	"fmt"
	"go/constant"
	"go/types"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/olvrng/ggen"
	"github.com/olvrng/rbot/be/pkg/l"
	"github.com/olvrng/rbot/be/tools/genapi/defs"
	"github.com/olvrng/rbot/be/tools/genutil"
)

var ll = l.New()
var ls = ll.Sugar()

func Services(ng ggen.Engine, pkg *packages.Package, kinds []defs.Kind) (services []*defs.Service, _ error) {
	objects := ng.GetObjectsByPackage(pkg)
	for _, obj := range objects {
		ls.Debugf("  object %v: %v", obj.Name(), obj.Type())
		directives := ng.GetDirectives(obj)
		switch obj := obj.(type) {
		case *types.TypeName:
			ls.Debugf("  type %v", obj.Name())
			switch typ := obj.Type().(type) {
			case *types.Named:
				switch underlyingType := typ.Underlying().(type) {
				case *types.Interface:
					kind := parseKind(kinds, obj.Name())
					if kind == "" {
						ls.Debugf("ignore unrecognized interface %v", obj.Name())
						continue
					}
					methods, err := parseService(ng, underlyingType)
					if err != nil {
						return nil, ggen.Errorf(err, "service %v: %v", obj.Name(), err)
					}

					apiPath := strings.TrimPrefix(directives.GetArg("api:path"), "/")
					apiPathID := strings.Replace(apiPath, "/", "-", -1)
					service := &defs.Service{
						Kind:      kind,
						Name:      strings.TrimSuffix(obj.Name(), string(kind)),
						FullName:  obj.Name(),
						APIPath:   apiPath,
						APIPathID: apiPathID,
						Methods:   methods,
						Interface: obj,
					}
					services = append(services, service)
					for _, m := range methods {
						m.Service = service
					}
				}
			}
		}
	}
	return services, nil
}

func parseKind(kinds []defs.Kind, name string) defs.Kind {
	for _, kind := range kinds {
		suffix := string(kind)
		if strings.HasSuffix(name, suffix) {
			return kind
		}
	}
	return ""
}

func parseService(ng ggen.Engine, iface *types.Interface) ([]*defs.Method, error) {
	methods := make([]*defs.Method, 0, iface.NumMethods())
	for i, n := 0, iface.NumMethods(); i < n; i++ {
		method := iface.Method(i)
		if !method.Exported() {
			continue
		}
		m, err := parseMethod(ng, method)
		if err != nil {
			return nil, ggen.Errorf(err, "method %v: %v", method.Name(), err)
		}
		methods = append(methods, m)
	}
	return methods, nil
}

func parseMethod(ng ggen.Engine, method *types.Func) (_ *defs.Method, err error) {
	apiPath := ng.GetDirectives(method).GetArg("api:path")
	if apiPath == "" {
		apiPath = method.Name()
	}
	apiPath = strings.TrimPrefix(apiPath, "/")

	mtyp := method.Type()
	styp := mtyp.(*types.Signature)
	params := styp.Params()
	results := styp.Results()
	requests, responses, err := checkMethodSignature(method.Name(), params, results)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", method.Name(), err)
	}
	return &defs.Method{
		Name:     method.Name(),
		APIPath:  apiPath,
		Comment:  ng.GetComment(method).Text(),
		Method:   method,
		Request:  requests,
		Response: responses,
	}, nil
}

func checkMethodSignature(name string, params *types.Tuple, results *types.Tuple) (request, response *defs.Message, err error) {
	if params.Len() == 0 {
		err = ggen.Errorf(nil, "expect at least 1 param")
		return
	}
	if results.Len() == 0 {
		err = ggen.Errorf(nil, "expect at least 1 result")
		return
	}
	var requestItems, responseItems []*defs.ArgItem
	{
		t := params.At(0)
		if t.Type().String() != "context.Context" {
			err = ggen.Errorf(nil, "expect the first param is context.Context")
			return
		}
	}
	{
		t := results.At(results.Len() - 1)
		if t.Type().String() != "error" {
			err = ggen.Errorf(nil, "expect the last return value is error")
			return
		}
	}
	{
		// skip the first param (context.Context)
		for i, n := 1, params.Len(); i < n; i++ {
			arg, err2 := checkArg(params.At(i), n == 2)
			if err2 != nil {
				return nil, nil, ggen.Errorf(err2, "%v: %v", name, err2)
			}
			requestItems = append(requestItems, arg)
			if !arg.Inline && arg.Name == "" {
				return nil, nil, ggen.Errorf(err2, "%v: must provide name for param %v", name, arg.Type)
			}
		}
	}
	{
		// skip the last result (error)
		for i, n := 0, results.Len()-1; i < n; i++ {
			arg, err2 := checkArg(results.At(i), n == 2)
			if err2 != nil {
				return nil, nil, ggen.Errorf(err2, "%v: %v", name, err2)
			}
			responseItems = append(responseItems, arg)
		}
		if len(responseItems) > 1 {
			for _, arg := range responseItems {
				if arg.Name == "" || strings.HasPrefix(arg.Name, "_") {
					return nil, nil, ggen.Errorf(err, "%v: must provide name for result %v", name, arg.Type)
				}
			}
		}
	}
	request = &defs.Message{Items: requestItems}
	response = &defs.Message{Items: responseItems}
	return request, response, nil
}

func checkArg(v *types.Var, autoInline bool) (*defs.ArgItem, error) {
	arg := &defs.ArgItem{
		Inline: v.Name() == "_" || v.Name() == "" && autoInline,
		Name:   toTitle(v.Name()),
		Var:    v,
		Type:   v.Type(),
	}
	// when inline, the param must be struct or pointer to struct
	if arg.Inline {
		var err error
		arg.Struct, arg.Ptr, err = checkStruct(v.Type())
		if err != nil {
			return nil, fmt.Errorf("type must be a struct or a pointer to struct to be inline: %v", err)
		}
	}
	return arg, nil
}

func checkStruct(t types.Type) (_ *types.Struct, ptr bool, _ error) {
	p, ptr := t.(*types.Pointer)
	if ptr {
		t = p.Elem()
	}

underlying:
	switch typ := t.(type) {
	case *types.Pointer:
		return nil, false, fmt.Errorf("got double pointer (%v)", t)

	case *types.Named:
		t = typ.Underlying()
		goto underlying

	case *types.Struct:
		return typ, ptr, nil

	default:
		return nil, false, fmt.Errorf("got %v", typ)
	}
}

func toTitle(s string) string {
	s = strings.TrimPrefix(s, "_")
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func ParseEnumInPackage(ng ggen.Engine, pkg *packages.Package) (map[string]*defs.Enum, error) {
	mapEnum := make(map[string]*defs.Enum)
	objects := ng.GetObjectsByPackage(pkg)
	sort.Slice(objects, func(i, j int) bool { return objects[i].Pos() < objects[j].Pos() })

	// read all enums in the package
	for _, obj := range objects {
		_, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}
		if !obj.Exported() {
			continue
		}
		directive, ok := ng.GetDirectives(obj).Get("enum")
		if !ok {
			continue
		}
		name := obj.Name()
		if directive.Arg != "" {
			return nil, ggen.Errorf(nil, "invalid argument for enum %v", name)
		}
		basic, err := genutil.CheckType(obj.Type(), genutil.Named, genutil.Basic)
		if err != nil {
			return nil, ggen.Errorf(err, "enum %v.%v must be integer type", pkg.PkgPath, name)
		}
		kind := basic.(*types.Basic).Kind()
		switch kind {
		case types.Int, types.Uint64:
			// no-op
		default:
			return nil, ggen.Errorf(err, "enum %v.%v must be int or uint64 type (got %v)", pkg.PkgPath, name, basic.String())
		}
		mapEnum[name] = &defs.Enum{
			Name:     name,
			Type:     obj.Type().(*types.Named),
			Basic:    basic.(*types.Basic),
			MapValue: map[string]interface{}{},
			MapName:  map[interface{}]string{},
			MapLabel: map[string]map[string]string{},
			MapConst: map[string]*types.Const{},
		}
	}

	// read values for enum
	for _, obj := range objects {
		cnst, ok := obj.(*types.Const)
		if !ok {
			continue
		}
		ds := ng.GetDirectives(cnst)
		directive, ok := ng.GetDirectives(cnst).Get("enum")
		enum := validateEnumConstType(pkg.Types, mapEnum, cnst.Type())
		if ok != (enum != nil) {
			return nil, ggen.Errorf(nil, "enum constant must have directive [+enum] (%v.%v)", pkg.PkgPath, obj.Name())
		}
		if !ok {
			continue
		}
		nameList := directive.Arg
		names, ok := validateEnumNames(strings.Split(nameList, ","))
		if !ok {
			return nil, ggen.Errorf(nil, "invalid enum value %v (%v.%v)", nameList, pkg.PkgPath, obj.Name())
		}
		for _, name := range names {
			if _cnst, exists := enum.MapConst[name]; exists {
				return nil, ggen.Errorf(nil, "duplicated value of %v for enum %v.%v (%v and %v)",
					name, pkg.PkgPath, enum.Name, _cnst.Name(), cnst.Name())
			}

			var enumValue interface{}
			switch enum.Basic.Kind() {
			case types.Int:
				value, ok2 := constant.Int64Val(cnst.Val())
				if !ok2 {
					return nil, ggen.Errorf(nil, "invalid enum value %v (%v.%v)", cnst.Val(), pkg.PkgPath, cnst.Name())
				}
				enumValue = int(value)
			case types.Uint64:
				value, ok2 := constant.Uint64Val(cnst.Val())
				if !ok2 {
					return nil, ggen.Errorf(nil, "invalid enum value %v (%v.%v)", cnst.Val(), pkg.PkgPath, cnst.Name())
				}
				enumValue = value
			default:
				panic(fmt.Sprintf("unexpected kind %v", enum.Basic))
			}

			lable, err := getLabels(ds, name)
			if err != nil {
				return nil, err
			}
			enum.MapConst[name] = cnst
			enum.MapValue[name] = enumValue
			enum.MapLabel[name] = lable
			enum.Labels = addLableToList(enum.Labels, enum.MapLabel[name])
			enum.Names = append(enum.Names, name)

			// multiple name can exist for a value, only map value to the first name
			if _, exists := enum.MapName[enumValue]; !exists {
				enum.MapName[enumValue] = name
				enum.Values = append(enum.Values, enumValue)
			}
		}
	}
	return mapEnum, nil
}

func addLableToList(labels []string, mapLable map[string]string) []string {
	for key, _ := range mapLable {
		isExisted := false
		for _, value := range labels {
			if value == key && key != "" {
				isExisted = true
				break
			}
		}
		if !isExisted {
			labels = append(labels, key)
		}
	}
	return labels
}

func getLabels(ds ggen.Directives, name string) (map[string]string, error) {
	var mapLabel = make(map[string]string)
	for i := 0; i < len(ds); i++ {
		rootPoint := ds[i]
		if rootPoint.Arg == name {
			for {
				if i+1 == len(ds) {
					break
				}
				strs := strings.Split(ds[i+1].Raw, `:`)
				if strs[0] != "+enum" || len(strs) != 3 {
					break
				}
				firstNameCharacter := strs[1][:1]
				if "A" > firstNameCharacter || firstNameCharacter > "Z" {
					i++
					continue
				}
				mapLabel[strs[1]] = strs[2]
				i++
			}
			break
		}
	}
	return mapLabel, nil
}

func validateEnumConstType(pkg *types.Package, mapEnum map[string]*defs.Enum, typ types.Type) *defs.Enum {
	named, ok := typ.(*types.Named)
	if ok && named.Obj().Pkg() == pkg {
		return mapEnum[named.Obj().Name()]
	}
	return nil
}

var reEnumValue = regexp.MustCompile(`[A-z_][A-z0-9_]*`)

func validateEnumNames(values []string) ([]string, bool) {
	for i, value := range values {
		value = strings.TrimSpace(value)
		values[i] = value
		if !reEnumValue.MatchString(value) {
			return nil, false
		}
	}
	return values, true
}
