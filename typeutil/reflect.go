package typeutil

//go:generate peg -inline sigparse.peg

import (
	"fmt"
	"reflect"
)

type functionSignatureGrammarMixin struct {
}

func (self *functionSignatureGrammarMixin) Hello() string {
	return Dump(self)
}

type TypeDeclaration string

func (self TypeDeclaration) String() string {
	return string(self)
}

func (self TypeDeclaration) IsSameTypeAs(value any) bool {
	if self == `any` {
		return true
	}

	var t reflect.Type

	if vv, ok := value.(reflect.Value); ok && vv.IsValid() {
		t = vv.Type()
	} else if tt, ok := value.(reflect.Type); ok {
		t = tt
	} else {
		t = reflect.TypeOf(value)
	}

	if t != nil && self != `` {
		return (t.String() == self.String())
	}

	return false
}

// Returns the number of input and return arguments a given function has.
func FunctionArity(fn any) (int, int, error) {
	if IsFunction(fn) {
		var fnT = reflect.TypeOf(fn)

		return fnT.NumIn(), fnT.NumOut(), nil
	} else {
		return 0, 0, fmt.Errorf("expected function, got %T", fn)
	}
}

// Parse the given function signature string and return the function name, input, and output arguments.
// Example: "helloWorld(string) error" would return an ident of "helloWorld", a 1-element type declaration
// representing the "string" argument, and a 1-element returns array with the "error" return parameter.
func ParseSignatureString(signature string) (ident string, args []TypeDeclaration, returns []TypeDeclaration, perr error) {
	var grammar = &typeutilFunctionSignatureSpec{
		Buffer: signature,
		Pretty: true,
	}

	if err := grammar.Init(); err != nil {
		perr = err
		return
	}

	if err := grammar.Parse(); err != nil {
		perr = err
		return
	}

	var decls = make([]TypeDeclaration, 0)

	for _, token := range grammar.Tokens() {
		switch rule := token.pegRule; rule {
		case ruleKW_FUNC:
			ident = `(anonymous)`
		case ruleIDENT:
			ident = signature[token.begin:token.end]
		case ruleSIGNATURE:
			args = make([]TypeDeclaration, len(decls))
			copy(args, decls)
			decls = make([]TypeDeclaration, 0)
		case ruleRETURNS:
			returns = make([]TypeDeclaration, len(decls))
			copy(returns, decls)
			decls = nil
		case ruleDATATYPE:
			decls = append(decls, TypeDeclaration(signature[token.begin:token.end]))
		}
	}

	return
}

// Returns whether the given function's actual signature matches the given spec string (as parsed by
// ParseSignatureString).
func FunctionMatchesSignature(fn any, signature string) error {
	fn = ResolveValue(fn)
	var fnT = reflect.ValueOf(fn).Type()

	if fnT.Kind() != reflect.Func {
		return fmt.Errorf("expected function, got %T", fn)
	}

	if _, args, returns, err := ParseSignatureString(signature); err == nil {
		if len(args) != fnT.NumIn() {
			return fmt.Errorf("expected %d arguments, got %d", len(args), fnT.NumIn())
		}

		if len(returns) != fnT.NumOut() {
			return fmt.Errorf("expected %d return arguments, got %d", len(returns), fnT.NumOut())
		}

		for i, arg := range args {
			var fnArg = fnT.In(i)

			if !arg.IsSameTypeAs(fnArg) {
				return fmt.Errorf("argument %d type mismatch: expected %v, got %v", i, arg, fnArg)
			}
		}

		for i, arg := range returns {
			var fnArg = fnT.Out(i)

			if !arg.IsSameTypeAs(fnArg) {
				return fmt.Errorf("return argument %d type mismatch: expected %v, got %v", i, arg, fnArg)
			}
		}

		return nil
	} else {
		return err
	}
}
