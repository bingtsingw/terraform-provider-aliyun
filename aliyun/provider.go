package aliyun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"strings"
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIYUN_ACCESS_KEY", os.Getenv("ALIYUN_ACCESS_KEY")),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIYUN_SECRET_KEY", os.Getenv("ALIYUN_SECRET_KEY")),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIYUN_REGION", os.Getenv("ALIYUN_REGION")),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"aliyun_fc_version": resourceAliyunFCVersion(),
		},
		ConfigureContextFunc: providerConfigure,
	}

	return provider
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		AccessKey: strings.TrimSpace(d.Get("access_key").(string)),
		SecretKey: strings.TrimSpace(d.Get("secret_key").(string)),
		RegionId:  strings.TrimSpace(d.Get("region").(string)),
		Region:    Region(strings.TrimSpace(d.Get("region").(string))),
	}

	return config.Client(), nil
}
