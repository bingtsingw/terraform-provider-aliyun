package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math/rand"
)

func resourceAliyunCRUserInfoAuth() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAliyunCRUserInfoAuthRead,
		CreateContext: resourceAliyunCRUserInfoAuthCreate,
		UpdateContext: resourceAlicloudCRUserInfoAuthUpdate,
		DeleteContext: resourceAlicloudCRUserInfoAuthDelete,

		Schema: map[string]*schema.Schema{
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  false,
				Sensitive: true,
			},
			"access_key": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCRUserInfoAuthDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceAlicloudCRUserInfoAuthUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn, _ := cr.NewClientWithAccessKey(d.Get("region").(string), d.Get("access_key").(string), d.Get("secret_key").(string))

	payload := &crUserInfoRequestPayload{}
	payload.User.Password = d.Get("password").(string)
	serialized, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	request := cr.CreateUpdateUserInfoRequest()
	request.SetContent(serialized)

	_, err = conn.UpdateUserInfo(request)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAliyunCRUserInfoAuthRead(ctx, d, m)
}

func resourceAliyunCRUserInfoAuthCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn, _ := cr.NewClientWithAccessKey(d.Get("region").(string), d.Get("access_key").(string), d.Get("secret_key").(string))

	payload := &crUserInfoRequestPayload{}
	payload.User.Password = d.Get("password").(string)
	serialized, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	request := cr.CreateCreateUserInfoRequest()
	request.SetContent(serialized)

	_, err = conn.CreateUserInfo(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", rand.Int()))

	return resourceAliyunCRUserInfoAuthRead(ctx, d, m)
}

func resourceAliyunCRUserInfoAuthRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
