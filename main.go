package main

import (
	"github.com/getlantern/terraform-provider-alicloud/alicloud"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: alicloud.Provider})
}
