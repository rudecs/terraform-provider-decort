/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>

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

package location

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

var DefaultGridID int

func UtilityLocationGetDefaultGridID(ctx context.Context, m interface{}) (int, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}

	log.Debug("utilityLocationGetDefaultGridID: retrieving locations list")
	apiResp, err := c.DecortAPICall(ctx, "POST", LocationsListAPI, urlValues)
	if err != nil {
		return 0, err
	}

	locList := LocationsListResp{}
	err = json.Unmarshal([]byte(apiResp), &locList)
	if err != nil {
		return 0, err
	}

	if len(locList) == 0 {
		DefaultGridID = 0
		return 0, fmt.Errorf("utilityLocationGetDefaultGridID: retrieved 0 length locations list")
	}

	DefaultGridID = locList[0].GridID
	log.Debugf("utilityLocationGetDefaultGridID: default location GridID %d, name %s", DefaultGridID, locList[0].Name)

	return DefaultGridID, nil
}
