/*
	Copyright 2022 Loophole Labs

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

package generator

import (
	"github.com/loopholelabs/polyglot-go/pkg/utils"
	"text/template"

	"github.com/loopholelabs/polyglot-go/internal/version"
	"github.com/loopholelabs/polyglot-go/templates"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type Generator struct {
	options      *protogen.Options
	templ        *template.Template
	CustomFields func() string
	CustomEncode func() string
	CustomDecode func() string
}

func New() *Generator {
	var g *Generator
	templ := template.Must(template.New("main").Funcs(template.FuncMap{
		"CamelCase":          utils.CamelCaseFullName,
		"CamelCaseName":      utils.CamelCaseName,
		"MakeIterable":       utils.MakeIterable,
		"Counter":            utils.Counter,
		"FirstLowerCase":     utils.FirstLowerCase,
		"FirstLowerCaseName": utils.FirstLowerCaseName,
		"FindValue":          FindValue,
		"GetKind":            GetKind,
		"GetLUTEncoder":      GetLUTEncoder,
		"GetLUTDecoder":      GetLUTDecoder,
		"GetEncodingFields":  GetEncodingFields,
		"GetDecodingFields":  GetDecodingFields,
		"GetKindLUT":         GetKindLUT,
		"CustomFields": func() string {
			return g.CustomFields()
		},
		"CustomEncode": func() string {
			return g.CustomEncode()
		},
		"CustomDecode": func() string {
			return g.CustomDecode()
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
	return g.templ.ExecuteTemplate(genFile, "base.templ", map[string]interface{}{
		"pluginVersion":   version.Version,
		"sourcePath":      protoFile.Desc.Path(),
		"package":         packageName,
		"requiredImports": RequiredImports,
		"enums":           protoFile.Desc.Enums(),
		"messages":        protoFile.Desc.Messages(),
		"header":          header,
	})
}
