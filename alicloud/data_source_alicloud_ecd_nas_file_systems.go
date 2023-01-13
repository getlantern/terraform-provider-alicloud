package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/getlantern/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcdNasFileSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdNasFileSystemsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Deleted", "Deleting", "Invalid", "Pending", "Running", "Stopped"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"systems": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_system_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_system_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metered_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_target_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_target_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nas_file_system_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"office_site_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"office_site_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_acl": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcdNasFileSystemsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeNASFileSystems"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("office_site_id"); ok {
		request["OfficeSiteId"] = v
	}
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var nasFileSystemNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nasFileSystemNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_nas_file_systems", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.FileSystems", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.FileSystems", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nasFileSystemNameRegex != nil && !nasFileSystemNameRegex.MatchString(fmt.Sprint(item["FileSystemName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["FileSystemId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["FileSystemStatus"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"capacity":             fmt.Sprint(object["Capacity"]),
			"create_time":          object["CreateTime"],
			"description":          object["Description"],
			"id":                   fmt.Sprint(object["FileSystemId"]),
			"file_system_id":       fmt.Sprint(object["FileSystemId"]),
			"file_system_type":     object["FileSystemType"],
			"metered_size":         fmt.Sprint(object["MeteredSize"]),
			"mount_target_domain":  object["MountTargetDomain"],
			"mount_target_status":  object["MountTargetStatus"],
			"nas_file_system_name": object["FileSystemName"],
			"office_site_id":       object["OfficeSiteId"],
			"office_site_name":     object["OfficeSiteName"],
			"status":               object["FileSystemStatus"],
			"storage_type":         object["StorageType"],
			"support_acl":          object["SupportAcl"],
			"zone_id":              object["ZoneId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["FileSystemName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("systems", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
