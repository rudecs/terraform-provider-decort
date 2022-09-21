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

package cloudapi

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/account"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/bservice"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/disks"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/image"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/k8s"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/kvmvm"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/lb"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/pfw"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/rg"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/snapshot"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/vins"
)

func NewRersourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"decort_resgroup":          rg.ResourceResgroup(),
		"decort_kvmvm":             kvmvm.ResourceCompute(),
		"decort_disk":              disks.ResourceDisk(),
		"decort_vins":              vins.ResourceVins(),
		"decort_pfw":               pfw.ResourcePfw(),
		"decort_k8s":               k8s.ResourceK8s(),
		"decort_k8s_wg":            k8s.ResourceK8sWg(),
		"decort_snapshot":          snapshot.ResourceSnapshot(),
		"decort_account":           account.ResourceAccount(),
		"decort_bservice":          bservice.ResourceBasicService(),
		"decort_bservice_group":    bservice.ResourceBasicServiceGroup(),
		"decort_image":             image.ResourceImage(),
		"decort_image_virtual":     image.ResourceImageVirtual(),
		"decort_lb":                lb.ResourceLB(),
		"decort_lb_backend":        lb.ResourceLBBackend(),
		"decort_lb_backend_server": lb.ResourceLBBackendServer(),
		"decort_lb_frontend":       lb.ResourceLBFrontend(),
		"decort_lb_frontend_bind":  lb.ResourceLBFrontendBind(),
	}
}
