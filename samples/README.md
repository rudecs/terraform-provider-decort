# Примеры применения ресурсов terraform-provider-decort
Каждый файл снабжен комментариями, которые кратко описывают возможности и параметры ресурса.  
Для успешной работы необходим установленный terraform.
## Ресурсы в примерах
- data:
  - grid
  - grid_list
  - image
  - image_list
  - image_list_stacks
- resources:
  - image
  - virtual_image
  - cdrom_image
  - delete_images

## Как пользоваться примерами
1. Установить terraform
2. Установить terraform-provider-decort с помощью команды `terraform init` (выполняется автоматически), либо вручную.
3. Заменить параметр *controller_url* на ваш.
4. Заменить параметр *oauth2* на ваш.
5. Добавить ключи 
*DECORT_APP_SECRET* и *DECORT_APP_ID* 
в качестве переменных окружения, либо 
можно добавить `app_id` и `app_secret` 
в блок `provider`,что небезопасно, т.к. данные
могут быть похищены при передачи файла.
