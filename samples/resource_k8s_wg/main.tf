/*
Пример использования
Ресурсов worker group
Ресурсы позволяет:
1. Создавать
2. Редактировать
3. Удалять

*/



#Расскомментируйте этот код,
#и внесите необходимые правки в версию и путь,
#чтобы работать с установленным вручную (не через hashicorp provider registry) провайдером
/*
terraform {
  required_providers {
    decort = {
      source = "terraform.local/local/decort"
      version = "1.0.0"
    }
  }
}
*/

provider "decort" {
  authenticator  = "oauth2"
  oauth2_url     = "https://sso.digitalenergy.online"
  controller_url = "https://mr4.digitalenergy.online"
  app_id         = ""
  app_secret     = ""
}


resource "decort_k8s_wg" "wg" {
  #id экземпляра k8s
  #обязательный параметр
  #тип - число
  k8s_id = 1234 //это значение должно быть и результат вызова decort_k8s.cluster.id

  #имя worker group
  #обязательный параметр
  #тип - строка
  name = "workers-2"

  #количество worker node для создания
  #опциональный параметр
  #тип - число
  #по - умолчанию - 1
  num = 2

  #количество cpu для 1 worker node
  #опциональный параметр
  #тип - число
  #по - умолчанию - 1
  cpu = 1

  #количество RAM для одной worker node в Мбайтах
  #опциональный параметр
  #тип - число
  #по-умолчанию - 1024
  ram = 1024

  #размер загрузочного диска для worker node, в Гбайтах
  #опциональный параметр
  #тип - число
  #по - умолчанию - 0
  #если установлен параметр 0, то размер диска будет равен размеру образа
  disk = 10
}


output "test_wg" {
  value = decort_k8s_wg.wg
}
