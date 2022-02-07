#node default { }

node 'default' {

	exec { 'creating osgp download folder':
		command => '/bin/mkdir -p /home/dev/Downloads/osgp',
		before => Exec['chown osgp download folder'],
	}

	exec { 'chown osgp download folder':
		command => '/bin/chown -R dev:dev /home/dev/Downloads',
	}
}
