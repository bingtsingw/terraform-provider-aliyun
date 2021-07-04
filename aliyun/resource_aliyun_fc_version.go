package aliyun

import (
	"context"
	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliyunFCVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAliyunFCVersionRead,
		CreateContext: resourceAliyunFCVersionCreate,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliyunFCVersionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).fcconn

	input := fc.PublishServiceVersionInput{
		ServiceName: d.Get("service_name").(*string),
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(*string)
	}

	v, err := conn.PublishServiceVersion(&fc.PublishServiceVersionInput{
		ServiceName: d.Get("service_name").(*string),
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

	v, err := conn.ListServiceVersions(&fc.ListServiceVersionsInput{
		ServiceName: d.Get("service_name").(*string),
		StartKey:    &id,
		Limit:       &limit,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	versions := v.Versions

	if len(versions) != 1 {
		return diag.Errorf("get version error")
	}

	if err := d.Set("description", versions[0].Description); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
