FROM postgres:latest

COPY db/create_user_table.sql /docker-entrypoint-initdb.d/1.sql
