---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "decort_account_deleted_list Data Source - terraform-provider-decort"
subcategory: ""
description: |-
  
---

# decort_account_deleted_list (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **id** (String) The ID of this resource.
- **page** (Number) Page number
- **size** (Number) Page size
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- **items** (List of Object) (see [below for nested schema](#nestedatt--items))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **default** (String)
- **read** (String)


<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- **account_id** (Number)
- **account_name** (String)
- **acl** (List of Object) (see [below for nested schema](#nestedobjatt--items--acl))
- **created_time** (Number)
- **deleted_time** (Number)
- **status** (String)
- **updated_time** (Number)

<a id="nestedobjatt--items--acl"></a>
### Nested Schema for `items.acl`

Read-Only:

- **explicit** (Boolean)
- **guid** (String)
- **right** (String)
- **status** (String)
- **type** (String)
- **user_group_id** (String)

