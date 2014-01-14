package main

import (
	"io"
	"strings"
	"text/template"
	"github.com/0xfaded/go-testgen"
)

type Test struct{}

var comment = template.Must(template.New("Comment").Parse(
`// Test {{ .Comment }}
`))

var defs string =
`	type I int
	is := []int{1, 2, 3}
	f := func(int, bool) int { return 1 }
	ft := func(I, bool) int { return 1 }
	v1 := func(...int) int { return 1 }
	v2 := func(int, ...int) int { return 1 }
	vt := func(I, ...I) int { return 1 }
	e := func() {}
	s := func() int { return 1}
	m := func() (int, int) { return 1, 1 }
	mt := func() (int, I) { return 1, 1 }
`
var underscores string =
`	_ = is
	_ = f
	_ = ft
	_ = v1
	_ = v2
	_ = vt
	_ = e
	_ = s
	_ = m
	_ = mt
`
var body = template.Must(template.New("Body").Parse(
defs + `

	env := makeEnv()
	env.Types["I"] = reflect.TypeOf(I(0))
	env.Vars["is"] = reflect.ValueOf(&is)
	env.Funcs["f"] = reflect.ValueOf(f)
	env.Funcs["ft"] = reflect.ValueOf(ft)
	env.Funcs["v1"] = reflect.ValueOf(v1)
	env.Funcs["v2"] = reflect.ValueOf(v2)
	env.Funcs["vt"] = reflect.ValueOf(vt)
	env.Funcs["e"] = reflect.ValueOf(e)
	env.Funcs["s"] = reflect.ValueOf(s)
	env.Funcs["m"] = reflect.ValueOf(m)
	env.Funcs["mt"] = reflect.ValueOf(mt)

{{ if .Errors }}{{ if .TestErrs }}
	expectCheckError(t, `+"`{{ .Expr }}`"+`, env,{{ range .Errors }}
		`+"`{{ . }}`"+`,{{ end }}
	){{ else }}	_ = env{{ end }}
{{ else }}
	expectType(t, `+"`{{ .Expr }}`"+`, env, reflect.TypeOf({{ .Expr }})){{ end }}
`))

func (*Test) Package() string {
	return "eval"
}

func (*Test) Prefix() string {
	return "CheckCallExpr"
}

func (*Test) Imports() map[string]string {
	return map[string]string { "reflect": "" }
}

func (*Test) Dimensions() []testgen.Dimension {
	funs := []testgen.Element{
		{"NoArg", "s"},
		{"Fixed", "f"},
		{"FixedTyped", "ft"},
		{"Variadic1", "v1"},
		{"Variadic2", "v2"},
		{"VariadicTyped", "vt"},
	}
	arg := []testgen.Element{
		{"X", ""},
		{"Int", "1"},
		{"Float", "1.5"},
		{"Bool", "true"},
		{"IntTyped", "I(1)"},
		{"Ints", "is..."},
		{"EmptyFunc", "e()"},
		{"SingleFunc", "s()"},
		{"MultiFunc", "m()"},
		{"MultiFuncMixedTypes", "mt()"},
	}
	return []testgen.Dimension{
		funs,
		arg,
		arg,
	}
}

func (*Test) Comment(w io.Writer, elts ...testgen.Element) error {
	fun := elts[0].Name
	sep := "("
	for _, elt := range elts[1:] {
		if elt.Value != "" {
			fun += sep + elt.Value.(string)
			sep = ", "
		}
	}
	if sep == "(" {
		fun += "("
	}
	fun += ")"

	vars := map[string] interface{} {
		"Comment": fun,
	}

	return comment.Execute(w, vars)
}

func (*Test) Body(w io.Writer, elts ...testgen.Element) error {
	expr := elts[0].Value.(string)
	sep := "("
	for _, elt := range elts[1:] {
		if elt.Value != "" {
			expr += sep + elt.Value.(string)
			sep = ", "
		}
	}
	if sep == "(" {
		expr += "("
	}
	expr += ")"

	compileErrs, err := compileExprWithDefs(expr, defs + underscores)
	if err != nil {
		return err
	}

	testErrs := true
	for i := range compileErrs {
		if strings.Index(compileErrs[i], "syntax error") != -1 {
			testErrs = false
		}
		compileErrs[i] = strings.Replace(compileErrs[i], "I", "eval.I", -1)
	}

	// gc is stripping duplicate errors on variadic functions. This is understandable
	// but we want to catch all errors.
	n0, n1, n2 := elts[0].Name, elts[1].Name, elts[2].Name
	if len(compileErrs) > 0 && n0[:3] == "Var" && n1 == n2 && (n1 == "Float" || n1 == "Bool" || n1 == "IntTyped" || n1 == "SingleFunc") {
		compileErrs = append(compileErrs[:1], append([]string{compileErrs[0]}, compileErrs[1:]...)...)
	}
	// Duplicates of these errors get stripped on all functions
	if n1 == n2 && (n1 == "EmptyFunc" || n1 == "MultiFunc" || n1 == "MultiFuncMixedTypes") {
		compileErrs = append(compileErrs[:1], append([]string{compileErrs[0]}, compileErrs[1:]...)...)
	}

	vars := map[string] interface{} {
		"Expr": expr,
		"Errors": compileErrs,
		"TestErrs": testErrs,
	}

	return body.Execute(w, &vars)
}

