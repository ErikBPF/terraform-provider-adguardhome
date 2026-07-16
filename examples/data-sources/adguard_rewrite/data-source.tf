# get a DNS rewrite rule
data "adguardhome_rewrite" "test" {
  domain = "example.org"
  answer = "1.2.3.4"
}
