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

package typescript

import (
	"github.com/loopholelabs/polyglot/generator/typescript/templates"
	"github.com/loopholelabs/polyglot/utils"
	"github.com/loopholelabs/polyglot/version"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"

	"bytes"
	"strings"
	"text/template"
)

type Generator struct {
	options      *protogen.Options
	templ        *template.Template
	CustomFields func() string
	CustomEncode func() string
	CustomDecode func() string

	dependencies map[string]struct{}
}

func New() *Generator {
	var g *Generator
	templ := template.Must(template.New("main").Funcs(template.FuncMap{
		"CamelCase":          utils.CamelCaseFullName,
		"CamelCaseName":      utils.CamelCaseName,
		"CamelCaseFullName":  utils.CamelCaseFullName,
		"MakeIterable":       utils.MakeIterable,
		"Counter":            utils.Counter,
		"FirstLowerCase":     utils.FirstLowerCase,
		"FirstLowerCaseName": utils.FirstLowerCaseName,
		"FindValue":          findValue,
		"GetKind": func(kind protoreflect.Kind) string {
			return getKind(g.dependencies, kind)
		},
		"GetLUTEncoder": func(kind protoreflect.Kind) string {
			return getLUTEncoder(g.trackDependency, kind)
		},
		"GetLUTDecoder": func(kind protoreflect.Kind) string {
			return getLUTDecoder(g.trackDependency, kind)
		},
		"GetEncodingFields": func(fields protoreflect.FieldDescriptors) encodingFields {
			return getEncodingFields(g.trackDependency, fields)
		},
		"GetDecodingFields": func(fields protoreflect.FieldDescriptors) decodingFields {
			return getDecodingFields(g.trackDependency, fields)
		},
		"GetKindLUT": func(kind protoreflect.Kind) string {
			return getKindLUT(g.trackDependency, kind)
		},
		"LowercaseCamelCase":     utils.LowercaseCamelCase,
		"LowercaseCamelCaseName": utils.LowercaseCamelCaseName,
		"CustomFields": func() string {
			return g.CustomFields()
		},
		"CustomEncode": func() string {
			return g.CustomEncode()
		},
		"CustomDecode": func() string {
			return g.CustomDecode()
		},
		"TrimSuffix": strings.TrimSuffix,
		"TrackDependency": func(dep string) string {
			return g.trackDependency(dep)
		},
	}).ParseFS(templates.FS, "*"))
	g = &Generator{
		options: &protogen.Options{
			ParamFunc:         func(name string, value string) error { return nil },
			ImportRewriteFunc: func(path protogen.GoImportPath) protogen.GoImportPath { return path },
		},
		templ:        templ,
		CustomEncode: func() string { return "" },
		CustomDecode: func() string { return "" },
		CustomFields: func() string { return "" },

		dependencies: map[string]struct{}{},
	}
	return g
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

	for _, f := range plugin.Files {
		if !f.Generate {
			continue
		}
		genFile := plugin.NewGeneratedFile(FileName(f.GeneratedFilenamePrefix), f.GoImportPath)

		packageName := string(f.Desc.Package().Name())
		if packageName == "" {
			packageName = string(f.GoPackageName)
		}

		err = g.ExecuteTemplate(genFile, f, packageName, true)
		if err != nil {
			return nil, err
		}
	}

	return plugin.Response(), nil
}

func (g *Generator) ExecuteTemplate(
	genFile *protogen.GeneratedFile,
	protoFile *protogen.File,
	packageName string,
	header bool,
) error {
	var bodyBuf bytes.Buffer
	if err := g.templ.ExecuteTemplate(&bodyBuf, "body.templ", map[string]interface{}{
		"enums":    protoFile.Desc.Enums(),
		"messages": protoFile.Desc.Messages(),
	}); err != nil {
		return err
	}

	var headBuf bytes.Buffer
	if err := g.templ.ExecuteTemplate(&headBuf, "head.templ", map[string]interface{}{
		"pluginVersion": version.Version(),
		"sourcePath":    protoFile.Desc.Path(),
		"package":       packageName,
		"header":        header,
		"dependencies":  g.dependencies,
	}); err != nil {
		return err
	}

	_, err := genFile.Write(append(headBuf.Bytes(), bodyBuf.Bytes()...))
	return err
}

func (g *Generator) trackDependency(dep string) string {
	g.dependencies[dep] = struct{}{}

	return dep
}
