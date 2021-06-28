package defs

import (
	"go/types"
)

type Kind string

const (
	KindService = "Service"
)

type Service struct {
	Kind      Kind
	Name      string
	FullName  string
	APIPath   string
	APIPathID string
	Methods   []*Method

	Interface *types.TypeName

	APIPath2 string
}

type Method struct {
	Service  *Service
	Name     string
	APIPath  string
	Comment  string
	Request  *Message
	Response *Message

	Method *types.Func
}

type Message struct {
	Items ArgItems
}

type ArgItems []*ArgItem

type ArgItem struct {
	Inline bool
	Name   string
	Type   types.Type
	Var    *types.Var
	Ptr    bool
	Struct *types.Struct
}

type Enum struct {
	Name string

	// sorted values as appear in code
	Values []interface{}
	Names  []string
	Labels []string

	MapValue map[string]interface{} // int or uint64
	MapName  map[interface{}]string
	MapLabel map[string]map[string]string

	Type     *types.Named
	Basic    *types.Basic
	MapConst map[string]*types.Const
}

type NodeType int

const (
	NodeNone = iota
	NodeField
	NodeStartInline
	NodeEndInline
)

type WalkFunc func(node NodeType, name string, field *types.Var, tag string) error

func (args ArgItems) Walk(fn WalkFunc) error {
	for _, arg := range args {
		if arg.Inline {
			s := arg.Struct
			if err := fn(NodeStartInline, arg.Name, arg.Var, ""); err != nil {
				return err
			}
			for i, n := 0, s.NumFields(); i < n; i++ {
				field := s.Field(i)
				if err := fn(NodeField, field.Name(), field, s.Tag(i)); err != nil {
					return err
				}
			}
			if err := fn(NodeEndInline, arg.Name, arg.Var, ""); err != nil {
				return err
			}
		} else {
			if err := fn(NodeField, arg.Name, arg.Var, ""); err != nil {
				return err
			}
		}
	}
	return nil
}
