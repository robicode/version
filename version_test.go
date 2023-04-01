package version

import (
	"reflect"
	"testing"
)

// structure used in most tests
type VersionTest struct {
	Version                   string
	ExpectedVersion           string
	ExpectedSegments          []string
	ExpectedCanonicalSegments []string
	ExpectedNumericSegments   []string
	ExpectedStringSegments    []string
	ExpectedResponse          bool
}

// Table tests used in most of the Test_XXX functions
var versionTests = []VersionTest{
	{
		Version:                   "1",
		ExpectedVersion:           "1",
		ExpectedSegments:          []string{"1"},
		ExpectedCanonicalSegments: []string{"1"},
		ExpectedNumericSegments:   []string{"1"},
		ExpectedStringSegments:    []string{},
		ExpectedResponse:          true},
	{
		Version:                   "1.2",
		ExpectedVersion:           "1.2",
		ExpectedSegments:          []string{"1", "2"},
		ExpectedCanonicalSegments: []string{"1", "2"},
		ExpectedNumericSegments:   []string{"1", "2"},
		ExpectedStringSegments:    []string{},
		ExpectedResponse:          true},
	{
		Version:          "1.",
		ExpectedResponse: false},
	{
		Version:                   "1.3.1",
		ExpectedVersion:           "1.3.1",
		ExpectedSegments:          []string{"1", "3", "1"},
		ExpectedCanonicalSegments: []string{"1", "3", "1"},
		ExpectedNumericSegments:   []string{"1", "3", "1"},
		ExpectedStringSegments:    []string{},
		ExpectedResponse:          true},
	{
		Version:          "1.5-",
		ExpectedResponse: false},
	{
		Version:                   "2.3-0-0",
		ExpectedVersion:           "2.3.pre.0.pre.0",
		ExpectedSegments:          []string{"2", "3", "pre", "0", "pre", "0"},
		ExpectedCanonicalSegments: []string{"2", "3", "pre", "0", "pre"},
		ExpectedNumericSegments:   []string{"2", "3"},
		ExpectedStringSegments:    []string{"pre", "0", "pre", "0"},
		ExpectedResponse:          true},
	{
		Version:                   "1.5-3",
		ExpectedResponse:          true,
		ExpectedSegments:          []string{"1", "5", "pre", "3"},
		ExpectedCanonicalSegments: []string{"1", "5", "pre", "3"},
		ExpectedNumericSegments:   []string{"1", "5"},
		ExpectedStringSegments:    []string{"pre", "3"},
		ExpectedVersion:           "1.5.pre.3",
	},
}

// New creates a newly prepared *Version.
func Test_New(t *testing.T) {
	for _, test := range versionTests {
		v, err := New(test.Version)
		if test.ExpectedResponse {
			if err != nil {
				t.Error("expected err to be nil")
			}
			if v == nil {
				t.Error("expected v not to be nil")
			}
		} else {
			if err == nil {
				t.Error("expected New() to return an error")
			}
			if v != nil {
				t.Error("expected v to be nil")
			}
		}
	}
}

// isCorrect validates the format of the version string.
func Test_IsCorrect(t *testing.T) {
	for _, test := range versionTests {
		if isCorrect(test.Version) != test.ExpectedResponse {
			if test.ExpectedResponse {
				t.Error("expected ", test.Version, "to be correct")
			} else {
				t.Error("expected ", test.Version, "to be incorrect")
			}
		}
	}
}

// Compare Compares this version with +other+ returning -1, 0, or 1 if the
// other version is larger, the same, or smaller than this
// one. Attempts to compare to something that's not a
// <tt>Gem::Version</tt> return +nil+.
func Test_Compare(t *testing.T) {
	version1, _ := New("1.1")
	older, _ := New("1.0")

	if version1.Compare(older) != 1 {
		t.Error("expected second Version to be older")
		t.Fail()
		return
	}

	newer, _ := New("1.2")

	if version1.Compare(newer) != -1 {
		t.Error("newer *Version should be newer and thus Compare should return -1")
		t.Fail()
		return
	}

	same, _ := New("1.1")

	if version1.Compare(same) != 0 {
		t.Error("Compare should == 0, indicating the same version")
		t.Fail()
		return
	}

	version1, err := New("1.1-1-1")
	if err != nil {
		t.Error("1.1-1-1 should be a valid version")
		t.Fail()
		return
	}

	version2, err := New("1.1-1-2")
	if err != nil {
		t.Error("1.1-1-2 should be a valid version")
		t.Fail()
		return
	}

	if version1.Compare(version2) != -1 {
		t.Error(version2.Version(), "should be greater than", version1.Version())
		t.Fail()
		return
	}

	version2, err = New("1.1-2-1")
	if err != nil {
		t.Error("1.1-2-1 should be a valid version")
		t.Fail()
		return
	}

	if version1.Compare(version2) != -1 {
		t.Error(version2.Version(), "should be greater than", version1.Version())
		t.Fail()
		return
	}

	version2, err = New("1.0-3-1")
	if err != nil {
		t.Error("1.0-3-1 should be a valid version")
		t.Fail()
		return
	}

	if version1.Compare(version2) != 1 {
		t.Error(version2.Version(), "should be less than", version1.Version())
		t.Fail()
		return
	}

	version4, err := New("1.12")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if version1.Compare(version4) != -1 {
		t.Error("version4 should be > version1")
	}
}

// splitSegments splits the segments into integer and alphanumeric arrays.
func Test_SplitSegments(t *testing.T) {
	for _, test := range versionTests {
		if !test.ExpectedResponse {
			continue
		}

		v, err := New(test.Version)
		if err != nil {
			t.Error("testing bug: version should be valid for test:", test.Version)
			t.Fail()
			return
		}

		numerics, stringset := v.splitSegments()

		for i, segment := range numerics {
			if segment != test.ExpectedNumericSegments[i] {
				t.Error("expected numeric segment", i, "to be", test.ExpectedNumericSegments[i], "but was", segment)
				t.Fail()
				return
			}
		}

		for i, segment := range stringset {
			if segment != test.ExpectedStringSegments[i] {
				t.Error("expected string segment", i, "to be", test.ExpectedStringSegments[i], "but was", segment)
				t.Fail()
				return
			}
		}
	}
}

// canonicalSegments is like segments, but with trailing zero segments removed.
func Test_CanonicalSegments(t *testing.T) {
	for _, test := range versionTests {
		if !test.ExpectedResponse {
			continue
		}

		v, _ := New(test.Version)
		segments := v.canonicalSegments()
		validSegments := test.ExpectedCanonicalSegments

		if len(segments) != len(validSegments) {
			t.Errorf("expected length of segments (%d) and validSegments (%d) to match", len(segments), len(validSegments))
			t.Error("canonicalSegments:", segments)
			return
		}

		for i, segment := range segments {
			if segment != validSegments[i] {
				t.Error("segment mismatch:: expecetd segment ", i, " to be ", validSegments[i], " but was ", segment)
			}
		}
	}
}

// segments splits the version string into its component parts.
func Test_Segments(t *testing.T) {
	for _, test := range versionTests {
		if !test.ExpectedResponse {
			continue
		}

		v, _ := New(test.Version)
		segments := v.segments()
		validSegments := test.ExpectedSegments

		if len(segments) != len(validSegments) {
			t.Errorf("expected length of segments (%d) and validSegments (%d) to match", len(segments), len(validSegments))
			t.Error("segments:", segments)
			return
		}

		for i, segment := range segments {
			if segment != validSegments[i] {
				t.Error("segment mismatch:: expecetd segment ", i, " to be ", validSegments[i], " but was ", segment)
			}
		}
	}
}

// Version returns the string representation of the *Version.
func Test_Version(t *testing.T) {
	for _, test := range versionTests {
		if !test.ExpectedResponse {
			continue
		}

		v, err := New(test.Version)
		if err != nil {
			t.Error("testing bug: test.Version should be a valid version for test: ", test.Version)
			t.Fail()
			return
		}

		if v.Version() != test.ExpectedVersion {
			t.Error("expected Version() to match ", test.ExpectedVersion, " but was ", v.Version())
			t.Fail()
			return
		}
	}
}

// Return a new version object where the next to the last revision
// number is one greater (e.g., 5.3.1 => 5.4).
//
// Pre-release (alpha) parts, e.g, 5.3.1.b.2 => 5.4, are ignored.
func Test_Bump(t *testing.T) {
	version, err := New("1.1.13.4-3")
	if err != nil {
		t.Error("version '1.1.13.4-3' is a valid version")
		t.Fail()
		return
	}

	result, err := version.Bump()
	if err != nil {
		t.Error("expected no error but received", err)
		t.Fail()
		return
	}

	if result == nil {
		t.Error("result should not be nil")
		t.Fail()
		return
	}

	if result.Version() != "1.1.14" {
		t.Error("expected Bump() to return a Version of '1.1.14' but got", result.Version())
		t.Fail()
		return
	}

	version, err = New("5.3.1.b.2")
	if err != nil {
		t.Error("version '5.3.1.b.2' is a valid version")
		t.Fail()
		return
	}

	result, err = version.Bump()
	if err != nil {
		t.Error("expected no error but received", err)
		t.Fail()
		return
	}

	if result == nil {
		t.Error("result should not be nil")
		t.Fail()
		return
	}

	if result.Version() != "5.4" {
		t.Error("result should be a Version of '5.4' but was", version.Version())
		t.Fail()
		return
	}
}

// IsPrerelease returns whether the Version is prerelease.
// A version is considered a prerelease if it contains a letter.
func Test_IsPrerelease(t *testing.T) {
	version, err := New("1.3-1")
	if err != nil {
		t.Error("expected '1.3-1' to be a valid version")
		t.Fail()
		return
	}

	if !version.IsPrerelease() {
		t.Error("expected version", version.Version(), "to be prerelease")
		t.Fail()
		return
	}

	version, err = New("1.13.2.1")
	if err != nil {
		t.Error("expected a valid Version from '1.13.2.1' but got error", err)
		t.Fail()
		return
	}

	if version.IsPrerelease() {
		t.Error("expected version '1.13.2.1' to not be prerelease but was")
		t.Fail()
		return
	}
}

// The release for this version (e.g. 1.2.0.a -> 1.2.0).
// Non-prerelease versions return themselves.
func Test_Release(t *testing.T) {
	version, err := New("1.2.5")
	if err != nil || version == nil {
		t.Error("new returned error:", err)
		t.Fail()
		return
	}

	v2 := version.Release()
	if v2 == nil {
		t.Error("expected v2 not to be nil")
		t.Fail()
		return
	}

	if v2 != version {
		t.Error("expected v2 to == version as version is not prerelease")
		t.Fail()
		return
	}

	version, err = New("1.2.1-2")
	if err != nil || version == nil {
		t.Error("new returned error:", err)
		t.Fail()
		return
	}

	v2 = version.Release()
	if v2.Version() != "1.2.1" {
		t.Error("expected v2.Version() to be '1.2.1' but was", v2.Version())
		t.Fail()
		return
	}
}

// A Version is only Eql() to another version if it's specified to the
// same precision. Version "1.0" is not the same as version "1".
func Test_Eql(t *testing.T) {
	version, err := New("1")
	if err != nil || version == nil {
		t.Error("new returned error:", err)
		t.Fail()
		return
	}

	version2, err := New("1.0")
	if err != nil || version2 == nil {
		t.Error("new returned error:", err)
		t.Fail()
		return
	}

	if version.Eql(version2) {
		t.Error("versions should not match:", version.version, "and", version2.version)
		t.Fail()
		return
	}

	version2, err = New("1")
	if err != nil || version2 == nil {
		t.Error("new returned error:", err)
		t.Fail()
		return
	}

	if !version.Eql(version2) {
		t.Error("versions should be Eql():", version.version, "and", version2.version)
		t.Fail()
		return
	}
}

// A recommended version for use with a ~> Requirement
func Test_ApproximateRecommendation(t *testing.T) {
	version, err := New("1.3.1-4")
	if err != nil || version == nil {
		t.Error("new returned error:", err)
		t.Fail()
		return
	}

	if version.ApproximateRecommendation() != "~> 1.3.a" {
		t.Error("expected ApproximateREcommendation() to return ~> 1.3.a but was", version.ApproximateRecommendation())
		t.Fail()
		return
	}

	version, err = New("1.3.5.7")
	if err != nil || version == nil {
		t.Error("new returned error:", err)
		t.Fail()
		return
	}

	if version.ApproximateRecommendation() != "~> 1.3" {
		t.Error("expected ApproximateREcommendation() to be '~> 1.3' but was", version.ApproximateRecommendation())
		t.Fail()
		return
	}
}

// reverseSlice sorts the slice s in reverse order.
func Test_ReverseSlice(t *testing.T) {
	slice := []string{"a", "b", "c"}

	newSlice := []string{"c", "b", "a"}

	reverseSlice(slice)

	for i, value := range slice {
		if value != newSlice[i] {
			t.Error("expected slices to be equal after reverse, but", value, "!=", newSlice[i])
			t.Fail()
			return
		}
	}
}

// deleteArrayElement deletes the given element from the given []string.
func Test_DeleteArrayElement(t *testing.T) {
	array := []string{"1", "2", "3", "4"}

	array2 := deleteArrayElement(array, 2)

	if len(array2) != len(array)-1 {
		t.Error("expected length of array2 to be 1 less than array")
		t.Fail()
		return
	}

	for _, value := range array2 {
		if value == array[2] {
			t.Error(array[2], "should not be in the result array")
			t.Fail()
			return
		}
	}
}

// strArrayEqual tests whether two []string slices are equal.
func Test_StrArrayEqual(t *testing.T) {
	array1 := []string{"one", "two", "three"}
	notequal := []string{"one", "two", "Three"}
	equal := []string{"one", "two", "three"}

	if strArraysEqual(array1, notequal) {
		t.Error("expected arrays not to be equal. Array:", array1, "notequal:", notequal)
		t.Fail()
		return
	}

	if !strArraysEqual(array1, equal) {
		t.Error("expected arrays to be equal. Array:", array1, "Equal:", equal)
		t.Fail()
		return
	}
}

// extractKind determines the underlying reflect.Kind of a string.
// Since wwe only deal with ints and strings, test just those two cases.
func Test_ExtractKind(t *testing.T) {
	s := "123"

	if extractKind(s) != reflect.Int {
		t.Error("expected extractKind to return reflect.Int but did not")
		t.Fail()
		return
	}

	s = "123abc"

	if extractKind(s) != reflect.String {
		t.Error("expected extractKind to return reflect.String but did not")
		t.Fail()
		return
	}
}
