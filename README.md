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
  # see https://github.com/getredash/redash/blob/v25.1/redash/query_runner/pg.py#L149-L153
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
make redash-upgrade-db
make testacc
```

## Run a development binary

```sh
docker compose up -d
make redash-setup
make redash-upgrade-db
cp etc/redash.tf.sample redash.tf
make tf-plan
make tf-apply
```

**NOTE:**
* local Redash URL: http://localhost:5001
* email: `admin@example.com`
* password: `password`
* mail server URL: http://localhost:10081

## Related Links

* https://github.com/winebarrel/redash-go
