# terraform-provider-redash

[![CI](https://github.com/winebarrel/terraform-provider-redash/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/terraform-provider-redash/actions/workflows/ci.yml)
[![terraform docs](https://img.shields.io/badge/terraform-docs-%35835CC?logo=terraform)](https://registry.terraform.io/providers/winebarrel/redash/latest/docs)

Terraform provider for [Redash](https://redash.io/).

## Usage

```tf
terraform {
  required_providers {
    redash = {
      source = "winebarrel/redash"
    }
  }
}

provider "redash" {
  url     = "http://localhost:5001" # default: $REDASH_URL
  api_key = "..."                   # default: $REDASH_API_KEY
}

resource "redash_data_source" "postgres" {
  name = "postgres"
  type = "pg"
  options = jsonencode({
    dbname = "postgres"
    host   = "postgres"
    port   = 5432
    user   = "postgres"
  })
}

resource "redash_query" "select_one" {
  data_source_id = redash_data_source.postgres.id
  name           = "select one"
  query          = "select 1"
}
```

## Tests

```sh
docker compose up -d
make redash-setup
make testacc
```

## Run a development binary

```sh
docker compose up -d
make redash-setup
cp etc/redash.tf.sample redash.tf
make tf-plan
make tf-apply
```

## Related Links

* https://github.com/winebarrel/redash-go
