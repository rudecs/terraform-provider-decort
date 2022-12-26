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

package vins

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	"github.com/rudecs/terraform-provider-decort/internal/dc"
	"github.com/rudecs/terraform-provider-decort/internal/status"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVinsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	rgId, rgOk := d.GetOk("rg_id")
	accountId, accountIdOk := d.GetOk("account_id")
	if !rgOk && !accountIdOk {
		return diag.Errorf("resourceVinsCreate: no valid accountId or resource group ID specified")
	}

	if rgOk {
		urlValues.Add("name", d.Get("name").(string))
		urlValues.Add("rgId", strconv.Itoa(rgId.(int)))
		if ipcidr, ok := d.GetOk("ipcidr"); ok {
			urlValues.Add("ipcidr", ipcidr.(string))
		}

		//extnet v1
		urlValues.Add("extNetId", strconv.Itoa(d.Get("ext_net_id").(int)))
		if extIp, ok := d.GetOk("ext_ip_addr"); ok {
			urlValues.Add("extIp", extIp.(string))
		}

		//extnet v2
		if extNetResp, ok := d.GetOk("ext_net"); ok {
			extNetSl := extNetResp.([]interface{})
			extNet := extNetSl[0].(map[string]interface{})
			urlValues.Add("vinsId", d.Id())
			urlValues.Add("netId", strconv.Itoa(extNet["ext_net_id"].(int)))
			urlValues.Add("extIp", extNet["ext_net_ip"].(string))
		}

		if desc, ok := d.GetOk("desc"); ok {
			urlValues.Add("desc", desc.(string))
		}
		urlValues.Add("preReservationsNum", strconv.Itoa(d.Get("pre_reservations_num").(int)))
		id, err := c.DecortAPICall(ctx, "POST", VinsCreateInRgAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(id)
	} else if accountIdOk {
		urlValues.Add("name", d.Get("name").(string))
		urlValues.Add("accountId", strconv.Itoa(accountId.(int)))
		if gid, ok := d.GetOk("gid"); ok {
			urlValues.Add("gid", strconv.Itoa(gid.(int)))
		}
		if ipcidr, ok := d.GetOk("ipcidr"); ok {
			urlValues.Add("ipcidr", ipcidr.(string))
		}
		if desc, ok := d.GetOk("desc"); ok {
			urlValues.Add("desc", desc.(string))
		}
		urlValues.Add("preReservationsNum", strconv.Itoa(d.Get("pre_reservations_num").(int)))
		id, err := c.DecortAPICall(ctx, "POST", VinsCreateInAccountAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(id)
	}

	warnings := dc.Warnings{}
	urlValues = &url.Values{}
	if ipRes, ok := d.GetOk("ip"); ok {
		ipsSlice := ipRes.([]interface{})
		for _, ipInterfase := range ipsSlice {
			ip := ipInterfase.(map[string]interface{})
			urlValues = &url.Values{}
			urlValues.Add("vinsId", d.Id())
			urlValues.Add("type", ip["type"].(string))
			if ipAddr, ok := ip["ip_addr"]; ok {
				urlValues.Add("ipAddr", ipAddr.(string))
			}
			if macAddr, ok := ip["mac_addr"]; ok {
				urlValues.Add("mac", macAddr.(string))
			}
			if computeId, ok := ip["compute_id"]; ok {
				urlValues.Add("computeId", strconv.Itoa(computeId.(int)))
			}
			_, err := c.DecortAPICall(ctx, "POST", VinsIpReserveAPI, urlValues)
			if err != nil {
				warnings.Add(err)
			}
		}
	}

	urlValues = &url.Values{}
	if natRule, ok := d.GetOk("nat_rule"); ok {
		addedNatRules := natRule.([]interface{})
		if len(addedNatRules) > 0 {
			for _, natRuleInterface := range addedNatRules {
				urlValues = &url.Values{}
				natRule := natRuleInterface.(map[string]interface{})

				urlValues.Add("vinsId", d.Id())
				urlValues.Add("intIp", natRule["int_ip"].(string))
				urlValues.Add("intPort", strconv.Itoa(natRule["int_port"].(int)))
				urlValues.Add("extPortStart", strconv.Itoa(natRule["ext_port_start"].(int)))
				urlValues.Add("extPortEnd", strconv.Itoa(natRule["ext_port_end"].(int)))
				urlValues.Add("proto", natRule["proto"].(string))

				_, err := c.DecortAPICall(ctx, "POST", VinsNatRuleAddAPI, urlValues)
				if err != nil {
					warnings.Add(err)
				}
			}
		}
	}

	defer resourceVinsRead(ctx, d, m)
	return warnings.Get()
}

func resourceVinsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	warnings := dc.Warnings{}

	vins, err := utilityVinsCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hasChangeState := false
	if vins.Status == status.Destroyed {
		d.SetId("")
		d.Set("vins_id", 0)
		return resourceVinsCreate(ctx, d, m)
	} else if vins.Status == status.Deleted {
		hasChangeState = true

		urlValues.Add("vinsId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", VinsRestoreAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	}
	urlValues = &url.Values{}

	isEnabled := d.Get("enable").(bool)
	if vins.Status == status.Disabled && isEnabled {
		hasChangeState = true

		urlValues.Add("vinsId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", VinsEnableAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	} else if vins.Status == status.Enabled && !isEnabled {
		hasChangeState = true

		urlValues.Add("vinsId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", VinsDisableAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	}
	if hasChangeState {
		vins, err = utilityVinsCheckPresence(ctx, d, m)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	flattenVins(d, *vins)
	return warnings.Get()
}

func isContainsIp(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["ip_addr"].(string) == elConv["ip_addr"].(string) {
			return true
		}
	}
	return false
}

func isContinsNatRule(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["int_ip"].(string) == elConv["int_ip"].(string) &&
			elOldConv["int_port"].(int) == elConv["int_port"].(int) &&
			elOldConv["ext_port_start"].(int) == elConv["ext_port_start"].(int) {
			return true
		}
	}
	return false
}

func resourceVinsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	warnings := dc.Warnings{}

	enableOld, enableNew := d.GetChange("enable")
	if enableOld.(bool) && !enableNew.(bool) {
		urlValues.Add("vinsId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", VinsDisableAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	} else if !enableOld.(bool) && enableNew.(bool) {
		urlValues.Add("vinsId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", VinsEnableAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	}

	//extnet v1
	oldExtNetId, newExtNedId := d.GetChange("ext_net_id")
	if oldExtNetId.(int) != newExtNedId.(int) {
		log.Debugf("resourceVinsUpdate: changing ViNS ID %s - ext_net_id %d -> %d", d.Id(), oldExtNetId.(int), newExtNedId.(int))

		extnetParams := &url.Values{}
		extnetParams.Add("vinsId", d.Id())

		if oldExtNetId.(int) > 0 {
			// there was preexisting external net connection - disconnect ViNS
			_, err := c.DecortAPICall(ctx, "POST", VinsExtNetDisconnectAPI, extnetParams)
			if err != nil {
				warnings.Add(err)
			}
		}

		if newExtNedId.(int) > 0 {
			// new external network connection requested - connect ViNS
			extnetParams.Add("netId", fmt.Sprintf("%d", newExtNedId.(int)))
			extNetIp, ok := d.GetOk("ext_net_ip")
			if ok && extNetIp.(string) != "" {
				urlValues.Add("Ip", extNetIp.(string))
			}
			_, err := c.DecortAPICall(ctx, "POST", VinsExtNetConnectAPI, extnetParams)
			if err != nil {
				warnings.Add(err)
			}
		}
	}

	urlValues = &url.Values{}
	if d.HasChange("ip") {
		deletedIps := make([]interface{}, 0)
		addedIps := make([]interface{}, 0)

		oldIpInterface, newIpInterface := d.GetChange("ip")
		oldIpSlice := oldIpInterface.([]interface{})
		newIpSlice := newIpInterface.([]interface{})

		for _, el := range oldIpSlice {
			if !isContainsIp(newIpSlice, el) {
				deletedIps = append(deletedIps, el)
			}
		}

		for _, el := range newIpSlice {
			if !isContainsIp(oldIpSlice, el) {
				addedIps = append(addedIps, el)
			}
		}

		if len(deletedIps) > 0 {
			for _, ipInterfase := range deletedIps {
				urlValues = &url.Values{}
				ip := ipInterfase.(map[string]interface{})

				urlValues.Add("vinsId", d.Id())
				if ip["ip_addr"].(string) != "" {
					urlValues.Add("ipAddr", ip["ip_addr"].(string))
				}
				if ip["mac_addr"].(string) != "" {
					urlValues.Add("mac", ip["mac_addr"].(string))
				}
				_, err := c.DecortAPICall(ctx, "POST", VinsIpReleaseAPI, urlValues)
				if err != nil {
					warnings.Add(err)
				}
			}
		}

		if len(addedIps) > 0 {
			for _, ipInterfase := range addedIps {
				urlValues = &url.Values{}
				ip := ipInterfase.(map[string]interface{})

				urlValues.Add("vinsId", d.Id())
				urlValues.Add("type", ip["type"].(string))
				if ip["ip_addr"].(string) != "" {
					urlValues.Add("ipAddr", ip["ip_addr"].(string))
				}
				if ip["mac_addr"].(string) != "" {
					urlValues.Add("mac", ip["mac_addr"].(string))
				}
				if ip["compute_id"].(int) != 0 {
					urlValues.Add("computeId", strconv.Itoa(ip["compute_id"].(int)))
				}
				_, err := c.DecortAPICall(ctx, "POST", VinsIpReserveAPI, urlValues)
				if err != nil {
					warnings.Add(err)
				}
			}
		}
	}

	if d.HasChange("nat_rule") {
		deletedNatRules := make([]interface{}, 0)
		addedNatRules := make([]interface{}, 0)

		oldNatRulesInterface, newNatRulesInterface := d.GetChange("nat_rule")
		oldNatRulesSlice := oldNatRulesInterface.([]interface{})
		newNatRulesSlice := newNatRulesInterface.([]interface{})

		for _, el := range oldNatRulesSlice {
			if !isContinsNatRule(newNatRulesSlice, el) {
				deletedNatRules = append(deletedNatRules, el)
			}
		}

		for _, el := range newNatRulesSlice {
			if !isContinsNatRule(oldNatRulesSlice, el) {
				addedNatRules = append(addedNatRules, el)
			}
		}

		if len(deletedNatRules) > 0 {
			for _, natRuleInterface := range deletedNatRules {
				urlValues = &url.Values{}
				natRule := natRuleInterface.(map[string]interface{})

				urlValues.Add("vinsId", d.Id())
				urlValues.Add("ruleId", strconv.Itoa(natRule["rule_id"].(int)))
				log.Debug("NAT_RULE_DEL_WITH: ", urlValues.Encode())
				_, err := c.DecortAPICall(ctx, "POST", VinsNatRuleDelAPI, urlValues)
				if err != nil {
					warnings.Add(err)
				}
			}
		}

		if len(addedNatRules) > 0 {
			for _, natRuleInterface := range addedNatRules {
				urlValues = &url.Values{}
				natRule := natRuleInterface.(map[string]interface{})

				urlValues.Add("vinsId", d.Id())
				urlValues.Add("intIp", natRule["int_ip"].(string))
				urlValues.Add("intPort", strconv.Itoa(natRule["int_port"].(int)))
				urlValues.Add("extPortStart", strconv.Itoa(natRule["ext_port_start"].(int)))
				if natRule["ext_port_end"].(int) != 0 {
					urlValues.Add("extPortEnd", strconv.Itoa(natRule["ext_port_end"].(int)))
				}
				if natRule["proto"].(string) != "" {
					urlValues.Add("proto", natRule["proto"].(string))
				}

				log.Debug("NAT_RULE_ADD_WITH: ", urlValues.Encode())
				_, err := c.DecortAPICall(ctx, "POST", VinsNatRuleAddAPI, urlValues)
				if err != nil {
					warnings.Add(err)
				}
			}
		}
	}

	if oldRestart, newRestart := d.GetChange("vnfdev_restart"); oldRestart == false && newRestart == true {
		urlValues = &url.Values{}
		urlValues.Add("vinsId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", VinsVnfdevRestartAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	}

	if oldRedeploy, newRedeploy := d.GetChange("vnfdev_redeploy"); oldRedeploy == false && newRedeploy == true {
		urlValues = &url.Values{}
		urlValues.Add("vinsId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", VinsVnfdevRedeployAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	}

	defer resourceVinsRead(ctx, d, m)
	return warnings.Get()
}

func resourceVinsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	urlValues := &url.Values{}
	c := m.(*controller.ControllerCfg)

	urlValues.Add("vinsId", d.Id())
	urlValues.Add("force", strconv.FormatBool(d.Get("force").(bool)))
	urlValues.Add("permanently", strconv.FormatBool(d.Get("permanently").(bool)))
	_, err := c.DecortAPICall(ctx, "POST", VinsDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func extNetSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ext_net_id": {
			Type:     schema.TypeInt,
			Default:  0,
			Optional: true,
		},
		"ext_net_ip": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  "",
		},
	}
}

func ipSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"ip_addr": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"mac_addr": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"compute_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

func natRuleSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"int_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"int_port": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"ext_port_start": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"ext_port_end": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"proto": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
			Computed:     true,
		},
		"rule_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func resourceVinsSchemaMake() map[string]*schema.Schema {
	rets := dataSourceVinsSchemaMake()
	rets["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	rets["rg_id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	}
	rets["account_id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Computed: true,
	}
	rets["ext_net_id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  -1,
	}
	rets["ext_ip_addr"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	}
	rets["ipcidr"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	rets["pre_reservations_num"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  32,
	}
	rets["gid"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Computed: true,
	}
	rets["enable"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
	rets["permanently"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	rets["force"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	rets["ext_net"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: extNetSchemaMake(),
		},
	}
	rets["ip"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: ipSchemaMake(),
		},
	}
	rets["nat_rule"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: natRuleSchemaMake(),
		},
	}
	rets["desc"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Optional user-defined text description of this ViNS.",
	}
	rets["restore"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	rets["vnfdev_restart"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	rets["vnfdev_redeploy"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	rets["vins_id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Unique ID of the ViNS. If ViNS ID is specified, then ViNS name, rg_id and account_id are ignored.",
	}

	return rets
}

func ResourceVins() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceVinsCreate,
		ReadContext:   resourceVinsRead,
		UpdateContext: resourceVinsUpdate,
		DeleteContext: resourceVinsDelete,

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

		Schema: resourceVinsSchemaMake(),
	}
}
