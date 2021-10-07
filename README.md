# terraform-provider-decort
Terraform provider for Digital Energy Cloud Orchestration Technology (DECORT) platform

NOTE: provider rc-1.40 is designed for DECORT API 3.7.x. For older API versions please use:
- DECORT API 3.6.x versions - provider version rc-1.10
- DECORT API versions prior to 3.6.0 - Terraform DECS provider (https://github.com/rudecs/terraform-provider-decs)

With this provider you can manage Compute instances, disks, virtual network segments and resource 
groups in DECORT platform, as well as query the platform for information about existing resources. 
This provider supports Import operations on pre-existing resources.

See user guide at https://github.com/rudecs/terraform-provider-decort/wiki

For a quick start follow these steps (assuming that your build host is running Linux; this provider builds on Windows as well, however, some paths may differ from what is mentioned below).

1. Obtain the latest GO compiler. As of beginning 2021 it is recommended to use v.1.16.3 but as new Terraform versions are released newer Go compiler may be required, so check official Terraform repository regularly for more information.
```
    cd /tmp
    wget https://golang.org/dl/go1.16.3.linux-amd64.tar.gz
    tar xvf ./go1.16.3.linux-amd64.tar.gz
    sudo mv go /usr/local
```

2. Add the following environment variables' declarations to shell startup script:
```
    export GOPATH=/opt/gopkg:~/
    export GOROOT=/usr/local/go
    export PATH=$PATH:$GOROOT/bin
```

3. Clone Terraform Plugin SDK framework repository to $GOPKG/src/github.com/hashicorp
```
    mkdir -p $GOPKG/src/github.com/hashicorp
    cd $GOPKG/src/github.com/hashicorp
    git clone https://github.com/hashicorp/terraform-plugin-sdk.git
```

4. Clone jwt-go package repository to $GOPKG/src/github.com/dgrijalva/jwt-go:
```
    mkdir -p $GOPKG/src/github.com/dgrijalva
    cd $GOPKG/src/github.com/dgrijalva
    git clone https://github.com/dgrijalva/jwt-go.git
```

5. Clone Terraform DECORT provider repository to $GOPKG/src/github.com/terraform-provider-decort
```
    cd $GOPKG/src/github.com
    git clone https://github.com/rudecs/terraform-provider-decort.git
```

6. Build Terraform DECORT provider:
```
    cd $GOPKG/src/github.com/terraform-provider-decort
    go build -o terraform-provider-decort
```