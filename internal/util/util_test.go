/*
 * Copyright (c) 2008-2021, Hazelcast, Inc. All Rights Reserved.
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

package util_test

import (
	"testing"

	pubcluster "github.com/hazelcast/hazelcast-go-client/cluster"
	"github.com/hazelcast/hazelcast-go-client/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestGetAddresses(t *testing.T) {
	host := "127.0.0.1"
	portRange := pubcluster.PortRange{
		Min: 5701,
		Max: 5703,
	}
	expectedAddrs := []pubcluster.Address{
		pubcluster.NewAddress(host, 5701),
		pubcluster.NewAddress(host, 5702),
		pubcluster.NewAddress(host, 5703),
	}
	addrs := util.GetAddresses(host, portRange)
	assert.Equal(t, addrs, expectedAddrs)
}