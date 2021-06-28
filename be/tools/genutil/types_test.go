package genutil

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

const testPath = "github.com/olvrng/rbot/tools/genutil/testdata"

var initialized bool
var testPkg *packages.Package
var testScope *types.Scope

func initOnce(t *testing.T) {
	if initialized {
		return
	}
	initialized = true

	cfg := &packages.Config{Mode: packages.LoadAllSyntax}
	pkgs, err := packages.Load(cfg, testPath)
	require.NoError(t, err)
	testPkg = pkgs[0]
	require.Equal(t, testPkg.PkgPath, testPath)
	testScope = testPkg.Types.Scope()
}

func TestExtractType(t *testing.T) {
	initOnce(t)
	namedSlice := getType(t, "NamedSliceOfPtrNamedStruct")

	t.Run("extract type", func(t *testing.T) {
		typs, inner := ExtractType(namedSlice.Underlying())
		expected := []types.Type{Slice, Pointer, Named}
		require.EqualValues(t, expected, typs)

		namedStruct := getType(t, "NamedStruct")
		require.Equal(t, namedStruct, inner)
	})
	t.Run("extract type skip named", func(t *testing.T) {
		typs, inner := ExtractTypeUnwrapNamed(namedSlice)
		expected := []types.Type{Named, Slice, Pointer, Named}
		require.EqualValues(t, expected, typs)

		namedStruct := getType(t, "NamedStruct")
		require.Equal(t, namedStruct, inner)
	})
	t.Run("extract type all", func(t *testing.T) {
		typs, inner := ExtractTypeAll(namedSlice)
		expected := []types.Type{Named, Slice, Pointer, Named, Struct}
		require.EqualValues(t, expected, typs)

		_, ok := inner.(*types.Struct)
		require.True(t, ok, "should be struct")
	})
}

func TestCheckType(t *testing.T) {
	initOnce(t)
	namedSlice := getType(t, "NamedSliceOfPtrNamedStruct")

	t.Run("no argument", func(t *testing.T) {
		result, err := CheckType(namedSlice)
		require.NoError(t, err)
		require.Equal(t, namedSlice, result)
	})
	t.Run("named (nil)", func(t *testing.T) {
		result, err := CheckType(namedSlice, Named)
		require.NoError(t, err)

		_, ok := result.(*types.Slice)
		require.True(t, ok, "should be slice")
	})
	t.Run("named match", func(t *testing.T) {
		c := testScope.Lookup("C")
		require.NotNil(t, c)
		namedSliceC := c.Type()
		require.Contains(t, namedSliceC.String(), "NamedSliceOfPtrNamedStruct")

		result, err := CheckType(namedSliceC, namedSliceC)
		require.NoError(t, err)

		_, ok := result.(*types.Slice)
		require.True(t, ok, "should be slice")
	})
	t.Run("named not match (error)", func(t *testing.T) {
		namedStruct := getType(t, "NamedStruct")

		_, err := CheckType(namedSlice, namedStruct)
		require.EqualError(t, err, "must be named type o.o/backend/tools/pkg/genutil/testdata.NamedStruct")
	})
	t.Run("named slice of ptr named struct", func(t *testing.T) {
		result, err := CheckType(namedSlice, Named, Slice, Pointer, Named, Struct)
		require.NoError(t, err)
		_ = result
	})
	t.Run("named slice of ptr named named (error)", func(t *testing.T) {
		result, err := CheckType(namedSlice, Named, Slice, Pointer, Named, Named)
		require.EqualError(t, err, "must be named slice of pointer to named named type")
		_ = result
	})
}

func TestCompatible(t *testing.T) {
	aType := testScope.Lookup("A").Type()
	bType := testScope.Lookup("B").Type()
	cType := testScope.Lookup("C").Type()
	dType := testScope.Lookup("D").Type()
	namedInt := testScope.Lookup("NamedInt").Type()
	bareInt := testScope.Lookup("I").Type()

	mustCheckType(t, aType, Slice, Pointer, Named, Struct)
	mustCheckType(t, bType, Slice, Pointer, Named, Struct)
	mustCheckType(t, cType, Named, Slice, Pointer, Named, Struct)
	mustCheckType(t, dType, Named, Slice, Pointer, Named, Struct)
	mustCheckType(t, namedInt, Named, Basic)
	mustCheckType(t, bareInt, Basic)

	require.True(t, Compatible(aType, bType), "a and b should be compatible")
	require.True(t, Compatible(aType, cType), "a and c should be compatible")
	require.True(t, Convertible(aType, bType), "a and b should be convertible")
	require.True(t, Convertible(aType, cType), "a and c should be convertible")
	require.False(t, Convertible(cType, dType), "c and d should not be convertible")
	require.False(t, Compatible(namedInt, bareInt), "NamedInt and int should not be compatible")
	require.True(t, Convertible(namedInt, bareInt), "NamedInt and int should be convertible")
}

func getType(t *testing.T, name string) *types.Named {
	obj := testScope.Lookup(name)
	require.NotNilf(t, obj, "type %v not found", name)

	_, ok := obj.(*types.TypeName)
	require.True(t, ok, "should be type name")

	typ := obj.Type()
	named, ok := typ.(*types.Named)
	require.True(t, ok, "should be named type")
	return named
}

func mustCheckType(t *testing.T, typ types.Type, typs ...types.Type) {
	_, err := CheckType(typ, typs...)
	require.NoError(t, err)
}
