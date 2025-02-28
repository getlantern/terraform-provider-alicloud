package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/getlantern/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBBackupPolicyCreate,
		Read:   resourceAlicloudPolarDBBackupPolicyRead,
		Update: resourceAlicloudPolarDBBackupPolicyUpdate,
		Delete: resourceAlicloudPolarDBBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"preferred_backup_period": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
				// terraform does not support ValidateFunc of TypeList attr
				// ValidateFunc: validateAllowedStringValue([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}),
				Optional: true,
				Computed: true,
			},

			"preferred_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},
			"backup_retention_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_retention_policy_on_cluster_deletion": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL", "LATEST", "NONE"}, false),
			},
		},
	}
}

func resourceAlicloudPolarDBBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("db_cluster_id").(string))

	return resourceAlicloudPolarDBBackupPolicyUpdate(d, meta)
}

func resourceAlicloudPolarDBBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polardbService := PolarDBService{client}
	object, err := polardbService.DescribeBackupPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_cluster_id", d.Id())
	d.Set("backup_retention_period", object.BackupRetentionPeriod)
	d.Set("preferred_backup_period", strings.Split(object.PreferredBackupPeriod, ","))
	d.Set("preferred_backup_time", object.PreferredBackupTime)
	d.Set("backup_retention_policy_on_cluster_deletion", object.BackupRetentionPolicyOnClusterDeletion)

	return nil
}

func resourceAlicloudPolarDBBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	polardbService := PolarDBService{client}

	if d.HasChange("preferred_backup_period") || d.HasChange("preferred_backup_time") ||
		d.HasChange("backup_retention_policy_on_cluster_deletion") {
		periodList := expandStringList(d.Get("preferred_backup_period").(*schema.Set).List())
		preferredBackupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
		preferredBackupTime := d.Get("preferred_backup_time").(string)
		var backupRetentionPolicyOnClusterDeletion string
		if v, ok := d.GetOk("backup_retention_policy_on_cluster_deletion"); ok && v.(string) != "" {
			backupRetentionPolicyOnClusterDeletion = v.(string)
		}
		// wait instance running before modifying
		if err := polardbService.WaitForCluster(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := polardbService.ModifyDBBackupPolicy(d.Id(), preferredBackupTime, preferredBackupPeriod, backupRetentionPolicyOnClusterDeletion); err != nil {
				if IsExpectedErrors(err, OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudPolarDBBackupPolicyRead(d, meta)
}

func resourceAlicloudPolarDBBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	//  Terraform can not destroy it..
	return nil
}
