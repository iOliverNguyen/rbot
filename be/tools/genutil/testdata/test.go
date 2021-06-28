package testdata

type NamedInt int
type NamedStruct struct{}
type NamedPtrNamedStruct *NamedStruct
type NamedSliceOfPtrNamedStruct []*NamedStruct
type MapNamedIntOfPtrNamedStruct map[NamedInt]*NamedStruct

type NamedStruct2 struct{}
type NamedSliceOfPtrNamedStruct2 []*NamedStruct2

// A, B and C should be compatible

var I int
var A []*NamedStruct
var B []*NamedStruct
var C NamedSliceOfPtrNamedStruct
var D NamedSliceOfPtrNamedStruct2
