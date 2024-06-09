SELECT * FROM `user_users` uu
LEFT JOIN `user_roles_groups_users` urgu ON(urgu.user_user_id = uu.id)
LEFT JOIN `user_roles_groups` urg ON(urgu.user_roles_group_id = urg.id)
LEFT JOIN `user_roles` ur ON(urg.user_role_id = ur.id)
LEFT JOIN `user_groups` ug ON(urg.user_group_id = ug.id)
WHERE uu.`id` = 2

