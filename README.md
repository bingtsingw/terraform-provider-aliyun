# terraform-provider-aliyun

## how to debug

```shell
vim ~/.terraformrc
```

```terraform
provider_installation {
  dev_overrides {
    "bingtsingw/aliyun" = "/path/to/project/terraform-provider-aliyun/bin"
  }

  direct {}
}
```
