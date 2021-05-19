SELECT lists.id, lists.title, lists.description
FROM lists
         INNER JOIN users_lists ON lists.id = users_lists.list_id
WHERE users_lists.user_id = $1;