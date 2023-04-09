// A Requirement is a set of one or more version restrictions. It supports a
// few (<tt>=, !=, >, <, >=, <=, ~></tt>) different restriction operators.
//
// See Gem::Version for a description on how versions and requirements work
// together in RubyGems.
package requirement

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/robicode/version"
)

var (
	// Order of operations is important!
	Ops = []string{
		"!=",
		">=",
		"<=",
		"~>",
		">",
		"<",
		"=",
	}

	ops = []opFunc{
		{Op: "!=", Func: notEquals},
		{Op: ">=", Func: gte},
		{Op: "<=", Func: lte},
		{Op: "~>", Func: tildeGT},
		{Op: ">", Func: greaterThan},
		{Op: "<", Func: lessThan},
		{Op: "=", Func: equals},
	}

	quoted  string = strings.Join(Ops, "|")
	pattern string = fmt.Sprintf("\\A\\s*(%s)?\\s*(%s)\\s*\\z", quoted, version.VersionPattern)
)

type operationFunc func(rs *RequirementSpecifier, v *version.Version) bool

type opFunc struct {
	Op   string
	Func operationFunc
}

// A RequirementSpecifier specifies a single requirement as an operator and version
// pair.
type RequirementSpecifier struct {
	Operator string
	Version  *version.Version
}

// Default requirements
var defaultRequirement RequirementSpecifier = RequirementSpecifier{
	Operator: ">=",
	Version:  version.New2("0"),
}

var defaultPrereleaseRequirement RequirementSpecifier = RequirementSpecifier{
	Operator: ">=",
	Version:  version.New2("0.a"),
}

// main struct
type Requirement struct {
	requirements []*RequirementSpecifier
}

func DefaultRequirement() *RequirementSpecifier {
	return &defaultRequirement
}

func DefaultPrereleaseRequirement() *RequirementSpecifier {
	return &defaultPrereleaseRequirement
}

func New(requirements ...string) (*Requirement, error) {
	var reqs []*RequirementSpecifier

	for _, value := range requirements {
		req, err := parse(value)
		if err != nil {
			return nil, err
		}

		reqs = append(reqs, req)
	}

	return &Requirement{
		requirements: reqs,
	}, nil
}

// parse parses +requirement+, returning an *RequirementSpecifier.
// Returns nil and an error on failure.
func parse(requirement string) (*RequirementSpecifier, error) {
	reg := regexp.MustCompile(pattern)
	var operator string

	if !reg.MatchString(requirement) {
		return nil, fmt.Errorf("unable to parse requirement: '%s' with regex: %s", requirement, pattern)
	}

	var parts []string
	for _, op := range Ops {
		if strings.Contains(requirement, op) {
			parts = append(parts, op)
			p1 := strings.Replace(requirement, op, "", 1)
			parts = append(parts, p1)
			break
		}
	}

	if len(parts) != 2 {
		return nil, errors.New("requirement should be an operator and version (e.g. '> 3.0')")
	}

	for i, value := range ops {
		if parts[0] == ops[i].Op {
			operator = value.Op
			break
		}
	}

	if operator == "" {
		return nil, fmt.Errorf("invalid operator '%s'", operator)
	}

	if operator == ">=" && parts[1] == "0" {
		return DefaultRequirement(), nil
	} else if operator == ">=" && parts[1] == "0.a" {
		return DefaultPrereleaseRequirement(), nil
	}

	ver, err := version.New(parts[1])
	if err != nil {
		return nil, err
	}

	return &RequirementSpecifier{
		Operator: operator,
		Version:  ver,
	}, nil
}

func (r *Requirement) Concat(requirements ...string) *Requirement {
	for _, req := range requirements {
		for _, registeredReq := range r.requirements {
			splitReq, err := parse(req)
			if err != nil {
				return nil
			}

			if registeredReq.Operator == splitReq.Operator || splitReq.Version.Compare(registeredReq.Version) == 0 {
				continue
			}

			r.requirements = append(r.requirements, splitReq)
		}
	}

	return r
}

// HasNone returns true if this *Requirement has no requirements.
func (r *Requirement) HasNone() bool {
	if len(r.requirements) == 1 {
		return r.requirements[0].Operator == DefaultRequirement().Operator && r.requirements[0].Version.Compare(DefaultPrereleaseRequirement().Version) == 0
	}

	return false
}

// Exact returns true if the requirement is for only an exact version.
func (r *Requirement) Exact() bool {
	if len(r.requirements) != 1 {
		return false
	}

	return r.requirements[0].Operator == "="
}

// AsList returns the list of requirements as a []string.
func (r *Requirement) AsList() []string {
	var list []string

	for _, req := range r.requirements {
		list = append(list, fmt.Sprintf("%s %s", req.Operator, req.Version.Version()))
	}

	return list
}

// IsPrerelease returns true if any of the requirements are
// prerelease.
func (r *Requirement) IsPrerelease() bool {
	for _, req := range r.requirements {
		if req.Version.IsPrerelease() {
			return true
		}
	}

	return false
}

// IsSatisfiedBy returns true if the given *Version satisfies all requirements
// of the *Requirement.
func (r *Requirement) IsSatisfiedBy(v *version.Version) bool {
	for _, requirement := range r.requirements {
		if !requirement.IsSatisfiedBy(v) {
			return false
		}
	}

	return true
}

func (r *Requirement) IsSpecific() bool {
	if len(r.requirements) > 1 {
		return true
	}

	if len(r.requirements) > 0 {
		req := r.requirements[0]
		return req.Operator != ">" && req.Operator != ">="
	}

	return true
}

// ToString returns the requirements as a string.
func (r *Requirement) ToString() string {
	var _strings []string

	for _, value := range r.requirements {
		_strings = append(_strings, value.ToString())
	}

	return strings.Join(_strings, ", ")
}

// IsSatisfiedBy returns true if a given *Version satisfies this requirement.
func (rs *RequirementSpecifier) IsSatisfiedBy(v *version.Version) bool {
	for _, value := range ops {
		if value.Op == rs.Operator {
			if !value.Func(rs, v) {
				return false
			}
		}
	}

	return true
}

// ToString returns the requirement specifier as a string.
func (rs *RequirementSpecifier) ToString() string {
	return fmt.Sprintf("%s %s", rs.Operator, rs.Version.Version())
}

// Requirement operators
func equals(rs *RequirementSpecifier, v *version.Version) bool {
	return rs.Version.Compare(v) == 0
}

func notEquals(rs *RequirementSpecifier, v *version.Version) bool {
	return rs.Version.Compare(v) != 0
}

func greaterThan(rs *RequirementSpecifier, v *version.Version) bool {
	return rs.Version.Compare(v) == -1
}

func lessThan(rs *RequirementSpecifier, v *version.Version) bool {
	return rs.Version.Compare(v) == 1
}

func gte(rs *RequirementSpecifier, v *version.Version) bool {
	return rs.Version.Compare(v) == 0 || rs.Version.Compare(v) == -1
}

func lte(rs *RequirementSpecifier, v *version.Version) bool {
	return rs.Version.Compare(v) == 0 || rs.Version.Compare(v) == 1
}

func tildeGT(rs *RequirementSpecifier, v *version.Version) bool {
	r, _ := rs.Version.Bump()

	return gte(rs, v) && r.Release().Compare(v) == 1
}
