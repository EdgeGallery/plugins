#!/bin/bash
./activemq start
apachectl -k start
cd ../../../tomcat8/bin
puppet apply /home/Config/puppet/manifests/postgresql.pp
sleep 100
sudo -u postgres /usr/bin/psql -p 5432 -f /home/dev/Sources/OSGP/Config/sql/fix_encoding.sql
puppet apply /home/Config/puppet/manifests/init-db.pp
sleep 10
./catalina.sh run
