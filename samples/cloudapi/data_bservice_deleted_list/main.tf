/*
Пример использования
Получение списка удаленных basic service
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

data "decort_bservice_deleted_list" "bsdl" {
  #id аккаунта для фильтрации данных
  #опциональный параметр
  #тип - число
  #если не задан - выводятся все доступные данные
  #account_id = 11111

  #id ресурсной группы, используется для фильтрации
  #опциональный параметр
  #тип - число
  #если не задан - выводятся все доступные данные
  #rg_id      = 11111

  #номер страницы для отображения
  #опциональный параметр
  #тип - число
  #если не задан - выводятся все доступные данные
  #page = 2

  #размер страницы
  #опциональный параметр
  #тип - число
  #если не задан - выводятся все доступные данные
  #size = 3

}

output "test" {
  value = data.decort_bservice_deleted_list.bsdl
}
