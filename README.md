
## woodchuck

## Start evironment

    $ cd ./manifest
    $ docker-compose up -d

## Start API

    $ cd ./source
    $ export WOODCHUCK_LOG_LEVEL=Debug; export WOODCHUCK_VIEWS_DIR=../conf/services/; go run api.go 

## MINIO

    http://localhost:9001/

## API requests

### Basis service

#### Create bucket without versioning

    curl -v -X POST  http:/localhost:9101/ -H 'Content-Type: application/json' -H 'X-Woodchuck-Service: Base' -d '{"Type":"Bucket","Attributes":{"name":"wo-versioning","versioning":false}}'



