#!/bin/bash
./activemq start
apachectl -k start
cd ../../../tomcat8/bin
puppet apply /home/Config/puppet/manifests/postgresql.pp
sleep 100
sudo -u postgres /usr/bin/psql -p 5432 -f /home/dev/Sources/OSGP/Config/sql/fix_encoding.sql
puppet apply /home/Config/puppet/manifests/init-db.pp
nohup ./catalina.sh run &
sleep 5
psql -U osp_admin -h localhost -d osgp_core -f /home/dev/Sources/OSGP/Config/sql/create-test-org.sql | grep 'INSERT 0 1' &> /dev/null
while [ $? != 0 ]
do
	sleep 5
	psql -U osp_admin -h localhost -d osgp_core -f /home/dev/Sources/OSGP/Config/sql/create-test-org.sql | grep 'INSERT 0 1'> /dev/null
done
sleep 30
fuser -n tcp -k 8080
mv ../tmp/* ../webapps/
./catalina.sh run 
