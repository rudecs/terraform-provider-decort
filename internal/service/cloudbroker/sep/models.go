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

package sep

///Sep Models
type SepConsumptionInd struct {
	DiskCount     int `json:"disk_count"`
	DiskUsage     int `json:"disk_usage"`
	SnapshotCount int `json:"snapshot_count"`
	SnapshotUsage int `json:"snapshot_usage"`
	Usage         int `json:"usage"`
	UsageLimit    int `json:"usage_limit"`
}

type SepConsumptionTotal struct {
	CapacityLimit int `json:"capacity_limit"`
	SepConsumptionInd
}

type SepConsumption struct {
	Total  SepConsumptionTotal          `json:"total"`
	Type   string                       `json:"type"`
	ByPool map[string]SepConsumptionInd `json:"byPool"`
}

type SepDiskList []int

type Sep struct {
	Ckey       string        `json:"_ckey"`
	Meta       []interface{} `json:"_meta"`
	ConsumedBy []int         `json:"consumedBy"`
	Desc       string        `json:"desc"`
	Gid        int           `json:"gid"`
	Guid       int           `json:"guid"`
	Id         int           `json:"id"`
	Milestones int           `json:"milestones"`
	Name       string        `json:"name"`
	ObjStatus  string        `json:"objStatus"`
	ProvidedBy []int         `json:"providedBy"`
	TechStatus string        `json:"techStatus"`
	Type       string        `json:"type"`
	Config     SepConfig     `json:"config"`
}

type SepConfig map[string]interface{}

type SepList []Sep
type SepPool map[string]interface{}
