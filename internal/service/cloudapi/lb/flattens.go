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

package lb

func flattenNode(node Node) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	n := map[string]interface{}{
		"backend_ip":  node.BackendIp,
		"compute_id":  node.ComputeId,
		"frontend_ip": node.FrontendIp,
		"guid":        node.GUID,
		"mgmt_ip":     node.MGMTIp,
		"network_id":  node.NetworkId,
	}

	temp = append(temp, n)

	return temp
}

func flattendBindings(bs []Binding) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0, len(bs))
	for _, b := range bs {
		t := map[string]interface{}{
			"address": b.Address,
			"guid":    b.GUID,
			"name":    b.Name,
			"port":    b.Port,
		}
		temp = append(temp, t)
	}
	return temp
}

func flattenFrontends(fs []Frontend) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0, len(fs))
	for _, f := range fs {
		t := map[string]interface{}{
			"backend":  f.Backend,
			"bindings": flattendBindings(f.Bindings),
			"guid":     f.GUID,
			"name":     f.Name,
		}
		temp = append(temp, t)
	}

	return temp
}

func flattenServers(servers []Server) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0, len(servers))
	for _, server := range servers {
		t := map[string]interface{}{
			"address":         server.Address,
			"check":           server.Check,
			"guid":            server.GUID,
			"name":            server.Name,
			"port":            server.Port,
			"server_settings": flattenServerSettings(server.ServerSettings),
		}

		temp = append(temp, t)
	}
	return temp
}

func flattenServerSettings(defSet ServerSettings) []map[string]interface{} {
	temp := map[string]interface{}{
		"downinter": defSet.DownInter,
		"fall":      defSet.Fall,
		"guid":      defSet.GUID,
		"inter":     defSet.Inter,
		"maxconn":   defSet.MaxConn,
		"maxqueue":  defSet.MaxQueue,
		"rise":      defSet.Rise,
		"slowstart": defSet.SlowStart,
		"weight":    defSet.Weight,
	}

	res := make([]map[string]interface{}, 0)
	res = append(res, temp)
	return res
}

func flattenLBBackends(backends []Backend) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0, len(backends))
	for _, item := range backends {
		t := map[string]interface{}{
			"algorithm":               item.Algorithm,
			"guid":                    item.GUID,
			"name":                    item.Name,
			"server_default_settings": flattenServerSettings(item.ServerDefaultSettings),
			"servers":                 flattenServers(item.Servers),
		}

		temp = append(temp, t)
	}
	return temp
}
