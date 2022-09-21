/*
Пример использования
Ресурса load balancer frontend bind (привязка фронтенда балансировщика нагрузок)
Ресурс позволяет:
1. Создавать привязку
2. Редактировать привязку
3. Удалять привязку

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

resource "decort_lb_frontend_bind" "lb" {
  #id балансировщика нагрузок
  #обязательный параметр
  #тип - число
  lb_id = 668

  #имя фронтенда для создания привязки
  #обязательный параметр
  #тип - строка
  frontend_name = "testFrontend"

  #наименование привязки
  #обязательный параметр
  #тип - строка
  name = "testBinding"

  #адрес привязки фронтенда
  #обязательный параметр
  #тип - строка
  address = "111.111.111.111"

  #порт для привязки фронтенда
  #обязательный параметр
  #тип - число
  port = 1111

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

output "test" {
  value = decort_lb_frontend_bind.lb
}
