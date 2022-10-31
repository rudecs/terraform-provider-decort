/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

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

/*
Terraform DECORT provider - manage resources provided by DECORT (Digital Energy Cloud
Orchestration Technology) with Terraform by Hashicorp.

Source code: https://github.com/rudecs/terraform-provider-decort

Please see README.md to learn where to place source code so that it
builds seamlessly.

Documentation: https://github.com/rudecs/terraform-provider-decort/wiki
*/

package kvmvm

import (
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func parseOsUsers(logins []OsUserRecord) []interface{} {
	var result = make([]interface{}, len(logins))

	for index, value := range logins {
		elem := make(map[string]interface{})

		elem["guid"] = value.Guid
		elem["login"] = value.Login
		elem["password"] = value.Password
		elem["public_key"] = value.PubKey
		result[index] = elem
		log.Debugf("parseOsUsers: parsed element %d - login %q", index, value.Login)
	}

	return result
}

func osUsersSubresourceSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"guid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "GUID of this guest OS user.",
		},

		"login": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Login name of this guest OS user.",
		},

		"password": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "Password of this guest OS user.",
		},

		"public_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "SSH public key of this guest OS user.",
		},
	}

	return rets
}
