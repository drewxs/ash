package evaluator

import "ash/object"

var builtins = map[string]*object.Builtin{
	"len":   object.GetBuiltinByName("len"),
	"print": object.GetBuiltinByName("print"),
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	"rest":  object.GetBuiltinByName("rest"),
	"push":  object.GetBuiltinByName("push"),
}
