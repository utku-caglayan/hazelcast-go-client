/*
 * Copyright (c) 2008-2022, Hazelcast, Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"errors"
	"strconv"

	"github.com/cucumber/godog"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/serialization"
)

// ctxKey is the key used to store the available godogs in the context.Context.
type (
	ctxKey    struct{}
	Resources struct {
		client *hazelcast.Client
	}
)

func WithResources(ctx context.Context, testContext Resources) context.Context {
	return context.WithValue(ctx, ctxKey{}, testContext)
}

func ResourcesFromContext(ctx context.Context) Resources {
	return ctx.Value(ctxKey{}).(Resources)
}

func TableToValues(t *godog.Table) ([][]interface{}, error) {
	rows := t.Rows
	if len(rows) == 0 {
		return nil, errors.New("empty table")
	}
	typeIdentifiers := rows[0].Cells
	// 2 for maps, 1-> key 1-> value
	serializers := make([]serializer, len(typeIdentifiers))
	for i, ti := range typeIdentifiers {
		s, err := typeNameToSerializer(ti.Value)
		if err != nil {
			return nil, err
		}
		serializers[i] = s
	}
	values := rows[1:]
	sValues := make([][]interface{}, len(values))
	var err error
	for i, r := range values {
		sValues[i] = make([]interface{}, len(typeIdentifiers))
		for j, c := range r.Cells {
			sValues[i][j], err = serializers[j](c.Value)
			if err != nil {
				return nil, err
			}
		}
	}
	return sValues, nil
}

type serializer func(string) (interface{}, error)

func typeNameToSerializer(name string) (serializer, error) {
	var serializer func(s string) (interface{}, error)
	switch name {
	case "string":
		serializer = func(s string) (interface{}, error) {
			return s, nil
		}
	case "int":
		serializer = func(s string) (interface{}, error) {
			return strconv.Atoi(s)
		}
	case "int32":
		serializer = func(s string) (interface{}, error) {
			return strconv.ParseInt(s, 10, 32)
		}
	case "int64", "long":
		serializer = func(s string) (interface{}, error) {
			return strconv.ParseInt(s, 10, 64)
		}
	case "float64":
		serializer = func(s string) (interface{}, error) {
			return strconv.ParseFloat(s, 64)
		}
	case "float32":
		serializer = func(s string) (interface{}, error) {
			return strconv.ParseFloat(s, 32)
		}
	case "json":
		serializer = func(s string) (interface{}, error) {
			return serialization.JSON(s), nil
		}
	default:
		return nil, errors.New("unknown type")
	}
	return serializer, nil
}
