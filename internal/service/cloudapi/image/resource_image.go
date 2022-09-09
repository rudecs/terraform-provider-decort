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
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourceImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceImageCreate: called for image %s", d.Get("name").(string))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("url", d.Get("url").(string))
	urlValues.Add("gid", strconv.Itoa(d.Get("gid").(int)))
	urlValues.Add("boottype", d.Get("boot_type").(string))
	urlValues.Add("imagetype", d.Get("type").(string))

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
	/* uncomment then OK
	imageId, err := c.DecortAPICall(ctx, "POST", imageCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	*/
	//innovation
	res, err := c.DecortAPICall(ctx, "POST", imageCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	i := make([]interface{}, 0)
	err = json.Unmarshal([]byte(res), &i)
	if err != nil {
		return diag.FromErr(err)
	}
	imageId := strconv.Itoa(int(i[1].(float64)))
	// end innovation

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

func resourceImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceImageRead: called for %s id: %s", d.Get("name").(string), d.Id())

	img, err := utilityImageCheckPresence(ctx, d, m)
	if img == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("unc_path", img.UNCPath)
	d.Set("ckey", img.CKey)
	d.Set("account_id", img.AccountId)
	d.Set("acl", img.Acl)
	d.Set("architecture", img.Architecture)
	d.Set("boot_type", img.BootType)
	d.Set("bootable", img.Bootable)
	d.Set("compute_ci_id", img.ComputeCiId)
	d.Set("deleted_time", img.DeletedTime)
	d.Set("desc", img.Description)
	d.Set("drivers", img.Drivers)
	d.Set("enabled", img.Enabled)
	d.Set("gid", img.GridId)
	d.Set("guid", img.GUID)
	d.Set("history", flattenHistory(img.History))
	d.Set("hot_resize", img.HotResize)
	d.Set("image_id", img.Id)
	d.Set("last_modified", img.LastModified)
	d.Set("link_to", img.LinkTo)
	d.Set("milestones", img.Milestones)
	d.Set("image_name", img.Name)
	d.Set("password", img.Password)
	d.Set("pool_name", img.Pool)
	d.Set("provider_name", img.ProviderName)
	d.Set("purge_attempts", img.PurgeAttempts)
	d.Set("res_id", img.ResId)
	d.Set("rescuecd", img.RescueCD)
	d.Set("sep_id", img.SepId)
	d.Set("shared_with", img.SharedWith)
	d.Set("size", img.Size)
	d.Set("status", img.Status)
	d.Set("tech_status", img.TechStatus)
	d.Set("type", img.Type)
	d.Set("username", img.Username)
	d.Set("version", img.Version)

	return nil
}

func resourceImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceImageDelete: called for %s, id: %s", d.Get("name").(string), d.Id())

	image, err := utilityImageCheckPresence(ctx, d, m)
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

	_, err = c.DecortAPICall(ctx, "POST", imageDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceImageEditName(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceImageEditName: called for %s, id: %s", d.Get("name").(string), d.Id())
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	urlValues.Add("name", d.Get("name").(string))
	_, err := c.DecortAPICall(ctx, "POST", imageEditNameAPI, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func resourceImageEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceImageEdit: called for %s, id: %s", d.Get("name").(string), d.Id())

	if d.HasChange("name") {
		err := resourceImageEditName(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceImageRead(ctx, d, m)
}

func ResourceImage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceImageCreate,
		ReadContext:   resourceImageRead,
		UpdateContext: resourceImageEdit,
		DeleteContext: resourceImageDelete,

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

		Schema: resourceImageSchemaMake(dataSourceImageExtendSchemaMake()),
	}
}
