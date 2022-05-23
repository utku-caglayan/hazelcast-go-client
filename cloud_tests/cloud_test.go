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
	"fmt"
	"reflect"
	"testing"

	"github.com/cucumber/godog"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/internal/it"
	"github.com/hazelcast/hazelcast-go-client/sql"
)

func thereIsRandomMapWithEntries(c context.Context, name string, t *godog.Table) (context.Context, error) {
	var (
		randomMap *hazelcast.Map
		err       error
	)
	if c, randomMap, err = thereIsRandomMap(c, name); err != nil {
		return c, err
	}
	entries, err := TableToValues(t)
	if err != nil {
		return c, fmt.Errorf("cannot serialize table values: %w", err)
	}
	for _, e := range entries {
		if err = randomMap.Set(c, e[0], e[1]); err != nil {
			return nil, fmt.Errorf("cannot set entry(%s, %s) to map: %w", e[0], e[1], err)
		}
	}
	return c, nil
}

func thereIsRandomMap(c context.Context, name string) (context.Context, *hazelcast.Map, error) {
	r := ResourcesFromContext(c)
	uniqueName := it.NewUniqueObjectName("map")
	randomMap, err := r.client.GetMap(c, uniqueName)
	if err != nil {
		return c, randomMap, fmt.Errorf("cannot create map %s: %w", uniqueName, err)
	}
	// save randomMap to context as "name"
	c = context.WithValue(c, name, randomMap)
	return c, randomMap, nil
}

func iCreateMappingForMap(c context.Context, mapName string, str *godog.DocString) (context.Context, error) {
	r := ResourcesFromContext(c)
	m := c.Value(mapName).(*hazelcast.Map)
	result, err := r.client.SQL().Execute(c, fmt.Sprintf(str.Content, m.Name()))
	if err != nil {
		return c, fmt.Errorf("can not create mapping: %w", err)
	}
	defer result.Close()
	return c, nil
}

func iExecuteSQLStatementForMapSaveResult(c context.Context, mapName, varName string, s *godog.DocString) (context.Context, error) {
	r := ResourcesFromContext(c)
	m := c.Value(mapName).(*hazelcast.Map)
	result, err := r.client.SQL().Execute(c, fmt.Sprintf(s.Content, m.Name()))
	if err != nil {
		return c, fmt.Errorf("can not execute SQL statement: %w", err)
	}
	c = context.WithValue(c, varName, result)
	return c, nil
}

func assertSQLResult(ctx context.Context, resultVar string, t *godog.Table) error {
	values, err := TableToValues(t)
	if err != nil {
		return fmt.Errorf("can not serialize table: %w", err)
	}
	r, ok := ctx.Value(resultVar).(sql.Result)
	if !ok {
		return errors.New("unexpected result type")
	}
	defer r.Close()
	if !r.IsRowSet() {
		return errors.New("unexpected SQL result type")
	}
	rm, err := r.RowMetadata()
	if err != nil {
		return errors.New("can not access row metadata")
	}
	if rm.ColumnCount() != len(values[0]) {
		return errors.New("unexpected result column count")
	}
	it, err := r.Iterator()
	if err != nil {
		return errors.New("can not instantiate result iterator")
	}
	var rows [][]interface{}
	for it.HasNext() {
		r, err := it.Next()
		if err != nil {
			return err
		}
		var row []interface{}
		for i := 0; i < rm.ColumnCount(); i++ {
			v, err := r.Get(i)
			if err != nil {
				return fmt.Errorf("can not access row value: %w", err)
			}
			row = append(row, v)
		}
		rows = append(rows, row)
	}
	for i, row := range rows {
		for j := range row {
			if !reflect.DeepEqual(row[j], values[i][j]) {
				return fmt.Errorf("result value mismatch, actual: %#v, expected: %#v", row[i], values[i][j])
			}
		}
	}
	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		TestSuiteInitializer: InitializeSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeSuite(suiteContext *godog.TestSuiteContext) {
	// todo create cluster here
	suiteContext.ScenarioContext().Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c, err := hazelcast.StartNewClient(ctx)
		if err != nil {
			return ctx, fmt.Errorf("can not create hazelcast client %w", err)
		}
		return WithResources(ctx, Resources{client: c}), nil
	})
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^there are following entries in a random map "([^"]*)"$`, thereIsRandomMapWithEntries)
	sc.Step(`^I create a mapping for "([^"]*)"$`, iCreateMappingForMap)
	sc.Step(`^I execute statement for "([^"]*)" with result "([^"]*)"$`, iExecuteSQLStatementForMapSaveResult)
	sc.Step(`^"([^"]*)" should be$`, assertSQLResult)
}
