/*
 * Copyright 2021 Huawei Technologies Co., Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nodes

import (
	_const "profile-manager/common/const"
	"profile-manager/nodes/fledge"
	"profile-manager/nodes/kuiper"
	"profile-manager/nodes/nodeif"
	"profile-manager/nodes/tdengine"
)

func GetProfileNode(name string) nodeif.Node {
	switch name {
	case _const.FledgeStr:
		return &fledge.FlEdge{}
	case _const.KuiperStr:
		return &kuiper.Kuiper{}
	case _const.TdEngineStr:
		return &tdengine.TdEngine{}
	}
	return nil
}
