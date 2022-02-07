node 'default' {

	exec { 'Update APT':
		command => '/usr/bin/apt-get update'
	}

	package { 'software-properties-common':
		ensure => latest,
		require => Exec['Update APT']
	}

	exec { 'install-ppa-ansible':
		command => '/usr/bin/apt-add-repository ppa:ansible/ansible',
		require => Package['software-properties-common']
	}

	exec { 'Update PPA':
		command => '/usr/bin/apt-get update',
		require => Exec['install-ppa-ansible']
	}

	package { 'ansible':
		ensure => latest,
		require => 'Exec[Update PPA]'
	}

	class { 'python':
		pip => true,
		require => Package['ansible']
	}

	package { 'python-lxml':
		ensure => latest,
		require => Class['python']
	}

	ensure_packages(['ansible-lint', 'boto3'], {
		ensure   => present,
		provider => 'pip',
		require  => [ Package['python-pip'], ],
	})
}
