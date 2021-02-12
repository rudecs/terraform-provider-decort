/*
Copyright (c) 2019-2021 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Sergey Shubin, <sergey.shubin@digitalenergy.online>, <svs1370@gmail.com>

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

package decort

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func makeSshKeysConfig(arg_list []interface{}) (sshkeys []SshKeyConfig, count int) {
	count = len(arg_list)
	if count < 1 {
		return nil, 0
	}

	sshkeys = make([]SshKeyConfig, count)
	var subres_data map[string]interface{}
	for index, value := range arg_list {
		subres_data = value.(map[string]interface{})
		sshkeys[index].User = subres_data["user"].(string)
		sshkeys[index].SshKey = subres_data["public_key"].(string)
		sshkeys[index].UserShell = subres_data["shell"].(string)
	}

	return sshkeys, count
}

func makeSshKeysArgString(sshkeys []SshKeyConfig) string {
	// Prepare a string with username and public ssh key value in a format recognized by cloud-init utility.
	// It is designed to be passed as "userdata" argument of virtual machine create API call.
	// The following format is expected:
	// '{"users": [{"ssh-authorized-keys": ["SSH_PUBCIC_KEY_VALUE"], "shell": "SHELL_VALUE", "name": "USERNAME_VALUE"}, {...}, ]}'

	/*
		`%s\n
		  - name: %s\n
			ssh-authorized-keys:
			- %s\n
			shell: /bin/bash`
	*/
	if len(sshkeys) < 1 {
		return ""
	}

	out := `{"users": [`
	const UserdataTemplate = `%s{"ssh-authorized-keys": ["%s"], "shell": "%s", "name": "%s"}, `
	const out_suffix = `]}`
	for _, elem := range sshkeys {
		out = fmt.Sprintf(UserdataTemplate, out, elem.SshKey, elem.UserShell, elem.User)
	}
	out = fmt.Sprintf("%s %s", out, out_suffix)
	return out
}

func sshSubresourceSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"user": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the guest OS user of a new compute, for which the following SSH key will be authorized.",
		},

		"public_key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Public SSH key to authorize to the specified guest OS user on the compute being created.",
		},

		"shell": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "/bin/bash",
			Description: "Guest user shell. This parameter is optional, default is /bin/bash.",
		},
	}

	return rets
}
