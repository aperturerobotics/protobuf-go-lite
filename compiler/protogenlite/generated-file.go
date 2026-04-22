package protogenlite

import (
	"fmt"
	"go/format"
	"sort"
	"strconv"
	"strings"
)

// GeneratedFile accumulates a generated Go source file.
type GeneratedFile struct {
	filename   string
	importPath GoImportPath
	plugin     *Plugin
	out        strings.Builder
	imports    map[GoImportPath]string
	aliases    map[string]struct{}
}

func newGeneratedFile(plugin *Plugin, filename string, importPath GoImportPath) *GeneratedFile {
	g := &GeneratedFile{
		filename:   filename,
		importPath: importPath,
		plugin:     plugin,
		imports:    make(map[GoImportPath]string),
		aliases:    make(map[string]struct{}),
	}
	if pkg := plugin.packageNameFor(importPath); pkg != "" {
		g.aliases[pkg] = struct{}{}
	}
	return g
}

// P appends a line to the generated output.
func (g *GeneratedFile) P(args ...any) {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			g.out.WriteString(v)
		case GoIdent:
			g.out.WriteString(g.QualifiedGoIdent(v))
		case int:
			g.out.WriteString(strconv.Itoa(v))
		case int32:
			g.out.WriteString(strconv.FormatInt(int64(v), 10))
		case int64:
			g.out.WriteString(strconv.FormatInt(v, 10))
		case uint64:
			g.out.WriteString(strconv.FormatUint(v, 10))
		case bool:
			g.out.WriteString(strconv.FormatBool(v))
		default:
			g.out.WriteString(fmt.Sprint(v))
		}
	}
	g.out.WriteByte('\n')
}

// QualifiedGoIdent returns the identifier qualified for the current file.
func (g *GeneratedFile) QualifiedGoIdent(id GoIdent) string {
	if id.GoImportPath == "" || id.GoImportPath == g.importPath {
		return id.GoName
	}
	alias, ok := g.imports[id.GoImportPath]
	if !ok {
		alias = g.plugin.packageNameFor(id.GoImportPath)
		if alias == "" {
			alias = goPackageName(string(id.GoImportPath))
		}
		base := alias
		for i := 2; ; i++ {
			if _, exists := g.aliases[alias]; !exists {
				break
			}
			alias = base + strconv.Itoa(i)
		}
		g.imports[id.GoImportPath] = alias
		g.aliases[alias] = struct{}{}
	}
	return alias + "." + id.GoName
}

// Content returns the generated file contents.
func (g *GeneratedFile) Content() ([]byte, error) {
	src := g.out.String()
	if len(g.imports) != 0 {
		src = insertImports(src, g.imports)
	}
	if strings.HasSuffix(g.filename, ".go") {
		formatted, err := format.Source([]byte(src))
		if err != nil {
			return nil, err
		}
		return formatted, nil
	}
	return []byte(src), nil
}

func insertImports(src string, imports map[GoImportPath]string) string {
	type importItem struct {
		path  GoImportPath
		alias string
	}
	items := make([]importItem, 0, len(imports))
	for path, alias := range imports {
		items = append(items, importItem{path: path, alias: alias})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].alias == items[j].alias {
			return items[i].path < items[j].path
		}
		return items[i].alias < items[j].alias
	})

	lines := strings.Split(src, "\n")
	var out strings.Builder
	inserted := false
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		out.WriteString(line)
		out.WriteByte('\n')
		if inserted || !strings.HasPrefix(line, "package ") {
			continue
		}
		if i+1 < len(lines) && lines[i+1] == "" {
			i++
		}
		out.WriteByte('\n')
		out.WriteString("import (\n")
		for _, item := range items {
			out.WriteString(item.alias)
			out.WriteString(" ")
			out.WriteString(strconv.Quote(string(item.path)))
			out.WriteByte('\n')
		}
		out.WriteString(")\n\n")
		inserted = true
	}
	return out.String()
}
