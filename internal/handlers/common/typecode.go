// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"fmt"
	"math"
	"time"

	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/must"
)

//go:generate ../../../bin/stringer -linecomment -type typeCode

// typeCode represents BSON type codes.
// BSON type codes represent corresponding codes in BSON specification.
// They could be used to query fields with particular type values using $type operator.
// Type code `number` is added to support MongoDB surrogate alias `number` which matches double, int and long type values.
type typeCode int32

const (
	typeCodeDouble    = typeCode(1)  // double
	typeCodeString    = typeCode(2)  // string
	typeCodeObject    = typeCode(3)  // object
	typeCodeArray     = typeCode(4)  // array
	typeCodeBinData   = typeCode(5)  // binData
	typeCodeObjectID  = typeCode(7)  // objectId
	typeCodeBool      = typeCode(8)  // bool
	typeCodeDate      = typeCode(9)  // date
	typeCodeNull      = typeCode(10) // null
	typeCodeRegex     = typeCode(11) // regex
	typeCodeInt       = typeCode(16) // int
	typeCodeTimestamp = typeCode(17) // timestamp
	typeCodeLong      = typeCode(18) // long
	// Not implemented.
	typeCodeDecimal = typeCode(19)  // decimal
	typeCodeMinKey  = typeCode(-1)  // minKey
	typeCodeMaxKey  = typeCode(127) // maxKey
	// Not actual type code. `number` matches double, int and long.
	typeCodeNumber = typeCode(-128) // number
)

// newTypeCode returns typeCde and error by given code.
func newTypeCode(code int32) (typeCode, error) {
	c := typeCode(code)
	switch c {
	case typeCodeDouble, typeCodeString, typeCodeObject, typeCodeArray,
		typeCodeBinData, typeCodeObjectID, typeCodeBool, typeCodeDate,
		typeCodeNull, typeCodeRegex, typeCodeInt, typeCodeTimestamp, typeCodeLong, typeCodeNumber:
		return c, nil
	case typeCodeDecimal, typeCodeMinKey, typeCodeMaxKey:
		return 0, NewErrorMsg(ErrNotImplemented, fmt.Sprintf(`Type code %v not implemented`, code))
	default:
		return 0, NewErrorMsg(ErrBadValue, fmt.Sprintf(`Invalid numerical type code: %d`, code))
	}
}

// hasSameTypeElements returns true if types.Array elements has the same type.
// MongoDB consider int32, int64 and float64 that could be converted to int as the same type.
func hasSameTypeElements(array *types.Array) bool {
	var prev string
	for i := 0; i < array.Len(); i++ {
		var cur string

		element := must.NotFail(array.Get(i))

		switch element := element.(type) {
		case *types.Document:
			cur = "object"
		case *types.Array:
			cur = "array"
		case float64:
			if element != math.Trunc(element) || math.IsNaN(element) || math.IsInf(element, 0) {
				cur = "double"
			} else {
				// float that could be converted to int should be compared as int
				cur = "int"
			}
		case string:
			cur = "string"
		case types.Binary:
			cur = "binData"
		case types.ObjectID:
			cur = "objectId"
		case bool:
			cur = "bool"
		case time.Time:
			cur = "date"
		case types.NullType:
			cur = "null"
		case types.Regex:
			cur = "regex"
		case int32:
			cur = "int"
		case types.Timestamp:
			cur = "timestamp"
		case int64:
			cur = "int"
		default:
			return false
		}

		if prev == "" {
			prev = cur
			continue
		}

		if prev != cur {
			return false
		}
	}

	return true
}

// aliasToTypeCode matches string type aliases to the corresponding typeCode value.
var aliasToTypeCode = map[string]typeCode{}

func init() {
	for _, i := range []typeCode{
		typeCodeDouble, typeCodeString, typeCodeObject, typeCodeArray,
		typeCodeBinData, typeCodeObjectID, typeCodeBool, typeCodeDate, typeCodeNull,
		typeCodeRegex, typeCodeInt, typeCodeTimestamp, typeCodeLong, typeCodeNumber,
	} {
		aliasToTypeCode[i.String()] = i
	}
}