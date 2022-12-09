package k8s

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceK8sListDeletedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	k8sList, err := utilityK8sListCheckPresence(ctx, d, m, K8sListDeletedAPI)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	flattenK8sList(d, k8sList)

	return nil
}

func dataSourceK8sListDeletedSchemaMake() map[string]*schema.Schema {
	k8sListDeleted := createK8sListSchema()
	delete(k8sListDeleted, "includedeleted")
	return k8sListDeleted
}

func DataSourceK8sListDeleted() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceK8sListDeletedRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceK8sListDeletedSchemaMake(),
	}
}
