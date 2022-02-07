node 'default' {

	$version = "5.15.7"

	# ActiveMQ is used as message broker.
	exec { 'wget activemq':
		command => "/usr/bin/wget -q -P /home/dev/Downloads/osgp - http://archive.apache.org/dist/activemq/${version}/apache-activemq-${version}-bin.tar.gz",
		onlyif => '/usr/bin/test ! -f /home/dev/Tools/activemq',
		before => Exec['unpack activemq'],
		returns => [0, 4],
	}

	exec { 'unpack activemq':
		command => "/bin/tar xzf /home/dev/Downloads/osgp/apache-activemq-${version}-bin.tar.gz -C /home/dev/Tools",
		creates => "/home/dev/Tools/apache-activemq-${version}",
		onlyif => '/usr/bin/test ! -f /home/dev/Tools/activemq',
		require => Exec['wget activemq']
	}

	file { 'create activemq link':
		ensure => link,
		path => '/home/dev/Tools/activemq',
		target => "/home/dev/Tools/apache-activemq-${version}",
		require => Exec['unpack activemq']
	}

	exec { 'Configure activemq':
		command => "/bin/mv /home/dev/Tools/apache-activemq-${version}/conf/activemq.xml /home/dev/Tools/apache-activemq-${version}/conf/activemq.xml.original ; /bin/cp /home/dev/Sources/OSGP/Config/apache-activemq/activemq.xml /home/dev/Tools/apache-activemq-${version}/conf/",
		require => File['create activemq link']
	}

	file { '/etc/osp':
    ensure => 'directory',
		mode =>  "0755",
  }

	file { 'create activemq ssl folder':
		path => '/etc/osp/activemq',
		ensure => directory,
		mode =>  "0755",
		require => [
		    Exec['unpack activemq'],
				File['/etc/osp'],]
	}

	file { 'create activemq client.ks link':
		ensure => link,
		path => '/etc/osp/activemq/client.ks',
		target => '/home/dev/Tools/activemq/conf/client.ks',
		require => File['create activemq ssl folder']
	}

	file { 'create activemq client.ts link':
		ensure => link,
		path => '/etc/osp/activemq/client.ts',
		target => '/home/dev/Tools/activemq/conf/client.ts',
		require => File['create activemq ssl folder']
	}

}
