/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>, <svs1370@gmail.com>

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
This file is part of Terraform (by Hashicorp) provider for Digital Energy Cloud Orchestration
Technology platfom.

Visit https://github.com/rudecs/terraform-provider-decort for full source code package and updates.
*/

package decort

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func flattenImage(d *schema.ResourceData, image *Image) {
	d.Set("name", image.Name)
	d.Set("drivers", image.Drivers)
	d.Set("url", image.Url)
	d.Set("gid", image.Gid)
	d.Set("image_id", image.ImageId)
	d.Set("boot_type", image.Boottype)
	d.Set("image_type", image.Imagetype)
	d.Set("bootable", image.Bootable)
	d.Set("sep_id", image.SepId)
	d.Set("unc_path", image.UNCPath)
	d.Set("link_to", image.LinkTo)
	d.Set("status", image.Status)
	d.Set("tech_status", image.TechStatus)
	d.Set("version", image.Version)
	d.Set("size", image.Size)
	d.Set("enabled", image.Enabled)
	d.Set("computeci_id", image.ComputeciId)
	d.Set("pool_name", image.PoolName)
	d.Set("username", image.Username)
	d.Set("username_dl", image.UsernameDL)
	d.Set("password", image.Password)
	d.Set("password_dl", image.PasswordDL)
	d.Set("account_id", image.AccountId)
	d.Set("guid", image.Guid)
	d.Set("milestones", image.Milestones)
	d.Set("provider_name", image.ProviderName)
	d.Set("purge_attempts", image.PurgeAttempts)
	d.Set("reference_id", image.ReferenceId)
	d.Set("res_id", image.ResId)
	d.Set("res_name", image.ResName)
	d.Set("rescuecd", image.Rescuecd)
	d.Set("architecture", image.Architecture)
	d.Set("hot_resize", image.Hotresize)
	d.Set("history", flattenHistory(image.History))
	d.Set("last_modified", image.LastModified)
	d.Set("meta", flattenMeta(image.Meta))
	d.Set("desc", image.Desc)
	d.Set("shared_with", image.SharedWith)
	return
}

func dataSourceImageRead(d *schema.ResourceData, m interface{}) error {
	image, err := utilityImageCheckPresence(d, m)
	if err != nil {

		return err
	}
	d.SetId(strconv.Itoa(image.Guid))
	flattenImage(d, image)

	return nil
}

func dataSourceImageSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the rescue disk",
		},
		"url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "URL where to download media from",
		},
		"gid": {
			Type:        schema.TypeInt,
			Computed:    true,
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
		"shared_with": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
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
		"drivers": {
			Type:     schema.TypeList,
			Computed: true,
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
			Computed:    true,
			Description: "Does this machine supports hot resize",
		},
		"username": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Optional username for the image",
		},
		"password": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Optional password for the image",
		},
		"account_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "AccountId to make the image exclusive",
		},
		"username_dl": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "username for upload binary media",
		},
		"password_dl": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "password for upload binary media",
		},
		"sep_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "storage endpoint provider ID",
		},
		"pool_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "pool for image create",
		},
		"architecture": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "binary architecture of this image, one of X86_64 of PPC64_LE",
		},
		"image_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "image id",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Whether to completely delete the image",
		},
		"bootable": {
			Type:        schema.TypeBool,
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
			Computed: true,
		},
		"computeci_id": {
			Type:     schema.TypeInt,
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
		"last_modified": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"desc": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceImageRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceImageSchemaMake(),
	}
}
