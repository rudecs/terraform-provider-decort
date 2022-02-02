---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "decort_pfw Resource - terraform-provider-decort"
subcategory: ""
description: |-
  
---

# decort_pfw (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **compute_id** (Number) ID of compute instance.
- **local_base_port** (Number) Internal base port number.
- **proto** (String) Network protocol, either 'tcp' or 'udp'.
- **public_port_start** (Number) External start port number for the rule.

### Optional

- **id** (String) The ID of this resource.
- **public_port_end** (Number) End port number (inclusive) for the ranged rule.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- **local_ip** (String) IP address of compute instance.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **default** (String)
- **delete** (String)
- **read** (String)
- **update** (String)

