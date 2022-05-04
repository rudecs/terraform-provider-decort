/*
Пример использования
Ресурса sep
Ресурс позволяет:
1. Создавать sep.
2. Редактировать sep.
3. Удалять sep.
4. Конфигурировать sep.

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

resource "decort_sep" "s" {
  #grid id
  #обязательный параметр
  #тип - число
  gid = 212

  #sep name
  #обязательный параметр
  #тип - строка
  name = "test sep"

  #тип sep
  #обязательный параметр
  #тип - строка
  #возможные значения - des, dorado, tatlin, hitachi
  type = "des"

  #описание sep
  #необязательный параметр, используется при создании ресурса
  #тип - строка
  desc = "rrrrr"

  #конфигурация sep
  #необязательный параметр, мб применен при создании или редактировании sep
  #представляет собой json-строку
  #тип - строка
  #config = file("./config.json")

  #изменение поля в конфигурации
  #необязательный параметр, мб применен на уже созданном sep
  #тип - объект
  #внимание, во избежание конфликтов не использовать с полем config
  /*
  field_edit {
    #имя поля
    #обязательный параметр
    #тип - строка
    field_name  = "edgeuser_password"
    
    #значение поля
    #обязательный параметр
    #тип - json строка
    field_value = "mosk"

    #тип значения
    #обязательный параметр
    #тип - строка, возможные значения: list,dict,int,bool,str
    field_type  = "str"
  }
  */

  #доступность sep
  #необязательный параметр, мб применен на уже созданном ресурсе
  #тип - булево значение
  #enable             = false

  #использование нодами
  #необязательный параметр, используется при редактировании ресурса
  #тип - массив чисел
  #consumed_by        = []

  #обновление лимита объема
  #необязательный параметр, применяется на уж созданнном ресурсе
  #тип - булев тип
  #upd_capacity_limit = true

  #id provided nodes
  #необязательный параметр, применяется на уже созданном ресурсе
  #тип - массив чисел
  #provided_by = [16, 14, 15]

  #отключение nodes
  #необязательный параметр, применяется на уже созданном ресурсе
  #тип - булев тип
  #используется в связке с clear_physically
  #decommission      = true

  #физическое очищение nodes
  #необязательный параметр, используется при удалении ресурса
  #тип - булев тип
  #clear_physically = false

}

output "test" {
  value = decort_sep.s
}

output "config" {
  value = jsondecode(decort_sep.s.config)

}
