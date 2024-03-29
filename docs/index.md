---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "decort Provider"
subcategory: ""
description: |-
  
---

# decort Provider





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `authenticator` (String) Authentication mode to use when connecting to DECORT cloud API. Should be one of 'oauth2', 'legacy' or 'jwt'.
- `controller_url` (String) URL of DECORT Cloud controller to use. API calls will be directed to this URL.

### Optional

- `allow_unverified_ssl` (Boolean) If true, DECORT API will not verify SSL certificates. Use this with caution and in trusted environments only!
- `app_id` (String) Application ID to access DECORT cloud API in 'oauth2' authentication mode.
- `app_secret` (String) Application secret to access DECORT cloud API in 'oauth2' authentication mode.
- `jwt` (String) JWT to access DECORT cloud API in 'jwt' authentication mode.
- `oauth2_url` (String) OAuth2 application URL in 'oauth2' authentication mode.
- `password` (String) User password for DECORT cloud API operations in 'legacy' authentication mode.
- `user` (String) User name for DECORT cloud API operations in 'legacy' authentication mode.
