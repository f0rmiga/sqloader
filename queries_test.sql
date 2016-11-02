-- /selectUser
SELECT *
FROM tbl_user u
WHERE u.id = $1
-- /

-- SELECT * FROM tbl_user

-- /listUser
SELECT id, name FROM tbl_user u LIMIT 1000
-- /
