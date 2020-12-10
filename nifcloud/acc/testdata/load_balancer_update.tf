provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_instance" "basic" {
  instance_id       = "%s"
  availability_zone = "east-21"
  image_id          = "221"
  key_name          = nifcloud_key_pair.basic.key_name
  depends_on = [nifcloud_key_pair.basic]

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }
}

resource "nifcloud_load_balancer" "basic" {
  load_balancer_name = "%s"
  instance_port = 80
  load_balancer_port = 80
  accounting_type = "1"
  availability_zones = ["east-21"]
  ip_version = "v4"
  instances = [nifcloud_instance.basic.instance_id]
  network_volume = 10
  ssl_certificate_id = nifcloud_ssl_certificate.basic.id
  policy_type = "standard"
  depends_on = [nifcloud_instance.basic, nifcloud_ssl_certificate.basic]
}

resource "nifcloud_key_pair" "basic" {
  key_name   = "%s"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
}

resource "nifcloud_ssl_certificate" "basic" {
  certificate = <<EOT
%s
EOT
  key         = <<EOT
%s
EOT
  ca          = <<EOT
%s
EOT
  description = "memo"
}
