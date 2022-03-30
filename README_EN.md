# terraform-provider-decort
Terraform provider for Digital Energy Cloud Orchestration Technology (DECORT) platform

NOTE: provider rc-1.25 is designed for DECORT API 3.7.x. For older API versions please use:
- DECORT API 3.6.x versions - provider version rc-1.10
- DECORT API versions prior to 3.6.0 - Terraform DECS provider (https://github.com/rudecs/terraform-provider-decs)

## Features
- Work with Compute instances, 
- Work with disks, 
- Work with k8s,
- Work with image,
- Work with reource groups,
- Work with VINS,
- Work with pfw,
- Work with accounts,
- Work with snapshots.

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
