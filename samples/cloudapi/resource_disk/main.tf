/*
Пример использования
Получение списка доступных образов
*/
#Расскомментируйте этот код,
#и внесите необходимые правки в версию и путь,
#чтобы работать с установленным вручную (не через hashicorp provider registry) провайдером

terraform {
  required_providers {
    decort = {
      version = "1.1"
      source  = "digitalenergy.online/decort/decort"
    }
  }
}


provider "decort" {
  authenticator = "oauth2"
  #controller_url = <DECORT_CONTROLLER_URL>
  controller_url = "https://ds1.digitalenergy.online"
  #oauth2_url = <DECORT_SSO_URL>
  oauth2_url           = "https://sso.digitalenergy.online"
  allow_unverified_ssl = true
}

resource "decort_disk" "acl" {
  account_id  = 88366
  gid         = 212
  disk_name   = "super-disk-re"
  size_max    = 20
  restore     = true
  permanently = true
  reason      = "delete"
  shareable = false
  iotune {
    read_bytes_sec      = 0
    read_bytes_sec_max  = 0
    read_iops_sec       = 0
    read_iops_sec_max   = 0
    size_iops_sec       = 0
    total_bytes_sec     = 0
    total_bytes_sec_max = 0
    total_iops_sec      = 3000
    total_iops_sec_max  = 0
    write_bytes_sec     = 0
    write_bytes_sec_max = 0
    write_iops_sec      = 0
    write_iops_sec_max  = 0
  }

}

output "test" {
  value = decort_disk.acl
}
