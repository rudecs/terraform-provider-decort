# terraform-provider-decort

Terraform provider for Digital Energy Cloud Orchestration Technology (DECORT) platform

## Mapping of platform versions with provider versions

|  DECORT API version | Terraform provider version |
| ------ | ------ |
| 3.8.5 | 3.4.x |
| 3.8.0 - 3.8.4 | 3.3.1 |
| 3.7.x |  rc-1.25 |
| 3.6.x |  rc-1.10 |
| до 3.6.0 | [terraform-provider-decs](https://github.com/rudecs/terraform-provider-decs) |

## Working modes

The provider support two working modes:

- User mode,
- Administator mode.
  Use flag DECORT_ADMIN_MODE for swithcing beetwen modes.
  See user guide at https://github.com/rudecs/terraform-provider-decort/wiki

## Features

- Work with Compute instances,
- Work with disks,
- Work with k8s,
- Work with image,
- Work with reource groups,
- Work with VINS,
- Work with pfw,
- Work with accounts,
- Work with snapshots,
- Work with pcidevice.
- Work with sep,
- Work with vgpu,
- Work with bservice,
- Work with extnets,
- Work with locations,
- Work with load balancers.

This provider supports Import operations on pre-existing resources.

See user guide at https://github.com/rudecs/terraform-provider-decort/wiki

## Get Started

Two ways for starting:

1. Installing via binary packages
2. Manual installing

### Installing via binary packages

1. Download and install terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/aws-get-started
2. Create a file `main.tf` and add to it next section.

```terraform
provider "decort" {
  authenticator = "oauth2"
  #controller_url = <DECORT_CONTROLLER_URL>
  controller_url = "https://ds1.digitalenergy.online"
  #oauth2_url = <DECORT_SSO_URL>
  oauth2_url           = "https://sso.digitalenergy.online"
  allow_unverified_ssl = true
}
```

3. Execute next command

```
terraform init
```

The Provider will automatically install on your computer from the terrafrom registry.

### Manual installing

1. Download and install Go Programming Language: https://go.dev/dl/
2. Download and install terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/aws-get-started
3. Clone provider's repo:

```bash
git clone https://github.com/rudecs/terraform-provider-decort.git
```

4. Change directory to clone provider's and execute next command

```bash
go build -o terraform-provider-decort
```

If you have experience with _makefile_, you can change `Makefile`'s paramters and execute next command

```bash
make build
```

5. Now move compilled file to:  
   Linux:

```bash
~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
```

Windows:

```powershell
%APPDATA%\terraform.d\plugins\${host_name}/${namespace}/${type}/${version}/${target}
```

NOTE: for Windows OS `%APP_DATA%` is a cataloge, where will place terraform files.
Example:

- host_name - digitalenergy.online
- namespace - decort
- type - decort
- version - 1.2
- target - windows_amd64

6. After all, create a file `main.tf`.
7. Add to the file next code section

```terraform
terraform {
  required_providers {
    decort = {
      version = "1.2"
      source  = "digitalenergy.online/decort/decort"
    }
  }
}
```

`version`- field for provider's version
Required
String
Note: Versions in code section and in a repository must be equal!

`source` - path to repository with provider's version

```bash
${host_name}/${namespace}/${type}
```

NOTE: all paramters must be equal to the repository path!

8. Execute command in your terminal

```bash
terraform init
```

9. If everything all right - you got green message in your terminal!

More details about the provider's building process: https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers

## Examples and Samples

- Examples: https://github.com/rudecs/terraform-provider-decort/wiki
- Samples: see in repository `samples`

Terraform schemas in:

- See in repository `docs`

Good work!
