node 'default' {

	$homedir="/home/dev"

	package { 'awscli':
		ensure => present
	}

	class { 'python':
		pip => true
	}

	package { 'pip':
		ensure   => '20.3.4',
		provider => 'pip',
		require  => Class['python']
	}

	exec { 'pytest-runner':
		command => '/usr/bin/python -m pip install pytest-runner==5.2',
		require => Package['pip']
	}

	exec { 'awsudo':
		command => '/usr/bin/python -m pip install awsudo',
		require => Package['pip']
	}

	exec { "Add awsudo to path in ${homedir}.profile":
		command => "/bin/sed -i 's/:\$PATH/:\\/home\\/dev\\/.local\\/bin:\$PATH/g' ${homedir}/.profile",
		unless => "/bin/grep PATH= /.profile | /bin/grep .local\\/bin",
		require => Exec['awsudo']
	}
}
