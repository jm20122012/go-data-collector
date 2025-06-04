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