// Copyright 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

// +build gofuzz

package json

import (
	"bytes"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func FuzzStrictDecode(data []byte) int {
	gvk := &schema.GroupVersionKind{Version: "v1"}
	strictOpt := SerializerOptions{false, false, true}
	strictYamlOpt := SerializerOptions{true, false, true}
	strictPrettyOpt := SerializerOptions{false, true, true}
	nonstrictOpt := SerializerOptions{false, false, true}
	scheme := runtime.NewScheme()
	strictSer := NewSerializerWithOptions(DefaultMetaFactory, scheme, scheme, strictOpt)
	nonstrictSer := NewSerializerWithOptions(DefaultMetaFactory, scheme, scheme, nonstrictOpt)
	ysSer := NewSerializerWithOptions(DefaultMetaFactory, scheme, scheme, strictYamlOpt)
	psSer := NewSerializerWithOptions(DefaultMetaFactory, scheme, scheme, strictPrettyOpt)
	obj0, _, err0 := strictSer.Decode(data, gvk, nil)
	obj1, _, err1 := nonstrictSer.Decode(data, gvk, nil)
	obj2, _, err2 := ysSer.Decode(data, gvk, nil)
	obj3, _, err3 := psSer.Decode(data, gvk, nil)
	if obj0 == nil {
		if obj1 != nil {
			panic("NonStrict is stricter than Strict")
		}
		if obj2 != nil {
			panic("Yaml strict different from plain strict")
		}
		if obj3 != nil {
			panic("Pretty strict different from plain strict")
		}
		if err0 == nil || err1 == nil || err2 == nil || err3 == nil {
			panic("no error")
		}
		return 0
	}
	if err0 != nil {
		panic("got object and error")
	}
	if err2 != nil {
		panic("got object and error")
	}
	if err3 != nil {
		panic("got object and error")
	}
	var b0 bytes.Buffer
	err4 := strictSer.Encode(obj0, &b0)
	if err4 != nil {
		panic("Can't encode decoded data")
	}
	if !bytes.Equal(b0.Bytes(), data) {
		panic("Encoded data doesn't match original")
	}
	return 1
}

func FuzzNonStrictDecode(data []byte) int {
	gvk := &schema.GroupVersionKind{Version: "v1"}
	nonstrictOpt := SerializerOptions{false, false, true}
	nonstrictYamlOpt := SerializerOptions{true, false, true}
	nonstrictPrettyOpt := SerializerOptions{false, true, true}
	scheme := runtime.NewScheme()
	nonstrictSer := NewSerializerWithOptions(DefaultMetaFactory, scheme, scheme, nonstrictOpt)
	ynsSer := NewSerializerWithOptions(DefaultMetaFactory, scheme, scheme, nonstrictYamlOpt)
	pnsSer := NewSerializerWithOptions(DefaultMetaFactory, scheme, scheme, nonstrictPrettyOpt)
	obj0, _, err0 := nonstrictSer.Decode(data, gvk, nil)
	if err0 != nil {
		return 0
	}

	var b0 bytes.Buffer
	err1 := nonstrictSer.Encode(obj0, &b0)
	if err1 != nil {
		panic("Can't encode decoded data")
	}
	if !bytes.Equal(b0.Bytes(), data) {
		panic("Encoded data doesn't match original")
	}

	b0.Reset()
	err2 := ynsSer.Encode(obj0, &b0)
	if err2 != nil {
		panic("Can't encode decoded data")
	}

	b0.Reset()
	err3 := pnsSer.Encode(obj0, &b0)
	if err3 != nil {
		panic("Can't encode decoded data")
	}
	return 1
}
