/*
	Copyright 2023 Loophole Labs

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

package golang

import (
	"flag"
	"fmt"
	"strings"

	"github.com/loopholelabs/polyglot/generator/proto/templates"
	"github.com/loopholelabs/polyglot/utils"
	"github.com/loopholelabs/polyglot/version"

	polygen "github.com/loopholelabs/polyglot/generator/golang"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"

	"text/template"
)

type Generator struct {
	options             *protogen.Options
	templ               *template.Template
	links               *map[protoreflect.FullName]*PackageInfo
	currentFile         *protogen.File
	usedNames           map[string]bool
	useNullablePointers *bool
}

func New() *Generator {
	var g *Generator

	var flags flag.FlagSet
	useNullablePointers := flags.Bool("nullable_pointers", false, "Use nullable pointers for optional fields")

	g = &Generator{
		options: &protogen.Options{
			ParamFunc:         flags.Set,
			ImportRewriteFunc: func(path protogen.GoImportPath) protogen.GoImportPath { return path },
		},
		templ:               nil,
		usedNames:           nil,
		useNullablePointers: useNullablePointers,
	}
	templ := template.Must(template.New("main").Funcs(template.FuncMap{
		"CamelCase":          g.camelCase,
		"CamelCaseFull":      utils.CamelCaseFullName,
		"MakeIterable":       utils.MakeIterable,
		"Counter":            utils.Counter,
		"FirstLowerCase":     utils.FirstLowerCase,
		"FirstLowerCaseName": utils.FirstLowerCaseName,
		"FindValue": func(field protoreflect.FieldDescriptor) string {
			return FindValue(g, field)
		},
		"GetKind":       polygen.GetKind,
		"GetLUTEncoder": polygen.GetLUTEncoder,
		"GetLUTDecoder": polygen.GetLUTDecoder,
		"GetEncodingFields": func(fields protoreflect.FieldDescriptors) EncodingFields {
			return GetEncodingFields(g, fields)
		},
		"GetDecodingFields": GetDecodingFields,
		"GetKindLUT":        polygen.GetKindLUT,
		"Params":            utils.Params,
		"Debug": func(input protoreflect.FullName) string {
			println(string(input))
			return string(input)
		},
		"ProtoName": g.protoName,
		"UseNullablePointers": func() bool {
			return *useNullablePointers
		},
	}).ParseFS(templates.FS, "*"))
	g.templ = templ
	return g
}

func (g *Generator) camelCase(n protoreflect.Name) string {
	name := utils.CamelCaseName(n)
	if g.usedNames[string(name)] {
		name = name + "_"
	}
	return name
}

func (g *Generator) protoName(n protoreflect.FullName) string {
	link := (*g.links)[n]
	if link == nil {
		return utils.GoCamelCase(string(n))
	}
	name := utils.GoCamelCase(strings.TrimPrefix(string(n), string(link.packagePath)+"."))
	if g.usedNames[name] {
		name = name + "_"
	}
	if link.importPath == string(g.currentFile.GoImportPath) {
		return name
	}
	return fmt.Sprintf("%s.%s", link.goPackageName, name)
}

func (*Generator) UnmarshalRequest(buf []byte) (*pluginpb.CodeGeneratorRequest, error) {
	req := new(pluginpb.CodeGeneratorRequest)
	return req, proto.Unmarshal(buf, req)
}

func (*Generator) MarshalResponse(res *pluginpb.CodeGeneratorResponse) ([]byte, error) {
	return proto.Marshal(res)
}

func (g *Generator) Generate(req *pluginpb.CodeGeneratorRequest) (res *pluginpb.CodeGeneratorResponse, err error) {
	plugin, err := g.options.New(req)
	if err != nil {
		return nil, err
	}

	g.usedNames = map[string]bool{
		"Reset":               true,
		"String":              true,
		"ProtoMessage":        true,
		"Marshal":             true,
		"Unmarshal":           true,
		"ExtensionRangeArray": true,
		"ExtensionMap":        true,
		"Descriptor":          true,
	}

	links := make(map[protoreflect.FullName]*PackageInfo)
	g.links = &links
	for _, f := range plugin.Files {
		for i := 0; i < f.Desc.Enums().Len(); i++ {
			enum := f.Desc.Enums().Get(i)
			(*g.links)[enum.FullName()] = &PackageInfo{
				packagePath:   string(f.Desc.Package()),
				goPackageName: string(f.GoPackageName),
				importPath:    string(f.GoImportPath),
			}
		}
		for i := 0; i < f.Desc.Messages().Len(); i++ {
			message := f.Desc.Messages().Get(i)
			(*g.links)[message.FullName()] = &PackageInfo{
				packagePath:   string(f.Desc.Package()),
				goPackageName: string(f.GoPackageName),
				importPath:    string(f.GoImportPath),
			}
			TreeShake(g.links, f, message)
		}
	}

	for _, f := range plugin.Files {
		g.currentFile = f
		deps := make(map[string]string)
		checked := make(map[string]bool)
		for i := 0; i < f.Desc.Messages().Len(); i++ {
			CheckDependencies(&deps, &checked, f.Desc.Messages().Get(i), g)
		}
		if !f.Generate {
			continue
		}
		genFile := plugin.NewGeneratedFile(polygen.FileName(f.GeneratedFilenamePrefix), f.GoImportPath)

		packageName := string(f.Desc.Package().Name())
		if packageName == "" {
			packageName = string(f.GoPackageName)
		}

		err = g.templ.ExecuteTemplate(genFile, "base.templ", map[string]interface{}{
			"pluginVersion": version.Version(),
			"sourcePath":    f.Desc.Path(),
			"package":       packageName,
			"enums":         f.Desc.Enums(),
			"messages":      f.Desc.Messages(),
			"imports":       deps,
		})
		if err != nil {
			return nil, err
		}
	}

	return plugin.Response(), nil
}

func CheckDependencies(deps *map[string]string, checked *map[string]bool, message protoreflect.MessageDescriptor, g *Generator) {
	if (*checked)[string(message.FullName())] {
		return
	}
	(*checked)[string(message.FullName())] = true
	for j := 0; j < message.Fields().Len(); j++ {
		field := message.Fields().Get(j)
		if field.Kind() == protoreflect.MessageKind {
			link := (*g.links)[field.Message().FullName()]
			if link == nil {
				continue
			}
			if link.importPath != string(g.currentFile.GoImportPath) {
				(*deps)[string(link.goPackageName)] = string(link.importPath)
			}
			CheckDependencies(deps, checked, field.Message(), g)
		}
	}
}

func TreeShake(links *map[protoreflect.FullName]*PackageInfo, f *protogen.File, message protoreflect.MessageDescriptor) {
	for i := 0; i < message.Enums().Len(); i++ {
		enum := message.Enums().Get(i)
		(*links)[enum.FullName()] = &PackageInfo{
			packagePath:   string(f.Desc.Package()),
			goPackageName: string(f.GoPackageName),
			importPath:    string(f.GoImportPath),
		}
	}
	for i := 0; i < message.Messages().Len(); i++ {
		message := message.Messages().Get(i)
		(*links)[message.FullName()] = &PackageInfo{
			packagePath:   string(f.Desc.Package()),
			goPackageName: string(f.GoPackageName),
			importPath:    string(f.GoImportPath),
		}
		if (*links)[message.FullName()] == nil {
			TreeShake(links, f, message)
		}
	}
}

type PackageInfo struct {
	packagePath   string
	goPackageName string
	importPath    string
}
