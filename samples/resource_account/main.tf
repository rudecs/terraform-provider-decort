/*
Пример использования
Ресурса cdrom image
Ресурс позволяет:
1. Создавать образ
2. Редактировать образ
3. Удалять образ

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

resource "decort_account" "a" {
  #id аккаунта
  #обязательный параметр
  #тип - число
  account_id         = 11111
  account_name       = "new_my_account"
  username           = "isername@decs3o"
  enable             = true
  send_access_emails = true
  /*users {
    user_id     = "username_2@decs3o"
    access_type = "R"
  }
  users {
    user_id     = "username_1@decs3o"
    access_type = "R"
  }*/
  resource_limits {
    cu_m = 1024
  }


}

output "test" {
  value = decort_account.a
}
