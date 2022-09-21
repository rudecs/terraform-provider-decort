/*
Пример использования
Работа с ресурсом kvmvm (compute)
Ресурс позволяет:
1. Создавать compute
2. Редактировать compute
3. Удалять compute
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

resource "decort_kvmvm" "comp" {
  #имя compute
  #обязательный параметр
  #мб изменен
  #тип - строка
  name = "test-tf-compute-update-new"

  #id resource group
  #обязательный параметр
  #тип - число
  rg_id = 1111

  #тип драйвера для compute
  #обязательный параметр
  #тип - строка
  driver = "KVM_X86"

  #число cpu
  #обязательный параметр
  #тип - число
  cpu = 1

  #кол-во оперативной памяти, МБ
  #обязательный параметр
  #тип - число
  ram = 2048

  #id образа диска для создания compute
  #обязательный параметр
  #тип - число
  image_id = 111

  #размер загрузочного диска
  #обязательный параметр
  #тип - число
  boot_disk_size = 20

  #описание compute
  #опциональный параметр
  #тип - строка
  description = "test update description in tf words update"


}

output "test" {
  value = decort_kvmvm.comp
}
