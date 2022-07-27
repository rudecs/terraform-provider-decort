/*
Пример использования
Ресурса snapshot
Ресурс позволяет:
1. Создавать snapshot
2. Удалять snapshot
3. Откатывать snapshot

*/

#Расскомментируйте этот код,
#и внесите необходимые правки в версию и путь,
#чтобы работать с установленным вручную (не через hashicorp provider registry) провайдером
/*
terraform {
  required_providers {
    decort = {
      version = "1.1"
      source  = "digitalenergy.online/decort/decort"
    }
  }
}
*/

provider "decort" {
  authenticator = "oauth2"
  #controller_url = <DECORT_CONTROLLER_URL>
  controller_url = "https://mr4.digitalenergy.online"
  #oauth2_url = <DECORT_SSO_URL>
  oauth2_url           = "https://sso.digitalenergy.online"
  allow_unverified_ssl = true
}


resource "decort_snapshot" "s" {
  #обязательный параметр
  #id вычислительной мощности
  #тип - число
  compute_id = 24074

  #обязательный параметр
  #наименование snapshot
  #тип - строка
  label = "test_ssht_3"

  #опциональный параметр
  #флаг отката
  #тип - булев тип
  #по-уолчанию - false
  #если флаг был измеен с false на true, то произойдет откат
  #rollback = false
}

output "test" {
  value = decort_snapshot.s
}
