// A conversion of the Rubygems Gem::Version class to Go.
//
// Some notes:
//
// I found out while writing this that adding slashes and modifiers
// to regexps in Go (e.g. /\[a-Z]+/i) will always cause match failure.
//
// A small note on the regular expressions: The Ruby version uses
// atomic groups, which re2 does not support. I opted to just remove
// the atomic operator and use the standard library. However, if you
// would like to use the original expression, another developer has
// contributed a package which has the same API as the standard library.
// To use it is simple:
//   - go get -u github.com/h2so5/goback/regexp
//   - Replace the regexp import with the above library.
//   - Replace the VersionPattern expression with the original regexp from Ruby(gems):
//     `[0-9]+(?>\.[0-9a-zA-Z]+)*(-[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?`
package version

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// The Version struct processes string versions into comparable
// values. A version string should normally be a series of numbers
// separated by periods. Each part (digits separated by periods) is
// considered its own number, and these are used for sorting. So for
// instance, 3.10 sorts higher than 3.2 because ten is greater than
// two.
//
// If any part contains letters (currently only a-z are supported) then
// that version is considered prerelease. Versions with a prerelease
// part in the Nth part sort less than versions with N-1
// parts. Prerelease parts are sorted alphabetically using the normal
// Ruby string sorting rules. If a prerelease part contains both
// letters and numbers, it will be broken into multiple parts to
// provide expected sort behavior (1.0.a10 becomes 1.0.a.10, and is
// greater than 1.0.a9).
//
// Prereleases sort between real releases (newest to oldest):
//
// 1. 1.0
// 2. 1.0.b1
// 3. 1.0.a.2
// 4. 0.9
//
// For further documentation and background, consult the Ruby Gem::Version docs.
type Version struct {
	version string
}

var (
	VersionPattern         = `[0-9]+(\.[0-9a-zA-Z]+)*(-[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?`
	VersionPatternAnchored = fmt.Sprintf(`\A\s*(%s)?\s*\z`, VersionPattern)
)

// New creates a new *Version with the given version string.
func New(version string) (*Version, error) {
	if !isCorrect(version) {
		return nil, fmt.Errorf("malformed version number string: '%s'", version)
	}

	ver := version

	if regexp.MustCompile(`/\A\s*\z/`).MatchString(version) {
		ver = "0"
	}

	ver = strings.TrimSpace(ver)
	ver = strings.ReplaceAll(ver, "-", ".pre.")

	return &Version{
		version: ver,
	}, nil
}

// New2 returns a new *Version with the given version string or nil on error
func New2(version string) *Version {
	v, err := New(version)
	if err != nil {
		return nil
	}

	return v
}

// Return a new version object where the next to the last revision
// number is one greater (e.g., 5.3.1 => 5.4).
//
// Pre-release (alpha) parts, e.g, 5.3.1.b.2 => 5.4, are ignored.
func (v *Version) Bump() (*Version, error) {
	segments := v.segments()

	for i, segment := range segments {
		if regexp.MustCompile(`[a-zA-Z]+`).MatchString(segment) {
			segments = segments[0:i]
		}
	}

	if len(segments) > 1 {
		segments = segments[:len(segments)-1]
	}

	num, err := strconv.Atoi(segments[len(segments)-1])
	if err != nil {
		panic("last element string should contain an int")
	}

	num = num + 1
	segments[len(segments)-1] = strconv.Itoa(num)

	version := strings.Join(segments, ".")

	ver, err := New(version)
	if err != nil {
		return nil, err
	}

	return ver, nil
}

// isCorrect validates the format of the version string.
func isCorrect(version string) bool {
	re := regexp.MustCompile(VersionPatternAnchored)

	return re.MatchString(version)
}

// segments splits the version string into its component parts.
func (v *Version) segments() []string {
	results := regexp.MustCompile(`[0-9]+|[a-zA-Z]+`).FindAllString(v.version, -1)
	if len(results) > 0 {
		return results
	}

	return []string{}
}

// IsPrerelease returns whether the Version is prerelease.
// A version is considered a prerelease if it contains a letter.
func (v *Version) IsPrerelease() bool {
	return regexp.MustCompile(`[a-zA-Z]`).MatchString(v.version)
}

// The release for this version (e.g. 1.2.0.a -> 1.2.0).
// Non-prerelease versions return themselves.
func (v *Version) Release() *Version {
	if !v.IsPrerelease() {
		return v
	}

	segments := v.segments()

	for i, segment := range segments {
		if regexp.MustCompile(`[a-zA-Z]+`).MatchString(segment) {
			segments = segments[0:i]
		}
	}

	newVersion, err := New(strings.Join(segments, "."))
	if err != nil {
		return nil
	}

	return newVersion

	// return strings.Join(segments, ".")
}

// A Version is only Eql() to another version if it's specified to the
// same precision. Version "1.0" is not the same as version "1".
func (v *Version) Eql(other *Version) bool {
	return v.version == other.version
}

// A recommended version for use with a ~> Requirement
func (v *Version) ApproximateRecommendation() string {
	segments := v.segments()

	for i, segment := range segments {
		if regexp.MustCompile(`[a-zA-Z]+`).MatchString(segment) {
			segments = segments[0:i]
		}
	}

	if len(segments) > 2 {
		segments = segments[0:2]
	}

	if len(segments) < 2 {
		segments = append(segments, "0")
	}

	recommendation := "~> " + strings.Join(segments, ".")

	if v.IsPrerelease() {
		recommendation += ".a"
	}

	return recommendation
}

// splitSegments splits the segments into integer and alphanumeric arrays.
func (v *Version) splitSegments() ([]string, []string) {
	var stringStart int
	segments := v.segments()

	for i, v := range segments {
		if regexp.MustCompile(`[a-zA-Z]+`).MatchString(v) {
			stringStart = i
			break
		}
	}

	if stringStart == 0 {
		stringStart = len(segments)
	}

	stringElements := segments[stringStart:]
	numericSegments := segments[0:stringStart]

	return numericSegments, stringElements
}

// reverseSlice sorts the slice s in reverse order.
func reverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// canonicalSegments is like segments, but with trailing zero segments removed.
func (v *Version) canonicalSegments() []string {
	var flattened []string

	numerics, stringset := v.splitSegments()

	reverseSlice(numerics)
	reverseSlice(stringset)

	for _, v := range numerics {
		value, _ := strconv.Atoi(v)

		if value == 0 {
			numerics = deleteArrayElement(numerics, 0)
		} else {
			break
		}
	}

	for _, v := range stringset {
		value, err := strconv.Atoi(v)
		if err != nil {
			break
		}

		if value == 0 {
			stringset = deleteArrayElement(stringset, 0)
		} else {
			break
		}
	}

	reverseSlice(numerics)
	reverseSlice(stringset)

	flattened = append(flattened, numerics...)
	flattened = append(flattened, stringset...)

	return flattened
}

// Version returns the version as a string
func (v *Version) Version() string {
	return v.version
}

// deleteArrayElement deletes the given element from the given []string.
func deleteArrayElement(arr []string, elem int) []string {
	if len(arr) == 0 {
		return arr
	}

	newArr := []string{}

	for i, v := range arr {
		if i != elem {
			newArr = append(newArr, v)
		}
	}

	return newArr
}

// Compare Compares this version with +other+ returning -1, 0, or 1 if the
// other version is larger, the same, or smaller than this
// one. Attempts to compare to something that's not a
// <tt>Gem::Version</tt> return +nil+.
func (v *Version) Compare(o *Version) int {
	l := v.canonicalSegments()
	r := o.canonicalSegments()

	if v.version == o.version || strArraysEqual(l, r) {
		return 0
	}

	lsz := len(l)
	rsz := len(r)

	var limit int

	if lsz > rsz {
		limit = lsz
	} else {
		limit = rsz
	}

	limit -= 1

	for i := 0; i <= limit; i++ {
		var li string
		if i < lsz {
			li = l[i]
		} else {
			li = "0"
		}

		var ri string
		if i < rsz {
			ri = r[i]
		} else {
			ri = "0"
		}

		if li == ri {
			continue
		}

		if extractKind(li) == reflect.String && extractKind(ri) == reflect.Int {
			return -1
		}

		if extractKind(li) == reflect.Int && extractKind(ri) == reflect.String {
			return 1
		}

		if extractKind(li) == reflect.String && extractKind(ri) == reflect.String {
			continue
		}

		lint, _ := strconv.Atoi(li)
		rint, _ := strconv.Atoi(ri)

		if lint > rint {
			return 1
		}

		return -1
	}

	return 0
}

// extractKind determines the underlying reflect.Kind of a string.
// Since wwe only deal with ints and strings, test just those two cases.
func extractKind(s string) reflect.Kind {
	if regexp.MustCompile(`[a-zA-Z]+`).MatchString(s) {
		return reflect.String
	}

	return reflect.Int
}

// strArrayEqual tests whether two []string slices are equal.
func strArraysEqual(sa1, sa2 []string) bool {
	if len(sa1) != len(sa2) {
		return false
	}

	var eq bool = true

	for i, v := range sa1 {
		if v != sa2[i] {
			eq = false
		}
	}

	return eq
}
