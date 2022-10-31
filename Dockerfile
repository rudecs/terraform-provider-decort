FROM docker.io/hashicorp/terraform:latest

WORKDIR /opt/decort/tf/
COPY provider.tf ./
COPY terraform-provider-decort ./terraform.d/plugins/digitalenergy.online/decort/decort/3.1.1/linux_amd64/
RUN terraform init

WORKDIR /tf
COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh", "/bin/terraform"]
