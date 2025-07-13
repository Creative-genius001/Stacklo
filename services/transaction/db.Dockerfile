FROM postgres:latest

COPY db/create_transaction_table.sql /docker-entrypoint-initdb.d/1.sql
