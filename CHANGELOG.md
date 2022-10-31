### New data sources

- decort_disk_snapshot_list
- decort_disk_snapshot
- decort_disk_list_deleted
- decort_disk_list_unattached
- decort_disk_list_types
- decort_disk_list_types_detailed

### New resources

- decort_disk_snapshot

### New features

- add dockerfile for creating an image for the tf provider
- change behaviour to disks: check the disk status during update the tf state
- add disk block to kvmvm resource

### New articles on wiki

- [Сборка terraform провайдера в образ](https://github.com/rudecs/terraform-provider-decort/wiki/04.05-Сборка-terraform-провайдера-в-образ)
- [Массовое создание ресурсов. Мета аргументы](https://github.com/rudecs/terraform-provider-decort/wiki/05.04-Массовое-создание-ресурсов.-Мета-аргументы)
- [Удаление ресурсов](https://github.com/rudecs/terraform-provider-decort/wiki/05.05-Удаление-ресурсов)
- [Управление снимком диска](https://github.com/rudecs/terraform-provider-decort/wiki/07.01.19-Resource-функция-decort_disk_snapshot-управление-снимком-диска)
- [Получение списка типов для диска](https://github.com/rudecs/terraform-provider-decort/wiki/06.01.39-Data-функция-decort_disk_list_types-получение-списка-типов-диска)
- [Расширенное получение списка поддерживаемых типов](https://github.com/rudecs/terraform-provider-decort/wiki/06.01.40-Data-функция-decort_disk_list_types_detailed-расширенное-получение-информации-о-поддерживаемых-типах-дисков)
- [Получение информации об удаленных дисках](https://github.com/rudecs/terraform-provider-decort/wiki/06.01.41-Data-функция-decort_disk_list_deleted-получение-информации-об-удаленных-дисках)
- [Получение информации о неподключенных дисках](https://github.com/rudecs/terraform-provider-decort/wiki/06.01.42-Data-функция-decort_disk_list_unattached-получение-информации-о-неподключенных-дисках)
- [Получение списка снимков состояния диска](https://github.com/rudecs/terraform-provider-decort/wiki/06.01.43-Data-функция-decort_disk_snapshot_list-получение-списка-снимков-состояния-диска)
- [Получение информацуии о снимке состояния диска](https://github.com/rudecs/terraform-provider-decort/wiki/06.01.44-Data-функция-decort_disk_snapshot-получение-информации-о-снимке-состояния)

### Update articles

- [Управление дисковыми ресурсами.](https://github.com/rudecs/terraform-provider-decort/wiki/07.01.03-Resource-функция-decort_disk-управление-дисковыми-ресурсами)
- [Управление виртуальными серверами, создаваемыми на базе системы виртуализации KVM](https://github.com/rudecs/terraform-provider-decort/wiki/07.01.01-Resource-функция-decort_kvmvm-управление-виртуальными-машинами-на-базе-KVM)
