/*
Пример использования
Ресурс конфигурации sep 
Ресурс позволяет:
1. Получить конфигурацию
2. Изменять конфигурацию
3. Изменять отдельные поля конфигурации
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

resource "decort_sep_config" "sc" {
  #id sep
  #обязательный параметр
  #тип - число
  sep_id = 1111

  #конфигурация
  #необязательное поле, используется для изменения конфигурации
  #тип - json-строка
  #config = file("./config.json")

  #редактироваие поля
  #неоябазательный параметр, используется при редактировании ресурса
  #тип - объект
  /*
  field_edit {
    #имя поля
    #обязательный параметр
    #тип - строка
    field_name = "edgeuser_password"

    #значение поля
    #обязательный параметр
    #тип - строка
    field_value = "shshs"

    #тип поля
    #обязательный параметр
    #тип - строка
    #возможные значения - int,bool, str, dict, list
    field_type = "str"
  }
  */
}

output "sep_config" {
  value = decort_sep_config.sc
}

output "sep_config_json" {
  value = jsondecode(decort_sep_config.sc.config)
}
