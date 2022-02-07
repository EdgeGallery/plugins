#node default { }

node 'default' {

	# Apache HTTPD is used as webserver with SSL and AJP.
	package { 'apache2':
		ensure => installed,
		before => Exec['enable mod_ssl','enable proxy_ajp', 'enable deflate'],
	}

	exec { 'enable mod_ssl':
		command => '/usr/sbin/a2enmod ssl',
	}

	exec { 'enable proxy_ajp':
		command => '/usr/sbin/a2enmod proxy proxy_ajp',
	}
	
	exec { 'enable deflate':
		command => '/usr/sbin/a2enmod -f deflate'
        }
}
