/*
Пример использования
Ресурса cdrom image
Ресурс позволяет:
1. Создавать basic service
2. Редактировать basic service
3. Удалять basic service
4. Создавать снимки состояний basic service
5. Совершать восстановление по снимкам состояний
6. Удалять снимки состояний

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

resource "decort_bservice" "b" {
  #имя сервиса
  #обязательный параметр
  #тип - строка
  service_name = "my_test_bservice_sn"

  #id ресурсной группы
  #обязательный параметр
  #тип - число
  rg_id = 11111

  #доступность сервиса
  #необязательный параметр
  #тип - булев тип
  #используется при редактировании ресурса
  #по-умолачанию - false
  #enable = true

  #снимок состояния
  #необязательный параметр
  #тип - объект
  #используется при редактировании ресурса
  #может быть несколько в ресурсе
  /*
  snapshots {
    #имя снимка состояния
    #обязательный параметр
    #тип - строка
    label = "test_snapshot"

    #восстановление сервиса из снимка состояния
    #необязательный параметр
    #тип - булев тип
    #по-умолчанию - false
    #восстановление происходит только при переключении с false на true
    rollback = false
  }
    snapshots {
    label = "test_snapshot_1"
  }
  */

  #старт сервиса
  #необязательный параметр
  #тип - булев тип
  #используется при редактировании ресурса
  #по-умолачанию - false
  #start        = false

  #восстановление сервиса после удаления
  #необязательный параметр
  #тип - булев тип
  #используется при редактировании ресурса
  #по-умолачанию - false
  #restore      = true

  #мгновенное удаление сервиса без права восстановления
  #необязательный параметр
  #тип - булев тип
  #используется при удалении ресурса
  #по-умолачанию - false
  #permanently  = true

  #id сервиса, позволяет сформировать .tfstate, если сервис есть в платформе
  #необязательный параметр
  #тип - булев тип
  #используется при создании ресурса
  #service_id   = 11111


}

output "test" {
  value = decort_bservice.b
}