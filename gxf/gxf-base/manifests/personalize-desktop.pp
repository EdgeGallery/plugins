#node default { }

node 'default' {

	exec { 'add osgp wallpaper for Ubuntu':
		command => '/bin/cp -f /home/dev/Sources/OSGP/Config/vagrant/OSGP-solar.png /usr/share/backgrounds/warty-final-ubuntu.png',
		onlyif => '/usr/bin/test -f /usr/share/backgrounds/warty-final-ubuntu.png'
	}

	exec { 'add osgp wallpaper for Mate':
		command => '/bin/cp -f /home/dev/Sources/OSGP/Config/vagrant/OSGP-solar.png /usr/share/backgrounds/ubuntu-mate-common/Ubuntu-Mate-Cold.jpg',
		onlyif => '/usr/bin/test -f /usr/share/backgrounds/ubuntu-mate-common/Ubuntu-Mate-Cold.jpg'
	}
	
	exec { 'add begin osgp wallpaper for Mate':
		command => '/bin/cp -f /home/dev/Sources/OSGP/Config/vagrant/OSGP-solar.png /usr/share/backgrounds/ubuntu-mate-common/Ubuntu-Mate-Cold-lightdm.jpg',
		onlyif => '/usr/bin/test -f /usr/share/backgrounds/ubuntu-mate-common/Ubuntu-Mate-Cold-lightdm.jpg'
	}
	
	exec { 'add desktop shortcuts':
		user => 'dev',
		command => '/bin/sh /home/dev/Sources/OSGP/Config/scripts/create_desktop_shortcuts.sh',
		require => Exec['chmod create_desktop_shortcuts.sh','chmod create_application_shortcut.sh','chmod create_weblink_shortcut.sh'],
	}
	
	exec { 'chmod create_desktop_shortcuts.sh':
		command => '/bin/chmod a+rwx /home/dev/Sources/OSGP/Config/scripts/create_desktop_shortcuts.sh',
	}
	
	exec { 'chmod create_application_shortcut.sh':
		command => '/bin/chmod a+rwx /home/dev/Sources/OSGP/Config/scripts/create_application_shortcut.sh',
	}
	
	exec { 'chmod create_weblink_shortcut.sh':
		command => '/bin/chmod a+rwx /home/dev/Sources/OSGP/Config/scripts/create_weblink_shortcut.sh',
	}

}
