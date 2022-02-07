#node default { }

node 'default' {

	file { '/etc/osp':
    ensure => 'directory',
  }
}
