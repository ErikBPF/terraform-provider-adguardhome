# manage a DNS rewrite rule
resource "adguardhome_rewrite" "test" {
  domain = "example.com"
  answer = "4.3.2.1"
}
