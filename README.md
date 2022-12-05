
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

#### Delete bucket

    curl -v -X DELETE http://localhost:9101/wo-versioning

#### List objects by bucket

    curl -v -X GET http://localhost:9101/wo-versioning/

#### Create object

    curl -v -X POST http:/localhost:9101/wo-versioning/  -H 'Content-Type: multipart/form-data' -H 'X-Woodchuck-Service: Base' --form metadata='{"type":"Object","attributes":{"name":"request.xml","prefix":"files/"}}' --form filedata=@./test/request_v1.xml

#### Get data object

    culr -v -X GET http://localhost:9101/wo-versioning/files/request.xml

#### Delete object

    curl -v -X DELETE http://localhost:9101/wo-versioning/files/request.xml


