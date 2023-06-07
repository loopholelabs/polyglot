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

package rust

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Dependencies struct {
	Enums bool
	Maps  bool
}

func DependencyAnalysis(file *protogen.File) *Dependencies {
	dependencies := &Dependencies{
		Enums: false,
		Maps:  false,
	}

	if len(file.Enums) > 0 {
		dependencies.Enums = true
	}
	for _, message := range file.Messages {
		for _, field := range message.Fields {
			if field.Desc.Kind() == protoreflect.MessageKind {
				if field.Desc.Message().IsMapEntry() {
					dependencies.Maps = true
				}
				if field.Desc.Message().Fields().Len() > 0 {
					dependencies = traverseFields(field.Message, dependencies)
				}
			}
		}
	}
	return dependencies
}

func traverseFields(message *protogen.Message, dependencies *Dependencies) *Dependencies {
	for _, field := range message.Fields {
		if field.Desc.Kind() == protoreflect.MessageKind {
			if field.Desc.Message().IsMapEntry() {
				dependencies.Maps = true
			}
			if field.Desc.Message().Fields().Len() > 0 {
				dependencies = traverseFields(field.Message, dependencies)
			}
		}
	}
	return dependencies
}
