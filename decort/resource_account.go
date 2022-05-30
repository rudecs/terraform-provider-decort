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
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceAccountCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceAccountCreate")

	if accountId, ok := d.GetOk("account_id"); ok {
		if exists, err := resourceAccountExists(d, m); exists {
			if err != nil {
				return err
			}
			d.SetId(strconv.Itoa(accountId.(int)))
			err = resourceAccountRead(d, m)
			if err != nil {
				return err
			}

			return nil
		}
		return errors.New("provided sep id does not exist")
	}

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("name", d.Get("account_name").(string))
	urlValues.Add("username", d.Get("username").(string))

	if emailaddress, ok := d.GetOk("emailaddress"); ok {
		urlValues.Add("emailaddress", emailaddress.(string))
	}
	if sendAccessEmails, ok := d.GetOk("send_access_emails"); ok {
		urlValues.Add("sendAccessEmails", strconv.FormatBool(sendAccessEmails.(bool)))
	}
	if resLimits, ok := d.GetOk("resource_limits"); ok {
		resLimit := resLimits.([]interface{})[0]
		resLimitConv := resLimit.(map[string]interface{})
		if resLimitConv["cu_m"] != nil {
			maxMemCap := int(resLimitConv["cu_m"].(float64))
			if maxMemCap == 0 {
				urlValues.Add("maxMemoryCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxMemoryCapacity", strconv.Itoa(maxMemCap))
			}
		}
		if resLimitConv["cu_d"] != nil {
			maxDiskCap := int(resLimitConv["cu_d"].(float64))
			if maxDiskCap == 0 {
				urlValues.Add("maxVDiskCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxVDiskCapacity", strconv.Itoa(maxDiskCap))
			}
		}
		if resLimitConv["cu_c"] != nil {
			maxCPUCap := int(resLimitConv["cu_c"].(float64))
			if maxCPUCap == 0 {
				urlValues.Add("maxCPUCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxCPUCapacity", strconv.Itoa(maxCPUCap))
			}

		}
		if resLimitConv["cu_i"] != nil {
			maxNumPublicIP := int(resLimitConv["cu_i"].(float64))
			if maxNumPublicIP == 0 {
				urlValues.Add("maxNumPublicIP", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxNumPublicIP", strconv.Itoa(maxNumPublicIP))
			}

		}
		if resLimitConv["cu_np"] != nil {
			maxNP := int(resLimitConv["cu_np"].(float64))
			if maxNP == 0 {
				urlValues.Add("maxNetworkPeerTransfer", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxNetworkPeerTransfer", strconv.Itoa(maxNP))
			}

		}
		if resLimitConv["gpu_units"] != nil {
			gpuUnits := int(resLimitConv["gpu_units"].(float64))
			if gpuUnits == 0 {
				urlValues.Add("gpu_units", strconv.Itoa(-1))
			} else {
				urlValues.Add("gpu_units", strconv.Itoa(gpuUnits))
			}
		}
	}

	accountId, err := controller.decortAPICall("POST", accountCreateAPI, urlValues)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(accountId)
	d.Set("account_id", accountId)

	err = resourceAccountRead(d, m)
	if err != nil {
		return err
	}

	d.SetId(id.String())

	return nil
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSepRead")

	acc, err := utilityAccountCheckPresence(d, m)
	if acc == nil {
		d.SetId("")
		return err
	}

	d.Set("dc_location", acc.DCLocation)
	d.Set("resources", flattenAccResources(acc.Resources))
	d.Set("ckey", acc.CKey)
	d.Set("meta", flattenMeta(acc.Meta))
	d.Set("acl", flattenAccAcl(acc.Acl))
	d.Set("company", acc.Company)
	d.Set("companyurl", acc.CompanyUrl)
	d.Set("created_by", acc.CreatedBy)
	d.Set("created_time", acc.CreatedTime)
	d.Set("deactivation_time", acc.DeactiovationTime)
	d.Set("deleted_by", acc.DeletedBy)
	d.Set("deleted_time", acc.DeletedTime)
	d.Set("displayname", acc.DisplayName)
	d.Set("guid", acc.GUID)
	d.Set("account_id", acc.ID)
	d.Set("account_name", acc.Name)
	d.Set("resource_limits", flattenRgResourceLimits(acc.ResourceLimits))
	d.Set("send_access_emails", acc.SendAccessEmails)
	d.Set("service_account", acc.ServiceAccount)
	d.Set("status", acc.Status)
	d.Set("updated_time", acc.UpdatedTime)
	d.Set("version", acc.Version)
	d.Set("vins", acc.Vins)
	d.Set("vinses", acc.Vinses)
	d.Set("computes", flattenAccComputes(acc.Computes))
	d.Set("machines", flattenAccMachines(acc.Machines))

	return nil
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceAccountDelete")

	account, err := utilityAccountCheckPresence(d, m)
	if account == nil {
		if err != nil {
			return err
		}
		return nil
	}

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
	urlValues.Add("permanently", strconv.FormatBool(d.Get("permanently").(bool)))

	_, err = controller.decortAPICall("POST", accountDeleteAPI, urlValues)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceAccountExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Debugf("resourceAccountExists")

	account, err := utilityAccountCheckPresence(d, m)
	if account == nil {
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func resourceAccountEdit(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceAccountEdit")
	c := m.(*ControllerCfg)

	urlValues := &url.Values{}
	if d.HasChange("enable") {
		api := accountDisableAPI
		enable := d.Get("enable").(bool)
		if enable {
			api = accountEnableAPI
		}
		urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))

		_, err := c.decortAPICall("POST", api, urlValues)
		if err != nil {
			return err
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("account_name") {
		urlValues.Add("name", d.Get("account_name").(string))
		urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
		_, err := c.decortAPICall("POST", accountUpdateAPI, urlValues)
		if err != nil {
			return err
		}

		urlValues = &url.Values{}
	}
	if d.HasChange("resource_limits") {
		resLimit := d.Get("resource_limits").([]interface{})[0]
		resLimitConv := resLimit.(map[string]interface{})

		if resLimitConv["cu_m"] != nil {
			maxMemCap := int(resLimitConv["cu_m"].(float64))
			if maxMemCap == 0 {
				urlValues.Add("maxMemoryCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxMemoryCapacity", strconv.Itoa(maxMemCap))
			}
		}
		if resLimitConv["cu_d"] != nil {
			maxDiskCap := int(resLimitConv["cu_d"].(float64))
			if maxDiskCap == 0 {
				urlValues.Add("maxVDiskCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxVDiskCapacity", strconv.Itoa(maxDiskCap))
			}
		}
		if resLimitConv["cu_c"] != nil {
			maxCPUCap := int(resLimitConv["cu_c"].(float64))
			if maxCPUCap == 0 {
				urlValues.Add("maxCPUCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxCPUCapacity", strconv.Itoa(maxCPUCap))
			}

		}
		if resLimitConv["cu_i"] != nil {
			maxNumPublicIP := int(resLimitConv["cu_i"].(float64))
			if maxNumPublicIP == 0 {
				urlValues.Add("maxNumPublicIP", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxNumPublicIP", strconv.Itoa(maxNumPublicIP))
			}

		}
		if resLimitConv["cu_np"] != nil {
			maxNP := int(resLimitConv["cu_np"].(float64))
			if maxNP == 0 {
				urlValues.Add("maxNetworkPeerTransfer", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxNetworkPeerTransfer", strconv.Itoa(maxNP))
			}

		}
		if resLimitConv["gpu_units"] != nil {
			gpuUnits := int(resLimitConv["gpu_units"].(float64))
			if gpuUnits == 0 {
				urlValues.Add("gpu_units", strconv.Itoa(-1))
			} else {
				urlValues.Add("gpu_units", strconv.Itoa(gpuUnits))
			}
		}

		urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
		_, err := c.decortAPICall("POST", accountUpdateAPI, urlValues)
		if err != nil {
			return err
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("send_access_emails") {
		urlValues.Add("sendAccessEmails", strconv.FormatBool(d.Get("send_access_emails").(bool)))
		urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
		_, err := c.decortAPICall("POST", accountUpdateAPI, urlValues)
		if err != nil {
			return err
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("restore") {
		restore := d.Get("restore").(bool)
		if restore {
			urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
			_, err := c.decortAPICall("POST", accountRestoreAPI, urlValues)
			if err != nil {
				return err
			}

			urlValues = &url.Values{}
		}
	}

	if d.HasChange("users") {
		deletedUsers := make([]interface{}, 0)
		addedUsers := make([]interface{}, 0)
		updatedUsers := make([]interface{}, 0)

		old, new := d.GetChange("users")
		oldConv := old.([]interface{})
		newConv := new.([]interface{})
		for _, el := range oldConv {
			if !isContainsUser(newConv, el) {
				deletedUsers = append(deletedUsers, el)
			}
		}
		for _, el := range newConv {
			if !isContainsUser(oldConv, el) {
				addedUsers = append(addedUsers, el)
			} else {
				if isChangedUser(oldConv, el) {
					updatedUsers = append(updatedUsers, el)
				}
			}
		}

		if len(deletedUsers) > 0 {
			for _, user := range deletedUsers {
				userConv := user.(map[string]interface{})
				urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
				urlValues.Add("userId", userConv["user_id"].(string))
				urlValues.Add("recursivedelete", strconv.FormatBool(userConv["recursive_delete"].(bool)))
				_, err := c.decortAPICall("POST", accountDeleteUserAPI, urlValues)
				if err != nil {
					return err
				}

				urlValues = &url.Values{}
			}
		}

		if len(addedUsers) > 0 {
			for _, user := range addedUsers {
				userConv := user.(map[string]interface{})
				urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
				urlValues.Add("userId", userConv["user_id"].(string))
				urlValues.Add("accesstype", strings.ToUpper(userConv["access_type"].(string)))
				_, err := c.decortAPICall("POST", accountAddUserAPI, urlValues)
				if err != nil {
					return err
				}

				urlValues = &url.Values{}
			}
		}

		if len(updatedUsers) > 0 {
			for _, user := range updatedUsers {
				userConv := user.(map[string]interface{})
				urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
				urlValues.Add("userId", userConv["user_id"].(string))
				urlValues.Add("accesstype", strings.ToUpper(userConv["access_type"].(string)))
				_, err := c.decortAPICall("POST", accountUpdateUserAPI, urlValues)
				if err != nil {
					return err
				}

				urlValues = &url.Values{}
			}
		}

	}

	return nil
}

func isContainsUser(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["user_id"].(string) == elConv["user_id"].(string) {
			return true
		}
	}
	return false
}

func isChangedUser(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["user_id"].(string) == elConv["user_id"].(string) &&
			(strings.ToUpper(elOldConv["access_type"].(string)) != strings.ToUpper(elConv["access_type"].(string)) ||
				elOldConv["recursive_delete"].(bool) != elConv["recursive_delete"].(bool)) {
			return true
		}
	}
	return false
}

func resourceAccountSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "account name",
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "username of owner the account",
		},
		"emailaddress": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "email",
		},
		"send_access_emails": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "if true send emails when a user is granted access to resources",
		},
		"users": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"user_id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"access_type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"recursive_delete": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
		"restore": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "restore a deleted account",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "whether to completely delete the account",
		},
		"enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "enable/disable account",
		},
		"resource_limits": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cu_c": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
					"cu_d": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
					"cu_i": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
					"cu_m": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
					"cu_np": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
					"gpu_units": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"account_id": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"dc_location": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"resources": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"current": {
						Type:     schema.TypeList,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cpu": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"disksize": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"extips": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"exttraffic": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"gpu": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"ram": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
					"reserved": {
						Type:     schema.TypeList,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cpu": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"disksize": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"extips": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"exttraffic": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"gpu": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"ram": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"meta": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"can_be_deleted": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"explicit": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"right": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"user_group_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"company": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"companyurl": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"deactivation_time": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"deleted_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deleted_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"displayname": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"service_account": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vins": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"computes": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"started": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"stopped": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"machines": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"halted": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"running": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"vinses": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountEdit,
		Delete: resourceAccountDelete,
		Exists: resourceAccountExists,

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

		Schema: resourceAccountSchemaMake(),
	}
}
