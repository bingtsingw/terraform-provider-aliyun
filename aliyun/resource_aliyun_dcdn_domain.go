package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
)

func resourceAliyunDcdnDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAliyunDcdnDomainRead,
		CreateContext: resourceAliyunDcdnDomainCreate,
		UpdateContext: resourceAliyunDcdnDomainUpdate,
		DeleteContext: resourceAliyunDcdnDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"domestic", "global", "overseas"}, false),
				Default:      "domestic",
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      80,
							ValidateFunc: validation.IntInSlice([]int{443, 80}),
						},
						"priority": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "20",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ipaddr", "domain", "oss"}, false),
						},
						"weight": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "10",
						},
					},
				},
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceAliyunDcdnDomainDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).dcdnconn

	request := dcdn.CreateDeleteDcdnDomainRequest()
	request.DomainName = d.Id()

	_, err := conn.DeleteDcdnDomain(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceAliyunDcdnDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(Client).dcdnconn

	if d.HasChange("scope") {
		request := dcdn.CreateModifyDCdnDomainSchdmByPropertyRequest()
		request.DomainName = d.Id()
		request.Property = fmt.Sprintf(`{"coverage":"%s"}`, d.Get("scope").(string))
		_, err := conn.ModifyDCdnDomainSchdmByProperty(request)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	updateDomain := false
	request := dcdn.CreateUpdateDcdnDomainRequest()
	if d.HasChange("resource_group_id") {
		updateDomain = true
		request.ResourceGroupId = d.Get("resource_group_id").(string)
	}
	if d.HasChange("sources") {
		updateDomain = true
		sources, err := convertSourcesToString(d.Get("sources").(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		request.Sources = sources
	}

	if updateDomain {
		request.DomainName = d.Id()
		_, err := conn.UpdateDcdnDomain(request)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAliyunDcdnDomainRead(ctx, d, m)
}

func resourceAliyunDcdnDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(Client).dcdnconn

	domain := d.Get("domain_name").(string)

	request := dcdn.CreateAddDcdnDomainRequest()
	request.DomainName = domain

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("scope"); ok {
		request.Scope = v.(string)
	}

	sources, err := convertSourcesToString(d.Get("sources").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}
	request.Sources = sources

	_, err = conn.AddDcdnDomain(request)
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request := dcdn.CreateDescribeDcdnDomainDetailRequest()
		request.DomainName = domain

		res, err := conn.DescribeDcdnDomainDetail(request)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error creating dcdn: %s", err))
		}

		if res.DomainDetail.DomainStatus != "online" {
			return resource.RetryableError(fmt.Errorf("dcdn creation is processing"))
		}

		if res.DomainDetail.DomainStatus == "online" {
			return nil
		}

		return resource.NonRetryableError(fmt.Errorf("error creating dcdn: unkown state"))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain)

	return resourceAliyunDcdnDomainRead(ctx, d, m)
}

func resourceAliyunDcdnDomainRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).dcdnconn

	request := dcdn.CreateDescribeDcdnDomainDetailRequest()
	request.DomainName = d.Id()

	res, err := conn.DescribeDcdnDomainDetail(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("domain_name", d.Id())
	d.Set("resource_group_id", res.DomainDetail.ResourceGroupId)
	d.Set("scope", res.DomainDetail.Scope)
	sources := make([]map[string]interface{}, 0)
	for _, val := range res.DomainDetail.Sources.Source {
		sources = append(sources, map[string]interface{}{
			"content":  val.Content,
			"port":     val.Port,
			"priority": val.Priority,
			"type":     val.Type,
			"weight":   val.Weight,
		})
	}
	if err := d.Set("sources", sources); err != nil {
		return diag.FromErr(err)
	}
	d.Set("sources", sources)
	d.Set("cname", res.DomainDetail.Cname)

	return diags
}

func convertSourcesToString(v []interface{}) (string, error) {
	arrayMaps := make([]interface{}, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		arrayMaps[i] = map[string]interface{}{
			"Content":  item["content"],
			"Port":     item["port"],
			"Priority": item["priority"],
			"Type":     item["type"],
			"Weight":   item["weight"],
		}
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", err
	}
	return string(maps), nil
}
