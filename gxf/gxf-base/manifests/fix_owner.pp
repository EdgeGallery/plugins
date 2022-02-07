#node default { }

node 'default' {

	$homedir='/home/dev'

	exec { "Correcting ${homedir} owner to dev:dev":
		command => "/bin/chown -R dev:dev ${homedir}"
	}

}
