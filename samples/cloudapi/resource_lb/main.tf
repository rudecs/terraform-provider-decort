/*
Пример использования
Ресурса load balancer
Ресурс позволяет:
1. Создавать load balancer
2. Редактировать load balancer
3. Удалять load balancer

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

resource "decort_lb" "lb" {
  #id ресурсной группы для со
  #обязательный параметр
  #тип - число
  rg_id = 1111

  #наименование load balancer
  #обязательный параметр
  #тип - строка
  name = "tf-test-lb"

  #id внешней сети
  #обязательный параметр
  #тип - число
  extnet_id = 6

  #id виртуальной сети
  #обязательный параметр
  #тип - число
  vins_id = 758

  #флаг запуска load balancer
  #обязательный параметр
  #тип - булев тип
  #по умолчанию - true
  #если load balancer был в статусе "stopped" (start = false), 
  #то для успешного старта, он должен быть доступен (enable = true)
  start = true

  #описание
  #опциональный параметр
  #тип - строка
  #desc      = "temp super lb for testing tf provider"

  #флаг доступности load balancer
  #необязательный параметр
  #тип - булев тип
  #enable = true

  #флаг перезапуска load balancer
  #необязательный параметр
  #тип - булев тип
  #перезагрузка срабатывает только при изменении флага с false на true
  #restart      = false

  #флаг сброса конфигурации load balancer
  #необязательный параметр
  #тип - булев тип
  #сброс срабатывает только при изменении флага с false на true
  #config_reset = false

  #флаг моментального удаления load balancer
  #необязательный параметр
  #тип - булев тип
  #по умолчанию - false
  #применяется при выполнении команды terraform destroy
  #permanently = false

  #флаг восстановления load balancer
  #необязательный параметр
  #тип - булев тип
  #восстановить можно load balancer, удаленным с флагом permanently = false
  #restore = true


  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

output "test" {
  value = decort_lb.lb
}
