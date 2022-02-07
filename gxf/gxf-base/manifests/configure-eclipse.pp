#node default { }

node 'default' {

	$homedir='/home/dev'

	# Extract OSGP workspace
	exec { 'unpack OSGP eclipse workspace':
		command => "/bin/mkdir -p ${homedir}/workspace && /bin/tar xzf ${homedir}/Sources/OSGP/Config/code-format-settings/eclipse/osgp.tar.gz -C ${homedir}/workspace",
		onlyif => "/usr/bin/test ! -d ${homedir}/workspace/OSGP",
	}

}
