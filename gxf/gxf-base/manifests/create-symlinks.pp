#node default { }

node 'default' {

	exec { 'create symlinks for apps and logging config files':
		command => '/bin/sh -c "/home/dev/Sources/OSGP/Config/scripts/create-symlinks.sh"',
	}

}
