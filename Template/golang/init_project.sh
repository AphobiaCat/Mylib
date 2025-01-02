#!/bin/bash

if [ -z "$1" ]; then
    echo "param 1 cant be null"
    exit 1
fi


echo "build target name:$1"

sed -i "1 c\TAR=$1" ./Makefile

sed -i "3 c\var Database_Name string = \"$1\"" ./src/public/env.go.dev
sed -i "3 c\var Database_Name string = \"$1\"" ./src/public/env.go.product

service_file=$1.service

echo "[Unit]
Description=$2
After=network.target

[Service]
ExecStart=/usr/bin/$1
Restart=always
User=nobody
Group=nogroup

[Install]
WantedBy=multi-user.target
" > $service_file



echo -n "#" > install.sh
echo -n "!" >> install.sh
echo -n "/bin/bash" >> install.sh

echo "
cp ./${service_file} /etc/systemd/system/${service_file}

cp -rf ./$1 /usr/bin/

systemctl stop ${service_file}
systemctl start ${service_file}
systemctl enable ${service_file}
" >> install.sh

chmod 777 install.sh
