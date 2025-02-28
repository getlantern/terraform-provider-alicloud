package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/getlantern/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudServiceMeshServiceMesh() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudServiceMeshServiceMeshCreate,
		Read:   resourceAlicloudServiceMeshServiceMeshRead,
		Update: resourceAlicloudServiceMeshServiceMeshUpdate,
		Delete: resourceAlicloudServiceMeshServiceMeshDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"load_balancer": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_server_public_eip": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"pilot_public_eip": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"api_server_loadbalancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pilot_public_loadbalancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				ForceNew: true,
			},
			"mesh_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_log": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"project": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"control_plane_log": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"project": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"audit": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"project": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"customized_zipkin": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"kiali": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"opa": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"limit_cpu": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"limit_memory": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"log_level": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"request_cpu": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"request_memory": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"outbound_traffic_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"ALLOW_ANY", "REGISTRY_ONLY"}, false),
						},
						"pilot": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http10_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"trace_sampling": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
								},
							},
						},
						"proxy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limit_cpu": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"limit_memory": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"request_cpu": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"request_memory": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"sidecar_injector": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_injection_policy_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"enable_namespaces_by_default": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"limit_cpu": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"limit_memory": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"request_cpu": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"request_memory": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"telemetry": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"tracing": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_locality_lb": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"network": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vswitche_list": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"service_mesh_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"edition": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Default", "Pro"}, false),
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"cluster_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "enterprise", "ultimate"}, false),
			},
			"cluster_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"extra_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cr_aggregation_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudServiceMeshServiceMeshCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateServiceMesh"
	request := make(map[string]interface{})
	conn, err := client.NewServicemeshClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("version"); ok {
		request["IstioVersion"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("service_mesh_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("edition"); ok {
		request["Edition"] = v
	}

	if v, ok := d.GetOk("network"); ok {
		for _, networkMap := range v.([]interface{}) {
			if networkArg, ok := networkMap.(map[string]interface{}); ok {
				if v, ok := networkArg["vpc_id"]; ok {
					request["VpcId"] = v
				}
				if v, ok := networkArg["vswitche_list"]; ok {
					request["VSwitches"] = convertListToJsonString(v.([]interface{}))
				}
			}
		}
	}

	if v, ok := d.GetOk("load_balancer"); ok {
		for _, loadBalancerMap := range v.([]interface{}) {
			if loadBalancerArg, ok := loadBalancerMap.(map[string]interface{}); ok {
				if v, ok := loadBalancerArg["api_server_public_eip"]; ok {
					request["ApiServerPublicEip"] = v
				}
				if v, ok := loadBalancerArg["pilot_public_eip"]; ok {
					request["PilotPublicEip"] = v
				}
			}
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		for _, meshConfigMap := range v.([]interface{}) {
			if meshConfigArg, ok := meshConfigMap.(map[string]interface{}); ok {

				if v, ok := meshConfigArg["customized_zipkin"]; ok {
					request["CustomizedZipkin"] = v
				}
				if v, ok := meshConfigArg["tracing"]; ok {
					request["Tracing"] = v
				}
				if v, ok := meshConfigArg["telemetry"]; ok {
					request["Telemetry"] = v
				}
				if v, ok := meshConfigArg["enable_locality_lb"]; ok {
					request["EnableLocalityLB"] = v
				}
				if pilot, ok := meshConfigArg["pilot"]; ok {
					for _, pilotMap := range pilot.([]interface{}) {
						if pilotArg, ok := pilotMap.(map[string]interface{}); ok {
							if v, ok := pilotArg["trace_sampling"]; ok {
								request["TraceSampling"] = v
							}
						}
					}
				}
				if AccessLog, ok := meshConfigArg["access_log"]; ok {
					for _, AccessLogMap := range AccessLog.([]interface{}) {
						if AccessLogArg, ok := AccessLogMap.(map[string]interface{}); ok {
							if v, ok := AccessLogArg["enabled"]; ok {
								request["AccessLogEnabled"] = v
							}
							if v, ok := AccessLogArg["project"]; ok {
								request["AccessLogProject"] = v
							}
						}
					}
				}
				if ControlPlaneLog, ok := meshConfigArg["control_plane_log"]; ok {
					for _, ControlPlaneLogMap := range ControlPlaneLog.([]interface{}) {
						if ControlPlaneLogArg, ok := ControlPlaneLogMap.(map[string]interface{}); ok {
							if v, ok := ControlPlaneLogArg["enabled"]; ok {
								request["ControlPlaneLogEnabled"] = v
							}
							if v, ok := ControlPlaneLogArg["project"]; ok {
								request["ControlPlaneLogProject"] = v
							}
						}
					}
				}
				if proxy, ok := meshConfigArg["proxy"]; ok {
					for _, proxyMap := range proxy.([]interface{}) {
						if proxyArg, ok := proxyMap.(map[string]interface{}); ok {
							if v, ok := proxyArg["request_memory"]; ok {
								request["ProxyRequestMemory"] = v
							}
							if v, ok := proxyArg["request_cpu"]; ok {
								request["ProxyRequestCPU"] = v
							}
							if v, ok := proxyArg["limit_memory"]; ok {
								request["ProxyLimitMemory"] = v
							}
							if v, ok := proxyArg["limit_cpu"]; ok {
								request["ProxyLimitCPU"] = v
							}
						}
					}
				}
				if opa, ok := meshConfigArg["opa"]; ok {
					for _, opaMap := range opa.([]interface{}) {
						if opaArg, ok := opaMap.(map[string]interface{}); ok {
							if v, ok := opaArg["enabled"]; ok {
								request["OpaEnabled"] = v
							}
							if v, ok := opaArg["log_level"]; ok {
								request["OPALogLevel"] = v
							}
							if v, ok := opaArg["request_cpu"]; ok {
								request["OPARequestCPU"] = v
							}
							if v, ok := opaArg["request_memory"]; ok {
								request["OPARequestMemory"] = v
							}
							if v, ok := opaArg["limit_cpu"]; ok {
								request["OPALimitCPU"] = v
							}
							if v, ok := opaArg["limit_memory"]; ok {
								request["OPALimitMemory"] = v
							}
						}

					}
				}
				if audit, ok := meshConfigArg["audit"]; ok {
					for _, auditMap := range audit.([]interface{}) {
						if auditArg, ok := auditMap.(map[string]interface{}); ok {
							if v, ok := auditArg["enabled"]; ok {
								request["EnableAudit"] = v
							}
							if v, ok := auditArg["project"]; ok {
								request["AuditProject"] = v
							}
							if v, ok := auditArg["enabled"]; ok {
								request["OpaEnabled"] = v
							}
						}
					}
				}
				if kiali, ok := meshConfigArg["kiali"]; ok {
					for _, kialiMap := range kiali.([]interface{}) {
						if kialiArg, ok := kialiMap.(map[string]interface{}); ok {
							if v, ok := kialiArg["enabled"]; ok {
								request["KialiEnabled"] = v
							}
						}
					}
				}
			}
		}
	}
	if v, ok := d.GetOk("cluster_spec"); ok {
		request["ClusterSpec"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ERR404", "InvalidActiveState.ACK"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_service_mesh_service_mesh", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServiceMeshId"]))
	servicemeshService := ServicemeshService{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, servicemeshService.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudServiceMeshServiceMeshUpdate(d, meta)
}
func resourceAlicloudServiceMeshServiceMeshRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	servicemeshService := ServicemeshService{client}
	object, err := servicemeshService.DescribeServiceMeshServiceMesh(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_service_mesh_service_mesh servicemeshService.DescribeServiceMeshServiceMesh Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("service_mesh_name", object["ServiceMeshInfo"].(map[string]interface{})["Name"])
	d.Set("edition", object["ServiceMeshInfo"].(map[string]interface{})["Profile"])

	if spec, ok := object["Spec"]; ok {
		if specArg, ok := spec.(map[string]interface{}); ok && len(specArg) > 0 {
			loadBalancerSli := make([]map[string]interface{}, 0)
			if loadBalancer, ok := specArg["LoadBalancer"]; ok {
				if loadBalancerArg, ok := loadBalancer.(map[string]interface{}); ok && len(loadBalancerArg) > 0 {
					loadBalancerMap := make(map[string]interface{})
					loadBalancerMap["pilot_public_eip"] = loadBalancerArg["PilotPublicEip"]
					loadBalancerMap["pilot_public_loadbalancer_id"] = loadBalancerArg["PilotPublicLoadbalancerId"]
					loadBalancerMap["api_server_loadbalancer_id"] = loadBalancerArg["ApiServerLoadbalancerId"]
					loadBalancerMap["api_server_public_eip"] = loadBalancerArg["ApiServerPublicEip"]
					loadBalancerSli = append(loadBalancerSli, loadBalancerMap)
				}
			}
			d.Set("load_balancer", loadBalancerSli)

			meshConfigSli := make([]map[string]interface{}, 0)
			if meshConfig, ok := specArg["MeshConfig"]; ok {
				meshConfigMap := make(map[string]interface{})
				if meshConfigArg, ok := meshConfig.(map[string]interface{}); ok && len(meshConfigArg) > 0 {
					accessLogSli := make([]map[string]interface{}, 0)
					if accessLog, ok := meshConfigArg["AccessLog"]; ok {
						if accessLogArg, ok := accessLog.(map[string]interface{}); ok && len(accessLogArg) > 0 {
							accessLogMap := make(map[string]interface{})
							accessLogMap["enabled"] = accessLogArg["Enabled"]
							accessLogMap["project"] = accessLogArg["Project"]
							accessLogSli = append(accessLogSli, accessLogMap)
						}
					}
					meshConfigMap["access_log"] = accessLogSli
					controlPlaneLogSli := make([]map[string]interface{}, 0)
					if controlPlaneLog, ok := meshConfigArg["ControlPlaneLogInfo"]; ok {
						if controlPlaneLogArg, ok := controlPlaneLog.(map[string]interface{}); ok && len(controlPlaneLogArg) > 0 {
							controlPlaneLogMap := make(map[string]interface{})
							controlPlaneLogMap["enabled"] = controlPlaneLogArg["Enabled"]
							controlPlaneLogMap["project"] = controlPlaneLogArg["Project"]
							controlPlaneLogSli = append(controlPlaneLogSli, controlPlaneLogMap)
						}
					}
					meshConfigMap["control_plane_log"] = controlPlaneLogSli
					auditSli := make([]map[string]interface{}, 0)
					if audit, ok := meshConfigArg["Audit"]; ok {
						if auditArg, ok := audit.(map[string]interface{}); ok && len(auditArg) > 0 {
							auditMap := make(map[string]interface{})
							auditMap["enabled"] = auditArg["Enabled"]
							auditMap["project"] = auditArg["Project"]
							auditSli = append(auditSli, auditMap)
							meshConfigMap["audit"] = auditSli
						}
					}

					meshConfigMap["customized_zipkin"] = meshConfigArg["CustomizedZipkin"]
					meshConfigMap["enable_locality_lb"] = meshConfigArg["EnableLocalityLB"]

					kialiSli := make([]map[string]interface{}, 0)
					if kiali, ok := meshConfigArg["Kiali"]; ok {
						if kialiArg, ok := kiali.(map[string]interface{}); ok && len(kialiArg) > 0 {
							kialiMap := make(map[string]interface{})
							kialiMap["enabled"] = kialiArg["Enabled"]
							kialiSli = append(kialiSli, kialiMap)
						}
					}
					meshConfigMap["kiali"] = kialiSli

					opaSli := make([]map[string]interface{}, 0)
					if opa, ok := meshConfigArg["OPA"]; ok {
						opaMap := make(map[string]interface{})
						if opaArg, ok := opa.(map[string]interface{}); ok && len(opaArg) > 0 {
							opaMap["enabled"] = opaArg["Enabled"]
							opaMap["limit_cpu"] = opaArg["LimitCPU"]
							opaMap["limit_memory"] = opaArg["LimitMemory"]
							opaMap["log_level"] = opaArg["LogLevel"]
							opaMap["request_cpu"] = opaArg["RequestCPU"]
							opaMap["request_memory"] = opaArg["RequestMemory"]
						}
						opaSli = append(opaSli, opaMap)
					}
					meshConfigMap["opa"] = opaSli
					meshConfigMap["outbound_traffic_policy"] = meshConfigArg["OutboundTrafficPolicy"]

					pilotSli := make([]map[string]interface{}, 0)
					if pilot := meshConfigArg["Pilot"]; ok {
						if pilotArg, ok := pilot.(map[string]interface{}); ok && len(pilotArg) > 0 {
							pilotMap := make(map[string]interface{})
							pilotMap["http10_enabled"] = pilotArg["Http10Enabled"]
							pilotMap["trace_sampling"] = pilotArg["TraceSampling"]
							pilotSli = append(pilotSli, pilotMap)
						}
					}
					meshConfigMap["pilot"] = pilotSli

					extraConfigSli := make([]map[string]interface{}, 0)
					if raw, ok := meshConfigArg["ExtraConfiguration"]; ok {
						if extraConfigArg, ok := raw.(map[string]interface{}); ok && len(extraConfigArg) > 0 {
							extraConfigMap := make(map[string]interface{})
							if v, ok := extraConfigArg["CRAggregationConfiguration"].(map[string]interface{}); ok {
								extraConfigMap["cr_aggregation_enabled"] = v["Enabled"]
								extraConfigSli = append(extraConfigSli, extraConfigMap)
							}
						}
					}
					d.Set("extra_configuration", extraConfigSli)

					proxySli := make([]map[string]interface{}, 0)
					if proxy, ok := meshConfigArg["Proxy"]; ok {
						if proxyArg, ok := proxy.(map[string]interface{}); ok && len(proxyArg) > 0 {
							proxyMap := make(map[string]interface{})
							proxyMap["limit_cpu"] = proxyArg["LimitCPU"]
							proxyMap["limit_memory"] = proxyArg["LimitMemory"]
							proxyMap["request_cpu"] = proxyArg["RequestCPU"]
							proxyMap["request_memory"] = proxyArg["RequestMemory"]
							proxySli = append(proxySli, proxyMap)
						}
					}
					meshConfigMap["proxy"] = proxySli

					sidecarInjectorSli := make([]map[string]interface{}, 0)
					if sidecarInjector, ok := meshConfigArg["SidecarInjector"]; ok {
						if sidecarInjectorArg, ok := sidecarInjector.(map[string]interface{}); ok && len(sidecarInjectorArg) > 0 {
							sidecarInjectorMap := make(map[string]interface{})
							sidecarInjectorMap["auto_injection_policy_enabled"] = sidecarInjectorArg["AutoInjectionPolicyEnabled"]
							sidecarInjectorMap["enable_namespaces_by_default"] = sidecarInjectorArg["EnableNamespacesByDefault"]
							sidecarInjectorMap["limit_cpu"] = sidecarInjectorArg["LimitCPU"]
							sidecarInjectorMap["limit_memory"] = sidecarInjectorArg["LimitMemory"]
							sidecarInjectorMap["request_cpu"] = sidecarInjectorArg["RequestCPU"]
							sidecarInjectorMap["request_memory"] = sidecarInjectorArg["RequestMemory"]
							sidecarInjectorSli = append(sidecarInjectorSli, sidecarInjectorMap)
						}
					}
					meshConfigMap["sidecar_injector"] = sidecarInjectorSli
					meshConfigMap["telemetry"] = meshConfigArg["Telemetry"]
					meshConfigMap["tracing"] = meshConfigArg["Tracing"]
					meshConfigSli = append(meshConfigSli, meshConfigMap)
				}
			}
			d.Set("mesh_config", meshConfigSli)

			networkSli := make([]map[string]interface{}, 0)
			if network, ok := specArg["Network"]; ok {
				if networkArg, ok := network.(map[string]interface{}); ok && len(networkArg) > 0 {
					networkMap := make(map[string]interface{})
					networkMap["vswitche_list"] = networkArg["VSwitches"]
					networkMap["vpc_id"] = networkArg["VpcId"]
					networkSli = append(networkSli, networkMap)
				}
			}
			d.Set("network", networkSli)
		}
	}
	d.Set("status", object["ServiceMeshInfo"].(map[string]interface{})["State"])
	d.Set("version", object["ServiceMeshInfo"].(map[string]interface{})["Version"])

	d.Set("cluster_spec", object["ClusterSpec"])
	clusters := make([]interface{}, 0)
	if v, ok := object["Clusters"].([]interface{}); ok {
		clusters = append(clusters, v...)
	}
	d.Set("cluster_ids", clusters)
	return nil
}
func resourceAlicloudServiceMeshServiceMeshUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	servicemeshService := ServicemeshService{client}
	var response map[string]interface{}
	d.Partial(true)
	conn, err := client.NewServicemeshClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	updateMeshFeatureReq := map[string]interface{}{
		"ServiceMeshId": d.Id(),
	}
	if d.HasChange("mesh_config") {
		update = true
	}

	if v, ok := d.GetOk("mesh_config"); ok {
		for _, meshConfigMap := range v.([]interface{}) {
			if meshConfigArg, ok := meshConfigMap.(map[string]interface{}); ok {

				if v, ok := meshConfigArg["customized_zipkin"]; ok {
					updateMeshFeatureReq["CustomizedZipkin"] = v
				}
				if v, ok := meshConfigArg["outbound_traffic_policy"]; ok {
					updateMeshFeatureReq["OutboundTrafficPolicy"] = v
				}
				if proxy, ok := meshConfigArg["proxy"]; ok {
					for _, proxyMap := range proxy.([]interface{}) {
						if proxyArg, ok := proxyMap.(map[string]interface{}); ok {
							if v, ok := proxyArg["request_memory"]; ok {
								updateMeshFeatureReq["ProxyRequestMemory"] = v
							}
							if v, ok := proxyArg["request_cpu"]; ok {
								updateMeshFeatureReq["ProxyRequestCPU"] = v
							}
							if v, ok := proxyArg["limit_memory"]; ok {
								updateMeshFeatureReq["ProxyLimitMemory"] = v
							}
							if v, ok := proxyArg["limit_cpu"]; ok {
								updateMeshFeatureReq["ProxyLimitCPU"] = v
							}
						}
					}
				}
				if AccessLog, ok := meshConfigArg["access_log"]; ok {
					for _, AccessLogMap := range AccessLog.([]interface{}) {
						if AccessLogArg, ok := AccessLogMap.(map[string]interface{}); ok {
							if v, ok := AccessLogArg["enabled"]; ok {
								updateMeshFeatureReq["AccessLogEnabled"] = v
							}
							if v, ok := AccessLogArg["project"]; ok {
								updateMeshFeatureReq["AccessLogProject"] = v
							}
						}
					}
				}
				if sidecarInjector, ok := meshConfigArg["sidecar_injector"]; ok && !d.IsNewResource() {
					for _, sidecarInjectorMap := range sidecarInjector.([]interface{}) {
						if sidecarInjectorArg, ok := sidecarInjectorMap.(map[string]interface{}); ok {
							if v, ok := sidecarInjectorArg["auto_injection_policy_enabled"]; ok {
								updateMeshFeatureReq["AutoInjectionPolicyEnabled"] = v
							}
							if v, ok := sidecarInjectorArg["enable_namespaces_by_default"]; ok {
								updateMeshFeatureReq["EnableNamespacesByDefault"] = v
							}
							if v, ok := sidecarInjectorArg["limit_cpu"]; ok {
								updateMeshFeatureReq["SidecarInjectorLimitCPU"] = v
							}
							if v, ok := sidecarInjectorArg["limit_memory"]; ok {
								updateMeshFeatureReq["SidecarInjectorLimitMemory"] = v
							}
							if v, ok := sidecarInjectorArg["request_cpu"]; ok {
								updateMeshFeatureReq["SidecarInjectorRequestCPU"] = v
							}
							if v, ok := sidecarInjectorArg["request_memory"]; ok {
								updateMeshFeatureReq["SidecarInjectorRequestMemory"] = v
							}
						}
					}
				}

				if AccessLog, ok := meshConfigArg["mesh_config"]; ok {
					for _, AccessLogMap := range AccessLog.([]interface{}) {
						if AccessLogArg, ok := AccessLogMap.(map[string]interface{}); ok {
							if v, ok := AccessLogArg["enabled"]; ok {
								updateMeshFeatureReq["AccessLogEnabled"] = v
							}
						}
					}
				}
				if v, ok := meshConfigArg["tracing"]; ok {
					updateMeshFeatureReq["Tracing"] = v
				}
				if v, ok := meshConfigArg["telemetry"]; ok {
					updateMeshFeatureReq["Telemetry"] = v
				}
				if pilot, ok := meshConfigArg["pilot"]; ok {
					for _, pilotMap := range pilot.([]interface{}) {
						if pilotArg, ok := pilotMap.(map[string]interface{}); ok {
							if v, ok := pilotArg["trace_sampling"]; ok {
								updateMeshFeatureReq["TraceSampling"] = v
							}
							if v, ok := pilotArg["http10_enabled"]; ok {
								updateMeshFeatureReq["Http10Enabled"] = v
							}
						}
					}
				}
				if opa, ok := meshConfigArg["opa"]; ok {
					for _, opaMap := range opa.([]interface{}) {
						if opaArg, ok := opaMap.(map[string]interface{}); ok {
							if v, ok := opaArg["enabled"]; ok {
								updateMeshFeatureReq["OpaEnabled"] = v
							}
							if v, ok := opaArg["log_level"]; ok {
								updateMeshFeatureReq["OPALogLevel"] = v
							}
							if v, ok := opaArg["request_cpu"]; ok {
								updateMeshFeatureReq["OPARequestCPU"] = v
							}
							if v, ok := opaArg["request_memory"]; ok {
								updateMeshFeatureReq["OPARequestMemory"] = v
							}
							if v, ok := opaArg["limit_cpu"]; ok {
								updateMeshFeatureReq["OPALimitCPU"] = v
							}
							if v, ok := opaArg["limit_memory"]; ok {
								updateMeshFeatureReq["OPALimitMemory"] = v
							}
						}
					}
				}
				if audit, ok := meshConfigArg["audit"]; ok {
					for _, auditMap := range audit.([]interface{}) {
						if auditArg, ok := auditMap.(map[string]interface{}); ok {
							if v, ok := auditArg["enabled"]; ok {
								updateMeshFeatureReq["EnableAudit"] = v
							}
							if v, ok := auditArg["project"]; ok {
								updateMeshFeatureReq["AuditProject"] = v
							}
						}
					}
				}
				if kiali, ok := meshConfigArg["kiali"]; ok {
					for _, kialiMap := range kiali.([]interface{}) {
						if kialiArg, ok := kialiMap.(map[string]interface{}); ok {
							if v, ok := kialiArg["enabled"]; ok {
								updateMeshFeatureReq["KialiEnabled"] = v
							}
						}
					}
				}
			}
		}
	}
	if !d.IsNewResource() && d.HasChange("cluster_spec") {
		update = true
		if v, ok := d.GetOk("cluster_spec"); ok {
			updateMeshFeatureReq["ClusterSpec"] = v
		}
	}

	if update {
		action := "UpdateMeshFeature"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-11"), StringPointer("AK"), nil, updateMeshFeatureReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateMeshFeatureReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, servicemeshService.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	update = false
	UpdateMeshCRAggregationReq := map[string]interface{}{
		"ServiceMeshId": d.Id(),
	}
	if d.HasChange("extra_configuration") {
		update = true
		if extraConfig, ok := d.GetOk("extra_configuration"); ok {
			for _, extraConfigMap := range extraConfig.([]interface{}) {
				if extraConfigArg, ok := extraConfigMap.(map[string]interface{}); ok {
					if v, ok := extraConfigArg["cr_aggregation_enabled"]; ok {
						UpdateMeshCRAggregationReq["Enabled"] = v
					}
				}
			}
		}
	}
	if update {
		action := "UpdateMeshCRAggregation"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-01-11"), StringPointer("AK"), UpdateMeshCRAggregationReq, nil, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateMeshFeatureReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, servicemeshService.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	update = false
	if !d.IsNewResource() && d.HasChange("version") {
		update = true
	}
	UpgradeEditionReq := map[string]interface{}{
		"ServiceMeshId": d.Id(),
	}
	if v, ok := d.GetOk("version"); ok {
		UpgradeEditionReq["ExpectedVersion"] = v
	}
	if update {
		action := "UpgradeMeshEditionPartially"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-01-11"), StringPointer("AK"), UpgradeEditionReq, nil, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateMeshFeatureReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, servicemeshService.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("cluster_ids") {
		oraw, nraw := d.GetChange("cluster_ids")
		removed := oraw.([]interface{})
		created := nraw.([]interface{})
		if len(removed) > 0 {
			for _, item := range removed {
				removeClusterReq := map[string]interface{}{
					"ServiceMeshId": d.Id(),
				}
				removeClusterReq["ClusterId"] = item
				action := "RemoveClusterFromServiceMesh"
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-11"), StringPointer("AK"), nil, removeClusterReq, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, removeClusterReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		if len(created) > 0 {
			for _, item := range created {
				addClusterReq := map[string]interface{}{
					"ServiceMeshId": d.Id(),
				}
				addClusterReq["ClusterId"] = item
				action := "AddClusterIntoServiceMesh"
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-11"), StringPointer("AK"), nil, addClusterReq, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, addClusterReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	d.Partial(false)

	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, servicemeshService.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudServiceMeshServiceMeshRead(d, meta)
}
func resourceAlicloudServiceMeshServiceMeshDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	servicemeshService := ServicemeshService{client}
	action := "DeleteServiceMesh"
	var response map[string]interface{}
	conn, err := client.NewServicemeshClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ServiceMeshId": d.Id(),
	}

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ErrorPermitted.ClustersNotEmpty", "RelatedResourceReused"}) || NeedRetry(err) {
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
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, servicemeshService.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
