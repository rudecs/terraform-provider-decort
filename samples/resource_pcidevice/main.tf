/*
Пример использования
Ресурса pdidevice
Ресурс позволяет:
1. Создавать устройство
2. Редактировать устройство
3. Удалять устройство

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

resource "decort_pcidevice" "pd" {
  #имя устройства
  #обязательный параметр
  #тип - строка
  name = "test_device"

  #путь до устройства
  #обязательный параметр
  #тип - строка
  hw_path = "0000:07:00.0"

  #описание устройства
  #обязательный параметр
  #тип - строка
  description = "test desc"

  #id ресурсной группы устройства
  #обязательный параметр
  #тип - число
  rg_id = 38138

  #id стака устройства
  #обязательный параметр
  #тип - число
  stack_id = 11

  #доступность устройства
  #опциональный параметр
  #может использоваться на созданном ресурсе
  #тип - булево значение
  #enable = false

  #принудительное удаение устройства
  #опциональный параметр
  #используется при удалении ресурса
  #тип - булево значение
  #force = true


  #id устройства
  #опциональный параметр
  #позволяет "восстановить" состояние ресурса терраформа на локальной машине
  #тип - число
  #device_id = 86


}

output "test" {
  value = decort_pcidevice.pd
}
