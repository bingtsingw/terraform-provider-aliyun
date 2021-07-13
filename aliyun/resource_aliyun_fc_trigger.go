package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAliyunFCTrigger() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAliyunFCTriggerRead,
		CreateContext: resourceAliyunFCTriggerCreate,
		UpdateContext: resourceAlicloudFCTriggerUpdate,
		DeleteContext: resourceAlicloudFCTriggerDelete,

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"function": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validation.StringLenBetween(1, 128),
			},
			"name_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 122),
			},

			"role": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) == "timer"
				},
			},

			"source_arn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"config": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"config_mns": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"config"},
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{fc.TRIGGER_TYPE_HTTP, fc.TRIGGER_TYPE_LOG, fc.TRIGGER_TYPE_OSS, fc.TRIGGER_TYPE_TIMER, fc.TRIGGER_TYPE_MNS_TOPIC, fc.TRIGGER_TYPE_CDN_EVENTS}, false),
			},
			"qualifier": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LATEST",
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trigger_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudFCTriggerDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).fcconn

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return diag.FromErr(err)
	}

	request := &fc.DeleteTriggerInput{
		ServiceName:  StringPointer(parts[0]),
		FunctionName: StringPointer(parts[1]),
		TriggerName:  StringPointer(parts[2]),
	}

	_, err = conn.DeleteTrigger(request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ServiceNotFound", "FunctionNotFound", "TriggerNotFound"}) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceAlicloudFCTriggerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(Client).fcconn

	updateInput := &fc.UpdateTriggerInput{}

	if d.HasChange("role") {
		updateInput.InvocationRole = StringPointer(d.Get("role").(string))
	}
	if d.HasChange("config") {
		var config interface{}
		if err := json.Unmarshal([]byte(d.Get("config").(string)), &config); err != nil {
			return diag.FromErr(err)
		}
		updateInput.TriggerConfig = config
	}
	if d.HasChange("qualifier") {
		updateInput.Qualifier = StringPointer(d.Get("qualifier").(string))
	}

	if updateInput != nil {
		parts, err := ParseResourceId(d.Id(), 3)
		if err != nil {
			return diag.FromErr(err)
		}
		updateInput.ServiceName = StringPointer(parts[0])
		updateInput.FunctionName = StringPointer(parts[1])
		updateInput.TriggerName = StringPointer(parts[2])

		_, err = conn.UpdateTrigger(updateInput)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAliyunFCTriggerRead(ctx, d, m)
}

func resourceAliyunFCTriggerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(Client).fcconn

	serviceName := d.Get("service").(string)
	fcName := d.Get("function").(string)
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		name = resource.PrefixedUniqueId(v.(string))
	} else {
		name = resource.UniqueId()
	}

	var config interface{}

	if d.Get("type").(string) == string(fc.TRIGGER_TYPE_MNS_TOPIC) {
		if v, ok := d.GetOk("config_mns"); ok {
			if err := json.Unmarshal([]byte(v.(string)), &config); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		if v, ok := d.GetOk("config"); ok {
			if err := json.Unmarshal([]byte(v.(string)), &config); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	object := fc.TriggerCreateObject{
		TriggerName:    StringPointer(name),
		TriggerType:    StringPointer(d.Get("type").(string)),
		InvocationRole: StringPointer(d.Get("role").(string)),
		TriggerConfig:  config,
		Qualifier:      StringPointer(d.Get("qualifier").(string)),
	}

	if v, ok := d.GetOk("source_arn"); ok && v.(string) != "" {
		object.SourceARN = StringPointer(v.(string))
	}
	request := &fc.CreateTriggerInput{
		ServiceName:         StringPointer(serviceName),
		FunctionName:        StringPointer(fcName),
		TriggerCreateObject: object,
	}

	response, err := conn.CreateTrigger(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", serviceName, COLON_SEPARATED, fcName, COLON_SEPARATED, *response.TriggerName))

	return resourceAliyunFCTriggerRead(ctx, d, m)
}

func resourceAliyunFCTriggerRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).fcconn

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return diag.FromErr(err)
	}

	service, function, name := parts[0], parts[1], parts[2]

	trigger, err := conn.GetTrigger(&fc.GetTriggerInput{
		ServiceName:  &service,
		FunctionName: &function,
		TriggerName:  &name,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("service", parts[0])
	d.Set("function", parts[1])
	d.Set("name", trigger.TriggerName)
	d.Set("trigger_id", trigger.TriggerID)
	d.Set("role", trigger.InvocationRole)
	d.Set("source_arn", trigger.SourceARN)
	d.Set("qualifier", trigger.Qualifier)

	data, err := trigger.RawTriggerConfig.MarshalJSON()
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("type").(string) == string(fc.TRIGGER_TYPE_MNS_TOPIC) {
		if err := d.Set("config_mns", string(data)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("config", string(data)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.Set("type", trigger.TriggerType)
	d.Set("last_modified", trigger.LastModifiedTime)

	return diags
}
