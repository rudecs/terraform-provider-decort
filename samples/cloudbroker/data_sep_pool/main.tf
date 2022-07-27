/*
Пример использования
Получение данных sep pool

*/
#Расскомментируйте код ниже,
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

data "decort_sep_pool" "sp" {
  #id sep
  #обязательный параметр
  #тип - число
  sep_id = 1111
  #sep pool name
  #обязательный параметр
  #тип - строка
  pool_name = "sep_pool"
}

output "test" {
  value = data.decort_sep_pool.sp
}

output "pool" {
  value = jsondecode(data.decort_sep_pool.sp.pool)
}
