#bash

export BUCKET_WO_VERSIONING=wo-versioning;
export BUCKET_VERSIONING=versioning;
export FILE_XML_V1="$(realpath ./request_v1.xml)"
export FILE_XML_V2="$(realpath ./request_v2.xml)"
export FILE_PDF="$(realpath ./ea_tutorial.pdf)"

echo $BUCKET_WO_VERSIONING
echo $BUCKET_VERSIONING
echo $FILE_XML_V1
echo $FILE_XML_V2
echo $FILE_PDF
echo "Start testing...."

cd ../model/bucket/
go test
cd ../.././test/

cd ../model/object
go test
cd ../../test/
read -p "Нажмите ENTER для продолжения"
