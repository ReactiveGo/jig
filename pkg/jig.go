package pkg

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/reactivego/jig/templ"
)

// jig contains everything from templ.Template and transient stuff not needed after analyzing all ast files.
type jig struct {
	templ.Generic
	Pos    token.Pos
	End    token.Pos
	Source string

	// common template is added to Needs [] of every other jig.
	common bool
}

func newJig(packageName string, cgroup *ast.CommentGroup) *jig {
	jig := &jig{}
	jig.PackageName = packageName
	jig.Pos = cgroup.End()
	for _, comment := range cgroup.List {
		if comment.Text == jigCommon {
			jig.common = true
			continue
		}

		kvmatch := reCommentPragma.FindStringSubmatch(comment.Text)
		if len(kvmatch) == 3 {
			switch kvmatch[1] {
			case jigTemplate:
				jig.Name = kvmatch[2]
			case jigNeeds:
				needs := strings.Split(kvmatch[2], ",")
				for _, need := range needs {
					jig.Needs = append(jig.Needs, strings.TrimSpace(need))
				}
			case jigEmbeds:
				embeds := strings.Split(kvmatch[2], ",")
				for _, embed := range embeds {
					embed = strings.TrimSpace(embed)
					jig.Embeds = append(jig.Embeds, embed)
					// If it embeds it, then it also needs it.
					jig.Needs = append(jig.Needs, embed)
				}
			case jigRequiredVars:
				requiredVars := strings.Split(kvmatch[2], ",")
				for _, required := range requiredVars {
					required = strings.TrimSpace(required)
					jig.RequiredVars = append(jig.RequiredVars, required)
				}
			}
		}
	}

	// Extract template vars e.g. [["<Foo>" "Foo"] ["<Bar>" "Bar"]]
	tplvars := reTemplateVar.FindAllStringSubmatch(jig.Name, -1)

	// Assign the template vars used in this template.
	// e.g. for tplvars == [["<Foo>" "Foo"] ["<Bar>" "Bar"]] assign Vars = ["Foo","Bar"]
	for _, tplvar := range tplvars {
		jig.Vars = append(jig.Vars, tplvar[1])
	}

	return jig
}

func (jig *jig) Close(pos token.Pos) {
	if jig != nil {
		jig.End = pos
	}
}

func (jig *jig) ContainsSourceRange(pos, end token.Pos) bool {
	return pos > jig.Pos && end < jig.End
}

func (jig *jig) AddSource(source string) {
	if jig.Source != "" {
		jig.Source += "\n"
	}
	jig.Source += source
	jig.Source += "\n"
}
