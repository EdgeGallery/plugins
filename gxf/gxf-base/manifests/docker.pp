include apt

node 'dev-box' {

	# Add repo for docker
	apt::source { 'docker':
		comment	 => 'The official Docker repo',
		location => 'https://apt.dockerproject.org/repo',
		release  => 'ubuntu-xenial',
		repos    => 'main',
		key	 => {
			id      => '58118E89F3A912897C070ADBF76221572C52609D',
			server  => 'ha.pool.sks-keyservers.net'
		}
	}

	# Add linux docker prerequisites
        package { "linux-image-extra-${::kernelrelease}":
                ensure => installed
        }

        package { 'linux-image-extra-virtual':
                ensure => installed
        }

        package { 'docker-engine':
                ensure => installed
        }

	# Add user to docker group to allow connectivity
	exec {"nrpe nagios membership":
		unless => "/usr/bin/getent group docker|/usr/bin/cut -d: -f4|/bin/grep -q dev",
		command => "/usr/sbin/usermod -a -G docker dev",
    		require => Package['docker-engine'],
  	}
}
