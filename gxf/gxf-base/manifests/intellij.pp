#node default { }

node 'dev-box' {

        $homedir="/home/dev"
        $version="2017.2.4"

        # IntellIJ
	exec { 'Download IntellIJ':
		command => "/usr/bin/wget -q -O ${homedir}/Downloads/ideaIC-${version}.tar.gz https://download.jetbrains.com/idea/ideaIC-${version}.tar.gz",
		creates => "${homedir}/Downloads/ideaIC-${version}.tar.gz",
		timeout => 1800,		
		returns => [0, 4],
	}

	exec { 'Extract IntellIJ':
		command => "/bin/mkdir -p ${homedir}/Tools/ideaIC-${version} && /bin/tar xzf ${homedir}/Downloads/ideaIC-${version}.tar.gz -C ${homedir}/Tools/ideaIC-${version}",
		onlyif => "/usr/bin/test ! -d ${homedir}/Tools/ideaIC-${version}",
		require => Exec['Download IntellIJ'], 
	}

        exec { 'Move inner idea-IC dir':
		command => "/bin/mv ${homedir}/Tools/ideaIC-${version}/idea-IC-*/* ${homedir}/Tools/ideaIC-${version} && /bin/rm -rf ${homedir}/Tools/ideaIC-${version}/idea-IC-*",
		onlyif => "/usr/bin/test -d ${homedir}/Tools/ideaIC-${version}/idea-IC-*",
		require => Exec['Extract IntellIJ']
	}

	file { 'Create IntellIJ link':
		ensure => link,
		path => "${homedir}/Tools/ideaIC",
		target => "${homedir}/Tools/ideaIC-${version}",
		require => Exec['Move inner idea-IC dir']
	}

        file { 'Create destop link':
                path => "${homedir}/Desktop/intellij.desktop",
                ensure => present,
                mode => 744,
                owner => "dev",
                group => "dev",
                content => '[Desktop Entry]
Type=Application
Terminal=false
Icon=/home/dev/Tools/ideaIC/bin/idea.png
Name=IntellIJ
Exec=/home/dev/Tools/ideaIC/bin/idea.sh
Categories=Utility;
'
        }

}

