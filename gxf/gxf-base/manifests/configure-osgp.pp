#node default { }

node 'default' {

	# Make sure to set the SunPKCS11 security provider for web-device-simulator.
	exec { 'Create web-device-simulator.properties file':
		command => '/usr/bin/touch /etc/osp/web-device-simulator.properties',
	}

	exec { 'Set property in web-device-simulator.properties file':
		command => '/bin/echo oslp.security.provider=SunPKCS11-NSS > /etc/osp/web-device-simulator.properties',
		require => Exec['Create web-device-simulator.properties file'],
	}

	exec { 'Change owner and group to dev.dev for web-device-simulator.properties file':
		command => '/bin/chown dev:dev /etc/osp/web-device-simulator.properties',
		require => Exec['Set property in web-device-simulator.properties file'],
	}
}
