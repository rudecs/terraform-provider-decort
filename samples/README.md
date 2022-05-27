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
  - snapshot_list
  - pcidevice_list
  - pcidevice
  - sep
  - sep_list
  - sep_disk_list
  - sep_config
  - sep_pool
  - sep_consumption
  - vgpu
  - disk_list
  - rg_list
  - account_list
  - account_computes_list
  - account_disks_list
  - account_vins_list
  - account_audits_list
  - account
  - account_rg_list
  - account_counsumed_units
  - account_counsumed_units_by_type
  - account_reserved_units
  - account_templates_list
- resources:
  - image
  - virtual_image
  - cdrom_image
  - delete_images
  - k8s
  - k8s_wg
  - snapshot
  - pcidevice
  - sep
  - sep_config
  - account

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
