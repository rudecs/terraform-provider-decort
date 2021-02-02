# terraform-provider-decort
Terraform provider for Digital Energy Cloud Orchestration Technology (DECORT)

NOTE: this provider is designed for DECORT API 3.5.x. For API versions prior to 3.5.x please use
Terraform DECS provider (https://github.com/rudecs/terraform-provider-decs).

With this provider you can manage Computes and resource groups in DECORT platform, as well as query the platform for
information about existing resources.

See user guide at https://github.com/rudecs/terraform-provider-decort/wiki

For a quick start follow these steps.
1. Obtain the latest GO compiler. As of November 2019 it is recommended to use v.1.13.x but as new Terraform versions are released newer Go compiler may be required, so check official Terraform repository regularly for more information.
```
    cd /tmp
    wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
    tar xvf ./go1.13.3.linux-amd64.tar.gz
    sudo mv go /usr/local
    # add the following environment variables' declaration to shell startup
    export GOPATH=/opt/gopkg:~/
    export GOROOT=/usr/local/go
    export PATH=$PATH:$GOROOT/bin
```

2. Clone Terraform framework repository to $GOPKG/src/github.com/hashicorp/terraform
```
    mkdir -p $GOPKG/src/github.com/hashicorp
    cd $GOPKG/src/github.com/hashicorp
    git clone https://github.com/hashicorp/terraform.git
```

3. Clone jwt-go package repository to $GOPKG/src/github.com/dgrijalva/jwt-go:
```
    mkdir -p $GOPKG/src/github.com/dgrijalva
    cd $GOPKG/src/github.com/dgrijalva
    git clone https://github.com/dgrijalva/jwt-go.git
```

4. Clone terraform-decs-provider repository to $GOPKG/src/github.com/terraform-provider-decort
```
    cd $GOPKG/src/github.com
    git clone https://github.com/rudecs/terraform-provider-decort.git
```

5. Build Terraform DECS provider:
```
    cd $GOPKG/src/github.com/terraform-provider-decort
    go build -o terraform-provider-decort
```