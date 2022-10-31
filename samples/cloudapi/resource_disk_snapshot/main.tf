/*
Пример использования
Ресурс снапшота диска
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

resource "decort_disk_snapshot" "ds" {
  #Номер диска
  #обязательный параметр
  #тип - число
  disk_id = 20100

  #Ярлык диска
  #обязательный параметр
  #тип - строка
  label = "label"

  #флаг rollback
  #опциональный параметр
  #тип - bool
  rollback = false
}

output "test" {
  value = decort_disk_snapshot.ds
}
