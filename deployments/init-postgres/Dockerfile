FROM postgres:13-alpine

COPY ./migrations/ddl.sql /migrations/ddl.sql
COPY ./migrations/dml.sql /migrations/dml.sql
COPY ./migrations/dal.sql /migrations/dal.sql

RUN chown -R postgres:postgres /migrations

USER postgres