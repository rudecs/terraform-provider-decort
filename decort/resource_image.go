/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>

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
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	log "github.com/sirupsen/logrus"
)

func resourceImageCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceImageCreate: called for image %s", d.Get("name").(string))

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("url", d.Get("url").(string))
	urlValues.Add("gid", strconv.Itoa(d.Get("gid").(int)))
	urlValues.Add("boottype", d.Get("boot_type").(string))
	urlValues.Add("imagetype", d.Get("image_type").(string))

	tstr := strings.Join(d.Get("drivers").([]string), ",")
	tstr = "[" + tstr + "]"
	urlValues.Add("drivers", tstr)

	if hotresize, ok := d.GetOk("hot_resize"); ok {
		urlValues.Add("hotresize", strconv.FormatBool(hotresize.(bool)))
	}
	if username, ok := d.GetOk("username"); ok {
		urlValues.Add("username", username.(string))
	}
	if password, ok := d.GetOk("password"); ok {
		urlValues.Add("password", password.(string))
	}
	if accountId, ok := d.GetOk("account_id"); ok {
		urlValues.Add("accountId", accountId.(string))
	}
	if accountId, ok := d.GetOk("account_id"); ok {
		urlValues.Add("accountId", strconv.Itoa(accountId.(int)))
	}
	if usernameDL, ok := d.GetOk("username_DL"); ok {
		urlValues.Add("usernameDL", usernameDL.(string))
	}
	if passwordDL, ok := d.GetOk("password_DL"); ok {
		urlValues.Add("passwordDL", passwordDL.(string))
	}
	if sepId, ok := d.GetOk("sep_id"); ok {
		urlValues.Add("sepId", strconv.Itoa(sepId.(int)))
	}
	if poolName, ok := d.GetOk("pool_name"); ok {
		urlValues.Add("poolName", poolName.(string))
	}
	if architecture, ok := d.GetOk("architecture"); ok {
		urlValues.Add("architecture", architecture.(string))
	}

	imageId, err := controller.decortAPICall("POST", imageCreateAPI, urlValues)
	if err != nil {
		return err
	}

	d.SetId(imageId)
	d.Set("image_id", imageId)

	image, err := utilityImageCheckPresence(d, m)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(image.ImageId))
	d.Set("image_id", image.ImageId)

	return nil
}

func resourceImageRead(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceImageRead: called for %s id: %s", d.Get("name").(string), d.Id())

	image, err := utilityImageCheckPresence(d, m)
	if image == nil {
		d.SetId("")
		return err
	}

	d.Set("image_id", image.ImageId)
	d.Set("name", image.Name)
	d.Set("url", image.Url)
	d.Set("gid", image.Gid)
	d.Set("boot_type", image.Boottype)
	d.Set("image_type", image.Imagetype)
	d.Set("drivers", image.Drivers)
	d.Set("hot_resize", image.Hotresize)
	d.Set("username", image.Username)
	d.Set("password", image.Password)
	d.Set("account_id", image.AccountId)
	d.Set("username_DL", image.UsernameDL)
	d.Set("password_DL", image.PasswordDL)
	d.Set("sep_id", image.SepId)
	d.Set("pool_name", image.PoolName)
	d.Set("architecture", image.Architecture)

	return nil
}

func resourceImageDelete(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceImageDelete: called for %s, id: %s", d.Get("name").(string), d.Id())

	image, err := utilityImageCheckPresence(d, m)
	if image == nil {
		if err != nil {
			return err
		}
		return nil
	}

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	if permanently, ok := d.GetOk("permanently"); ok {
		urlValues.Add("permanently", strconv.FormatBool(permanently.(bool)))
	}

	_, err = controller.decortAPICall("POST", imageDeleteAPI, urlValues)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceImageExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Debugf("resourceImageExists: called for %s, id: %s", d.Get("name").(string), d.Id())

	image, err := utilityImageCheckPresence(d, m)
	if image == nil {
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func resourceImageEditName(d *schema.ResourceDiff, m interface{}) error {
	log.Debugf("resourceImageEditName: called for %s, id: %s", d.Get("name").(string), d.Id())
	c := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	urlValues.Add("name", d.Get("name").(string))
	_, err := c.decortAPICall("POST", imageEditNameAPI, urlValues)
	if err != nil {
		return err
	}
	return nil
}

func resourceImageLink(d *schema.ResourceDiff, m interface{}) error {
	log.Debugf("resourceImageLink: called for %s, id: %s", d.Get("name").(string), d.Id())
	c := m.(*ControllerCfg)
	urlValues := &url.Values{}
	link := d.Get("link").(map[string]interface{})
	urlValues.Add("imageId", strconv.Itoa(link["image_id"].(int)))
	urlValues.Add("targetId", strconv.Itoa(link["target_id"].(int)))
	_, err := c.decortAPICall("POST", imageLinkAPI, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func resourceImageSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the rescue disk",
		},
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "URL where to download media from",
		},
		"gid": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "grid (platform) ID where this template should be create in",
		},
		"boot_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"bios", "uefi"}, false),
			Description:  "Boot type of image bios or uefi",
		},
		"image_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Image type linux, windows or other",
		},
		"drivers": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			Description: "List of types of compute suitable for image. Example: [ \"KVM_X86\" ]",
		},
		"hot_resize": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Does this machine supports hot resize",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional username for the image",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional password for the image",
		},
		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "AccountId to make the image exclusive",
		},
		"username_DL": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "username for upload binary media",
		},
		"password_DL": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "password for upload binary media",
		},
		"sep_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "storage endpoint provider ID",
		},
		"pool_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "pool for image create",
		},
		"architecture": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"X86_64", "PPC64_LE"}, false),
			Description:  "binary architecture of this image, one of X86_64 of PPC64_LE",
		},
		"image_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "image id",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to completely delete the image",
		},
		"virtual": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Create virtual image",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "name of the virtual image to create",
					},
					"v_image_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "",
					},
				},
			},
		},
		"link": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Link virtual image to another image in the platform",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"image_id": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "ID of the virtual image",
					},
					"target_id": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "ID of real image to link this virtual image to",
					},
				},
			},
		},
	}
}

func resourceImage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceImageCreate,
		Read:   resourceImageRead,
		Delete: resourceImageDelete,
		Exists: resourceImageExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &Timeout60s,
			Read:    &Timeout30s,
			Update:  &Timeout60s,
			Delete:  &Timeout60s,
			Default: &Timeout60s,
		},
		CustomizeDiff: customdiff.All(
			customdiff.IfValueChange("name", func(old, new, meta interface{}) bool {
				return !(old.(string) == new.(string))
			}, resourceImageEditName),
			customdiff.IfValueChange("link", func(old, new, meta interface{}) bool {
				o := old.(map[string]interface{})
				n := new.(map[string]interface{})
				if o["image_id"].(int) != n["image_id"].(int) && o["target_id"].(int) != n["target_id"].(int) {
					return true
				}
				return false
			}, resourceImageLink),
		),

		Schema: resourceImageSchemaMake(),
	}
}
