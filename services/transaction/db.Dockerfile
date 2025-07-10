FROM postgres:latest

COPY db/create_transactions_table.sql /docker-entrypoint-initdb.d/1.sql
