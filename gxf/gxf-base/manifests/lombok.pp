node 'default' {
	
	exec { 'wget lombok':
		command => "/usr/bin/wget -q -O /home/dev/Downloads/lombok.jar https://projectlombok.org/downloads/lombok.jar",
		creates => "/home/dev/Downloads/lombok.jar",
		timeout => 1800,
		returns => [0, 4],
	}

	exec { 'cp lombok':
		command => "/bin/cp /home/dev/Downloads/lombok.jar /home/dev/Tools/eclipse/",
		onlyif => "/usr/bin/test -L /home/dev/Tools/eclipse",
		require => Exec["wget lombok"],
	}

	exec { 'add lombok line':
		command =>'/bin/echo "-javaagent:/home/dev/Tools/eclipse/lombok.jar" >> /home/dev/Tools/eclipse/eclipse.ini',
		unless => '/bin/grep -qxF /home/dev/Tools/eclipse/eclipse.ini -e "-javaagent:/home/dev/Tools/eclipse/lombok.jar"',
		require => Exec['cp lombok'],
	}
	
}
