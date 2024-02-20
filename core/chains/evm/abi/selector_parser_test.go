// Sourced from https://github.com/ethereum/go-ethereum/blob/fe91d476ba3e29316b6dc99b6efd4a571481d888/accounts/abi/selector_parser_test.go

// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// https://docs.soliditylang.org/en/latest/grammar.html#a4.SolidityParser.parameterList
type parameter struct {
	typeName   string
	identifier string
}

func mkType(parameterList ...interface{}) []abi.ArgumentMarshaling {
	var result []abi.ArgumentMarshaling
	for i, p := range parameterList {
		name := fmt.Sprintf("name%d", i)

		if safeParameter, ok := p.(parameter); ok {
			name = safeParameter.identifier
			p = safeParameter.typeName
		}

		if typeName, ok := p.(string); ok {
			result = append(result, abi.ArgumentMarshaling{Name: name, Type: typeName, InternalType: typeName, Components: nil, Indexed: false})
		} else if components, ok := p.([]abi.ArgumentMarshaling); ok {
			result = append(result, abi.ArgumentMarshaling{Name: name, Type: "tuple", InternalType: "tuple", Components: components, Indexed: false})
		} else if components, ok := p.([][]abi.ArgumentMarshaling); ok {
			result = append(result, abi.ArgumentMarshaling{Name: name, Type: "tuple[]", InternalType: "tuple[]", Components: components[0], Indexed: false})
		} else {
			log.Fatalf("unexpected type %T", p)
		}
	}
	return result
}

// mkType := func(name string, typeOrComponents interface{}) abi.ArgumentMarshaling {
// 	if typeName, ok := typeOrComponents.(string); ok {
// 		return abi.ArgumentMarshaling{Name: name, Type: typeName, InternalType: typeName, Components: nil, Indexed: false}
// 	} else if components, ok := typeOrComponents.([]abi.ArgumentMarshaling); ok {
// 		return abi.ArgumentMarshaling{Name: name, Type: "tuple", InternalType: "tuple", Components: components, Indexed: false}
// 	} else if components, ok := typeOrComponents.([][]abi.ArgumentMarshaling); ok {
// 		return abi.ArgumentMarshaling{Name: name, Type: "tuple[]", InternalType: "tuple[]", Components: components[0], Indexed: false}
// 	}
// 	log.Fatalf("unexpected type %T", typeOrComponents)
// 	return abi.ArgumentMarshaling{}
// }

func TestParseSelector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		name  string
		args  []abi.ArgumentMarshaling
	}{
		{"noargs()", "noargs", []abi.ArgumentMarshaling{}},
		{"simple(uint256,uint256,uint256)", "simple", mkType("uint256", "uint256", "uint256")},
		{"other(uint256,address)", "other", mkType("uint256", "address")},
		{"withArray(uint256[],address[2],uint8[4][][5])", "withArray", mkType("uint256[]", "address[2]", "uint8[4][][5]")},
		{"singleNest(bytes32,uint8,(uint256,uint256),address)", "singleNest", mkType("bytes32", "uint8", mkType("uint256", "uint256"), "address")},
		{"multiNest(address,(uint256[],uint256),((address,bytes32),uint256))", "multiNest",
			mkType("address", mkType("uint256[]", "uint256"), mkType(mkType("address", "bytes32"), "uint256"))},
		{"arrayNest((uint256,uint256)[],bytes32)", "arrayNest", mkType([][]abi.ArgumentMarshaling{mkType("uint256", "uint256")}, "bytes32")},
		{"multiArrayNest((uint256,uint256)[],(uint256,uint256)[])", "multiArrayNest",
			mkType([][]abi.ArgumentMarshaling{mkType("uint256", "uint256")}, [][]abi.ArgumentMarshaling{mkType("uint256", "uint256")})},
		{"singleArrayNestAndArray((uint256,uint256)[],bytes32[])", "singleArrayNestAndArray",
			mkType([][]abi.ArgumentMarshaling{mkType("uint256", "uint256")}, "bytes32[]")},
		{"singleArrayNestWithArrayAndArray((uint256[],address[2],uint8[4][][5])[],bytes32[])", "singleArrayNestWithArrayAndArray",
			mkType([][]abi.ArgumentMarshaling{mkType("uint256[]", "address[2]", "uint8[4][][5]")}, "bytes32[]")},
	}
	for i, tt := range tests {
		selector, err := ParseSelector(tt.input)
		if err != nil {
			t.Errorf("test %d: failed to parse selector '%v': %v", i, tt.input, err)
		}
		if selector.Name != tt.name {
			t.Errorf("test %d: unexpected function name: '%s' != '%s'", i, selector.Name, tt.name)
		}

		if selector.Type != "function" {
			t.Errorf("test %d: unexpected type: '%s' != '%s'", i, selector.Type, "function")
		}
		if !reflect.DeepEqual(selector.Inputs, tt.args) {
			t.Errorf("test %d: unexpected args: '%v' != '%v'", i, selector.Inputs, tt.args)
		}
	}
}

func TestParseSelectorWithNames(t *testing.T) {
	t.Parallel()

	type testCases struct {
		description    string
		input          string
		expectedOutput abi.SelectorMarshaling
	}

	for _, tc := range []testCases{
		{
			description: "no_args",
			input:       "noargs()",
			expectedOutput: abi.SelectorMarshaling{
				Name:   "noargs",
				Type:   "function",
				Inputs: []abi.ArgumentMarshaling{},
			},
		},
		{
			description: "simple_named_args",
			input:       "simple(uint256 a, address b, byte c)",
			expectedOutput: abi.SelectorMarshaling{
				Name:   "simple",
				Type:   "function",
				Inputs: mkType(parameter{"uint256", "a"}, parameter{"address", "b"}, parameter{"byte", "c"}),
			},
		},
		// FAILING
		{
			description: "simple_named_args",
			input:       "simple    (uint256     a,     address b, byte c)",
			expectedOutput: abi.SelectorMarshaling{
				Name:   "simple",
				Type:   "function",
				Inputs: mkType(parameter{"uint256", "a"}, parameter{"address", "b"}, parameter{"byte", "c"}),
			},
		},
		{
			description: "tuple_named_args",
			input:       "addPerson((string name, uint16 age) person)",
			expectedOutput: abi.SelectorMarshaling{
				Name: "addPerson",
				Type: "function",
				Inputs: []abi.ArgumentMarshaling{
					{
						Name:         "person",
						Type:         "tuple",
						InternalType: "tuple",
						Components: []abi.ArgumentMarshaling{
							{
								Name: "name", Type: "string", InternalType: "string", Components: nil, Indexed: false,
							},
							{
								Name: "age", Type: "uint16", InternalType: "uint16", Components: nil, Indexed: false,
							},
						},
					},
				},
			},
		},
		// FAILING CASE
		{
			description: "explicit_tuple_named_args",
			input:       "addPerson(tuple(string name, uint16 age) person)",
			expectedOutput: abi.SelectorMarshaling{
				Name: "addPerson",
				Type: "function",
				Inputs: []abi.ArgumentMarshaling{
					{
						Name:         "person",
						Type:         "tuple",
						InternalType: "tuple",
						Components: []abi.ArgumentMarshaling{
							{
								Name: "name", Type: "string", InternalType: "string", Components: nil, Indexed: false,
							},
							{
								Name: "age", Type: "uint16", InternalType: "uint16", Components: nil, Indexed: false,
							},
						},
					},
				},
			},
		},

		// "function ",
		//   "function addPeople(tuple(string name, uint16 age)[] person)",

		// {"other(uint256 foo,    address bar )", "other", []abi.ArgumentMarshaling{mkType("foo", "uint256"), mkType("bar", "address")}},
		// {"withArray(uint256[] a, address[2] b, uint8[4][][5] c)", "withArray", []abi.ArgumentMarshaling{mkType("a", "uint256[]"), mkType("b", "address[2]"), mkType("c", "uint8[4][][5]")}},
		// {"singleNest(bytes32 d, uint8 e, (uint256,uint256) f, address g)", "singleNest", []abi.ArgumentMarshaling{mkType("d", "bytes32"), mkType("e", "uint8"), mkType("f", []abi.ArgumentMarshaling{mkType("name0", "uint256"), mkType("name1", "uint256")}), mkType("g", "address")}},
		// {"singleNest(bytes32 d, uint8 e, (uint256 first,   uint256 second ) f, address g)", "singleNest", []abi.ArgumentMarshaling{mkType("d", "bytes32"), mkType("e", "uint8"), mkType("f", []abi.ArgumentMarshaling{mkType("first", "uint256"), mkType("second", "uint256")}), mkType("g", "address")}},
	} {
		t.Run(tc.description, func(t *testing.T) {
			selector, err := ParseSelector(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, selector)
		})
	}
}

func TestParseSelectorErrors(t *testing.T) {
	type errorTestCases struct {
		description   string
		input         string
		expectedError string
	}

	for _, scenario := range []errorTestCases{
		{
			description:   "invalid name",
			input:         "123()",
			expectedError: "failed to parse selector identifier '123()': invalid token start. Expected: [a-zA-Z], received: 1",
		},
		{
			description:   "missing closing parenthesis",
			input:         "noargs(",
			expectedError: "failed to parse selector args 'noargs(': expected ')', got ''",
		},
		{
			description:   "missing opening parenthesis",
			input:         "noargs)",
			expectedError: "failed to parse selector args 'noargs)': expected '(', got )",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			_, err := ParseSelector(scenario.input)
			require.Error(t, err)
			assert.Equal(t, scenario.expectedError, err.Error())
		})
	}
}
