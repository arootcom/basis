
## woodchuck

## Start evironment

    $ cd ./manifest
    $ docker-compose up -d

## Start API

    $ cd ./source
    $ export WOODCHUCK_LOG_LEVEL=Debug; export WOODCHUCK_VIEWS_DIR=../conf/services/; go run api.go 

