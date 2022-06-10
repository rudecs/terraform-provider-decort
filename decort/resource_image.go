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
	"errors"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		urlValues.Add("poolName", poolName.(string))
	}
	if architecture, ok := d.GetOk("architecture"); ok {
		urlValues.Add("architecture", architecture.(string))
	}

	api := ""
	if isSync := d.Get("sync").(bool); !isSync {
		api = imageCreateAPI
	} else {
		api = imageSyncCreateAPI
	}
	imageId, err := controller.decortAPICall("POST", api, urlValues)
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
	d.Set("bootable", image.Bootable)
	//d.Set("image_id", image.ImageId)

	err = resourceImageRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func resourceImageRead(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceImageRead: called for %s id: %s", d.Get("name").(string), d.Id())

	image, err := utilityImageCheckPresence(d, m)
	if image == nil {
		d.SetId("")
		return err
	}

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
	d.Set("meta", flattenMeta(image.Meta))
	d.Set("hot_resize", image.Hotresize)
	d.Set("history", flattenHistory(image.History))
	d.Set("last_modified", image.LastModified)
	d.Set("desc", image.Desc)
	d.Set("shared_with", image.SharedWith)

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
	if reason, ok := d.GetOk("reason"); ok {
		urlValues.Add("reason", reason.(string))
	} else {
		urlValues.Add("reason", "")
	}
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

func resourceImageEdit(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceImageEdit: called for %s, id: %s", d.Get("name").(string), d.Id())
	c := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	urlValues.Add("name", d.Get("name").(string))

	urlValues.Add("username", d.Get("username").(string))
	urlValues.Add("password", d.Get("password").(string))
	urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
	urlValues.Add("bootable", strconv.FormatBool(d.Get("bootable").(bool)))
	urlValues.Add("hotresize", strconv.FormatBool(d.Get("hot_resize").(bool)))

	//_, err := c.decortAPICall("POST", imageEditAPI, urlValues)
	_, err := c.decortAPICall("POST", imageEditAPI, urlValues)
	if err != nil {
		err = resourceImageRead(d, m)
		if err != nil {
			return err
		}
		return nil
	}
	err = resourceImageRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func resourceImageChangeEnabled(d *schema.ResourceDiff, m interface{}) error {
	var api string

	c := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	if d.Get("enabled").(bool) {
		api = imageEnableAPI
	} else {
		api = imageDisableAPI
	}
	resp, err := c.decortAPICall("POST", api, urlValues)
	if err != nil {
		return err
	}
	res, err := strconv.ParseBool(resp)
	if err != nil {
		return err
	}
	if !res {
		return errors.New("Cannot enable/disable")
	}
	return nil
}

func resourceImageLink(d *schema.ResourceDiff, m interface{}) error {
	log.Debugf("resourceVirtualImageLink: called for %s, id: %s", d.Get("name").(string), d.Id())
	c := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	urlValues.Add("targetId", strconv.Itoa(d.Get("link_to").(int)))
	_, err := c.decortAPICall("POST", imageLinkAPI, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func resourceImageShare(d *schema.ResourceDiff, m interface{}) error {
	log.Debugf("resourceImageShare: called for %s, id: %s", d.Get("name").(string), d.Id())
	c := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	accIds := d.Get("shared_with").([]interface{})
	temp := ""
	l := len(accIds)
	for i, accId := range accIds {
		s := strconv.Itoa(accId.(int))
		if i != (l - 1) {
			s += ",\n"
		} else {
			s += "\n"
		}
		temp = temp + s
	}
	temp = "[" + temp + "]"
	urlValues.Add("accounts", temp)
	_, err := c.decortAPICall("POST", imageShareAPI, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func resourceImageChangeComputeci(d *schema.ResourceDiff, m interface{}) error {
	c := m.(*ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	computeci := d.Get("computeci_id").(int)

	api := ""

	if computeci == 0 {
		api = imageComputeciUnsetAPI
	} else {
		urlValues.Add("computeciId", strconv.Itoa(computeci))
		api = imageComputeciSetAPI
	}

	_, err := c.decortAPICall("POST", api, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func resourceImageUpdateNodes(d *schema.ResourceDiff, m interface{}) error {
	log.Debugf("resourceImageUpdateNodes: called for %s, id: %s", d.Get("name").(string), d.Id())
	c := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	enabledStacks := d.Get("enabled_stacks").([]interface{})
	temp := ""
	l := len(enabledStacks)
	for i, stackId := range enabledStacks {
		s := stackId.(string)
		if i != (l - 1) {
			s += ","
		}
		temp = temp + s
	}
	temp = "[" + temp + "]"
	urlValues.Add("enabledStacks", temp)
	_, err := c.decortAPICall("POST", imageUpdateNodesAPI, urlValues)
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
			Description: "Name of the rescue disk",
		},
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "URL where to download media from",
		},
		"gid": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "grid (platform) ID where this template should be create in",
		},
		"boot_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Boot type of image bios or uefi",
		},
		"image_type": {
			Type:        schema.TypeString,
			Required:    true,
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
			Optional:    true,
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
		"reason": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"last_modified": {
			Type:     schema.TypeInt,
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
		"sync": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Create image from a media identified by URL (in synchronous mode)",
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

func resourceImage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceImageCreate,
		Read:   resourceImageRead,
		Update: resourceImageEdit,
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
			customdiff.IfValueChange("enabled", func(old, new, meta interface{}) bool {
				return old.(bool) != new.(bool)
			}, resourceImageChangeEnabled),
			customdiff.IfValueChange("name", func(old, new, meta interface{}) bool {
				return old.(string) != new.(string) && old.(string) != ""
			}, resourceImageEditName),
			customdiff.IfValueChange("shared_with", func(old, new, meta interface{}) bool {
				o := old.([]interface{})
				n := new.([]interface{})

				if len(o) != len(n) {
					return true
				} else if len(o) == 0 {
					return false
				}

				count := 0
				for i, v := range n {
					if v.(int) == o[i].(int) {
						count++
					}
				}
				return count == 0
			}, resourceImageShare),
			customdiff.IfValueChange("computeci_id", func(old, new, meta interface{}) bool {
				return old.(int) != new.(int)
			}, resourceImageChangeComputeci),
			customdiff.IfValueChange("enabled_stacks", func(old, new, meta interface{}) bool {
				o := old.([]interface{})
				n := new.([]interface{})

				if len(o) != len(n) {
					return true
				} else if len(o) == 0 {
					return false
				}
				count := 0
				for i, v := range n {
					if v.(string) == o[i].(string) {
						count++
					}
				}
				return count == 0
			}, resourceImageUpdateNodes),
		),

		Schema: resourceImageSchemaMake(),
	}
}

func flattenMeta(m []interface{}) []string {
	output := []string{}
	for _, item := range m {
		switch d := item.(type) {
		case string:
			output = append(output, d)
		case int:
			output = append(output, strconv.Itoa(d))
		case int64:
			output = append(output, strconv.FormatInt(d, 10))
		case float64:
			output = append(output, strconv.FormatInt(int64(d), 10))
		default:
			output = append(output, "")
		}
	}
	return output
}

func flattenHistory(history []History) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	for _, item := range history {
		t := map[string]interface{}{
			"id":        item.Id,
			"guid":      item.Guid,
			"timestamp": item.Timestamp,
		}

		temp = append(temp, t)
	}
	return temp
}
