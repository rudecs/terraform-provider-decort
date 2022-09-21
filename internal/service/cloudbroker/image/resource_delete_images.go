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

/*
Terraform DECORT provider - manage resources provided by DECORT (Digital Energy Cloud
Orchestration Technology) with Terraform by Hashicorp.

Source code: https://github.com/rudecs/terraform-provider-decort

Please see README.md to learn where to place source code so that it
builds seamlessly.

Documentation: https://github.com/rudecs/terraform-provider-decort/wiki
*/

package image

import (
	"context"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourceCreateListImages(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := uuid.New()
	d.SetId(id.String())
	return nil
}

func resourceDeleteListImages(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceDeleteListImages: start deleting...")

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	imageIds := d.Get("image_ids").([]interface{})
	temp := ""
	l := len(imageIds)
	for i, imageId := range imageIds {
		s := strconv.Itoa(imageId.(int))
		if i != (l - 1) {
			s += ",\n"
		} else {
			s += "\n"
		}
		temp = temp + s
	}
	temp = "[" + temp + "]"

	urlValues.Add("reason", d.Get("reason").(string))
	urlValues.Add("permanently", strconv.FormatBool(d.Get("permanently").(bool)))
	urlValues.Add("imageIds", temp)

	_, err := c.DecortAPICall(ctx, "POST", imageDeleteImagesAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceReadListImages(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceDeleteImagesSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"image_ids": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "images ids for deleting",
		},
		"reason": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "reason for deleting the images",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "whether to completely delete the images",
		},
	}

}

func ResourceDeleteImages() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceCreateListImages,
		ReadContext:   resourceReadListImages,
		DeleteContext: resourceDeleteListImages,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout60s,
			Read:    &constants.Timeout30s,
			Update:  &constants.Timeout60s,
			Delete:  &constants.Timeout60s,
			Default: &constants.Timeout60s,
		},

		Schema: resourceDeleteImagesSchemaMake(),
	}
}
