-- name: WriteAvtechRecord :exec
INSERT INTO
    avtech_data(id, timestamp, temp_f, temp_c, device_id, device_type_id)
VALUES (DEFAULT, $1, $2, $3, $4, $5);