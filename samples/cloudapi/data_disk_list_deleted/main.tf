/*
Пример использования
Получение списка дисков со статусом DELETED
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

data "decort_disk_list_deleted" "dld" {
  #id аккаунта для получения списка дисков
  #опциональный параметр 
  #тип - число
  #account_id = 11111

  #тип диска
  #опциональный параметр 
  #тип - строка
  #возможные типы: "b" - boot_disk, "d" - data_disk
  #type = "d"

  #кол-во страниц для вывода
  #опицональный параметр 
  #тип - число
  #page = 1

  #размер страницы
  #опицональный параметр 
  #тип - число
  #size = 1
}

output "test" {
  value = data.decort_disk_list_deleted.dld
}
