package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAliyunDcdnDomainConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAliyunDcdnDomainConfigCreate,
		ReadContext:   resourceAliyunDcdnDomainConfigRead,
		DeleteContext: resourceAliyunDcdnDomainConfigDelete,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(5, 67),
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"function_args": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arg_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"arg_value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunDcdnDomainConfigDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).dcdnconn

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return diag.FromErr(err)
	}
	request := dcdn.CreateDescribeDcdnDomainConfigsRequest()
	request.DomainName = parts[0]
	request.FunctionNames = parts[1]

	res, err := conn.DescribeDcdnDomainConfigs(request)
	if err != nil {
		return diag.FromErr(err)
	}
	config := res.DomainConfigs.DomainConfig[0]

	deleteRequest := dcdn.CreateDeleteDcdnSpecificConfigRequest()
	deleteRequest.ConfigId = config.ConfigId
	deleteRequest.DomainName = parts[0]

	_, err = conn.DeleteDcdnSpecificConfig(deleteRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceAliyunDcdnDomainConfigRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).dcdnconn

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return diag.FromErr(err)
	}
	request := dcdn.CreateDescribeDcdnDomainConfigsRequest()
	request.DomainName = parts[0]
	request.FunctionNames = parts[1]

	res, err := conn.DescribeDcdnDomainConfigs(request)
	if err != nil {
		return diag.FromErr(err)
	}
	config := res.DomainConfigs.DomainConfig[0]

	var funArgs []map[string]string

	for _, args := range config.FunctionArgs.FunctionArg {
		if args.ArgName == "cert" || args.ArgName == "cert_id" || args.ArgName == "cert_name" || args.ArgName == "cert_type" || args.ArgName == "dkey" || args.ArgName == "pkey" || args.ArgName == "https" {
			continue
		}
		funArgs = append(funArgs, map[string]string{
			"arg_name":  args.ArgName,
			"arg_value": args.ArgValue,
		})
	}

	d.Set("domain_name", parts[0])
	d.Set("function_name", parts[1])
	d.Set("function_args", funArgs)

	return diags
}

func resourceAliyunDcdnDomainConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(Client).dcdnconn

	config := make([]map[string]interface{}, 1)
	functionArgs := d.Get("function_args").(*schema.Set).List()
	args := make([]map[string]interface{}, len(functionArgs))
	for key, value := range functionArgs {
		arg := value.(map[string]interface{})
		args[key] = map[string]interface{}{
			"argName":  arg["arg_name"],
			"argValue": arg["arg_value"],
		}
	}
	config[0] = map[string]interface{}{
		"functionArgs": args,
		"functionName": d.Get("function_name").(string),
	}
	functions, _ := json.Marshal(config)

	request := dcdn.CreateBatchSetDcdnDomainConfigsRequest()
	request.DomainNames = d.Get("domain_name").(string)
	request.Functions = string(functions)

	_, err := conn.BatchSetDcdnDomainConfigs(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s", request.DomainNames, d.Get("function_name").(string)))

	return resourceAliyunDcdnDomainConfigRead(ctx, d, m)
}
