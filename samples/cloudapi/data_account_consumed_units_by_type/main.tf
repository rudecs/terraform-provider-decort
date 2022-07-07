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

data "decort_account_consumed_units_by_type" "acubt" {
  #id аккаунта
  #обязательный параметр
  #тип - число
  account_id = 33333

  #тип вычислительной единицы
  #обязательный параметр
  #тип - строка
  #значения:
  #cu_c - кол-во виртуальных cpu ядер
  #cu_m - кол-во RAM в МБ
  #cu_d - кол-в используемой дисковой памяти, в ГБ
  #cu_i - кол-во публичных ip адресов
  #cu_np - кол-во полученного/отданного трафика, в ГБ
  #gpu_units - кол-во gpu ядер
  cu_type = "cu_с"
}

output "test" {
  value = data.decort_account_consumed_units_by_type.acubt
}
