/*
	Copyright 2021 Misir Jafarov

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

			http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package proxy

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type PathType string

const (
	PrefixPath PathType = "prefix"
	RegexpPath PathType = "regexp"
)

type UrlMatcher interface {
	Matches(u *url.URL) bool
}

type Host struct {
	name           string
	paths          []*Path
	defaultBackend string
}

type Path struct {
	pattern  string
	pathType PathType
	backend  string
}

func NewHost(name string, defaultBackend string) (host *Host) {
	host = &Host{
		name:           name,
		paths:          []*Path{},
		defaultBackend: defaultBackend,
	}
	return
}

func (h *Host) AddPath(pattern string, pathType PathType, backend string) {
	h.paths = append(h.paths, &Path{pattern, pathType, backend})
}

func (h *Host) Matches(u *url.URL) bool {
	return u.Hostname() == h.name
}

func (h *Host) ReplaceDefault(u *url.URL) (*url.URL, error) {
	rawurl := combineUrls(h.defaultBackend, u.Path)
	return url.Parse(copyMeta(rawurl, u))
}

var regexCache = make(map[string]*regexp.Regexp)

func complieRegexp(pattern string) (re *regexp.Regexp, err error) {
	if cached, ok := regexCache[pattern]; ok {
		return cached, nil
	}
	re, err = regexp.Compile(pattern)
	regexCache[pattern] = re
	return
}

func (p *Path) Matches(u *url.URL) bool {
	switch p.pathType {
	case PrefixPath:
		p1 := normalizePath(u.Path)
		p2 := normalizePath(p.pattern)
		return strings.HasPrefix(p1, p2)

	case RegexpPath:
		re, err := complieRegexp(p.pattern)
		if err != nil {
			log.Error("Failed to parse pattern '%s': %s", p.pattern, err)
			return false
		}
		if re == nil {
			return false
		}
		return re.MatchString(u.Path)

	default:
		return false
	}
}

func (p *Path) Replace(u *url.URL) (*url.URL, error) {
	switch p.pathType {
	case PrefixPath:
		path := u.Path[len(p.pattern):]
		rawurl := combineUrls(p.backend, path)
		return url.Parse(copyMeta(rawurl, u))

	case RegexpPath:
		re, err := complieRegexp(p.pattern)
		if err != nil || re == nil {
			return nil, fmt.Errorf("Failed to parse pattern '%s': %s", p.pattern, err)
		}
		groups := re.FindStringSubmatch(u.Path)
		rawurl := p.backend
		for i, group := range groups {
			rawurl = strings.ReplaceAll(rawurl, fmt.Sprintf("$%d", i), group)
		}
		return url.Parse(copyMeta(rawurl, u))

	default:
		return nil, fmt.Errorf("invalid path type '%s'", p.pathType)
	}
}

func normalizePath(p string) string {
	if p == "" {
		return p
	}
	if p[0] == '/' {
		p = "/" + p
	}
	if p[len(p)-1] != '/' {
		p += "/"
	}
	return p
}

func combineUrls(elem ...string) (r string) {
	r = ""
	for i, e := range elem {
		if i > 0 {
			r += "/"
		}
		if strings.HasSuffix(e, "/") {
			e = e[:len(e)-1]
		}
		if strings.HasPrefix(e, "/") {
			e = e[1:]
		}
		r += e
	}
	return
}

func copyMeta(rawurl string, u *url.URL) string {
	if u.RawQuery != "" {
		rawurl += "?" + u.RawQuery
	}
	if u.RawFragment != "" {
		rawurl += "#" + u.RawFragment
	}
	return rawurl
}
