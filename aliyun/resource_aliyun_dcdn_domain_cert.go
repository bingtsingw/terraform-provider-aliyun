package aliyun

import (
	"context"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			"cert_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cas", "free", "upload"}, false),
			},
			"ssl_pri": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"ssl_pub": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
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
	request.CertType = d.Get("cert_type").(string)
	request.SSLProtocol = "on"
	request.ForceSet = "1"
	if v, ok := d.GetOk("ssl_pub"); ok {
		request.SSLPub = v.(string)
	}
	if v, ok := d.GetOk("ssl_pri"); ok {
		request.SSLPri = v.(string)
	}

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
	d.Set("cert_type", certInfo.CertType)
	d.Set("ssl_pub", certInfo.SSLPub)

	return diags
}
