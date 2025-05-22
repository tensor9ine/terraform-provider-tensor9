terraform {
  required_providers {
    tensor9 = {
      source  = "tf-registry.alpha.tensor9.com/alpha/tensor9"
    }
  }
}

provider "tensor9" {
  api_key = "deadbeef"
  endpoint = "https://localhost:12345"
}

resource "tensor9_lofi_twin" "aws_s3_bucket_bucket_x_twin" {
  rsx_id = "bucket_x"
  vars = {}
  template = "{}"
  projection_id = "0000000000000001:0000000000000001:0000000000000001"
  template_fmt = "TerraformJson"
  schema = {}
}