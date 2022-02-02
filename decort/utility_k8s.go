package decort

import (
	"encoding/json"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func utilityK8sCheckPresence(d *schema.ResourceData, m interface{}) (*K8sRecord, error) {
	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", d.Id())

	resp, err := controller.decortAPICall("POST", K8sGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, nil
	}

	var k8s K8sRecord
	if err := json.Unmarshal([]byte(resp), &k8s); err != nil {
		return nil, err
	}

	return &k8s, nil
}
