#node default { }

node 'default' {

	# Delete OSGP download dir
	exec { 'delete OSGP download dir':
		command => '/bin/rm -rf /home/dev/Downloads/osgp',
	}
	
	# Changes the permissions for the unpacked archive folders to dev user.
	exec { 'chown /home/dev/Tools':
		command => '/bin/chown -R dev:dev /home/dev/Tools/*',
	}
}
