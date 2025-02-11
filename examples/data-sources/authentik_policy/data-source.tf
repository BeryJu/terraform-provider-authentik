# To get the ID of a policy by name

data "authentik_policy" "default-authentication-flow-password-stage" {
  name = "default-authentication-flow-password-stage"
}

# Then use `data.authentik_policy.default-authentication-flow-password-stage.id`
