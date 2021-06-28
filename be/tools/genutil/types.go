package genutil

import (
	"fmt"
	"go/types"

	"github.com/olvrng/ggen"
)

var (
	Array     types.Type = (*types.Array)(nil)
	Basic     types.Type = (*types.Basic)(nil)
	Interface types.Type = (*types.Interface)(nil)
	Map       types.Type = (*types.Map)(nil)
	Named     types.Type = (*types.Named)(nil)
	Pointer   types.Type = (*types.Pointer)(nil)
	Slice     types.Type = (*types.Slice)(nil)
	Struct    types.Type = (*types.Struct)(nil)
)

func UnwrapNamed(typ types.Type) types.Type {
	for {
		named, ok := typ.(*types.Named)
		if !ok {
			return typ
		}
		typ = named.Underlying()
	}
}

func ExtractNamed(typ types.Type) *types.Named {
	for typ != typ.Underlying() {
		named, ok := typ.(*types.Named)
		if ok {
			return named
		}
		typ = typ.Underlying()
	}
	return nil
}

func ExtractType(typ types.Type) (_ []types.Type, inner types.Type) {
	return extractType(false, typ, nil)
}

func ExtractTypeUnwrapNamed(typ types.Type) (typs []types.Type, inner types.Type) {
	for {
		named, ok := typ.(*types.Named)
		if !ok {
			break
		}
		typ = named.Underlying()
		typs = append(typs, Named)
	}
	typs, inner = extractType(false, typ, typs)
	return typs, inner
}

func ExtractTypeAll(typ types.Type) (_ []types.Type, inner types.Type) {
	return extractType(true, typ, nil)
}

func extractType(all bool, typ types.Type, typs []types.Type) (_ []types.Type, inner types.Type) {
	switch t := typ.(type) {
	case *types.Basic:
		return append(typs, Basic), t

	case *types.Interface:
		return append(typs, Interface), t

	case *types.Struct:
		return append(typs, Struct), t

	case *types.Pointer:
		return extractType(all, t.Elem(), append(typs, Pointer))

	case *types.Array:
		return extractType(all, t.Elem(), append(typs, Array))

	case *types.Slice:
		return extractType(all, t.Elem(), append(typs, Slice))

	case *types.Map:
		return extractType(all, t.Elem(), append(typs, Map, t.Key()))

	case *types.Named:
		if all {
			return extractType(all, t.Underlying(), append(typs, Named))
		}
		return append(typs, Named), t

	default:
		panic(fmt.Sprintf("unsupported type %v", t))
	}
}

func CheckType(typ types.Type, typs ...types.Type) (types.Type, error) {
	result, err := checkType(typ, typs...)
	if err != nil {
		return nil, ggen.Errorf(err, "must be %v", err)
	}
	return result, nil
}

func checkType(typ types.Type, typs ...types.Type) (types.Type, error) {
	if len(typs) == 0 {
		return typ, nil
	}
	switch typs[0] {
	case Basic:
		_, ok := typ.(*types.Basic)
		if len(typs) > 1 {
			return nil, ggen.Errorf(nil, "unexpected type after struct")
		}
		if !ok {
			return nil, ggen.Errorf(nil, "basic type")
		}
		return typ, nil

	case Named:
		named, ok := typ.(*types.Named)
		if !ok {
			return nil, ggen.Errorf(nil, "named type")
		}
		result, err := checkType(named.Underlying(), typs[1:]...)
		if err != nil {
			return nil, ggen.Errorf(err, "named %v", err)
		}
		return result, nil

	case Pointer:
		ptr, ok := typ.(*types.Pointer)
		if !ok {
			return nil, ggen.Errorf(nil, "pointer type")
		}
		result, err := checkType(ptr.Elem(), typs[1:]...)
		if err != nil {
			return nil, ggen.Errorf(err, "pointer to %v", err)
		}
		return result, nil

	case Slice:
		sl, ok := typ.(*types.Slice)
		if !ok {
			return nil, ggen.Errorf(nil, "slice type")
		}
		result, err := checkType(sl.Elem(), typs[1:]...)
		if err != nil {
			return nil, ggen.Errorf(err, "slice of %v", err)
		}
		return result, nil

	case Map:
		mp, ok := typ.(*types.Map)
		if !ok {
			return nil, ggen.Errorf(nil, "map type")
		}
		if len(typs) > 1 {
			if typs[1] != nil {
				_, err := checkType(mp.Key(), typs[1])
				if err != nil {
					return nil, ggen.Errorf(err, "map[%v] type", err)
				}
			}
			typs = typs[1:]
		}
		result, err := checkType(mp.Elem(), typs[1:]...)
		if err != nil {
			return nil, ggen.Errorf(err, "map[%v] of %v", typs[0], err)
		}
		return result, nil

	case Struct:
		if len(typs) > 1 {
			return nil, ggen.Errorf(nil, "unexpected type after struct")
		}
		st, ok := typ.(*types.Struct)
		if !ok {
			return nil, ggen.Errorf(nil, "struct type")
		}
		return st, nil
	}

	// exact basic type
	if basic, ok := typs[0].(*types.Basic); ok {
		if len(typs) > 1 {
			return nil, ggen.Errorf(nil, "unexpected type after basic")
		}
		if typ != basic {
			return nil, ggen.Errorf(nil, "basic type %v", basic.String())
		}
		return basic, nil
	}
	// exact named type
	if named, ok := typs[0].(*types.Named); ok {
		if len(typs) > 1 {
			return nil, ggen.Errorf(nil, "unexpected type after named")
		}
		if typ != named {
			return nil, ggen.Errorf(nil, "named type %v", named.String())
		}
		return named.Underlying(), nil
	}

	panic(fmt.Sprintf("invalid type %v", typs[0]))
}

func Compatible(typ1, typ2 types.Type) bool {
	if typ1 == typ2 {
		return true
	}
	typs1, inner1 := ExtractTypeUnwrapNamed(typ1)
	typs2, inner2 := ExtractTypeUnwrapNamed(typ2)
	if typs1[0] == Named && typs2[0] == Named {
		return false
	}
	if typs1[0] == Basic || typs1[1] == Basic {
		return false
	}
	return inner1 == inner2 && EqualTypes(TrimNamed(typs1), TrimNamed(typs2))
}

func Convertible(typ1, typ2 types.Type) bool {
	if typ1 == typ2 {
		return true
	}
	return Compatible(UnwrapNamed(typ1), UnwrapNamed(typ2))
}

func TrimNamed(typs []types.Type) []types.Type {
	for i := 0; i < len(typs); i++ {
		if typs[i] != Named {
			return typs[i:]
		}
	}
	return nil
}

func EqualTypes(typs1, typs2 []types.Type) bool {
	if len(typs1) != len(typs2) {
		return false
	}
	for i := 0; i < len(typs1); i++ {
		if typs1[i] != typs2[i] {
			return false
		}
	}
	return true
}
