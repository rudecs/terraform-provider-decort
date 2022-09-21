/*
Пример использования
Ресурса load balancer frontend
Ресурс позволяет:
1. Создавать frontend
2. Удалять frontend

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

resource "decort_lb_frontend" "lb" {
  #id балансировщика нагрузок
  #обязательный параметр
  #тип - число
  lb_id = 668

  #имя бекенда для создания фронтенда
  #обязательный параметр
  #тип - строка
  backend_name = "testBackend"

  #имя фронтенда
  #обязательный параметр
  #тип - строка
  name = "testFrontend"

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

output "test" {
  value = decort_lb_frontend.lb
}
