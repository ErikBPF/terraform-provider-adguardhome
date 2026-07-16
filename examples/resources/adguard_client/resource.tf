# manage a client
resource "adguardhome_client" "test" {
  name = "Test Client"
  ids  = ["192.168.100.15", "test-client"]
}
