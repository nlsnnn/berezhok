-- name: ListLocationPins :many
SELECT code, name_ru, sort_order FROM location_pins ORDER BY sort_order ASC;

-- name: GetLocationSelectedPins :many
SELECT lp.code, lp.name_ru, lp.sort_order
FROM location_selected_pins lsp
JOIN location_pins lp ON lp.code = lsp.pin_code
WHERE lsp.location_id = $1
ORDER BY lp.sort_order ASC;

-- name: DeleteLocationSelectedPins :exec
DELETE FROM location_selected_pins WHERE location_id = $1;

-- name: InsertLocationSelectedPin :exec
INSERT INTO location_selected_pins (location_id, pin_code)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: CreateLocationPin :one
INSERT INTO location_pins (code, name_ru, sort_order)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteLocationPin :exec
DELETE FROM location_pins WHERE code = $1;
