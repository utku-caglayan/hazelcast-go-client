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

package sql

import "context"

type Service interface {
	// Execute executes the given SQL statement.
	Execute(ctx context.Context, stmt Statement) (Result, error)
	// ExecuteQuery is a convenient method to execute a distributed query with the given parameter
	// values. You may define parameter placeholders in the query with the "?" character.
	// For every placeholder, a value must be provided.
	ExecuteQuery(ctx context.Context, query string, params ...interface{}) (Result, error)
}
