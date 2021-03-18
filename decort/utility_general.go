/*
Copyright (c) 2019-2021 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Sergey Shubin, <sergey.shubin@digitalenergy.online>, <svs1370@gmail.com>

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

	"strings"

)

func Jo2JSON(arg_str string) string {
	// DECORT API historically returns response in the form of Python dictionary, which generally
	// looks like JSON, but does not comply with JSON syntax.
	// For Golang JSON Unmarshal to work properly we need to pre-process API response as follows:   
	ret_string := strings.Replace(string(arg_str), "u'", "\"", -1)
	ret_string = strings.Replace(ret_string, "'", "\"", -1)
	ret_string = strings.Replace(ret_string, ": False", ": false", -1)
	ret_string = strings.Replace(ret_string, ": True", ": true", -1)
	ret_string = strings.Replace(ret_string, "null", "\"\"", -1)
	ret_string = strings.Replace(ret_string, "None", "\"\"", -1)

	// fix for incorrect handling of usage info
	// ret_string = strings.Replace(ret_string, "<", "\"", -1)
	// ret_string = strings.Replace(ret_string, ">", "\"", -1)
	return ret_string
}
