/*
Пример использования
Получение информации о k8s кластере
*/
#Расскомментируйте этот код,
#и внесите необходимые правки в версию и путь,
#чтобы работать с установленным вручную (не через hashicorp provider registry) провайдером

terraform {
  required_providers {
    decort = {
      version = "1.1"
      source  = "digitalenergy.online/decort/decort"
    }
  }
}


provider "decort" {
  authenticator = "oauth2"
  #controller_url = <DECORT_CONTROLLER_URL>
  controller_url = "https://ds1.digitalenergy.online"
  #oauth2_url = <DECORT_SSO_URL>
  oauth2_url           = "https://sso.digitalenergy.online"
  allow_unverified_ssl = true
}

data "decort_k8s_wf" "k8s_wg" {
  #id кластера
  #обязательный параметр
  #тип - число
  k8s_id = 49304

  #id группы воркеров
  #обязательный параметр
  #тип - число
  wg_id = 43329
}

output "output_k8s_wg" {
  value = data.decort_k8s.k8s
}
