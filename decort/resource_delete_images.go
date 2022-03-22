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

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceCreateListImages(d *schema.ResourceData, m interface{}) error {
	id := uuid.New()
	d.SetId(id.String())
	return nil
}

func resourceDeleteListImages(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceDeleteListImages: start deleting...")

	c := m.(*ControllerCfg)
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

	_, err := c.decortAPICall("POST", imageDeleteImagesAPI, urlValues)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceReadListImages(d *schema.ResourceData, m interface{}) error {
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

func resourceDeleteImages() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceCreateListImages,
		Read:   resourceReadListImages,
		Delete: resourceDeleteListImages,

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

		Schema: resourceDeleteImagesSchemaMake(),
	}
}
