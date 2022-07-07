/*
Пример использования
Ресурса delete images
Ресурс является служебным
Его можно использоваться для быстрого удаления нескольких образов
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

resource "decort_delete_images" "my_images" {
  #массив, содержащий набор id образов для удаления
  #обязательный параметр 
  #тип - массив чисел
  image_ids = [6125]

  #параметр удаления
  #опциональный тип
  #по-умолчанию - false
  #тип - булев тип
  permanently = true

  #причина удаления 
  #обязательный параметр 
  #тип - строка
  reason = "test delete"
}

output "test" {
  value = decort_delete_images.my_images
}

/*
Применение:
1. terraform plan
2. terraform apply
3. terraform destroy


Примечание:
Данный ресурс не поддерживает обновления параметров, поэтому, для переиспользования
необходимо удалить старое состояние и повторить шаги выше.
*/
