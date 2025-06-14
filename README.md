**NOTE:** All commands are run using taskfiles.  Refer to task definition for actual commands.

Timescale requires special handling for schema restores. Follow the steps below to correctly dump and restore a TimescaleDB schema

1) After base schema is created: 
`task db:dump-tsdb-schema-local`

2) Restore the schema to new DB (Make sure .env is populated correctly):
`task db:restore-tsdb-schema-remote`

3) Manually need to recreate hypertables (Make sure .env is populated: correctly)
`task db:create-hypertables-remote`

4) Use normal psql to seed the DB if required:
`psql -h host -U user -d database -f seed.sql`