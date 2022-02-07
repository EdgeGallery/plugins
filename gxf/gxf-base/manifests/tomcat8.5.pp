#node default { }

node 'default' {

	$homedir = '/home/dev'
	$version = '8.5.50'

	# Tomcat8 is used as application server.
	exec { 'wget tomcat 8.5':
		command => "/usr/bin/wget -q -P ${homedir}/Downloads/osgp - http://archive.apache.org/dist/tomcat/tomcat-8/v$version/bin/apache-tomcat-$version.tar.gz",
		before => Exec['wget postgresql jdbc','unpack tomcat','change permissions of tomcat conf files'],
		returns => [0, 4],
	}

	exec { 'unpack tomcat':
		command => "/bin/tar xzf ${homedir}/Downloads/osgp/apache-tomcat-$version.tar.gz -C ${homedir}/Tools",
		before => Exec['wget postgresql jdbc','change permissions of tomcat conf files'],
	}

	exec { 'change permissions of tomcat conf files':
		command => "/bin/chmod 644 ${homedir}/Tools/apache-tomcat-$version/conf/*",
	}

	file { 'create tomcat link in tools':
		path => "${homedir}/Tools/tomcat",
		ensure => link,
		target => "${homedir}/Tools/apache-tomcat-$version"
	}

	file { 'create tomcat8 link in tools':
		path => "${homedir}/Tools/tomcat8",
		ensure => link,
		target => "${homedir}/Tools/apache-tomcat-$version"
	}

	# A JDBC is needed by Tomcat to connect to PostgreSQL.
	exec { 'wget postgresql jdbc':
		command => "/usr/bin/wget -q -P ${homedir}/Tools/tomcat/lib - https://jdbc.postgresql.org/download/postgresql-42.2.9.jar",
		returns => [0, 4],
		require => File['create tomcat link in tools']
	}
}
