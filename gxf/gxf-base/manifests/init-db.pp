#node default { }

node 'default' {

	# Creates the uses and databases.
	exec { 'create users and databases':
		command => '/bin/sh -c "sudo -u postgres /usr/bin/psql -p 5432 -f /home/dev/Sources/OSGP/Config/sql/create-users-and-databases.sql"',
		before => Exec['backup pg_hba.conf','change pg_hba.conf','reload config'],
	}
	
	# Creates a backup of the pg_hba.conf file.
	exec { 'backup pg_hba.conf':
		command => '/bin/cp -p /etc/postgresql/latest/main/pg_hba.conf /etc/postgresql/latest/main/pg_hba.backup',
		before => Exec['reload config'],
	}

	# Makes sure 'psql' can be used like this example: psql -U osp_admin -d osgp_core  (then enter the password when prompted).
	exec { 'change pg_hba.conf':
		command => '/bin/sh -c "sed -i \'86i\local      all     osp_admin               md5\' /etc/postgresql/latest/main/pg_hba.conf"',
		before => Exec['reload config'],
	}

	# Apply the change to pg_hba.conf by reloading the configuration for PostgreSQL.
	exec { 'reload config':
		command => 'service postgresql reload',
		path => ['/usr/bin', '/usr/sbin', '/bin']
	}
}
