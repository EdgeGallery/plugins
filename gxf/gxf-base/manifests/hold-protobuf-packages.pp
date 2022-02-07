#node default { }

node 'default' {

	exec { 'hold libprotobuf7':
		command => '/bin/echo libprotobuf7 hold | sudo /usr/bin/dpkg --set-selections',
	}

	exec { 'hold libprotoc7':
		command => '/bin/echo libprotoc7 hold | sudo /usr/bin/dpkg --set-selections',
	}

	exec { 'hold protobuf-compiler':
		command => '/bin/echo protobuf-compiler hold | sudo /usr/bin/dpkg --set-selections',
	}

}
