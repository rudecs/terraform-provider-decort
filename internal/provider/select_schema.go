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

package provider

import (
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ca "github.com/rudecs/terraform-provider-decort/internal/provider/cloudapi"
	cb "github.com/rudecs/terraform-provider-decort/internal/provider/cloudbroker"
)

func selectSchema(isDatasource bool) map[string]*schema.Resource {
	adminMode, err := strconv.ParseBool(os.Getenv("DECORT_ADMIN_MODE"))
	if err != nil {
		adminMode = false
	}
	if isDatasource {
		return selectDataSourceSchema(adminMode)
	}
	return selectResourceSchema(adminMode)

}

func selectDataSourceSchema(adminMode bool) map[string]*schema.Resource {
	if adminMode {
		return cb.NewDataSourcesMap()
	}
	return ca.NewDataSourcesMap()
}

func selectResourceSchema(adminMode bool) map[string]*schema.Resource {
	if adminMode {
		return cb.NewRersourcesMap()
	}
	return ca.NewRersourcesMap()
}
