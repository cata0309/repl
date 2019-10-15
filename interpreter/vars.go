package interpreter

import (
	"fmt"
	"go/token"
	"strings"
)

type Var struct {
	Name  string
	Type  string
	Value string
}
type Vars map[string]Var

func (v Var) String() string {
	if v.Value != "" {
		return fmt.Sprintf("%s %s = %s", v.Name, v.Type, v.Value)
	}
	return fmt.Sprintf("%s %s", v.Name, v.Type)
}

func (s *Interpreter) addVar(v Var) {
	s.vars[strings.TrimSpace(v.Name)] = v
}

func (vs Vars) String() string {
	var sets []string
	for _, v := range vs {
		sets = append(sets, v.String())
	}
	return strings.Join(sets, "\n\t")
}

func isVarDecl(code string) bool {
	tokens, _ := tokenizerAndLiterizer(code)
	for _, t := range tokens {
		if t == token.DEFINE || t == token.VAR {
			return true
		}
	}
	return false
}
func NewVar(code string) Var {
	tokens, lits := tokenizerAndLiterizer(code)
	for _, t := range tokens {
		if t == token.DEFINE {
			return extractDataFromVarWithDefine(tokens, lits)
		} else if t == token.VAR {
			return extractDataFromVarWithVar(tokens, lits)
		}
	}
	return Var{}
}

func extractDataFromVarWithVar(tokens []token.Token, lits []string) Var {
	var idents []string
	for idx, tok := range tokens {
		if tok == token.VAR {
			continue
		}
		if tok == token.ASSIGN {
			continue
		}
		if lits[idx] == " " || lits[idx] == "\n" || lits[idx] == "" {
			continue
		}
		idents = append(idents, lits[idx])
	}

	if len(idents) == 2 {
		return Var{
			Name: idents[0], Value: idents[1],
		}
	} else if len(idents) == 3 {
		return Var{
			idents[0], idents[1], idents[2],
		}
	}
	return Var{}
}
func extractDataFromVarWithDefine(tokens []token.Token, lits []string) Var {
	var idents []string
	var valueIdx int
	for idx, tok := range tokens {
		if tok == token.IDENT {
			idents = append(idents, lits[idx])
			continue
		} else if tok == token.DEFINE {
			continue
		}
		if lits[idx] == " " || lits[idx] == "\n" {
			continue
		}
		valueIdx = idx
	}
	return Var{Name: idents[0], Value: lits[valueIdx]}
}