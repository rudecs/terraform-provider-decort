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

package image

import (
	"context"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourceImageVirtualCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceImageVirtualCreate: called for image %s", d.Get("name").(string))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("targetId", strconv.Itoa(d.Get("target_id").(int)))

	imageId, err := c.DecortAPICall(ctx, "POST", imageCreateVirtualAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(imageId)
	d.Set("image_id", imageId)

	_, err = utilityImageCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := resourceImageRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceImageVirtualEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceImageEdit: called for %s, id: %s", d.Get("name").(string), d.Id())

	if d.HasChange("name") {
		err := resourceImageEditName(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("link_to") {
		err := resourceImageVirtualLink(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceImageRead(ctx, d, m)
}

func resourceImageVirtualLink(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceVirtualImageLink: called for %s, id: %s", d.Get("name").(string), d.Id())
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	urlValues.Add("targetId", strconv.Itoa(d.Get("link_to").(int)))
	_, err := c.DecortAPICall(ctx, "POST", imageLinkAPI, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func ResourceImageVirtual() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceImageVirtualCreate,
		ReadContext:   resourceImageRead,
		UpdateContext: resourceImageVirtualEdit,
		DeleteContext: resourceImageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout600s,
			Read:    &constants.Timeout300s,
			Update:  &constants.Timeout300s,
			Delete:  &constants.Timeout300s,
			Default: &constants.Timeout300s,
		},

		Schema: resourceImageVirtualSchemaMake(dataSourceImageExtendSchemaMake()),
	}
}
