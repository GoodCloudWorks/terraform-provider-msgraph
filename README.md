# Terraform provider for Microsoft Graph

## ~/terraform.rc

```hcl
provider_installation {
  dev_overrides {
    "msgraph" = "GOBIN || ~/go/bin"
  }

  direct {}
}
```
