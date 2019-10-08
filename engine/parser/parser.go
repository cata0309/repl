package parser

import (
	"regexp"
	"strings"
)

type StmtType uint8

const (
	StmtTypeFunctionCall = iota
	StmtTypeImport
	StmtTypePrint
	StmtTypeComment
	StmtTypeExpr
	StmtTypeTypeDecl
	StmtTypeVarDecl
	StmtTypeFuncDecl
	StmtUnknown
	StmtEmpty
)

func Parse(code string) (StmtType, error) {
	if isEmpty(code) {
		return StmtEmpty, nil
	} else if isComment(code) {
		return StmtTypeComment, nil
	} else if isImport(code) {
		return StmtTypeImport, nil
	} else if isFunc(code) {
		return StmtTypeFuncDecl, nil
	} else if isTypeDecl(code) {
		return StmtTypeTypeDecl, nil
	} else if isPrint(code) {
		return StmtTypePrint, nil
	} else if isComment(code) {
		return StmtTypeComment, nil
	} else if isExpr(code) {
		return StmtTypeExpr, nil
	} else {
		return StmtUnknown, nil
	}
}
func ShouldContinue(code string) (int, bool) {
	var stillOpenChars int
	for _, c := range code {
		if c == '{' || c == '(' {
			stillOpenChars++
			continue
		}
		if c == '}' || c == ')' {
			stillOpenChars--
		}
	}
	return stillOpenChars, stillOpenChars > 0
}
func isEmpty(code string) bool {
	return len(code) == 0
}
func isComment(code string) bool {
	if len(code) < 2 {
		return false
	}
	if code[:2] == "//" || code[:2] == "/*" {
		return true
	}
	return false
}

func isShellCommand(code string) bool {
	if len(code) == 0 {
		return false
	}
	return code[0] == ':'
}

func isTypeDecl(code string) bool {
	matched, err := regexp.Match("type .+", []byte(code))
	if err != nil {
		return false
	}
	return matched
}
func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}
func isFunctionCall(code string) bool {
	m, err := regexp.Match("^[a-zA-Z0-9_.-]+\\(.*\\)", []byte(code))
	if err != nil {
		return false
	}
	return m && strings.ContainsAny(code, "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm")
}

func isExpr(code string) bool {
	if (strings.Contains(code, "=") && !strings.Contains(code, "==")) || strings.Contains(code, "var") || isFunctionCall(code) {
		return false
	}
	return true
}
func isFunc(code string) bool {
	matched, err := regexp.Match("^func.+", []byte(code))
	if err != nil {
		return false
	}
	return matched
}
func isImport(im string) bool {
	matched, err := regexp.Match("import .+", []byte(im))
	if err != nil {
		panic(err)
	}
	return matched
}
func isPrint(code string) bool {
	matched1, err := regexp.Match("^fmt.Print.*\\(.*\\)", []byte(code))
	if err != nil {
		panic(err)
	}
	matched2, err := regexp.Match("^print(ln|f).*", []byte(code))
	if err != nil {
		panic(err)
	}
	return matched1 || matched2
}
