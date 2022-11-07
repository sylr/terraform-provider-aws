---
subcategory: "RAM (Resource Access Manager)"
layout: "aws"
page_title: "AWS: aws_ram_resource_share"
description: |-
  Retrieve information about a RAM Resource Share
---

# Data Source: aws_ram_resource_share

`aws_ram_resource_share` Retrieve information about a RAM Resource Share.

## Example Usage

```terraform
data "aws_ram_resource_share" "example" {
  name           = "example"
  resource_owner = "SELF"
}
```

## Search by filters

```terraform
data "aws_ram_resource_share" "tag_filter" {
  name           = "MyResourceName"
  resource_owner = "SELF"

  filter {
    name   = "NameOfTag"
    values = ["exampleNameTagValue"]
  }
}
```

## Argument Reference

The following Arguments are supported

* `name` - (Required) Name of the resource share to retrieve.
* `resource_owner` (Required) Owner of the resource share. Valid values are `SELF` or `OTHER-ACCOUNTS`.

* `resource_share_status` (Optional) Specifies that you want to retrieve details of only those resource shares that have this status. Valid values are `PENDING`, `ACTIVE`, `FAILED`, `DELETING`, and `DELETED`.
* `filter` - (Optional) Filter used to scope the list e.g., by tags. See [related docs] (https://docs.aws.amazon.com/ram/latest/APIReference/API_TagFilter.html).
    * `name` - (Required) Name of the tag key to filter on.
    * `values` - (Required) Value of the tag key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the resource share.
* `id` - ARN of the resource share.
* `status` - Status of the RAM share.
* `owning_account_id` - ID of the AWS account that owns the resource share.
* `tags` - Tags attached to the RAM share
