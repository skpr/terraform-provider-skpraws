resource "skpraws_managed_login_branding" "test" {
  client_id    = "client-id"
  user_pool_id = "user-pool-id"
  settings     = file("${path.module}/settings.json")
}
