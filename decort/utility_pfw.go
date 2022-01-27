package decort

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func utilityPfwCheckPresence(d *schema.ResourceData, m interface{}) (*PfwRecord, error) {
	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	resp, err := controller.decortAPICall("POST", ComputePfwListAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, nil
	}

	idS := d.Id()
	id, err := strconv.Atoi(idS)
	if err != nil {
		return nil, err
	}

	var pfws []PfwRecord
	if err := json.Unmarshal([]byte(resp), &pfws); err != nil {
		return nil, err
	}

	for _, pfw := range pfws {
		if pfw.ID == id {
			return &pfw, nil
		}
	}

	return nil, nil
}
