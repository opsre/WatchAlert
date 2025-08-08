- 初始化SQL
>
> notice_template_examples.sql: 通知模版
>
> rule_template_groups.sql: 告警规则模版组
>
> rule_templates.sql: 告警规则模版
> 
> tenants.sql: 租户
> 
> tenants_linked_users.sql: 租户关联用户表
```shell
# mysql -h xxx:3306 -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < notice_template_examples.sql
# mysql -h xxx:3306 -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < rule_template_groups.sql
# mysql -h xxx:3306 -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < rule_templates.sql
# mysql -h xxx:3306 -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < tenants.sql
# mysql -h xxx:3306 -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < tenants_linked_users.sql
```