package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/getlantern/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudGaBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaBandwidthPackageAttachmentCreate,
		Read:   resourceAlicloudGaBandwidthPackageAttachmentRead,
		Update: resourceAlicloudGaBandwidthPackageAttachmentUpdate,
		Delete: resourceAlicloudGaBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"accelerators": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "BandwidthPackageAddAccelerator"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	request["AcceleratorId"] = d.Get("accelerator_id")
	request["BandwidthPackageId"] = d.Get("bandwidth_package_id")
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.BandwidthPackage", "StateError.Accelerator", "GreaterThanGa.IpSetBandwidth", "BandwidthIllegal.BandwidthPackage", "NotExist.BasicBandwidthPackage"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_bandwidth_package_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AcceleratorId"], ":", response["BandwidthPackageId"]))
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaBandwidthPackageAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaBandwidthPackageAttachmentRead(d, meta)
}

func resourceAlicloudGaBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	if !strings.Contains(d.Id(), ":") {
		d.SetId(fmt.Sprint(d.Get("accelerator_id"), ":", d.Id()))
	}
	object, err := gaService.DescribeGaBandwidthPackageAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_bandwidth_package_attachment gaService.DescribeGaBandwidthPackageAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("bandwidth_package_id", parts[1])
	d.Set("accelerator_id", parts[0])
	d.Set("accelerators", []string{parts[0]})
	d.Set("status", object["State"])

	return nil
}

func resourceAlicloudGaBandwidthPackageAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	update := false
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"AcceleratorId": parts[0],
	}

	if d.HasChange("bandwidth_package_id") {
		update = true
	}
	oldValue, newValue := d.GetChange("bandwidth_package_id")
	oldBandwidthPackageId := oldValue.(string)
	newBandwidthPackageId := newValue.(string)
	request["BandwidthPackageId"] = newBandwidthPackageId
	request["TargetBandwidthPackageId"] = oldBandwidthPackageId

	if update {
		action := "ReplaceBandwidthPackage"
		conn, err := client.NewGaplusClient()
		if err != nil {
			return WrapError(err)
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.BandwidthPackage", "StateError.Accelerator", "GreaterThanGa.IpSetBandwidth", "BandwidthIllegal.BandwidthPackage", "NotExist.BasicBandwidthPackage"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", parts[0], request["BandwidthPackageId"]))

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaBandwidthPackageAttachmentStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	return resourceAlicloudGaBandwidthPackageAttachmentRead(d, meta)
}

func resourceAlicloudGaBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	if !strings.Contains(d.Id(), ":") {
		d.SetId(fmt.Sprint(d.Get("accelerator_id"), ":", d.Id()))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "BandwidthPackageRemoveAccelerator"
	var response map[string]interface{}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"BandwidthPackageId": parts[1],
	}

	request["AcceleratorId"] = parts[0]
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.BandwidthPackage", "StateError.Accelerator", "BindExist.CrossDomain", "Exist.EndpointGroup", "Exist.IpSet", "BandwidthPackageCannotUnbind.HasCrossRegion", "BandwidthPackageCannotUnbind.IpSet", "BandwidthPackageCannotUnbind.EndpointGroup"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.BandwidthPackage", "Exist.EndpointGroup"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaBandwidthPackageAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
