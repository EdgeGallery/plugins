node 'default' {
	
        exec { 'Configure activemq for development':
		command => "/bin/sed -i 's/<broker xmlns=/<broker deleteAllMessagesOnStartup=\"true\" xmlns=/g' /home/dev/Tools/activemq/conf/activemq.xml"
        }
}
