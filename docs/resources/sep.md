---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "decort_sep Resource - decort"
subcategory: ""
description: |-
  
---

# decort_sep (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `gid` (Number) grid (platform) ID
- `name` (String) SEP name
- `type` (String) type of storage

### Optional

- `clear_physically` (Boolean) clear disks and images physically
- `config` (String) sep config string
- `consumed_by` (List of Number) list of consumer nodes IDs
- `decommission` (Boolean) unlink everything that exists from SEP
- `desc` (String) sep description
- `enable` (Boolean) enable SEP after creation
- `field_edit` (Block List, Max: 1) (see [below for nested schema](#nestedblock--field_edit))
- `provided_by` (List of Number) list of provider nodes IDs
- `sep_id` (Number) sep type des id
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `upd_capacity_limit` (Boolean) Update SEP capacity limit

### Read-Only

- `ckey` (String)
- `guid` (Number)
- `id` (String) The ID of this resource.
- `meta` (List of String)
- `milestones` (Number)
- `obj_status` (String)
- `tech_status` (String)

<a id="nestedblock--field_edit"></a>
### Nested Schema for `field_edit`

Required:

- `field_name` (String)
- `field_type` (String)
- `field_value` (String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `default` (String)
- `delete` (String)
- `read` (String)
- `update` (String)


