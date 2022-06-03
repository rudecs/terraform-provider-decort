# terraform-provider-decort
Terraform provider для платформы Digital Energy Cloud Orchestration Technology (DECORT)

Внимание: провайдер версии rc-1.25 разработан для DECORT API 3.7.x.  
Для более старых версий можно использовать:
- DECORT API 3.6.x - версия провайдера rc-1.10
- DECORT API до 3.6.0 - terraform DECS provider (https://github.com/rudecs/terraform-provider-decs)

## Возможности провайдера
- Работа с Compute instances, 
- Работа с disks, 
- Работа с k8s,
- Работа с image,
- Работа с reource groups,
- Работа с VINS,
- Работа с pfw,
- Работа с accounts,
- Работа с snapshots,
- Работа с pcidevice,
- Работа с sep,
- Работа с vgpu,
- Работа с bservice.

Вики проекта: https://github.com/rudecs/terraform-provider-decort/wiki

## Начало
Старт возможен по двум путям:  
1. Установка через собранные пакеты.
2. Ручная установка.

### Установка через собранные пакеты.
1. Скачайте и установите terraform по ссылке: https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/aws-get-started
2. Создайте файл `main.tf` и добавьте в него следующий блок.
```terraform
provider "decort" {
  authenticator = "oauth2"
  #controller_url = <DECORT_CONTROLLER_URL>
  controller_url = "https://ds1.digitalenergy.online"
  #oauth2_url = <DECORT_SSO_URL>
  oauth2_url           = "https://sso.digitalenergy.online"
  allow_unverified_ssl = true
}
```
3. Выполните команду
```
terraform init
```
Провайдер автоматически будет установлен на ваш компьютер из terraform registry.

### Ручная установка
1. Скачайте и установите Go по ссылке: https://go.dev/dl/
2. Скачайте и установите terraform по ссылке: https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/aws-get-started
3. Склонируйте репозиторий с провайдером, выполнив команду:
```bash
git clone https://github.com/rudecs/terraform-provider-decort.git
```
4. Перейдите в скачанную папку с провайдером и выполните команду
```bash
go build -o terraform-provider-decort
```
Если вы знаете как устроен _makefile_, то можно изменить в файле `Makefile` параметры под вашу ОС и выполнить команду
```bash
make build
```
5. Полученный файл необходимо поместить:  
Linux:
```bash
~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
```
Windows:
```powershell
%APPDATA%\terraform.d\plugins\${host_name}\${namespace}\${type}\${version}\${target}
```
ВНИМАНИЕ: для ОС Windows `%APP_DATA%` является каталогом, в котором будут помещены будущие файлы terraform.
Где:
- host_name - имя хоста, держателя провайдера, например, digitalenergy.online
- namespace - пространство имен хоста, например decort 
- type - тип провайдера, может совпадать с пространством имен, например, decort
- version - версия провайдера, например 1.2
- target - версия ОС, например windows_amd64
6. После этого, создайте файл `main.tf`.
7. Добавьте в него следующий блок
```terraform
terraform {
  required_providers {
    decort = {
      version = "1.2"
      source  = "digitalenergy.online/decort/decort"
    }
  }
}
```
В поле `version` указывается версия провайдера.  
Обязательный параметр  
Тип поля - строка  
ВНИМАНИЕ: Версии в блоке и в репозитории, в который был помещен провайдер должны совпадать!

В поле `source` помещается путь до репозитория с версией вида:
```bash
${host_name}/${namespace}/${type}
```
ВНИМАНИЕ: все параметры должны совпадать с путем репозитория, в котором помещен провайдер. 

8. В консоле выполнить команду 
```bash
terraform init
```

9. Если все прошло хорошо - ошибок не будет.  

Более подробно о сборке провайдера можно найти по ссылке: https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers

## Примеры работы
Примеры работы можно найти:
- На вики проекта: https://github.com/rudecs/terraform-provider-decort/wiki
- В папке `samples`  

Схемы к terraform'у доступны:
- В папке `docs`  

Хорошей работы!
