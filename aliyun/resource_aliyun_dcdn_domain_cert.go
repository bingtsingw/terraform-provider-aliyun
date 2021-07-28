package aliyun

import (
	"context"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliyunDcdnDomainCert() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAliyunDcdnDomainCertRead,
		CreateContext: resourceAliyunDcdnDomainCertCreate,
		DeleteContext: resourceAliyunDcdnDomainCertDelete,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cert_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunDcdnDomainCertDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).dcdnconn

	request := dcdn.CreateSetDcdnDomainCertificateRequest()
	request.DomainName = d.Id()
	request.SSLProtocol = "off"

	_, err := conn.SetDcdnDomainCertificate(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceAliyunDcdnDomainCertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(Client).dcdnconn
	domain := d.Get("domain_name").(string)

	request := dcdn.CreateSetDcdnDomainCertificateRequest()
	request.DomainName = domain
	request.CertName = d.Get("cert_name").(string)
	request.CertType = "cas"
	request.SSLProtocol = "on"
	request.ForceSet = "0"

	_, err := conn.SetDcdnDomainCertificate(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain)

	return resourceAliyunDcdnDomainCertRead(ctx, d, m)
}

func resourceAliyunDcdnDomainCertRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).dcdnconn

	request := dcdn.CreateDescribeDcdnDomainCertificateInfoRequest()
	request.DomainName = d.Id()

	res, err := conn.DescribeDcdnDomainCertificateInfo(request)
	if err != nil {
		return diag.FromErr(err)
	}

	certInfo := res.CertInfos.CertInfo[0]

	d.Set("domain_name", d.Id())
	d.Set("cert_name", certInfo.CertName)

	return diags
}
