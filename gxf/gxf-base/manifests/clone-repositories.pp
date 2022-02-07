#node default { }

node 'default' {

	exec { 'cloning Config repo':
		command => '/usr/bin/git clone https://github.com/OSGP/Config.git /home/dev/Sources/OSGP/Config',
		creates => '/home/dev/Sources/OSGP/Config'
	}

	exec { 'cloning open-smart-grid-platform repo':
		command => '/usr/bin/git clone https://github.com/OSGP/open-smart-grid-platform.git /home/dev/Sources/OSGP/open-smart-grid-platform',
		creates => '/home/dev/Sources/OSGP/open-smart-grid-platform'
	}

	exec { 'cloning Documentation repo':
		command => '/usr/bin/git clone https://github.com/OSGP/Documentation.git /home/dev/Sources/OSGP/Documentation',
		creates => '/home/dev/Sources/OSGP/Documentation'
	}

}
