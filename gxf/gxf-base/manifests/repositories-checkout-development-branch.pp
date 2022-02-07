#node default { }

node 'default' {

	exec { 'Config repo':
		command => '/bin/sh -c "cd /home/dev/Sources/OSGP/Config; /usr/bin/git checkout development"',
		returns => [0,128],
	}

	exec { 'open-smart-grid-platform repo':
		command => '/bin/sh -c "cd /home/dev/Sources/OSGP/open-smart-grid-platform; /usr/bin/git checkout development"',
		returns => [0,128],
	}

	exec { 'Documentation repo':
		command => '/bin/sh -c "cd /home/dev/Sources/OSGP/Documentation; /usr/bin/git checkout development"',
		returns => [0,128],
	}

}
