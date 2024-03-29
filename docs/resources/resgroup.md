---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "decort_resgroup Resource - decort"
subcategory: ""
description: |-
  
---

# decort_resgroup (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_id` (Number) Unique ID of the account, which this resource group belongs to.
- `name` (String) Name of this resource group. Names are case sensitive and unique within the context of a account.

### Optional

- `def_net_type` (String) Type of the network, which this resource group will use as default for its computes - PRIVATE or PUBLIC or NONE.
- `description` (String) User-defined text description of this resource group.
- `ext_ip` (String) IP address on the external netowrk to request when def_net_type=PRIVATE and ext_net_id is not 0
- `ext_net_id` (Number) ID of the external network for default ViNS. Pass 0 if def_net_type=PUBLIC or no external connection required for the defult ViNS when def_net_type=PRIVATE
- `ipcidr` (String) Address of the netowrk inside the private network segment (aka ViNS) if def_net_type=PRIVATE
- `quota` (Block List, Max: 1) Quota settings for this resource group. (see [below for nested schema](#nestedblock--quota))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `account_name` (String) Name of the account, which this resource group belongs to.
- `def_net_id` (Number) ID of the default network for this resource group (if any).
- `id` (String) The ID of this resource.

<a id="nestedblock--quota"></a>
### Nested Schema for `quota`

Optional:

- `cpu` (Number) Limit on the total number of CPUs in this resource group.
- `disk` (Number) Limit on the total volume of storage resources in this resource group, specified in GB.
- `ext_ips` (Number) Limit on the total number of external IP addresses this resource group can use.
- `ext_traffic` (Number) Limit on the total ingress network traffic for this resource group, specified in GB.
- `gpu_units` (Number) Limit on the total number of virtual GPUs this resource group can use.
- `ram` (Number) Limit on the total amount of RAM in this resource group, specified in MB.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `default` (String)
- `delete` (String)
- `read` (String)
- `update` (String)


