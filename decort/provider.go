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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	// "github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var decsController *ControllerCfg

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"authenticator": {
				Type:         schema.TypeString,
				Required:     true,
				StateFunc:    stateFuncToLower,
				ValidateFunc: validation.StringInSlice([]string{"oauth2", "legacy", "jwt"}, true), // ignore case while validating
				Description:  "Authentication mode to use when connecting to DECORT cloud API. Should be one of 'oauth2', 'legacy' or 'jwt'.",
			},

			"oauth2_url": {
				Type:        schema.TypeString,
				Optional:    true,
				StateFunc:   stateFuncToLower,
				DefaultFunc: schema.EnvDefaultFunc("DECORT_OAUTH2_URL", nil),
				Description: "OAuth2 application URL in 'oauth2' authentication mode.",
			},

			"controller_url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				StateFunc:   stateFuncToLower,
				Description: "URL of DECORT Cloud controller to use. API calls will be directed to this URL.",
			},

			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DECORT_USER", nil),
				Description: "User name for DECORT cloud API operations in 'legacy' authentication mode.",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DECORT_PASSWORD", nil),
				Description: "User password for DECORT cloud API operations in 'legacy' authentication mode.",
			},

			"app_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DECORT_APP_ID", nil),
				Description: "Application ID to access DECORT cloud API in 'oauth2' authentication mode.",
			},

			"app_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DECORT_APP_SECRET", nil),
				Description: "Application secret to access DECORT cloud API in 'oauth2' authentication mode.",
			},

			"jwt": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DECORT_JWT", nil),
				Description: "JWT to access DECORT cloud API in 'jwt' authentication mode.",
			},

			"allow_unverified_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, DECORT API will not verify SSL certificates. Use this with caution and in trusted environments only!",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"decort_resgroup": resourceResgroup(),
			"decort_kvmvm":    resourceCompute(),
			"decort_disk":     resourceDisk(),
			"decort_vins":     resourceVins(),
			"decort_pfw":      resourcePfw(),
			"decort_k8s":      resourceK8s(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"decort_account":  dataSourceAccount(),
			"decort_resgroup": dataSourceResgroup(),
			"decort_kvmvm":    dataSourceCompute(),
			"decort_image":    dataSourceImage(),
			"decort_disk":     dataSourceDisk(),
			"decort_vins":     dataSourceVins(),
			// "decort_pfw": dataSourcePfw(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func stateFuncToLower(argval interface{}) string {
	return strings.ToLower(argval.(string))
}

func stateFuncToUpper(argval interface{}) string {
	return strings.ToUpper(argval.(string))
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	decsController, err := ControllerConfigure(d)
	if err != nil {
		return nil, err
	}

	// initialize global default Grid ID - it will be needed to create some resource types, e.g. disks
	gridId, err := decsController.utilityLocationGetDefaultGridID()
	if err != nil {
		return nil, err
	}
	if gridId == 0 {
		return nil, fmt.Errorf("providerConfigure: invalid default Grid ID = 0")
	}

	return decsController, nil
}
