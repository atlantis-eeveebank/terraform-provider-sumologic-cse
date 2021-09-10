data "sumologiccse_users" "all" {}

data "sumologiccse_permissions" "all" {}

locals {
  roles = {
    Auditor = [
      "can_access_audit_logs",
      "can_comment_on_insights",
    ]
  }
}

resource "sumologiccse_role" "roles" {
  for_each    = local.roles
  name        = "${each.key} (Managed By Terraform)"
  permissions = each.value
}
