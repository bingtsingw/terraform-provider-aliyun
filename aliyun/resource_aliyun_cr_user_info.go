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

type crUserInfoRequestPayload struct {
	User struct {
		Password string `json:"Password"`
	} `json:"User"`
}

func resourceAliyunCRUserInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAliyunCRUserInfoRead,
		CreateContext: resourceAliyunCRUserInfoCreate,
		UpdateContext: resourceAlicloudCRUserInfoUpdate,
		DeleteContext: resourceAlicloudCRUserInfoDelete,

		Schema: map[string]*schema.Schema{
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  false,
				Sensitive: true,
			},
		},
	}
}

func resourceAlicloudCRUserInfoDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceAlicloudCRUserInfoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(Client).crconn

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

	return resourceAliyunCRUserInfoRead(ctx, d, m)
}

func resourceAliyunCRUserInfoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(Client).crconn

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

	return resourceAliyunCRUserInfoRead(ctx, d, m)
}

func resourceAliyunCRUserInfoRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
