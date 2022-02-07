#node default { }

node 'default' {

	exec { 'creating firmware folder':
		command => '/bin/mkdir -p /var/www/html/firmware',
		before => Exec['chown firmware folder','chmod firmware folder'],	
	}

	exec { 'chown firmware folder':
		command => '/bin/chown dev:dev /var/www/html/firmware',
	}

	exec { 'chmod firmware folder':
		command => '/bin/chmod 755 /var/www/html/firmware',
	}

}
