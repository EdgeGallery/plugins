#node default { }

node 'default' {

	# PostgreSQL.
	package { 'postgresql':
		ensure => installed,
	}

	exec { 'Create latest postgresql link':
		command => '/bin/ln -sf `ls -td -- /etc/postgresql/*` /etc/postgresql/latest',
		require => Package['postgresql'],
		onlyif => '/usr/bin/test ! -f /etc/postgresql/latest'
	}

	exec { 'Update locale':
		command => '/bin/sh -c "locale-gen en_US.UTF-8 && update-locale LANG=en_US.UTF-8"',
		before => Exec['Fix encoding for postgres 9.5'],
	}

	# The following is needed for postgres 9.5 as the default encoding is here LATIN1. The create database scripts will fail when this is not fixed.
	exec { 'Fix encoding for postgres 9.5':
		command => '/bin/sh -c "sudo -u postgres /usr/bin/psql -p 5432 -f /home/dev/Sources/OSGP/Config/sql/fix_encoding.sql"',
		require => Package['postgresql'],
		onlyif => '/usr/bin/test ! -f /etc/postgresql/9.5'
	}

	service { 'postgresql':
		ensure => 'running',
		enable => true,
		require => Package['postgresql'],
	}

	exec { 'Increase connections postgres':
		command => '/bin/sed -i \'s/max_connections = .*$/max_connections = 1000 /\' /etc/postgresql/latest/main/postgresql.conf',
		returns => [0, 4],
		require => [Package['postgresql'], Exec['Create latest postgresql link']],
		notify => Service['postgresql'],
	}
}
