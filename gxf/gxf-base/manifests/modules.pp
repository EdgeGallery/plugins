#class base::modules {
node 'default' {

#sudo puppet module install rcoleman/puppet_module
#	module { 'puppetlabs-git':
#		ensure => present
#	}

	exec { 'Installing module puppetlabs-stdlib':
		command => "puppet module install puppetlabs-stdlib --version 5.2.0 --ignore-dependencies --force",
		unless  => "puppet module list | grep \"puppetlabs-stdlib.*5\"",
		path    => ['/bin', '/usr/bin']
	}

	exec { 'Installing module puppetlabs-git':
		command => "puppet module install puppetlabs-git",
		unless  => "puppet module list | grep puppetlabs-git",
		path    => ['/bin', '/usr/bin']
	}

	exec { 'Installing module mjanser-eclipse':
		command => "puppet module install mjanser-eclipse",
		unless  => "puppet module list | grep mjanser-eclipse",
		path    => ['/bin', '/usr/bin']
	}

	exec { 'Installing module apt':
		command => "puppet module install puppetlabs-apt --version 2.4.0 --ignore-dependencies",
		unless  => "puppet module list | grep puppetlabs-apt",
		path    => ['/bin', '/usr/bin']
	}
        
	exec { 'Installing module python':
		command => "puppet module install stankevich-python",
		unless  => "puppet module list | grep stankevich-python",
		path    => ['/bin', '/usr/bin']
	}
}

