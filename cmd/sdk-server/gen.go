// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/OneOfOne/xast"
	"github.com/missionMeteora/apiserv"
)

var (
	_ = json.Marshal
	_ = fmt.Print

	cv reflect.Value
)

type MethodSignature struct {
	xast.Func
	Route   string
	ReqType string
}

func firstToLower(s string) string {
	r, sz := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[sz:]
}

func makeMethodSig(f *xast.Func) (ms *MethodSignature) {
	if !f.Exported() || ignoredFuncs[f.Name] {
		return
	}

	ms = &MethodSignature{
		Func: *f,
	}

	switch {
	case strings.HasPrefix(ms.Name, "Get"):
		ms.Route, ms.ReqType = ms.Name[3:], "GET"
	case strings.HasPrefix(ms.Name, "List"):
		ms.Route, ms.ReqType = ms.Name[4:], "GET"
	case strings.HasPrefix(ms.Name, "Delete"):
		ms.Route, ms.ReqType = ms.Name[6:], "DELETE"
	case strings.HasPrefix(ms.Name, "Create"):
		ms.Route, ms.ReqType = ms.Name[6:], "POST"
	case strings.HasPrefix(ms.Name, "Update"):
		ms.Route, ms.ReqType = ms.Name[6:], "PUT"
	default:
		return
	}

	ms.Route = firstToLower(ms.Route)

	return
}

var ignoredFuncs = map[string]bool{
	"RawRequest":       true,
	"RawRequestCtx":    true,
	"CurrentKey":       true,
	"AsUser":           true,
	"GetUserID":        true,
	"GetUserAPIKey":    true,
	"GetAPIVersion":    true,
	"CreateAdFromFile": true,
	"ListAdsFilter":    true,
	"GetAgencies":      true,

	"GetCampaignReport": true,
	"GetAdsReport":      true,
}

func main() {
	log.SetFlags(log.Lshortfile)
	// c := sdk.New("")
	d := xast.NewDumper(func(name, typ string) bool {
		switch typ {
		case "struct":
			return name == "Client"
		case "func":
			if strings.HasPrefix(name, "Client.") {
				n := name[7:]
				return ast.IsExported(n) && !ignoredFuncs[n]
			}
		}

		return true
	})

	if err := d.ProcessDir("github.com/missionMeteora/sdk"); err != nil {
		log.Fatal(err)
	}

	var ct *xast.Type
	for _, t := range d.Types {
		if t.Name == "Client" {
			ct = t
			break
		}
	}

	sort.Slice(ct.Methods, func(i, j int) bool {
		return ct.Methods[i].Name < ct.Methods[j].Name
	})

	var (
		allMethods []*MethodSignature
		b          strings.Builder
	)

	b.WriteString(header)
	for _, m := range ct.Methods {
		ms := makeMethodSig(m)
		if ms == nil {
			continue
		}
		switch ms.ReqType {
		case "GET", "DELETE":
			allMethods = append(allMethods, ms)
			getTmpl(&b, ms)
		case "POST", "PUT":
			allMethods = append(allMethods, ms)
			postTmpl(&b, ms)
		}
	}

	b.WriteString("\nfunc (ch *clientHandler) init() {\n")
	for _, ms := range allMethods {
		b.WriteByte('\t')
		fmt.Fprintf(&b, `ch.g.AddRoute("%s", "/%s`, ms.ReqType, ms.Route)
		for _, p := range ms.Params[1:] {
			switch p.Type {
			case "string", "time.Time":
				fmt.Fprintf(&b, "/:%s", p.Name)
			}
		}

		fmt.Fprintf(&b, `", ch.%s)`, ms.Name)
		b.WriteByte('\n')
	}
	b.WriteString("}\n")

	fmt.Println(b.String())
	// j, _ := json.MarshalIndent(d.Types[0].Methods, "", "  ")
	// println(fmt.Sprintf("%s", j))
}

func getTmpl(b *strings.Builder, ms *MethodSignature) {
	fmt.Fprintf(b, "\nfunc (ch *clientHandler) %s(ctx *apiserv.Context) apiserv.Response {\n", ms.Name)
	fmt.Fprintf(b, "\tc := ch.getClient(ctx)\n\tif ctx.Done() {\n\t\treturn nil\n\t}\n\n")

	if len(ms.Results) == 1 {
		fmt.Fprintf(b, "\terr := c.%s(context.Background()", ms.Name)
	} else {
		fmt.Fprintf(b, "\tdata, err := c.%s(context.Background()", ms.Name)
	}

	for _, f := range ms.Params[1:] {
		if f.Type == "time.Time" {
			fmt.Fprintf(b, `, sdk.DateToTime(ctx.Param("%s"))`, f.Name)
		} else {
			fmt.Fprintf(b, `, ctx.Param("%s")`, f.Name)
		}
	}

	b.WriteString(")\n")
	b.WriteString("\tif err != nil {\n\t\treturn apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)\n\t}\n\n")
	if len(ms.Results) == 1 {
		b.WriteString("\t return apiserv.RespOK\n")
	} else {
		b.WriteString("\treturn apiserv.NewJSONResponse(data)\n")
	}
	b.WriteString("}\n")
}

var _ apiserv.Context

func postTmpl(b *strings.Builder, ms *MethodSignature) {
	fmt.Fprintf(b, "\nfunc (ch *clientHandler) %s(ctx *apiserv.Context) apiserv.Response { // method:%s\n", ms.Name, ms.ReqType)
	fmt.Fprintf(b, "\tc := ch.getClient(ctx)\n\tif ctx.Done() {\n\t\treturn nil\n\t}\n\n")

	var params []string
	for _, p := range ms.Params[1:] {
		if p.Type == "string" {
			params = append(params, `, ctx.Param("`+p.Name+`")`)
			continue
		}
		if p.Type == "time.Time" {
			params = append(params, `, sdk.DateToTime(ctx.Param("`+p.Name+`"))`)
			continue

		}
		fmt.Fprintf(b, "\tvar %s *sdk.%s\n", p.Name, p.Type[1:]) // strip the *
		fmt.Fprintf(b, "\tif err := ctx.BindJSON(&%s); err != nil {\n", p.Name)
		b.WriteString("\t\treturn apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)\n\t}\n\n")

		params = append(params, ", "+p.Name)
	}

	if len(ms.Results) == 1 {
		fmt.Fprintf(b, "\terr := c.%s(context.Background()", ms.Name)
	} else {
		fmt.Fprintf(b, "\tdata, err := c.%s(context.Background()", ms.Name)
	}

	for _, p := range params {
		b.WriteString(p)
	}

	b.WriteString(")\n")
	b.WriteString("\tif err != nil {\n\t\treturn apiserv.NewJSONErrorResponse(http.StatusBadRequest, err)\n\t}\n\n")
	if len(ms.Results) == 1 {
		b.WriteString("\t return apiserv.RespOK\n")
	} else {
		b.WriteString("\treturn apiserv.NewJSONResponse(data)\n")
	}
	b.WriteString("}\n")
}

const header = `// this file is automatically generated, make sure you don't ovewrite your changes

package main

import (
	"net/http"
	"context"

	"github.com/missionMeteora/apiserv"
	"github.com/missionMeteora/sdk"
)

/*
// add this to main.go and handle the logic in it
// it may return errors using ctx
// remember to call ch.init(g)
// these are notes to myself, don't judge, I barely remember my name, ok?
type clientHandler struct{}
func (ch *clientHandler) getClient(ctx *apiserv.Context) *sdk.Client {
		return nil
}
*/

`
