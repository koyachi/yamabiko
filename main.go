/*
- echo-client-gen
- yamabiko 山彦
- kodama
*/
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	//parseHandlerFunc(src, "echo_sample.go")
	parseRoute(srcRoute, "route.go")
}

func parseRoute(src string, filename string) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, filename, src, parser.Mode(0))

	for _, d := range f.Decls {
		ast.Print(fset, d)
		fmt.Println()

		echoInstanceName := ""
		if function, ok := d.(*ast.FuncDecl); ok {
			for _, l := range function.Body.List {
				if assignStmt, ok := l.(*ast.AssignStmt); ok {
					callExpr := assignStmt.Rhs[0].(*ast.CallExpr)
					selectorExpr := callExpr.Fun.(*ast.SelectorExpr)
					x := selectorExpr.X.(*ast.Ident).Name
					sel := selectorExpr.Sel.Name
					tok := assignStmt.Tok
					if x == "echo" && sel == "New" && tok == token.DEFINE {
						v := assignStmt.Lhs[0].(*ast.Ident).Name
						fmt.Printf("found echo.New(). echo instance name = %s\n", v)
						echoInstanceName = v
					}
				}
				if exprStmt, ok := l.(*ast.ExprStmt); ok {
					if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
						f := callExpr.Fun.(*ast.SelectorExpr)
						if f.X.(*ast.Ident).Name == echoInstanceName {
							args := exprStmt.X.(*ast.CallExpr).Args
							pathLit := args[0].(*ast.BasicLit).Value
							//funcLit := args[1].(*ast.FuncLit)
							/*
								switch f.Sel.Name {
								case "GET":
								case "PUT":
								case "POST":
								case "DELETE":
								default:
									break
								}
							*/
							fmt.Printf("method = %s, path = %s\n", f.Sel.Name, pathLit)
						}
					}
				}
			}
		}
	}
}

var srcRoute = `
package main

import (
	"net/http"
	
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}`

func parseHandlerFunc(src string, filename string) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, filename, src, parser.Mode(0))
	//fmt.Printf("%v\n", f)
	//ast.Print(fset, f)
	//	return

	for _, d := range f.Decls {
		//ast.Print(fset, d)
		//fmt.Println()

		// TODO: inspect
		if function, ok := d.(*ast.FuncDecl); ok {
			var found = false
			for _, p := range function.Type.Params.List {
				if sel, ok := p.Type.(*ast.SelectorExpr); ok {
					if sel.X.(*ast.Ident).Name == "echo" && sel.Sel.Name == "Context" {
						fmt.Printf("-------------------------\n")
						fmt.Printf("Found! %s\n", function.Name.Name)
						found = true
					}
				}
			}
			if found {
				for _, l := range function.Body.List {
					if assignStmt, ok := l.(*ast.AssignStmt); ok {
						callExpr := assignStmt.Rhs[0].(*ast.CallExpr)
						selectorExpr := callExpr.Fun.(*ast.SelectorExpr)
						x := selectorExpr.X.(*ast.Ident).Name
						sel := selectorExpr.Sel.Name
						if x == "c" && sel == "Param" {
							fmt.Printf("param = %s\n", callExpr.Args[0].(*ast.BasicLit).Value)
						}
					}
				}
			}
		}

	}
}

var src = `package p

import (
	"github.com/labstack/echo"
	"net/http"
)

type response struct {
	ID string ` + "`json:\"id\"`" + `
	CreatedAt string ` + "`json:\"createdAt\"`" + `
}

func HandlerX(c echo.Context) {
	x := c.Param("x")
}

func HandlerY(c echo.Context) {
	y := c.Param("y")
}
// return c.JSON()がNG?

func HandlerAx(c echo.Context) error {
	x := c.Param("x")
	y := c.Param("y")
	return c.JSON(http.StatusOK, response{ID: "1", CreatedAt: "2017-06-14"}})
}

func HandlerAdash(c echo.Context, i string) error {
	x1 := c.Param("x1")
	y1 := c.Param("y1")
	return c.JSON(http.StatusOK, response{ID: "2", CreatedAt: "2017-06-14"}})
}

/*
func HandlerB(c echo.Context, ctx *usecases.Context) error {
	a := c.Param("a")
	b := c.Param("b")
	return c.JSON(http.StatusOK, struct{foo: "123"})
}
*/
`
