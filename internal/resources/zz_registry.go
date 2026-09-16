// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// apiDialect is how this API shapes the values the descriptors read: what a wrapped value
// looks like, and what it spells "inherit". The rules that decided which fields these apply
// to ran at generation time.
var apiDialect = &dialect{
	property:         propertyFields{value: "value", raw: "rawvalue", parsed: "parsed", source: "source"},
	inherit:          "INHERIT",
	inheritedSources: []string{"INHERITED", "DEFAULT"},
}

// APIVersion is the API snapshot these resources were generated from. The provider asks the
// server for this version, so a box running something else is refused rather than guessed at.
const APIVersion = "v25.10.5"

// Resources lists every generated resource the provider serves.
var Resources = []func() resource.Resource{
	NewAcmeDnsAuthenticator,
	NewCertificate,
	NewCloudsyncCredentials,
	NewCloudsyncTask,
	NewCronJob,
	NewDataset,
	NewGeneralConfig,
	NewFTPConfig,
	NewGroup,
	NewInitScript,
	NewMailConfig,
	NewISCSIAuth,
	NewISCSIGlobalConfig,
	NewISCSIExtent,
	NewISCSIInitiator,
	NewISCSIPortal,
	NewISCSITarget,
	NewISCSITargetExtent,
	NewNetworkConfig,
	NewNFSConfig,
	NewNFSShare,
	NewReplication,
	NewSMBConfig,
	NewService,
	NewSMBShare,
	NewSnmpConfig,
	NewSnapshotTask,
	NewSSHConfig,
	NewStaticRoute,
	NewTunable,
	NewUPSConfig,
	NewUser,
	NewVM,
	NewVMDevice,
	NewZvol,
}

// DataSources lists every generated data source.
var DataSources = []func() datasource.DataSource{
	NewAcmeDnsAuthenticatorDataSource,
	NewCertificateDataSource,
	NewCloudsyncCredentialsDataSource,
	NewCloudsyncTaskDataSource,
	NewCronJobDataSource,
	NewDatasetDataSource,
	NewGeneralConfigDataSource,
	NewFTPConfigDataSource,
	NewGroupDataSource,
	NewInitScriptDataSource,
	NewMailConfigDataSource,
	NewISCSIAuthDataSource,
	NewISCSIGlobalConfigDataSource,
	NewISCSIExtentDataSource,
	NewISCSIInitiatorDataSource,
	NewISCSIPortalDataSource,
	NewISCSITargetDataSource,
	NewISCSITargetExtentDataSource,
	NewNetworkConfigDataSource,
	NewNFSConfigDataSource,
	NewNFSShareDataSource,
	NewReplicationDataSource,
	NewSMBConfigDataSource,
	NewServiceDataSource,
	NewSMBShareDataSource,
	NewSnmpConfigDataSource,
	NewSnapshotTaskDataSource,
	NewSSHConfigDataSource,
	NewStaticRouteDataSource,
	NewTunableDataSource,
	NewUPSConfigDataSource,
	NewUserDataSource,
	NewVMDataSource,
	NewVMDeviceDataSource,
	NewZvolDataSource,
}

// ListResources lists the resources that support terraform query.
var ListResources = []func() list.ListResource{
	NewAcmeDnsAuthenticatorList,
	NewCertificateList,
	NewCloudsyncCredentialsList,
	NewCloudsyncTaskList,
	NewCronJobList,
	NewDatasetList,
	NewGroupList,
	NewInitScriptList,
	NewISCSIAuthList,
	NewISCSIExtentList,
	NewISCSIInitiatorList,
	NewISCSIPortalList,
	NewISCSITargetList,
	NewISCSITargetExtentList,
	NewNFSShareList,
	NewReplicationList,
	NewServiceList,
	NewSMBShareList,
	NewSnapshotTaskList,
	NewStaticRouteList,
	NewTunableList,
	NewUserList,
	NewVMList,
	NewVMDeviceList,
	NewZvolList,
}

// descriptors maps a resource type to its model, for code in this package that needs the
// shape of a call rather than a resource: the fake server the tests run against. Keyed by
// type rather than namespace, because two resources can be variants of one namespace.
var descriptors = map[string]*model{
	"acme_dns_authenticator": modelAcmeDnsAuthenticator,
	"certificate":            modelCertificate,
	"cloudsync_credentials":  modelCloudsyncCredentials,
	"cloudsync_task":         modelCloudsyncTask,
	"cron_job":               modelCronJob,
	"dataset":                modelDataset,
	"general_config":         modelGeneralConfig,
	"ftp_config":             modelFTPConfig,
	"group":                  modelGroup,
	"init_script":            modelInitScript,
	"mail_config":            modelMailConfig,
	"iscsi_auth":             modelISCSIAuth,
	"iscsi_global_config":    modelISCSIGlobalConfig,
	"iscsi_extent":           modelISCSIExtent,
	"iscsi_initiator":        modelISCSIInitiator,
	"iscsi_portal":           modelISCSIPortal,
	"iscsi_target":           modelISCSITarget,
	"iscsi_target_extent":    modelISCSITargetExtent,
	"network_config":         modelNetworkConfig,
	"nfs_config":             modelNFSConfig,
	"nfs_share":              modelNFSShare,
	"replication":            modelReplication,
	"smb_config":             modelSMBConfig,
	"service":                modelService,
	"smb_share":              modelSMBShare,
	"snmp_config":            modelSnmpConfig,
	"snapshot_task":          modelSnapshotTask,
	"ssh_config":             modelSSHConfig,
	"static_route":           modelStaticRoute,
	"tunable":                modelTunable,
	"ups_config":             modelUPSConfig,
	"user":                   modelUser,
	"vm":                     modelVM,
	"vm_device":              modelVMDevice,
	"zvol":                   modelZvol,
}

// Actions lists every generated action.
var Actions = []func() action.Action{
	NewReplicationRunAction,
	NewCloudsyncSyncAction,
	NewSnapshotTaskRunAction,
	NewScrubRunAction,
	NewAppRedeployAction,
	NewAppStartAction,
	NewAppStopAction,
	NewServiceControlAction,
	NewUiRestartAction,
}
