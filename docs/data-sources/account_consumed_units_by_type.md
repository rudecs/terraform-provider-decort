---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "decort_account_consumed_units_by_type Data Source - decort"
subcategory: ""
description: |-
  
---

# decort_account_consumed_units_by_type (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_id` (Number) ID of the account
- `cu_type` (String) cloud unit resource type

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `cu_result` (Number)
- `id` (String) The ID of this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `default` (String)
- `read` (String)


