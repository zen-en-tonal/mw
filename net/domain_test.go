package net

import "testing"

func TestParseDomain(t *testing.T) {
	actual, err := ParseDomain("www.example.co.jp")
	if err != nil {
		t.Error(err)
	}

	expacted := NewDomain("jp", "example.co", "www")
	if expacted != *actual {
		t.Error("assertion failed", "expected", expacted, "actual", actual)
	}
}

func TestParseDomainOnlyRoot(t *testing.T) {
	actual, err := ParseDomain("example.com")
	if err != nil {
		t.Error(err)
	}

	expacted := NewDomainFromRoot("com", "example")
	if expacted != *actual {
		t.Error("assertion failed", "expected", expacted, "actual", actual)
	}
}
