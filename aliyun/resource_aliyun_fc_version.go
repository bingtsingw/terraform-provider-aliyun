package aliyun

import (
	"context"
	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func resourceAliyunFCVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAliyunFCVersionRead,
		CreateContext: resourceAliyunFCVersionCreate,
		DeleteContext: resourceAliyunFCVersionDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunFCVersionDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).fcconn
	id := d.Id()
	serviceName := d.Get("service_name").(string)

	_, err := conn.DeleteServiceVersion(&fc.DeleteServiceVersionInput{
		ServiceName: &serviceName,
		VersionID:   &id,
	})

	if err != nil {
		if !strings.Contains(err.Error(), "\"HttpStatus\": 404,") {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return diags
}

func resourceAliyunFCVersionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).fcconn
	serviceName := d.Get("service_name").(string)

	input := fc.PublishServiceVersionInput{
		ServiceName: &serviceName,
	}

	if description, ok := d.GetOk("description"); ok {
		desc := description.(string)
		input.Description = &desc
	}

	v, err := conn.PublishServiceVersion(&fc.PublishServiceVersionInput{
		ServiceName: &serviceName,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*v.VersionID)

	diags = resourceAliyunFCVersionRead(ctx, d, m)

	return diags
}

func resourceAliyunFCVersionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).fcconn
	id := d.Id()
	var limit int32 = 1
	serviceName := d.Get("service_name").(string)

	v, err := conn.ListServiceVersions(&fc.ListServiceVersionsInput{
		ServiceName: &serviceName,
		StartKey:    &id,
		Limit:       &limit,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	versions := v.Versions

	if len(versions) == 1 {
		if err := d.Set("description", versions[0].Description); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
