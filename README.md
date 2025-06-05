# go-data-collector

## To save database state:
### After base schema is created:
`pg_dump -h host -U user -d database --schema-only > schema.sql`

### Populate with device data, then create seed
`pg_dump -h host -U user -d database --data-only > seed.sql`

## To load to new DB:
1) Start the new DB
2) `psql -h host -U user -d database -f schema.sql`
3) `psql -h host -U user -d database -f seed.sql`