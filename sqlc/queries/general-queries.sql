-- name: GetDeviceList :many
SELECT * FROM device_list;

-- name: GetDeviceListByDeviceTypeName :many
SELECT 
    dl.id,
    dl.device_name,
    dl.location,
    dl.ip_address,
    dl.device_type_id
FROM device_list dl 
INNER JOIN device_types dt ON dl.device_type_id = dt.id
WHERE dt.device_type = $1;

-- name: GetDeviceListByDeviceTypeId :many
SELECT 
    dl.device_name,
    dl.location,
    dl.ip_address,
    dl.device_type_id
FROM device_list dl 
WHERE dl.device_type_id = $1;

-- name: GetAllCollectorGroups :many
select * from collector_groups;

-- name: GetEnabledCollectorGroups :many
select * from collector_groups where enabled = true;

-- name: GetDevicesByCollectorGroupID :many
select * from device_list where collector_group_id = $1;

-- name: GetDevicesByCollectorGroupName :many
select dl.*, cg.group_name
from device_list dl
inner join collector_groups cg on dl.collector_group_id = cg.id
where cg.group_name = $1;

-- name: GetEnabledDevicesByCollectorGroupID :many
select * from device_list where collector_group_id = $1 and enabled = true;

-- name: GetEnabledDevicesByCollectorGroupName :many
select dl.*, cg.group_name
from device_list dl
inner join collector_groups cg on dl.collector_group_id = cg.id
where cg.group_name = $1 and dl.enabled = true;