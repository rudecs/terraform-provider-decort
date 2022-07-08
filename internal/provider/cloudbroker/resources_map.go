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

package cloudbroker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudbroker/image"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudbroker/pcidevice"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudbroker/sep"
)

func NewRersourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"decort_image":         image.ResourceImage(),
		"decort_virtual_image": image.ResourceVirtualImage(),
		"decort_cdrom_image":   image.ResourceCDROMImage(),
		"decort_delete_images": image.ResourceDeleteImages(),
		"decort_pcidevice":     pcidevice.ResourcePcidevice(),
		"decort_sep":           sep.ResourceSep(),
		"decort_sep_config":    sep.ResourceSepConfig(),
	}
}
