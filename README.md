# Redash Terraform Provider #
[![Actions Status][actions-image]][actions] [![Go Report Card][goreport-image]][goreport] [![Release][release-image]][releases] [![License][license-image]][license]

## Overview ##

Terraform provider for managing Redash configurations.

## Quick start ##

Assuming git is installed:

```bash
 host> git clone https://github.com/snowplow-devops/terraform-provider-redash
 host> cd terraform-provider-redash
 host> make test
 host> make
```

To remove all build files:

```bash
 host> make clean
```

To format the golang code in the source directory:

```bash
 host> make format
```

**Note:** Always run `format` before submitting any code.

**Note:** The `make test` command also generates a code coverage file which can be found at `build/coverage/coverage.html`.


## Installation

First download the pre-compiled binary for your platform from the release assets at the following links or generate the binaries locally using the provided `make` command:

https://github.com/snowplow-devops/terraform-provider-redash/releases/latest

From here you will need to move the binary into your Terraform plugins directory - depending on your platform / installation this might change but generally speaking they are located at:

* Darwin & Linux: `~/.terraform.d/plugins`
* Windows: `%APPDATA%\terraform.d\plugins`

## How to use?

### Setting up the provider

To actually start tracking Snowplow events from Terraform you will need to configure the `provider` and a `resource`:

Before using the provider, you will need an environment variable set with your personal API Key (which can be found in your Redash profile)
```bash
$ export REDASH_API_KEY="<YourPersonalAPIKeyHere>"
```

```hcl
# Minimal configuration
provider "redash" {
  redash_uri = "https://com.acme.redash"
}
```

With the provider configured, we can now use data sources and manage resources.

### Users ###
```hcl
data "redash_user" "rrunner" {
  id = 1
}

resource "redash_user" "wcoyote" {
  name   = "Wile E. Coyote"
  email  = "wcoyote@acme.com"
  groups = [32,1]
}

```

### Groups ###
```hcl
data "redash_group" "genuises" {
  id = 35
}

resource "redash_group" "runners" {
  name = "Beep Beep"
}
```

### Data Sources ###

Please note that the list of required/accepted options varies wildly by type. This is entirely dependent on the Redash installation that you are connecting to. For a detailed list of types and options, you can GET from the `/api/data_sources/types` endpoint on your Redash instance.

```hcl
data "redash_data_source" "acme_corp" {
  id = 123
}

resource "redash_data_source" "acme_corp" {
   name = "ACME Corporation Product Database"
   type = "redshift"

   options {
     host     = "newproducts.acme.com"
     port     = 5439
     dbname   = "products"
     user     = "wcoyote"
     password = "eth3LbeRt"
    }
}

resource "redash_group_data_source_attachment" "wcoyote_acme" {
  group_id       = redash_group.genuises.id
  data_source_id = redash_data_source.acme_corp.id
}
```

For more detailed documentation, please see the Terraform Provider documentaton at https://registry.terraform.io/providers/snowplow-devops/redash/latest

### Publishing

This is handled through CI/CD on Github Actions. However all binaries will be generated by using the `make` command for local publishing.

### Copyright and license

The Redash Go Client is copyright 2019-2020 Snowplow Analytics Ltd.

Licensed under the **[Apache License, Version 2.0][license]** (the "License");
you may not use this software except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[actions-image]: https://github.com/snowplow-devops/terraform-provider-redash/workflows/ci/badge.svg
[actions]: https://github.com/snowplow-devops/terraform-provider-redash/actions

[release-image]: https://img.shields.io/github/v/release/snowplow-devops/terraform-provider-redash?style=flat&color=6ad7e5
[releases]: https://github.com/snowplow-devops/terraform-provider-redash/releases

[license-image]: http://img.shields.io/badge/license-Apache--2-blue.svg?style=flat
[license]: http://www.apache.org/licenses/LICENSE-2.0

[goreport-image]: https://goreportcard.com/badge/github.com/snowplow-devops/terraform-provider-redash
[goreport]: https://goreportcard.com/report/github.com/snowplow-devops/terraform-provider-redash