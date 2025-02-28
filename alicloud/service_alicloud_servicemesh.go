package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/getlantern/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ServicemeshService struct {
	client *connectivity.AliyunClient
}

func (s *ServicemeshService) DescribeServiceMeshServiceMesh(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServicemeshClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeServiceMeshDetail"
	request := map[string]interface{}{
		"ServiceMeshId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-11"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ServiceMesh.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ServiceMesh:ServiceMesh", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ServiceMesh", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ServiceMesh", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ServicemeshService) DescribeServiceMeshDetail(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServicemeshClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeServiceMeshDetail"
	request := map[string]interface{}{
		"ServiceMeshId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-11"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.ServiceMesh", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ServiceMesh", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ServicemeshService) ServiceMeshServiceMeshStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeServiceMeshServiceMesh(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["ServiceMeshInfo"].(map[string]interface{})["State"]) == failState {
				return object, fmt.Sprint(object["ServiceMeshInfo"].(map[string]interface{})["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["ServiceMeshInfo.State"])))
			}
		}
		return object, fmt.Sprint(object["ServiceMeshInfo"].(map[string]interface{})["State"]), nil
	}
}

func (s *ServicemeshService) DescribeUserPermissions(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServicemeshClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeUserPermissions"
	request := map[string]interface{}{
		"SubAccountUserId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-11"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"StatusForbidden"}) {
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
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil

}

func (s *ServicemeshService) DescribeServiceMeshExtensionProvider(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServicemeshClient()
	if err != nil {
		return object, WrapError(err)
	}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return object, WrapError(err)
	}
	action := "DescribeExtensionProvider"
	request := map[string]interface{}{
		"ServiceMeshId": parts[0],
		"Type":          parts[1],
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-11"), StringPointer("AK"), nil, request, &runtime)
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

	resp, err := jsonpath.Get("$.ExtensionProviderArray", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ExtensionProviderArray", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ServiceMesh:ExtensionProvider", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Name"]) == parts[2] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("ServiceMesh:ExtensionProvider", id)), NotFoundWithResponse, response)
	}

	return object, nil
}
