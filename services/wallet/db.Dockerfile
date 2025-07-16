FROM postgres:latest

COPY db/create_wallet_table.sql /docker-entrypoint-initdb.d/1.sql
