/*
Пример использования
Ресурса cdrom image
Ресурс позволяет:
1. Создавать образ
2. Редактировать образ
3. Удалять образ

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
  #controller_url = "https://ds1.digitalenergy.online"
  controller_url = "https://gamma.dev.decs.online"
  #oauth2_url = <DECORT_SSO_URL>
  #oauth2_url           = "https://sso.digitalenergy.online"
  oauth2_url           = "https://sso-gamma.dev.decs.online:8443"
  allow_unverified_ssl = true
}

resource "decort_sep_des" "sd" {
  sep_id             = 11
  gid                = 214
  name               = "test sep"
  desc               = "rrrrr"
  enable             = false
  consumed_by        = []
  upd_capacity_limit = true
  #provided_by = [16, 14, 15]
  #decommision      = true
  #clear_physically = false

}

output "test" {
  //value = tolist(data.decort_sep_des.sl.config)[0].api_ips
  value = decort_sep_des.sd
}
