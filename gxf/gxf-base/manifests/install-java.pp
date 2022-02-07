#node default { }

node 'default' {

	# Will install the default JDK. On 16.04 this will install java-8-openjdk.
	package { 'default-jdk':
		ensure => installed,
	}

	# Enable the SunPKCS11 security provider.
	exec { 'Add NSS library to java.security':
		command => '/bin/sed -i "/SunPCSC/a security.provider.10=sun.security.pkcs11.SunPKCS11 \$\{java.home\}\/lib\/security\/nss.cfg" /usr/lib/jvm/java-8-openjdk-amd64/jre/lib/security/java.security',
		onlyif => '/usr/bin/test $(/bin/grep -c security.provider.10 /usr/lib/jvm/java-8-openjdk-amd64/jre/lib/security/java.security) -eq 0',
		require => Package['default-jdk']
	}

	#exec { 'Fix NSS library dir':
	#	command => '/bin/sed -i "s/@NSS_LIBDIR@/\/usr\/lib64/g" /usr/lib/jvm/java-8-openjdk-amd64/jre/lib/security/nss.cfg',
	#	onlyif => '/usr/bin/test -f /usr/lib/jvm/java-8-openjdk-amd64/jre/lib/security/nss.cfg',
	#	require => Exec['Add NSS library to java.security'],
	#}
        file_line { 'NSS library dir':
               path => '/usr/lib/jvm/java-8-openjdk-amd64/jre/lib/security/nss.cfg',
               line => 'nssLibraryDirectory = /usr/lib/x86_64-linux-gnu'
        }

}
