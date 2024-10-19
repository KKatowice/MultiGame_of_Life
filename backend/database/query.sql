-- name: ListUsers :one
SELECT userId FROM users WHERE users.name = (?);

-- name: CreateGame :exec
INSERT INTO games (roomId, hei, wid) VALUES (?,?,?);

-- name: CreateUser :exec
INSERT INTO users (userId, name) VALUES (?,?);

-- name: CreateLobby :exec
INSERT INTO lobby (userId, roomId, id) VALUES (?,?,?);

-- name: AllUser :many
SELECT * FROM users;

-- name: Get_wh :one
SELECT wid, hei FROM games WHERE games.roomId = (?);

/* -- name: DeleteUser :exec
DELETE FROM users WHERE userId = ?; */