---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "decort_kvmvm Resource - decort"
subcategory: ""
description: |-
  
---

# decort_kvmvm (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `boot_disk_size` (Number) This compute instance boot disk size in GB. Make sure it is large enough to accomodate selected OS image.
- `cpu` (Number) Number of CPUs to allocate to this compute instance.
- `driver` (String) Hardware architecture of this compute instance.
- `image_id` (Number) ID of the OS image to base this compute instance on.
- `name` (String) Name of this compute. Compute names are case sensitive and must be unique in the resource group.
- `ram` (Number) Amount of RAM in MB to allocate to this compute instance.
- `rg_id` (Number) ID of the resource group where this compute should be deployed.

### Optional

- `cloud_init` (String) Optional cloud_init parameters. Applied when creating new compute instance only, ignored in all other cases.
- `description` (String) Optional text description of this compute instance.
- `detach_disks` (Boolean)
- `extra_disks` (Set of Number) Optional list of IDs of extra disks to attach to this compute. You may specify several extra disks.
- `ipa_type` (String) compute purpose
- `is` (String) system name
- `network` (Block Set, Max: 8) Optional network connection(s) for this compute. You may specify several network blocks, one for each connection. (see [below for nested schema](#nestedblock--network))
- `permanently` (Boolean)
- `pool` (String) Pool to use if sepId is set, can be also empty if needed to be chosen by system.
- `sep_id` (Number) ID of SEP to create bootDisk on. Uses image's sepId if not set.
- `started` (Boolean) Is compute started.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `account_id` (Number) ID of the account this compute instance belongs to.
- `account_name` (String) Name of the account this compute instance belongs to.
- `boot_disk_id` (Number) This compute instance boot disk ID.
- `id` (String) The ID of this resource.
- `os_users` (List of Object) Guest OS users provisioned on this compute instance. (see [below for nested schema](#nestedatt--os_users))
- `rg_name` (String) Name of the resource group where this compute instance is located.

<a id="nestedblock--network"></a>
### Nested Schema for `network`

Required:

- `net_id` (Number) ID of the network for this connection.
- `net_type` (String) Type of the network for this connection, either EXTNET or VINS.

Optional:

- `ip_address` (String) Optional IP address to assign to this connection. This IP should belong to the selected network and free for use.

Read-Only:

- `mac` (String) MAC address associated with this connection. MAC address is assigned automatically.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `default` (String)
- `delete` (String)
- `read` (String)
- `update` (String)


<a id="nestedatt--os_users"></a>
### Nested Schema for `os_users`

Read-Only:

- `guid` (String)
- `login` (String)
- `password` (String)
- `public_key` (String)


