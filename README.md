# GO-MAHNEM

## Description 
My first golang project, simple web cwawler of mahnem.ru with rest API and nice UI

## TODO

- Add Repository (work with db, storing fetched users data)
- Add fetch strategies (base on dictionary, on int ranges and on built in  search)
- Add service using strategies
- Cover by tests

- Add Rest API to drive fetch process
- Add UI to manage fetch process
- Cover by tests

## Maintainer
Vladimir Dergachev (4dergachev@gmail.com)

## Deps
#### Compile
* github.com/spf13/viper
* github.com/PuerkitoBio/goquery
* github.com/jackc/pgx
* github.com/lib/pq (db migrate only)
* github.com/Masterminds/squirrel
* github.com/golang-migrate/migrate
* github.com/golang-migrate/migrate/database/postgres
* github.com/mitchellh/mapstructure

#### Test
* github.com/stretchr/testify