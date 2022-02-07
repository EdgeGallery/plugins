user { 'dev' : 
	name        		=> 'dev',
	ensure      		=> present, 
	shell       		=> '/bin/bash',
	# Password is dev, hashed using openssl passwd -1
	password   		=> '',
	home        		=> '/home/dev/',
	# Makes sure user has uid more than 500, ensure it can login via GUI
	system      		=> false,                 
	managehome  		=> true,
	comment     		=> 'The "dev" user',
	groups			=> ['sudo', 'vboxsf']
}

file { '/etc/lightdm/lightdm.conf':
	ensure => present,
	content => "[SeatDefaults]\nautologin-user=vagrant",
	before => Exec['make dev user default']
}

exec { 'make dev user default':
	command => '/bin/sed -i "s/vagrant/dev/g" /etc/lightdm/lightdm.conf',
	require => [File['/etc/lightdm/lightdm.conf'], User['dev']]
}

