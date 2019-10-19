package interpreter

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func Test_newSession(t *testing.T) {
	monkey.Patch(createTmpDir, func(wd string) (string, error) {
		return "somedir/Interpreter", nil
	})
	monkey.Patch(os.Chdir, func(string) error {
		return nil
	})
	monkey.Patch(getModuleNameOfCurrentProject, func(string) string {
		return "tmpmodule"
	})
	monkey.Unpatch(createTmpDir)
	monkey.Unpatch(os.Chdir)
	monkey.Unpatch(getModuleNameOfCurrentProject)
}

func Test_addCode(t *testing.T) {
	s := &Interpreter{}
	s.code = []string{}
	_, err := s.Eval("fmt.Println(")
	assert.NoError(t, err)
	_, err = s.Eval(`"Salam",`)
	assert.NoError(t, err)
	_, err = s.Eval(`)`)
	assert.NoError(t, err)
	assert.Equal(t, []string{"fmt.Println(\n\"Salam\",\n)"}, s.code)
}

func Test_removeLastCode(t *testing.T) {
	s := &Interpreter{}
	s.code = append(s.code, "some ok code", "some code caused error")
	s.removeLastCode()
	assert.Equal(t, []string{"some ok code"}, s.code)
}
func Test_removeTmpCodes(t *testing.T) {
	s := &Interpreter{}
	s.code = append(s.code, `a := 1+2`)
	s.code = append(s.code, `fmt.Println("aaa")`)
	s.tmpCodes = append(s.tmpCodes, 1)
	s.removeTmpCodes()
	assert.Equal(t, []string{"a := 1+2"}, s.code)
}

func Test_add_print(t *testing.T) {
	s := &Interpreter{}
	s.Eval(`fmt.Println("Salam")`)
	assert.Equal(t, []string{`fmt.Println("Salam")`}, s.code)
	assert.Equal(t, []int{0}, s.tmpCodes)
}

func Test_add_function_call(t *testing.T) {
	s := &Interpreter{}
	s.Eval(`someFunc("salam man be to yare ghadimi")`)
	assert.Equal(t, s.code, []string{"fmt.Printf(\"<%T> %+v\\n\", someFunc(\"salam man be to yare ghadimi\"), someFunc(\"salam man be to yare ghadimi\"))"})
}
func Test_add_continue_mode(t *testing.T) {
	s := &Interpreter{}
	s.Eval("fmt.Println(")
	s.Eval("2,")
	s.Eval(")")
	assert.Equal(t, []string{"fmt.Println(\n2,\n)"}, s.code)
}

func Test_checkIfErrIsNotDecl(t *testing.T) {
	assert.True(t, checkIfErrIsNotDecl(`"fmt" imported and not used`))
	assert.True(t, checkIfErrIsNotDecl(`a declared and not used`))
	assert.False(t, checkIfErrIsNotDecl("not able to compile"))
}
func Test_createReplaceRequireClause_with_moduleName(t *testing.T) {
	moduleName := "shell"
	localPath := "inja"
	assert.Equal(t, "replace shell => inja", createReplaceRequireClause(moduleName, localPath))
}

func Test_createReplaceRequireClause_without_moduleName(t *testing.T) {
	moduleName := ""
	localPath := "inja"
	assert.Equal(t, "", createReplaceRequireClause(moduleName, localPath))
}
func Test_wrapInPrint(t *testing.T) {
	assert.Equal(t, `fmt.Printf("<%T> %+v\n", 1+2, 1+2)`, wrapInPrint("1+2"))
	assert.Equal(t, `fmt.Printf("<%T> %+v\n", "Hello", "Hello")`, wrapInPrint(`"Hello"`))

}

func Test_getModuleNameOfCurrentProject_in_go_project(t *testing.T) {
	monkey.Patch(ioutil.ReadFile, func(string) ([]byte, error) {
		return []byte(`module somemodule
go 1.13`), nil
	})
	moduleName := getModuleNameOfCurrentProject("somedir")
	assert.Equal(t, moduleName, "somemodule")
	monkey.Unpatch(ioutil.ReadFile)
}

func Test_getModuleNameOfCurrentProject_not_in_go_project(t *testing.T) {
	monkey.Patch(ioutil.ReadFile, func(string) ([]byte, error) {
		return nil, os.ErrNotExist
	})
	moduleName := getModuleNameOfCurrentProject("somedir")
	assert.Equal(t, moduleName, "")
	monkey.Unpatch(ioutil.ReadFile)
}

func Test_Integration(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	i, err := NewSession(wd)
	assert.NoError(t, err)
	_, err = i.Eval("var x =2")
	assert.NoError(t, err)
	assert.Equal(t, Var{"x", "", "2"}, i.vars["x"])

	_, err = i.Eval("x := 3")
	assert.NoError(t, err)
	assert.Equal(t, Var{"x", "", "3"}, i.vars["x"])

	_, err = i.Eval("var z,y int= 3,4")
	assert.NoError(t, err)
	assert.Equal(t, Var{"z,y", "int", "3,4"}, i.vars["z,y"])

	_, err = i.Eval("var z=2")
	assert.NoError(t, err)
	assert.Equal(t, Var{"z", "", "2"}, i.vars["z"])
	assert.Equal(t, Var{"_,y", "int", "3,4"}, i.vars["_,y"])

	_, err = i.Eval("type user struct{")
	assert.NoError(t, err)
	assert.True(t, i.continueMode)

	_, err = i.Eval("Name string")
	assert.NoError(t, err)
	assert.True(t, i.continueMode)

	_, err = i.Eval("}")
	assert.NoError(t, err)
	assert.False(t, i.continueMode)
	assert.Equal(t, "type user struct{\nName string\n}", i.types["user"])

	_, err = i.Eval(`import "fmt"`)
	assert.NoError(t, err)
	assert.Equal(t, ImportDatas{{`"fmt"`, ""}}, i.imports)

	_, err = i.Eval(`import (`)
	assert.NoError(t, err)
	assert.True(t, i.continueMode)

	_, err = i.Eval(`"os"`)
	assert.NoError(t, err)
	assert.True(t, i.continueMode)

	_, err = i.Eval(`"exec"`)
	assert.NoError(t, err)
	assert.True(t, i.continueMode)

	_, err = i.Eval(")")
	assert.NoError(t, err)
	assert.False(t, i.continueMode)
	assert.Equal(t, ImportDatas{{`"fmt"`, ""}, {`"os"`, ""}, {`"exec"`, ""}}, i.imports)

	_, err = i.Eval(":vars")
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(i.vars, Vars{
		"_,y": Var{
			"_,y", "int", "3,4",
		},
		"z": Var{
			"z", "", "2",
		},
		"x": Var{
			"x", "", "3",
		},
	}))

	out, err := i.Eval(`:types`)
	assert.NoError(t, err)
	assert.Equal(t, i.typesForSource()+"\n", out)

	out, err = i.Eval(":help")
	assert.NoError(t, err)
	assert.Equal(t, helpText+"\n", out)

	out, err = i.Eval(`:imports`)
	assert.NoError(t, err)
	assert.Equal(t, i.imports.AsDump()+"\n", out)

	_, err = i.Eval("var x int = 2")
	assert.NoError(t, err)
	assert.Equal(t, Var{"x", "int", "2"}, i.vars["x"])

	_, err = i.Eval("x+=2")
	assert.NoError(t, err)

	out = i.eval()
	assert.Empty(t, out)
	assert.Equal(t, []string{"x+=2"}, i.code)

	_, err = i.Eval(":e x")
	assert.NoError(t, err)
	assert.Equal(t, []string{"x+=2", wrapInPrint("x")}, i.code)

	out = i.eval()
	assert.Equal(t, "<int> 4\n", out)

	out, err = i.Eval(":doc fmt.Println")
	assert.NoError(t, err)

	doc, err := goDoc("fmt.Println")
	assert.NoError(t, err)
	assert.Equal(t, string(doc)+"\n", out)

	_, err = i.Eval("func Name() string{}")
	assert.NoError(t, err)
	assert.Equal(t, "func Name() string{}", i.funcs["Name"])

	out, err = i.Eval(":funcs")
	assert.NoError(t, err)
	assert.Equal(t, "Name => func Name() string{}\n", out)

	out, err = i.Eval(":dump")
	assert.NoError(t, err)
}
