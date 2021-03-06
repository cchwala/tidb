// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"reflect"

	. "github.com/pingcap/check"
)

func (s *testConfigSuite) TestCloneConf(c *C) {
	c1, err := CloneConf(&defaultConf)
	c.Assert(err, IsNil)
	c2, err := CloneConf(c1)
	c.Assert(err, IsNil)
	c.Assert(reflect.DeepEqual(c1, c2), IsTrue)

	c1.Store = "abc"
	c1.Port = 2333
	c1.Log.EnableSlowLog = !c1.Log.EnableSlowLog
	c1.RepairTableList = append(c1.RepairTableList, "abc")
	c.Assert(c1.Store, Not(Equals), c2.Store)
	c.Assert(c1.Port, Not(Equals), c2.Port)
	c.Assert(c1.Log.EnableSlowLog, Not(Equals), c2.Log.EnableSlowLog)
	c.Assert(fmt.Sprintf("%v", c1.RepairTableList), Not(Equals), fmt.Sprintf("%v", c2.RepairTableList))
}

func (s *testConfigSuite) TestMergeConfigItems(c *C) {
	oriConf, _ := CloneConf(&defaultConf)
	oldConf, _ := CloneConf(oriConf)
	newConf, _ := CloneConf(oldConf)

	// allowed
	newConf.Performance.MaxProcs = 123
	newConf.Performance.MaxMemory = 123
	newConf.Performance.CrossJoin = false
	newConf.Performance.FeedbackProbability = 123
	newConf.Performance.QueryFeedbackLimit = 123
	newConf.Performance.PseudoEstimateRatio = 123
	newConf.OOMAction = "panic"
	newConf.MemQuotaQuery = 123
	newConf.TiKVClient.StoreLimit = 123

	// rejected
	newConf.Store = "tiflash"
	newConf.Port = 2333
	newConf.AdvertiseAddress = "1.2.3.4"
	newConf.Log.SlowThreshold = 2345

	as, rs := MergeConfigItems(oldConf, newConf)
	c.Assert(len(as), Equals, 9)
	c.Assert(len(rs), Equals, 4)
	for _, a := range as {
		_, ok := dynamicConfigItems[a]
		c.Assert(ok, IsTrue)
	}
	for _, a := range rs {
		_, ok := dynamicConfigItems[a]
		c.Assert(ok, IsFalse)
	}

	c.Assert(oldConf.Performance.MaxProcs, Equals, newConf.Performance.MaxProcs)
	c.Assert(oldConf.Performance.MaxMemory, Equals, newConf.Performance.MaxMemory)
	c.Assert(oldConf.Performance.CrossJoin, Equals, newConf.Performance.CrossJoin)
	c.Assert(oldConf.Performance.FeedbackProbability, Equals, newConf.Performance.FeedbackProbability)
	c.Assert(oldConf.Performance.QueryFeedbackLimit, Equals, newConf.Performance.QueryFeedbackLimit)
	c.Assert(oldConf.Performance.PseudoEstimateRatio, Equals, newConf.Performance.PseudoEstimateRatio)
	c.Assert(oldConf.OOMAction, Equals, newConf.OOMAction)
	c.Assert(oldConf.MemQuotaQuery, Equals, newConf.MemQuotaQuery)
	c.Assert(oldConf.TiKVClient.StoreLimit, Equals, newConf.TiKVClient.StoreLimit)

	c.Assert(oldConf.Store, Equals, oriConf.Store)
	c.Assert(oldConf.Port, Equals, oriConf.Port)
	c.Assert(oldConf.AdvertiseAddress, Equals, oriConf.AdvertiseAddress)
	c.Assert(oldConf.Log.SlowThreshold, Equals, oriConf.Log.SlowThreshold)
}
