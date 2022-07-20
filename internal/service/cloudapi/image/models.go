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

/*
type History struct {
	Guid      string `json:"guid"`
	Id        int    `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

type Image struct {
	ImageId       int           `json:"id"`
	Name          string        `json:"name"`
	Url           string        `json:"url"`
	Gid           int           `json:"gid"`
	Guid          int           `json:"guid"`
	Boottype      string        `json:"bootType"`
	Imagetype     string        `json:"type"`
	Drivers       []string      `json:"drivers"`
	Hotresize     bool          `json:"hotResize"`
	Bootable      bool          `json:"bootable"`
	Username      string        `json:"username"`
	Password      string        `json:"password"`
	AccountId     int           `json:"accountId"`
	UsernameDL    string        `json:"usernameDL"`
	PasswordDL    string        `json:"passwordDL"`
	SepId         int           `json:"sepId"`
	PoolName      string        `json:"pool"`
	Architecture  string        `json:"architecture"`
	UNCPath       string        `json:"UNCPath"`
	LinkTo        int           `json:"linkTo"`
	Status        string        `json:"status"`
	TechStatus    string        `json:"techStatus"`
	Size          int           `json:"size"`
	Version       string        `json:"version"`
	Enabled       bool          `json:"enabled"`
	ComputeciId   int           `json:"computeciId"`
	Milestones    int           `json:"milestones"`
	ProviderName  string        `json:"provider_name"`
	PurgeAttempts int           `json:"purgeAttempts"`
	ReferenceId   string        `json:"referenceId"`
	ResId         string        `json:"resId"`
	ResName       string        `json:"resName"`
	Rescuecd      bool          `json:"rescuecd"`
	Meta          []interface{} `json:"_meta"`
	History       []History     `json:"history"`
	LastModified  int64         `json:"lastModified"`
	Desc          string        `json:"desc"`
	SharedWith    []int         `json:"sharedWith"`
}
*/

type Image struct {
	AccountId    int      `json:"accountId"`
	Architecture string   `json:"architecture"`
	BootType     string   `json:"bootType"`
	Bootable     bool     `json:"bootable"`
	CDROM        bool     `json:"cdrom"`
	Description  string   `json:"desc"`
	Drivers      []string `json:"drivers"`
	HotResize    bool     `json:"hotResize"`
	Id           int      `json:"id"`
	LinkTo       int      `json:"linkTo"`
	Name         string   `json:"name"`
	Pool         string   `json:"pool"`
	SepId        int      `json:"sepId"`
	Size         int      `json:"size"`
	Status       string   `json:"status"`
	Type         string   `json:"type"`
	Username     string   `json:"username"`
	Virtual      bool     `json:"virtual"`
}

type ImageList []Image

type History struct {
	Guid      string `json:"guid"`
	Id        int    `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

type ImageExtend struct {
	UNCPath       string      `json:"UNCPath"`
	CKey          string      `json:"_ckey"`
	AccountId     int         `json:"accountId"`
	Acl           interface{} `json:"acl"`
	Architecture  string      `json:"architecture"`
	BootType      string      `json:"bootType"`
	Bootable      bool        `json:"bootable"`
	ComputeCiId   int         `json:"computeciId"`
	DeletedTime   int         `json:"deletedTime"`
	Description   string      `json:"desc"`
	Drivers       []string    `json:"drivers"`
	Enabled       bool        `json:"enabled"`
	GridId        int         `json:"gid"`
	GUID          int         `json:"guid"`
	History       []History   `json:"history"`
	HotResize     bool        `json:"hotResize"`
	Id            int         `json:"id"`
	LastModified  int         `json:"lastModified"`
	LinkTo        int         `json:"linkTo"`
	Milestones    int         `json:"milestones"`
	Name          string      `json:"name"`
	Password      string      `json:"password"`
	Pool          string      `json:"pool"`
	ProviderName  string      `json:"provider_name"`
	PurgeAttempts int         `json:"purgeAttempts"`
	ResId         string      `json:"resId"`
	RescueCD      bool        `json:"rescuecd"`
	SepId         int         `json:"sepId"`
	SharedWith    []int       `json:"sharedWith"`
	Size          int         `json:"size"`
	Status        string      `json:"status"`
	TechStatus    string      `json:"techStatus"`
	Type          string      `json:"type"`
	Username      string      `json:"username"`
	Version       string      `json:"version"`
}
