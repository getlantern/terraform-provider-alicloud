package alicloud

import (
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/getlantern/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCSEdgeKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCSEdgeKubernetesClustersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"connections": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_server_internet": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"api_server_intranet": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudCSEdgeKubernetesClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var allClusterTypes []cs.ClusterType

	var requestInfo *cs.Client
	invoker := NewInvoker()
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeClusters("")
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_edge_kubernetes_clusters", "DescribeClusters", DenverdinoAliyungo)
	}
	if debugOn() {
		addDebug("DescribeClusters", response, requestInfo)
	}
	allClusterTypes, _ = response.([]cs.ClusterType)

	var filteredClusterTypes []cs.ClusterType
	for _, v := range allClusterTypes {
		if v.ClusterType != cs.ClusterTypeManagedKubernetes {
			continue
		}
		if client.RegionId != string(v.RegionID) {
			continue
		}
		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r, err := regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
			if !r.MatchString(v.Name) {
				continue
			}
		}
		if ids, ok := d.GetOk("ids"); ok {
			var found bool
			for _, i := range expandStringList(ids.([]interface{})) {
				if v.ClusterID == i {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		filteredClusterTypes = append(filteredClusterTypes, v)
	}

	var filteredKubernetesCluster []cs.KubernetesCluster

	for _, v := range filteredClusterTypes {
		var kubernetesCluster cs.KubernetesCluster

		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return csClient.DescribeKubernetesCluster(v.ClusterID)
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_edge_kubernetes_clusters", "DescribeKubernetesCluster", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["Id"] = v.ClusterID
			addDebug("DescribeKubernetesCluster", response, requestInfo, requestMap)
		}
		kubernetesCluster = response.(cs.KubernetesCluster)

		if az, ok := d.GetOk("availability_zone"); ok && az != kubernetesCluster.ZoneId {
			continue
		}

		filteredKubernetesCluster = append(filteredKubernetesCluster, kubernetesCluster)
	}

	return csEdgeKubernetesClusterDescriptionAttributes(d, filteredKubernetesCluster, meta)
}

func csEdgeKubernetesClusterDescriptionAttributes(d *schema.ResourceData, clusterTypes []cs.KubernetesCluster, meta interface{}) error {
	detailEnabled := false
	if v, ok := d.GetOk("enable_details"); ok {
		detailEnabled = v.(bool)
	}

	var ids, names []string
	var s []map[string]interface{}
	for _, ct := range clusterTypes {
		mapping := map[string]interface{}{
			"id":   ct.ClusterID,
			"name": ct.Name,
		}

		if !detailEnabled {
			ids = append(ids, ct.ClusterID)
			names = append(names, ct.Name)
			s = append(s, mapping)
			continue
		}

		mapping["vpc_id"] = ct.VPCID
		mapping["security_group_id"] = ct.SecurityGroupID
		mapping["availability_zone"] = ct.ZoneId

		var workerNodes []map[string]interface{}

		invoker := NewInvoker()
		client := meta.(*connectivity.AliyunClient)
		pageNumber := 1
		for {
			var result []cs.KubernetesNodeType
			var pagination *cs.PaginationResult
			var requestInfo *cs.Client
			var response interface{}

			if err := invoker.Run(func() error {
				raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
					requestInfo = csClient
					nodes, paginationResult, err := csClient.GetKubernetesClusterNodes(ct.ClusterID, common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}, "")
					return []interface{}{nodes, paginationResult}, err
				})
				response = raw
				return err
			}); err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_edge_kubernetes_clusters", "GetKubernetesClusterNodes", DenverdinoAliyungo)
			}
			if debugOn() {
				requestMap := make(map[string]interface{})
				requestMap["Id"] = ct.ClusterID
				requestMap["Pagination"] = common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}
				addDebug("GetKubernetesClusterNodes", response, requestInfo, requestMap)
			}
			result, _ = response.([]interface{})[0].([]cs.KubernetesNodeType)
			pagination, _ = response.([]interface{})[1].(*cs.PaginationResult)

			if pageNumber == 1 && (len(result) == 0 || result[0].InstanceId == "") {
				err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					if err := invoker.Run(func() error {
						raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
							requestInfo = csClient
							nodes, _, err := csClient.GetKubernetesClusterNodes(ct.ClusterID, common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}, "")
							return nodes, err
						})
						response = raw
						return err
					}); err != nil {
						return resource.NonRetryableError(err)
					}
					if debugOn() {
						requestMap := make(map[string]interface{})
						requestMap["Id"] = ct.ClusterID
						requestMap["Pagination"] = common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}
						addDebug("GetKubernetesClusterNodes", response, requestInfo, requestMap)
					}
					tmp, _ := response.([]cs.KubernetesNodeType)
					if len(tmp) > 0 && tmp[0].InstanceId != "" {
						result = tmp
					}
					for _, stableState := range cs.NodeStableClusterState {
						// If cluster is in NodeStableClusteState, node list will not change
						if ct.State == stableState {
							return nil
						}
					}
					time.Sleep(5 * time.Second)
					return resource.RetryableError(Error("there is no any nodes in kubernetes cluster %s", d.Id()))
				})
				if err != nil {
					return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_managed_kubernetes_clusters", "GetKubernetesClusterNodes", DenverdinoAliyungo)
				}

			}

			for _, node := range result {
				subMapping := map[string]interface{}{
					"id":         node.InstanceId,
					"name":       node.InstanceName,
					"private_ip": node.IpAddress[0],
				}
				workerNodes = append(workerNodes, subMapping)
			}

			if len(result) < pagination.PageSize {
				break
			}
			pageNumber += 1
		}
		mapping["worker_nodes"] = workerNodes

		var response interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				endpoints, err := csClient.GetClusterEndpoints(ct.ClusterID)
				return endpoints, err
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_edge_kubernetes_clusters", "GetClusterEndpoints", DenverdinoAliyungo)
		}
		connection := make(map[string]string)
		if endpoints, ok := response.(cs.ClusterEndpoints); ok && endpoints.ApiServerEndpoint != "" {
			connection["api_server_internet"] = endpoints.ApiServerEndpoint
		}
		if endpoints, ok := response.(cs.ClusterEndpoints); ok && endpoints.IntranetApiServerEndpoint != "" {
			connection["api_server_intranet"] = endpoints.IntranetApiServerEndpoint
		}
		mapping["connections"] = connection

		request := vpc.CreateDescribeNatGatewaysRequest()
		request.VpcId = ct.VPCID

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, ct.VPCID, "DescribeNatGateways", AlibabaCloudSdkGoERROR)
		}
		addDebug("DescribeNatGateways", raw, request.RpcRequest, request)
		nat, _ := raw.(*vpc.DescribeNatGatewaysResponse)
		if nat != nil && len(nat.NatGateways.NatGateway) > 0 {
			mapping["nat_gateway_id"] = nat.NatGateways.NatGateway[0].NatGatewayId
		}

		ids = append(ids, ct.ClusterID)
		names = append(names, ct.Name)
		s = append(s, mapping)
	}

	d.Set("ids", ids)
	d.Set("names", names)
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
