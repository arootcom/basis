#bash

FILE="request.xml"
FILE_XML_V1="$(realpath ./request_v1.xml)"
FILE_XML_V2="$(realpath ./request_v2.xml)"

for BUCKET in "wo-versioning" "versioning"
#for BUCKET in "versioning"
do
    if [ "$BUCKET" == "versioning" ] ; then
        VERSIONING="true"
    elif [  "$BUCKET" == "wo-versioning" ] ; then
        VERSIONING="false"
    fi
    
    echo -en  "\n\033[0;34m Запрос несуществующей корзины \033[0m\n"
    curl -v http:/localhost:9101/$BUCKET | jq
    read -p "Нажмите ENTER для продолжения"

    echo -en "\n\033[0;34m Создание корзины $BUCKET для объектов \033[0m\n"
    curl -v -X POST  http:/localhost:9101/ \
        -H 'Content-Type: application/json' \
        -d "{\"name\":\"$BUCKET\",\"versioning\":$VERSIONING}" | jq
    read -p "Нажмите ENTER для продолжения"

    echo -en "\n\033[0;34m Запрос списка доступных корзин \033[0m\n"
    curl -v http:/localhost:9101/ | jq
    read -p "Нажмите ENTER для продолжения"

    echo -en "\n\033[0;34m Запрос данных о корзине \033[0m\n"
    curl -v http:/localhost:9101/$BUCKET | jq
    read -p "Нажмите ENTER для продолжения"

    for PREFIX in "" "signed/"
    do
        if [ "$PREFIX" == "signed/" ] ; then
            FILE_TYPE="DOSSIER_CHANGE_REQUEST_SIGNED"
        else
            FILE_TYPE="DOSSIER_CHANGE_REQUEST_FORM"
        fi

        echo -en "\n\033[0;34m Загрузка файла $FILE_XML_V1 \033[0m\n"
        echo -en "\033[0;34m В обект $PREFIX$FILE \033[0m\n"
        curl -v -X POST http:/localhost:9101/$BUCKET/ \
            -H 'Content-Type: multipart/form-data' \
            --form prefix="$PREFIX" \
            --form name="$FILE" \
            --form metadata="{\"Service\":\"REGISTER_OF_MEDICINES\",\"Type\":\"$FILE_TYPE\"}" \
            --form datafile=@$FILE_XML_V1 | jq
        read -p "Нажмите ENTER для продолжения"

        echo -en "\n\033[0;34m Список объектов в корзине \033[0m\n"
        curl -v http:/localhost:9101/$BUCKET/ | jq
        read -p "Нажмите ENTER для продолжения"

        echo -en "\n\033[0;34m Список версий объекта  \033[0m\n"
        curl -v -X OPTIONS http:/localhost:9101/$BUCKET/$PREFIX$FILE | jq
        read -p "Нажмите ENTER для продолжения"

        echo -en "\n\033[0;34m Данные файла объекта \033[0m\n"
        curl -v -X GET http:/localhost:9101/$BUCKET/$PREFIX$FILE
        read -p "Нажмите ENTER для продолжения"

        echo -en "\n\033[0;34m Загрузка новой версии файла \033[0m\n"
        curl -v -X PUT http:/localhost:9101/$BUCKET/$PREFIX$FILE  \
            -H 'Content-Type: multipart/form-data' \
            --form datafile=@$FILE_XML_V2
        read -p "Нажмите ENTER для продолжения"

        echo -en "\n\033[0;34m Список версий объекта  \033[0m\n"
        curl -v -X OPTIONS http:/localhost:9101/$BUCKET/$PREFIX$FILE | jq
        read -p "Нажмите ENTER для продолжения"


        if [ "$BUCKET" == "versioning" ] ; then
            for i in 1 2
            do
                echo "versionId = "
                read VERSION_ID
                
                echo -en "\n\033[0;34m Данные файла объекта \033[0m\n"
                curl -v -X GET http:/localhost:9101/$BUCKET/$PREFIX$FILE?versionId=$VERSION_ID
                read -p "Нажмите ENTER для продолжения"
            done
        elif [  "$BUCKET" == "wo-versioning" ] ; then
            echo -en "\n\033[0;34m Данные файла объекта \033[0m\n"
            curl -v -X GET http:/localhost:9101/$BUCKET/$PREFIX$FILE
            read -p "Нажмите ENTER для продолжения"
        fi
    done
done


for BUCKET in "wo-versioning" "versioning"
#for BUCKET in "versioning"
do
    for PREFIX in "" "signed/"
    do
        echo -en "\n\033[0;34m Удаление объекта \033[0m\n"
        curl -v -X DELETE http:/localhost:9101/$BUCKET/$PREFIX$FILE | jq
        read -p "Нажмите ENTER для продолжения"
    done

    echo -en "\n\033[0;34m Удаление корзины $BUCKET \033[0m\n"
    curl -v -X DELETE http:/localhost:9101/$BUCKET |jq
    read -p "Нажмите ENTER для продолжения"
done

exit

