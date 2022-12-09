/*
Пример использования
Получение данных о compute (виртулаьной машине)
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
  controller_url = "https://ds1.digitalenergy.online"
  #oauth2_url = <DECORT_SSO_URL>
  oauth2_url           = "https://sso.digitalenergy.online"
  allow_unverified_ssl = true
}

data "decort_kvmvm" "comp" {
  #Для получения информации о виртуальной машине
  #можно воспользоваться двумя методами:
  #1. получение информации по идентификатору машины - compute_id
  #2. получение информации по имени машины и идентификатору ресурсной группы - name и rg_id

  #id виртуальной машины
  #опциональный параметр
  #тип - число
  #compute_id = 11346

  #название машины
  #опциональный параметр
  #тип - строка
  #используется вместе с параметром rg_id
  #name = "test-rg-temp"

  #id ресурсной группы
  #опциональный параметр
  #тип - число
  #используется вместе с параметром name
  #rg_id = 1825
}

output "test" {
  value = data.decort_kvmvm.comp
}
