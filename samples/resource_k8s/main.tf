/*
Пример использования
Ресурсов k8s cluster
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


resource "decort_k8s" "cluster" {
  #имя кластера
  #обязательный параметр
  #при изменении - обновдяет имя кластера
  #тип - строка
  name = "tftest"

  #id resource group
  #обязательный параметр
  #тип - число
  rg_id = 776

  #id catalogue item 
  #обязательный параметр
  #тип - число
  k8sci_id = 9

  #имя для первой worker group, созданной в кластере
  #обязательный параметр
  #тип - строка
  wg_name = "workers"

  #настройка мастер node или nodes
  #опциональный параметр
  #максимальное кол-во элементов - 1
  #тип - список нод
  masters {
    #кол-во node
    #обязательный параметр
    #тип - число
    num = 1

    #кол-во cpu
    #обязательный параметр
    #тип - число
    cpu = 2


    #кол-во RAM в Мбайтах
    #обязательный параметр
    #тип - число
    ram = 2048


    #размер диска в Гбайтах
    #обязательный параметр
    #тип - число
    disk = 10
  }

  #настройка worker node или nodes
  #опциональный параметр
  #максимальное кол-во элементов - 1
  #тип - список нод
  workers {
    #кол-во node
    #обязательный параметр
    #тип - число
    num = 1

    #кол-во cpu
    #обязательный параметр
    #тип - число
    cpu = 2

    #кол-во RAM в Мбайтах
    #обязательный параметр
    #тип - число
    ram = 2048

    #размер диска в Гбайтах
    #обязательный параметр
    #тип - число
    disk = 10
  }
}

output "test_cluster" {
  value = decort_k8s.cluster
}

/*
output "kubeconfig"{
  value = decort_k8s.cluster.kubeconfig
}

*/
