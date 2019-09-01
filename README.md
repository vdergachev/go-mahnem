# GO-MAHNEM

## Description 
My first golang project, simple web cwawler of mahnem.ru with rest API and nice UI

## TODO

- Add fetch strategies 
    1. Base on dictionary
    2. Int ranges
    3. Built in search
    4. Scan existing prifiles for updates
- Add service using strategies
- Cover by tests (fetch html pages, db queries, service)

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

## References 
* https://github.com/avelino/awesome-go