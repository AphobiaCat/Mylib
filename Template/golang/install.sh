#!/bin/bash
cp ./test.service /etc/systemd/system/test.service

cp -rf ./test /usr/bin/

systemctl stop test.service
systemctl start test.service
systemctl enable test.service

