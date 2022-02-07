#node default { }

node 'default' {

	exec { 'creating source code folder':
		command => '/bin/mkdir -p /home/dev/Sources/OSGP',
		before => Exec['chown source code folder'],	
	}
	
	exec { 'chown source code folder':
		command => '/bin/chown -R dev:dev /home/dev/Sources',
	}
	
}
