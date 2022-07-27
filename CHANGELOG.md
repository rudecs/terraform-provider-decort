### Bug fixes
- fatal error when trying to retrieve compute boot disk if former does not have one
- ignored timeouts
- wrong handling of errors when attaching network interfaces and disks to kvmvm

### New features
- parameter iotune in disk
- migrated to terraform SDKv2
- admin mode (activated by environment variable DECORT\_ADMIN\_MODE) for resources: account, k8s, image, disk, resgroup, kvmvm, vins
- parameters sep\_id and pool in kvmvm
