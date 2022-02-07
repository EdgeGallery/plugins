#node default { }

node 'default' {

	package { 'gedit':
		ensure => present
	}

}
