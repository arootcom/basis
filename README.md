
## woodchuck

## Start evironment

    $ cd ./manifest
    $ docker-compose up -d

## Start API

    $ cd ./source
    $ export BASIS_LOG_LEVEL=Debug; export BASIS_VIEWS_DIR=../conf/services/; go run api.go 
