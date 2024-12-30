-- name: GetMouseStats :one
select mouseDistance, leftClick, rightClick from mouse where id = 1;

-- name: UpdateMouseStats :one
update mouse set mouseDistance = ?, leftClick = ?, rightClick = ? where id = 1 returning *;
