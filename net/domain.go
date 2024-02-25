package net

import (
	"errors"
	"regexp"
)

type Domain struct {
	Top  string // com
	Name string // github
	Sub  string // www
}

func NewDomain(top string, name string, sub string) Domain {
	return Domain{
		Top:  top,
		Name: name,
		Sub:  sub,
	}
}

func NewDomainFromRoot(top string, name string) Domain {
	return Domain{
		Top:  top,
		Name: name,
	}
}

func (d Domain) String() string {
	domain := d.RootDomain()
	if d.Sub != "" {
		domain = d.Sub + "." + domain
	}
	return domain
}

func (d Domain) RootDomain() string {
	return d.Name + "." + d.Top
}

func (d Domain) TrimSub() Domain {
	return Domain{
		Top:  d.Top,
		Name: d.Name,
	}
}

func (d Domain) IsGroupOn(other Domain) bool {
	return d.Top == other.Top && d.Name == other.Name
}

func ParseDomain(s string) (*Domain, error) {
	domainRegexp := *regexp.MustCompile(`^((?P<sub>[a-z0-9\-]+)\.)?(?P<name>[a-z0-9\-\.?]+)\.(?P<root>[a-z]{2,4})$`)
	match := domainRegexp.FindStringSubmatch(s)
	sub := match[2]
	name := match[3]
	top := match[4]
	if name == "" || top == "" {
		return nil, errors.New("")
	}
	return &Domain{
		Top:  top,
		Name: name,
		Sub:  sub,
	}, nil
}

func MustParseDomain(s string) Domain {
	d, err := ParseDomain(s)
	if err != nil {
		panic(err)
	}
	return *d
}
