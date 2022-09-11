package evaluator

import (
	"fmt"

	"github.com/vancanhuit/monkey/internal/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{
					Message: fmt.Sprintf(
						"wrong number of arguments. got=%d, want=1", len(args)),
				}
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{
					Value: int64(len(arg.Value)),
				}
			default:
				return &object.Error{
					Message: fmt.Sprintf("argument to `len` not supported, got %s", arg.Type()),
				}
			}
		},
	},
}
