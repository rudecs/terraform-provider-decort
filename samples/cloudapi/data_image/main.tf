/*
Пример использования
Получение информации об образе
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

data "decort_image" "image" {
  #id образа
  #обязательный параметр
  #тип - число
  image_id = 111

  #позывать ли информацию об удаленном образе
  #опциональный параметр
  #тип - булево значение
  #по умолчанию - false
  #show_all = false
}

output "test" {
  value = data.decort_image.image
}
