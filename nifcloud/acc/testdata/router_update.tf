provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_router" "basic" {
  name              = "%supd"
  description       = "memo-upd"
  availability_zone = "east-21"
  accounting_type   = "1"
  type              = "medium"
  security_group    = nifcloud_security_group.basic.group_name
  nat_table_id      = nifcloud_nat_table.basic.id
  route_table_id    = nifcloud_route_table.basic.id

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }

  network_interface {
    network_id      = nifcloud_private_lan.basic.id
    ip_address      = "192.168.1.254"
    dhcp            = true
    dhcp_config_id  = nifcloud_dhcp_config.basic.id
    dhcp_options_id = nifcloud_dhcp_option.basic.id
  }
}

resource "nifcloud_dhcp_config" "basic" {
    ipaddress_pool {
        ipaddress_pool_start = "192.168.1.50"
        ipaddress_pool_stop  = "192.168.1.100"
    }
}

resource "nifcloud_dhcp_option" "basic" {
    default_router      = "192.168.1.1"
    domain_name_servers = ["8.8.8.8", "8.8.4.4"]
}

resource "nifcloud_route_table" "basic" {
  route {
    cidr_block = "0.0.0.0/0"
    network_id = "net-COMMON_GLOBAL"
  }
}

resource "nifcloud_nat_table" "basic" {
  snat {
    rule_number                   = "1"
    protocol                      = "ALL"
    source_address                = "192.168.1.150"
    outbound_interface_network_id = "net-COMMON_GLOBAL"
  }
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name  = "%s"
  availability_zone = "east-21"
  cidr_block        = "192.168.1.0/24"
}

resource "nifcloud_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}
