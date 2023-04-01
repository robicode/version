package requirement

import (
	"testing"

	"github.com/robicode/version"
)

func Test_Equals(t *testing.T) {
	req, err := New("> 1.3.5")
	if err != nil || req == nil {
		t.Error("BUG: New() returned error for valid req:", err)
		t.Fail()
		return
	}

	v, err := version.New("1.3.5")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if !equals(req.requirements[0], v) {
		t.Error("expected versions to be equal")
		t.Fail()
		return
	}

	v2, err := version.New("1.2.1")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if equals(req.requirements[0], v2) {
		t.Error("expected versions not to be equal")
		t.Fail()
		return
	}
}

func Test_NotEquals(t *testing.T) {
	req, err := New("> 1.3.5")
	if err != nil || req == nil {
		t.Error("BUG: New() returned error for valid req:", err)
		t.Fail()
		return
	}

	v, err := version.New("1.3.5")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if notEquals(req.requirements[0], v) {
		t.Error("expected versions to be equal")
		t.Fail()
		return
	}

	v2, err := version.New("1.2.1")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if !notEquals(req.requirements[0], v2) {
		t.Error("expected versions not to be equal")
		t.Fail()
		return
	}
}

func Test_GreaterThan(t *testing.T) {
	req, err := New("> 1.3.5")
	if err != nil || req == nil {
		t.Error("BUG: New() returned error for valid req:", err)
		t.Fail()
		return
	}

	v, err := version.New("1.3.5")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if greaterThan(req.requirements[0], v) {
		t.Error("expected version not to be greater than req version")
		t.Fail()
		return
	}

	v2, err := version.New("1.4")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if !greaterThan(req.requirements[0], v2) {
		t.Error("expected v2 to be greater than requirement")
		t.Fail()
		return
	}
}

func Test_LessThan(t *testing.T) {
	req, err := New("< 1.3.5")
	if err != nil || req == nil {
		t.Error("BUG: New() returned error for valid req:", err)
		t.Fail()
		return
	}

	v, err := version.New("1.3.4")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if !lessThan(req.requirements[0], v) {
		t.Error("expected versions to be less than requirement")
		t.Fail()
		return
	}

	v2, err := version.New("1.5")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if lessThan(req.requirements[0], v2) {
		t.Error("expected versions to be greater")
		t.Fail()
		return
	}
}

func Test_Gte(t *testing.T) {
	req, err := New(">= 1.3.5")
	if err != nil || req == nil {
		t.Error("BUG: New() returned error for valid req:", err)
		t.Fail()
		return
	}

	v, err := version.New("1.3.5")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if !gte(req.requirements[0], v) {
		t.Error("expected matching versions to return true")
		t.Fail()
		return
	}

	v2, err := version.New("1.2.1")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if gte(req.requirements[0], v2) {
		t.Error("expected lesser versions not to true")
		t.Fail()
		return
	}

	v3, err := version.New("1.36")
	if err != nil {
		t.Error("version.New returned error:", err)
		t.Fail()
		return
	}

	if !gte(req.requirements[0], v3) {
		t.Error("expected v3 to be greater than requirement")
		t.Fail()
		return
	}
}

func Test_Lte(t *testing.T) {
	req, err := New(">= 1.3.5")
	if err != nil || req == nil {
		t.Error("BUG: New() returned error for valid req:", err)
		t.Fail()
		return
	}

	v, err := version.New("1.3.5")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if !lte(req.requirements[0], v) {
		t.Error("expected matching versions to return true")
		t.Fail()
		return
	}

	v2, err := version.New("1.2.1")
	if err != nil {
		t.Error("version.New returned error for valid version:", err)
		t.Fail()
		return
	}

	if !lte(req.requirements[0], v2) {
		t.Error("expected lesser versions not to true")
		t.Fail()
		return
	}

	v3, err := version.New("1.36")
	if err != nil {
		t.Error("version.New returned error:", err)
		t.Fail()
		return
	}

	if lte(req.requirements[0], v3) {
		t.Error("expected v3 to be greater than requirement")
		t.Fail()
		return
	}
}

func Test_TildeGT(t *testing.T) {
	req, err := New("~> 3.5")
	if err != nil {
		t.Error("expected New not to return error but got:", err)
		t.Fail()
		return
	}

	ver, err := version.New("3.5.1")
	if err != nil {
		t.Error("expected version.New not to return error but got:", err)
		t.Fail()
		return
	}

	if !tildeGT(req.requirements[0], ver) {
		t.Error("expected version to satisfy requirement")
		t.Fail()
		return
	}

	ver, err = version.New("3.4")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if tildeGT(req.requirements[0], ver) {
		t.Error("expected version to not match requirement")
		t.Fail()
		return
	}
}

func Test_New(t *testing.T) {
	req, err := New(">= 1.3.5")
	if err != nil {
		t.Error("expected err to be nil but got:", err)
		t.Fail()
		return
	}

	if req == nil {
		t.Error("expected req not to be nil but was nil")
		t.Fail()
		return
	}
}

func Test_IsSatisfiedBy(t *testing.T) {
	req, err := New(">= 1.3.5")
	if err != nil {
		t.Error("expected err to be nil but got:", err)
		t.Fail()
		return
	}

	if req == nil {
		t.Error("expected req not to be nil but was nil")
		t.Fail()
		return
	}

	ver, err := version.New("1.0.6")
	if err != nil || ver == nil {
		t.Error("expected version to be created but got error:", err)
		t.Fail()
		return
	}

	if req.IsSatisfiedBy(ver) {
		t.Error("expected requirement not to be satisfied by version")
		t.Fail()
		return
	}

	req, err = New("> 1.2", "< 1.4", "!= 1.3.3")
	if err != nil || req == nil {
		t.Error(err)
		t.Fail()
		return
	}

	ver, err = version.New("1.3.5")
	if err != nil || ver == nil {
		t.Error(err)
	}

	if !req.IsSatisfiedBy(ver) {
		t.Error("expected version to satisfy requirement")
		t.Fail()
		return
	}
}
