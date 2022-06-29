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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourceCDROMImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceCDROMImageCreate: called for image %s", d.Get("name").(string))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("url", d.Get("url").(string))
	urlValues.Add("gid", strconv.Itoa(d.Get("gid").(int)))

	tstr := d.Get("drivers").([]interface{})
	temp := ""
	l := len(tstr)
	for i, str := range tstr {
		s := "\"" + str.(string) + "\""
		if i != (l - 1) {
			s += ","
		}
		temp = temp + s
	}
	temp = "[" + temp + "]"
	urlValues.Add("drivers", temp)

	if username, ok := d.GetOk("username"); ok {
		urlValues.Add("username", username.(string))
	}
	if password, ok := d.GetOk("password"); ok {
		urlValues.Add("password", password.(string))
	}
	if accountId, ok := d.GetOk("account_id"); ok {
		urlValues.Add("accountId", strconv.Itoa(accountId.(int)))
	}
	if usernameDL, ok := d.GetOk("username_dl"); ok {
		urlValues.Add("usernameDL", usernameDL.(string))
	}
	if passwordDL, ok := d.GetOk("password_dl"); ok {
		urlValues.Add("passwordDL", passwordDL.(string))
	}
	if sepId, ok := d.GetOk("sep_id"); ok {
		urlValues.Add("sepId", strconv.Itoa(sepId.(int)))
	}
	if poolName, ok := d.GetOk("pool_name"); ok {
		urlValues.Add("pool_name", poolName.(string))
	}
	if architecture, ok := d.GetOk("architecture"); ok {
		urlValues.Add("architecture", architecture.(string))
	}

	imageId, err := c.DecortAPICall("POST", imageCreateCDROMAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(imageId)
	d.Set("image_id", imageId)

	image, err := utilityImageCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(image.ImageId))
	d.Set("bootable", image.Bootable)
	//d.Set("image_id", image.ImageId)

	diagnostics := resourceImageRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceCDROMImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceCDROMImageDelete: called for %s, id: %s", d.Get("name").(string), d.Id())

	image, err := utilityImageCheckPresence(d, m)
	if image == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))

	if permanently, ok := d.GetOk("permanently"); ok {
		urlValues.Add("permanently", strconv.FormatBool(permanently.(bool)))
	}

	_, err = c.DecortAPICall("POST", imageDeleteCDROMAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceCDROMImageSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the rescue disk",
		},
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "URL where to download ISO from",
		},
		"gid": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "grid (platform) ID where this template should be create in",
		},
		"boot_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Boot type of image bios or uefi",
		},
		"image_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image type linux, windows or other",
		},
		"drivers": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "List of types of compute suitable for image. Example: [ \"KVM_X86\" ]",
		},
		"meta": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "meta",
		},
		"hot_resize": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Does this machine supports hot resize",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional username for the image",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional password for the image",
		},
		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "AccountId to make the image exclusive",
		},
		"username_dl": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "username for upload binary media",
		},
		"password_dl": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "password for upload binary media",
		},
		"sep_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "storage endpoint provider ID",
		},
		"pool_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "pool for image create",
		},
		"architecture": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "binary architecture of this image, one of X86_64 of PPC64_LE",
		},
		"image_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "image id",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Whether to completely delete the image",
		},
		"bootable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Does this image boot OS",
		},
		"unc_path": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "unc path",
		},
		"link_to": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "status",
		},
		"tech_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "tech atatus",
		},
		"version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "version",
		},
		"size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "image size",
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"computeci_id": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"provider_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"purge_attempts": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"reference_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"res_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"res_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rescuecd": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"desc": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"shared_with": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"enabled_stacks": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"history": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"timestamp": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func ResourceCDROMImage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceCDROMImageCreate,
		ReadContext:   resourceImageRead,
		UpdateContext: resourceImageEdit,
		DeleteContext: resourceCDROMImageDelete,
		Exists:        resourceImageExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout60s,
			Read:    &constants.Timeout30s,
			Update:  &constants.Timeout60s,
			Delete:  &constants.Timeout60s,
			Default: &constants.Timeout60s,
		},

		Schema: resourceCDROMImageSchemaMake(),
	}
}
