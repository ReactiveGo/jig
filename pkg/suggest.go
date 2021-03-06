package pkg

import (
	"regexp"
	"strings"
)

// SuggestTypesToGenerate will suggest specialization signatures for types that are missing based on the error detected.
// e.g. if the error mentions and ObservableInt that is missing a field or method MapFloat32 then the specialization
// signature will become "ObservableInt MapFloat32" this signature can then be used by specialization code to generate
// the required type, field or method definition.
func (Package) SuggestTypesToGenerate(errs []error) []string {
	if len(errs) == 0 {
		return nil
	}

	var sigs []string
	sigmap := make(map[string]struct{})
	for _, err := range errs {
		errstr := err.Error()
		for _, exp := range reFixableErrors {
			matches := exp.FindStringSubmatch(errstr)
			if len(matches) == 5 || len(matches) == 6 {
				signature := strings.Join(matches[4:], " ")
				if _, present := sigmap[signature]; !present {
					sigmap[signature] = struct{}{}
					sigs = append(sigs, signature)
				}
				break
			}
		}
	}
	return sigs
}

var reFixableErrors = []*regexp.Regexp{
	regexp.MustCompile("^(.*):([0-9]*):([0-9]*): undeclared name: (.*)$"),
	regexp.MustCompile("^(.*):([0-9]*):([0-9]*): invalid operation: .* [(]value of type [*](.*)[)] has no field or method (.*)"),
	regexp.MustCompile("^(.*):([0-9]*):([0-9]*): invalid operation: .* [(]variable of type [*](.*)[)] has no field or method (.*)"),
	regexp.MustCompile("^(.*):([0-9]*):([0-9]*): invalid operation: .* [(]value of type (.*)[)] has no field or method (.*)"),
	regexp.MustCompile("^(.*):([0-9]*):([0-9]*): invalid operation: .* [(]variable of type (.*)[)] has no field or method (.*)"),

	// Errors are reported differently for Go 1.12
	regexp.MustCompile("^(.*):([0-9]*):([0-9]*): .* undefined [(]type [*](.*) has no field or method (.*)[)]"),
	regexp.MustCompile("^(.*):([0-9]*):([0-9]*): .* undefined [(]type (.*) has no field or method (.*)[)]"),
	regexp.MustCompile("^(.*):([0-9]*):([0-9]*): .* [(]variable of type [*](.*)[)] .*: missing method (.*)"),
	regexp.MustCompile("^(.*):([0-9]*):([0-9]*): .* [(]variable of type (.*)[)] .*: missing method (.*)"),
}
