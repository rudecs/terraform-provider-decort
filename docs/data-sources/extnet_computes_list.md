---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "decort_extnet_computes_list Data Source - decort"
subcategory: ""
description: |-
  
---

# decort_extnet_computes_list (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_id` (Number) filter by account ID

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `items` (List of Object) (see [below for nested schema](#nestedatt--items))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `default` (String)
- `read` (String)


<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `account_id` (Number)
- `account_name` (String)
- `extnets` (List of Object) (see [below for nested schema](#nestedobjatt--items--extnets))
- `id` (Number)
- `name` (String)
- `rg_id` (Number)
- `rg_name` (String)

<a id="nestedobjatt--items--extnets"></a>
### Nested Schema for `items.extnets`

Read-Only:

- `ipaddr` (String)
- `ipcidr` (String)
- `name` (String)
- `net_id` (Number)


