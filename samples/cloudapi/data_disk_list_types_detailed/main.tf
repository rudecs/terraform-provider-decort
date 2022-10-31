/*
Пример использования
Получение списка типов дисков, но детально
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


data "decort_disk_list_types_detailed" "dltd" {
  #Нет входных параметров

  #Выходной параметр
  #тип - лист типов
  # items {}

  #Выходной параметр
  #Список пулов
  # pools

  #Выходной параметр
  #Имя 
  # name

  #Выходной параметр
  #Список типов
  #types
}

output "test" {
  value = data.decort_disk_list_types_detailed.dltd
}
