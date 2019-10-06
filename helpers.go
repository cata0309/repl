package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func createReplaceRequireClause(moduleName, localPath string) string {
	if moduleName == "" {
		return ""
	}
	return fmt.Sprintf(`replace %s => %s`, moduleName, localPath)
}
func isComment(code string) bool {
	if code[:2] == "//" || code[:2] == "/*" {
		return true
	}
	return false
}
func multiplyString(s string, n int) string {
	var out string
	if n == 0 {
		return ""
	}
	for i := 0; i < n; i++ {
		out += s
	}
	return out
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
	if strings.Contains(code, "=") || strings.Contains(code, "var") || isFunctionCall(code) {
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
func goGet() error {
	return exec.Command("go", "get", "-u", "./...").Run()
}
func getModuleNameOfCurrentProject(workingDirectory string) string {
	bs, err := ioutil.ReadFile(workingDirectory + "/go.mod")
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		panic(err)
	}
	gomod := string(bs)
	moduleName := strings.Split(strings.Split(gomod, "\n")[0], " ")[1]
	return moduleName
}
func wrapInPrint(code string) string {
	return fmt.Sprintf(`fmt.Printf("<%%T> %%+v\n", %s, %s)`, code, code)
}

func checkIfErrIsNotDecl(err string) bool {
	return strings.Contains(err, "not used")
}
