---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "decort_extnet Data Source - terraform-provider-decort"
subcategory: ""
description: |-
  
---

# decort_extnet (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **net_id** (Number)

### Optional

- **id** (String) The ID of this resource.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- **check__ips** (List of String)
- **check_ips** (List of String)
- **ckey** (String)
- **default** (Boolean)
- **default_qos** (List of Object) (see [below for nested schema](#nestedatt--default_qos))
- **desc** (String)
- **dns** (List of String)
- **excluded** (List of String)
- **free_ips** (Number)
- **gateway** (String)
- **gid** (Number)
- **guid** (Number)
- **ipcidr** (String)
- **meta** (List of String) meta
- **milestones** (Number)
- **net_name** (String)
- **network** (String)
- **network_id** (Number)
- **pre_reservations_num** (Number)
- **prefix** (Number)
- **pri_vnf_dev_id** (Number)
- **reservations** (List of Object) (see [below for nested schema](#nestedatt--reservations))
- **shared_with** (List of Number)
- **status** (String)
- **vlan_id** (Number)
- **vnfs** (List of Object) (see [below for nested schema](#nestedatt--vnfs))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **default** (String)
- **read** (String)


<a id="nestedatt--default_qos"></a>
### Nested Schema for `default_qos`

Read-Only:

- **e_rate** (Number)
- **guid** (String)
- **in_burst** (Number)
- **in_rate** (Number)


<a id="nestedatt--reservations"></a>
### Nested Schema for `reservations`

Read-Only:

- **client_type** (String)
- **desc** (String)
- **domainname** (String)
- **hostname** (String)
- **ip** (String)
- **mac** (String)
- **type** (String)
- **vm_id** (Number)


<a id="nestedatt--vnfs"></a>
### Nested Schema for `vnfs`

Read-Only:

- **dhcp** (Number)

