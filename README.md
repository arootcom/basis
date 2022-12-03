
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

#### List buckets

    curl -v -X GET http://localhost:9101/

#### Create bucket

    curl -v -X POST  http:/localhost:9101/ -H 'Content-Type: application/json' -H 'X-Woodchuck-Service: Base' -d '{"Type":"Bucket","Attributes":{"name":"wo-versioning","versioning":false}}'

#### Get bucket

    curl -v -X GET http://localhost:9101/wo-versioning

### Delete bucket

    curl -v -X DELETE http://localhost:9101/wo-versioning
