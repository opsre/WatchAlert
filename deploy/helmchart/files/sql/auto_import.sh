#!/bin/bash
mysql -h ${MYSQL_HOST} -u root -p${MYSQL_ROOT_PASSWORD} --default-character-set=utf8mb4 -D ${MYSQL_DATABASE} < /sql/notice_template_examples.sql
mysql -h ${MYSQL_HOST} -u root -p${MYSQL_ROOT_PASSWORD} --default-character-set=utf8mb4 -D ${MYSQL_DATABASE} < /sql/rule_template_groups.sql
mysql -h ${MYSQL_HOST} -u root -p${MYSQL_ROOT_PASSWORD} --default-character-set=utf8mb4 -D ${MYSQL_DATABASE} < /sql/rule_templates.sql
mysql -h ${MYSQL_HOST} -u root -p${MYSQL_ROOT_PASSWORD} --default-character-set=utf8mb4 -D ${MYSQL_DATABASE} < /sql/tenants.sql
mysql -h ${MYSQL_HOST} -u root -p${MYSQL_ROOT_PASSWORD} --default-character-set=utf8mb4 -D ${MYSQL_DATABASE} < /sql/tenants_linked_users.sql
