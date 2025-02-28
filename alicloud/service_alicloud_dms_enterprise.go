package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/getlantern/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type DmsEnterpriseService struct {
	client *connectivity.AliyunClient
}

func (s *DmsEnterpriseService) DescribeDmsEnterpriseInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDmsenterpriseClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetInstance"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Host":     parts[0],
		"Port":     parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNoEnoughNumber"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DmsEnterpriseInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Instance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instance", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DmsEnterpriseService) DescribeDmsEnterpriseUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDmsenterpriseClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetUser"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Uid":      id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.User", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.User", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DmsEnterpriseService) DescribeDmsEnterpriseProxy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDmsenterpriseClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetProxy"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"ProxyId":  id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidParameterValid"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("DMSEnterprise:Proxy", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DmsEnterpriseService) DescribeDmsEnterpriseProxyAccess(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewDmsenterpriseClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"ProxyAccessId": id,
		"RegionId":      s.client.RegionId,
	}

	var response map[string]interface{}
	action := "GetProxyAccess"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameterValid"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ProxyAccess", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ProxyAccess", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DmsEnterpriseService) DmsEnterpriseProxyAccessStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDmsEnterpriseProxyAccess(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object[""]) == failState {
				return object, fmt.Sprint(object[""]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object[""])))
			}
		}
		return object, fmt.Sprint(object[""]), nil
	}
}

func (s *DmsEnterpriseService) InspectProxyAccessSecret(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewDmsenterpriseClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"ProxyAccessId": id,
		"RegionId":      s.client.RegionId,
	}

	var response map[string]interface{}
	action := "InspectProxyAccessSecret"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameterValid"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DmsEnterpriseService) DescribeDmsEnterpriseLogicDatabase(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewDmsenterpriseClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"DbId":     id,
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "GetLogicDatabase"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.LogicDatabase", response)
	success, _ := jsonpath.Get("$.Success", response)
	if err != nil && success.(bool) {
		return object, WrapErrorf(Error(GetNotFoundMessage("DmsEnterprise", id)), NotFoundWithResponse, response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DmsEnterpriseService) DmsEnterpriseLogicDatabaseStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDmsEnterpriseLogicDatabase(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object[""]) == failState {
				return object, fmt.Sprint(object[""]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object[""])))
			}
		}
		return object, fmt.Sprint(object[""]), nil
	}
}
