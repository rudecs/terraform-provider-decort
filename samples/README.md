# Примеры применения ресурсов terraform-provider-decort

Каждый файл снабжен комментариями, которые кратко описывают возможности и параметры ресурса.  
Для успешной работы необходим установленный terraform.

## Ресурсы в примерах

- cloudapi:
  - data:
    - image
    - image_list
    - image_list_stacks
    - snapshot_list
    - pcidevice_list
    - pcidevice
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
    - account_deleted_list
    - bservice_list
    - bservice_snapshot_list
    - bservice_deleted_list
    - bservice
    - bservice_group
    - extnet_default
    - extnet_list
    - extnet
    - extnet_computes_list
    - vins_list
    - locations_list
    - location_url
    - lb
    - lb_list
    - lb_list_deleted
    - disk_list_deleted
    - disk_list_unattached
    - disk_list_types
    - disk_list_types_detailed
    - disk_snapshot_list
    - disk_snapshot
  - resources:
    - image
    - virtual_image
    - cdrom_image
    - delete_images
    - k8s
    - k8s_wg
    - snapshot
    - pcidevice
    - account
    - bservice
    - bservice_group
    - lb
    - lb_frontend
    - lb_backend
    - lb_frontend_bind
    - lb_backend_server
    - disk_snapshot
- cloudbroker:
  - data:
    - grid
    - grid_list
    - image
    - image_list
    - image_list_stacks
    - pcidevice_list
    - pcidevice
    - sep
    - sep_list
    - sep_disk_list
    - sep_config
    - sep_pool
    - sep_consumption
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
    - account_deleted_list
    - vins_list
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
    - vins

## Как пользоваться примерами

1. Установить terraform
2. Установить terraform-provider-decort с помощью команды `terraform init` (выполняется автоматически), либо вручную.
3. Заменить параметр _controller_url_ на ваш.
4. Заменить параметр _oauth2_ на ваш.
5. Добавить ключи
   _DECORT_APP_SECRET_ и _DECORT_APP_ID_
   в качестве переменных окружения, либо
   можно добавить `app_id` и `app_secret`
   в блок `provider`,что небезопасно, т.к. данные
   могут быть похищены при передачи файла.
