/*
Пример использования
Получение списка load balancer (балансировщиков нагрузки)
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

data "decort_lb_list" "lbl" {
  #флаг влючения в результат удаленных балансироващиков нагрузки
  #опциональный параметр
  #тип - булев тип
  #значение по-умолчанию - false
  #если не задан - выводятся все доступные неудаленные балансировщики
  #includedeleted = true

  #номер страницы для отображения
  #опциональный параметр
  #тип - число
  #если не задан - выводятся все доступные данные
  #page           = 1

  #размер страницы
  #опциональный параметр
  #тип - число
  #если не задан - выводятся все доступные данные
  #size = 1
}

output "test" {
  value = data.decort_lb_list.lbl
}
