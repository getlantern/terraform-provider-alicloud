package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/getlantern/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type CloudfwService struct {
	client *connectivity.AliyunClient
}

func (s *CloudfwService) DescribeCloudFirewallControlPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudfwClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeControlPolicy"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AclUuid":     parts[0],
		"Direction":   parts[1],
		"CurrentPage": 1,
		"PageSize":    100,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if fmt.Sprint(response["Message"]) == "not buy user" {
			conn.Endpoint = String(connectivity.CloudFirewallOpenAPIEndpointControlPolicy)
			return resource.RetryableError(fmt.Errorf("%s", response))

		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Policys", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Policys", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudFirewall", id)), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CloudfwService) DescribeCloudFirewallAddressBook(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudfwClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeAddressBook"
	request := map[string]interface{}{
		"PageSize":    PageSizeLarge,
		"CurrentPage": 1,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
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
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Acls", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Acls", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudFirewall", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["GroupUuid"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudFirewall", id)), NotFoundWithResponse, response)
	}
	return
}
func (s *CloudfwService) DescribeCloudFirewallInstanceMember(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCloudfirewallClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"MemberUid": id,
	}

	var response map[string]interface{}
	action := "DescribeInstanceMembers"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidResource.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Members", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Members", response)
	}
	members := v.([]interface{})
	if len(members) < 1 {
		return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return members[0].(map[string]interface{}), nil
}
func (s *CloudfwService) CloudFirewallInstanceMemberStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallInstanceMember(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["MemberStatus"]) == failState {
				return object, fmt.Sprint(object["MemberStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["MemberStatus"])))
			}
		}
		return object, fmt.Sprint(object["MemberStatus"]), nil
	}
}

func (s *CloudfwService) DescribeCloudFirewallVpcFirewallCen(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCloudfwClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"VpcFirewallId": id,
	}

	var response map[string]interface{}
	action := "DescribeVpcFirewallCenDetail"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ErrorVpcFirewallExist"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VpcFirewallCen", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	status, err := jsonpath.Get("$.FirewallSwitchStatus", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	if status.(string) == "notconfigured" {
		return object, WrapErrorf(Error(GetNotFoundMessage("VpcFirewallCen", id)), NotFoundWithResponse, response)
	}
	return v.(map[string]interface{}), nil
}

func (s *CloudfwService) DescribeVpcFirewallCenList(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCloudfwClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"VpcFirewallId": id,
	}

	var response map[string]interface{}
	action := "DescribeVpcFirewallCenList"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.VpcFirewalls[0]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	totalCount, _ := jsonpath.Get("$.TotalCount", response)
	total, _ := totalCount.(json.Number).Int64()
	if err != nil && total == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VpcFirewallCen", id)), NotFoundWithResponse, response)
	}
	return v.(map[string]interface{}), nil
}

func (s *CloudfwService) CloudFirewallVpcFirewallCenStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallVpcFirewallCen(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["FirewallSwitchStatus"]) == failState {
				return object, fmt.Sprint(object["FirewallSwitchStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["FirewallSwitchStatus"])))
			}
		}
		return object, fmt.Sprint(object["FirewallSwitchStatus"]), nil
	}
}

func (s *CloudfwService) DescribeCloudFirewallVpcFirewall(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCloudfwClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"VpcFirewallId": id,
	}

	var response map[string]interface{}
	action := "DescribeVpcFirewallDetail"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ErrorFirewallNotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	switchStatus, _ := jsonpath.Get("$.FirewallSwitchStatus", response)
	switchStatusString := switchStatus.(string)
	if switchStatusString == "notconfigured" {
		err = WrapErrorf(Error(GetNotFoundMessage("Cloudfirewall", id)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *CloudfwService) DescribeVpcFirewallList(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCloudfwClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"VpcFirewallId": id,
	}

	var response map[string]interface{}
	action := "DescribeVpcFirewallList"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ErrorFirewallNotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	totalCount, _ := jsonpath.Get("$.TotalCount", response)
	total, _ := totalCount.(json.Number).Int64()
	if err != nil && total == 0 {
		err = WrapErrorf(Error(GetNotFoundMessage("CloudFirewall", id)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	v, err := jsonpath.Get("$.VpcFirewalls[0]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *CloudfwService) CloudFirewallVpcFirewallStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallVpcFirewall(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["FirewallSwitchStatus"]) == failState {
				return object, fmt.Sprint(object["FirewallSwitchStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["FirewallSwitchStatus"])))
			}
		}
		return object, fmt.Sprint(object["FirewallSwitchStatus"]), nil
	}
}

func (s *CloudfwService) DescribeCloudFirewallVpcFirewallControlPolicy(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCloudfirewallClient()
	if err != nil {
		return object, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"VpcFirewallId": parts[0],
		"AclUuid":       parts[1],
		"CurrentPage":   1,
		"PageSize":      PageSizeLarge,
	}

	var response map[string]interface{}
	action := "DescribeVpcFirewallControlPolicy"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.Policys", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Policys", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudFirewall", id)), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}
