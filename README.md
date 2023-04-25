# terraform-provider-redash

[![terraform docs](https://img.shields.io/badge/terraform-docs-%35835CC?logo=terraform)](https://registry.terraform.io/providers/winebarrel/redash/latest/docs)

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
  url     = "http://localhost:5001"
  api_key = "..."
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
  data_source_id = resource.redash_data_source.postgres.id
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
