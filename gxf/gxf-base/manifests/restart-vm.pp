node 'default' {

	exec { 'Restart': 
		command => '/sbin/shutdown -r +1'
	}
}


