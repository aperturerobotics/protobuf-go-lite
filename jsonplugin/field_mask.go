// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package jsonplugin

import "strings"

// FieldMask is the interface for field masks.
type FieldMask interface {
	GetPaths() []string
}

type path struct {
	parent  *path
	element string
}

func newPath(p string) path {
	i := strings.LastIndex(p, ".")
	if i == -1 {
		return path{
			element: p,
		}
	}
	parent := newPath(p[:i])
	return path{
		parent:  &parent,
		element: p[i+1:],
	}
}

func (p *path) Equals(other *path) bool {
	if p == nil || other == nil {
		return p == other
	}
	if p.element != other.element {
		return false
	}
	return p.parent.Equals(other.parent)
}

func (p *path) String() string {
	if p == nil {
		return ""
	}
	if p.parent == nil {
		return p.element
	}
	return p.parent.String() + "." + p.element
}

func (p *path) push(field string) *path {
	return &path{
		parent:  p,
		element: field,
	}
}

type pathSlice struct {
	paths []path
}

func newPathSlice(paths ...string) *pathSlice {
	m := &pathSlice{}
	m.paths = make([]path, len(paths))
	for i, p := range paths {
		m.paths[i] = newPath(p)
	}
	return m
}

func (m *pathSlice) add(path path) {
	m.paths = append(m.paths, path)
}

func (m *pathSlice) contains(search path) bool {
	for _, path := range m.paths {
		if path.Equals(&search) {
			return true
		}
	}
	return false
}

func (m *pathSlice) GetPaths() []string {
	if m == nil {
		return nil
	}
	paths := make([]string, len(m.paths))
	for i, path := range m.paths {
		paths[i] = path.String()
	}
	return paths
}
