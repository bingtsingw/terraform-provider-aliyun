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
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIYUN_ACCOUNT_ID", os.Getenv("ALIYUN_ACCOUNT_ID")),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIYUN_REGION", os.Getenv("ALIYUN_REGION")),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"aliyun_fc_version":         resourceAliyunFCVersion(),
			"aliyun_fc_trigger":         resourceAliyunFCTrigger(),
			"aliyun_cr_user_info":       resourceAliyunCRUserInfo(),
			"aliyun_cr_user_info_auth":  resourceAliyunCRUserInfoAuth(),
			"aliyun_dcdn_domain":        resourceAliyunDcdnDomain(),
			"aliyun_dcdn_domain_cert":   resourceAliyunDcdnDomainCert(),
			"aliyun_dcdn_domain_config": resourceAliyunDcdnDomainConfig(),
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
		AccountID: strings.TrimSpace(d.Get("account_id").(string)),
	}

	return config.Client(), nil
}
