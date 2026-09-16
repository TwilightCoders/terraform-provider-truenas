// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"math/big"
)

var modelCloudsyncCredentials = &model{
	typeName:     "cloudsync_credentials",
	namespace:    "cloudsync.credentials",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "cloudsync.credentials.create",
	updateMethod: "cloudsync.credentials.update",
	getMethod:    "cloudsync.credentials.get_instance",
	deleteMethod: "cloudsync.credentials.delete",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "name", api: "name", path: "name",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Human-readable name for the cloud credential.",
		},
		{
			name: "storage", api: "provider", path: "provider",
			kind: kindUnion, role: roleRequired,
			readable: true, updatable: true,
			description:   "Cloud provider configuration including type and authentication details.",
			discriminator: "type",
			children: []*node{
				{
					name: "azureblob", api: "AZUREBLOB", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `AZUREBLOB`.",
					children: []*node{
						{
							name: "account", api: "account", path: "provider.account",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "Azure Blob Storage account name for authentication.",
						},
						{
							name: "key", api: "key", path: "provider.key",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "Azure Blob Storage access key for authentication.",
						},
						{
							name: "endpoint", api: "endpoint", path: "provider.endpoint",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Custom Azure Blob Storage endpoint URL. Empty string for default endpoints. Defaults to `\"\"`.",
						},
					},
				},
				{
					name: "b2", api: "B2", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `B2`.",
					children: []*node{
						{
							name: "account", api: "account", path: "provider.account",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "Backblaze B2 account ID for authentication.",
						},
						{
							name: "key", api: "key", path: "provider.key",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "Backblaze B2 application key for authentication.",
						},
					},
				},
				{
					name: "box", api: "BOX", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `BOX`.",
					children: []*node{
						{
							name: "client_id", api: "client_id", path: "provider.client_id",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Box OAuth application client ID. Defaults to `\"\"`.",
						},
						{
							name: "client_secret", api: "client_secret", path: "provider.client_secret",
							kind: kindString, role: roleOptional,
							sensitive: true, writeOnly: true, readable: true,
							description: "Box OAuth application client secret.",
						},
						{
							name: "token", api: "token", path: "provider.token",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "Box OAuth access token for API authentication.",
						},
					},
				},
				{
					name: "dropbox", api: "DROPBOX", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `DROPBOX`.",
					children: []*node{
						{
							name: "client_id", api: "client_id", path: "provider.client_id",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Dropbox OAuth application client ID. Defaults to `\"\"`.",
						},
						{
							name: "client_secret", api: "client_secret", path: "provider.client_secret",
							kind: kindString, role: roleOptional,
							sensitive: true, writeOnly: true, readable: true,
							description: "Dropbox OAuth application client secret.",
						},
						{
							name: "token", api: "token", path: "provider.token",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "Dropbox OAuth access token for API authentication.",
						},
					},
				},
				{
					name: "ftp", api: "FTP", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `FTP`.",
					children: []*node{
						{
							name: "host", api: "host", path: "provider.host",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "FTP server hostname or IP address.",
						},
						{
							name: "port", api: "port", path: "provider.port",
							kind: kindInt, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: 21,
							description: "FTP server port number. Defaults to `21`.",
						},
						{
							name: "user", api: "user", path: "provider.user",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "FTP username for authentication.",
						},
						{
							name: "pass", api: "pass", path: "provider.pass",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "FTP password for authentication.",
						},
					},
				},
				{
					name: "google_cloud_storage", api: "GOOGLE_CLOUD_STORAGE", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `GOOGLE_CLOUD_STORAGE`.",
					children: []*node{
						{
							name: "service_account_credentials", api: "service_account_credentials", path: "provider.service_account_credentials",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "JSON service account credentials for Google Cloud Storage authentication.",
						},
					},
				},
				{
					name: "google_drive", api: "GOOGLE_DRIVE", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `GOOGLE_DRIVE`.",
					children: []*node{
						{
							name: "client_id", api: "client_id", path: "provider.client_id",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "OAuth client ID for Google Drive API access. Defaults to `\"\"`.",
						},
						{
							name: "client_secret", api: "client_secret", path: "provider.client_secret",
							kind: kindString, role: roleOptional,
							sensitive: true, writeOnly: true, readable: true,
							description: "OAuth client secret for Google Drive API access.",
						},
						{
							name: "token", api: "token", path: "provider.token",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "OAuth access token for Google Drive authentication.",
						},
						{
							name: "team_drive", api: "team_drive", path: "provider.team_drive",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Google Drive team drive ID or empty string for personal drive. Defaults to `\"\"`.",
						},
					},
				},
				{
					name: "google_photos", api: "GOOGLE_PHOTOS", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `GOOGLE_PHOTOS`.",
					children: []*node{
						{
							name: "client_id", api: "client_id", path: "provider.client_id",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "OAuth client ID for Google Photos API access. Defaults to `\"\"`.",
						},
						{
							name: "client_secret", api: "client_secret", path: "provider.client_secret",
							kind: kindString, role: roleOptional,
							sensitive: true, writeOnly: true, readable: true,
							description: "OAuth client secret for Google Photos API access.",
						},
						{
							name: "token", api: "token", path: "provider.token",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "OAuth access token for Google Photos authentication.",
						},
					},
				},
				{
					name: "http", api: "HTTP", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `HTTP`.",
					children: []*node{
						{
							name: "url", api: "url", path: "provider.url",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "HTTP URL for file access.",
						},
					},
				},
				{
					name: "hubic", api: "HUBIC", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `HUBIC`.",
					children: []*node{
						{
							name: "token", api: "token", path: "provider.token",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "OAuth access token for Hubic authentication.",
						},
					},
				},
				{
					name: "mega", api: "MEGA", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `MEGA`.",
					children: []*node{
						{
							name: "user", api: "user", path: "provider.user",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "MEGA username for authentication.",
						},
						{
							name: "pass", api: "pass", path: "provider.pass",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "MEGA password for authentication.",
						},
					},
				},
				{
					name: "onedrive", api: "ONEDRIVE", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `ONEDRIVE`.",
					children: []*node{
						{
							name: "client_id", api: "client_id", path: "provider.client_id",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "OAuth client ID for OneDrive API access. Defaults to `\"\"`.",
						},
						{
							name: "client_secret", api: "client_secret", path: "provider.client_secret",
							kind: kindString, role: roleOptional,
							sensitive: true, writeOnly: true, readable: true,
							description: "OAuth client secret for OneDrive API access.",
						},
						{
							name: "token", api: "token", path: "provider.token",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "OAuth access token for OneDrive authentication.",
						},
						{
							name: "drive_type", api: "drive_type", path: "provider.drive_type",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "Type of OneDrive to access.",
						},
						{
							name: "drive_id", api: "drive_id", path: "provider.drive_id",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "OneDrive drive identifier.",
						},
					},
				},
				{
					name: "pcloud", api: "PCLOUD", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `PCLOUD`.",
					children: []*node{
						{
							name: "client_id", api: "client_id", path: "provider.client_id",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "OAuth client ID for pCloud API access. Defaults to `\"\"`.",
						},
						{
							name: "client_secret", api: "client_secret", path: "provider.client_secret",
							kind: kindString, role: roleOptional,
							sensitive: true, writeOnly: true, readable: true,
							description: "OAuth client secret for pCloud API access.",
						},
						{
							name: "token", api: "token", path: "provider.token",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "OAuth access token for pCloud authentication.",
						},
						{
							name: "hostname", api: "hostname", path: "provider.hostname",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "pCloud hostname or empty string for default. Defaults to `\"\"`.",
						},
					},
				},
				{
					name: "s3", api: "S3", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `S3`.",
					children: []*node{
						{
							name: "access_key_id", api: "access_key_id", path: "provider.access_key_id",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "S3 access key ID for authentication.",
						},
						{
							name: "secret_access_key", api: "secret_access_key", path: "provider.secret_access_key",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "S3 secret access key for authentication.",
						},
						{
							name: "endpoint", api: "endpoint", path: "provider.endpoint",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "S3-compatible endpoint URL or empty string for AWS S3. Defaults to `\"\"`.",
						},
						{
							name: "region", api: "region", path: "provider.region",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "S3 region or empty string for default. Defaults to `\"\"`.",
						},
						{
							name: "skip_region", api: "skip_region", path: "provider.skip_region",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "Whether to skip region validation. Defaults to `false`.",
						},
						{
							name: "signatures_v2", api: "signatures_v2", path: "provider.signatures_v2",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "Whether to use AWS Signature Version 2. Defaults to `false`.",
						},
						{
							name: "max_upload_parts", api: "max_upload_parts", path: "provider.max_upload_parts",
							kind: kindInt, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: 10000,
							description: "Maximum number of parts for multipart uploads. Defaults to `10000`.",
						},
					},
				},
				{
					name: "sftp", api: "SFTP", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `SFTP`.",
					children: []*node{
						{
							name: "host", api: "host", path: "provider.host",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "SFTP server hostname or IP address.",
						},
						{
							name: "port", api: "port", path: "provider.port",
							kind: kindInt, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: 22,
							description: "SFTP server port number. Defaults to `22`.",
						},
						{
							name: "user", api: "user", path: "provider.user",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "SFTP username for authentication.",
						},
						{
							name: "pass", api: "pass", path: "provider.pass",
							kind: kindString, role: roleOptional,
							nullable: true, sensitive: true, writeOnly: true, readable: true,
							description: "SFTP password for authentication or `null` for key-based auth.",
						},
						{
							name: "private_key", api: "private_key", path: "provider.private_key",
							kind: kindInt, role: roleOptionalComputed,
							nullable: true, readable: true,
							description: "SSH private key ID for authentication or `null` for password auth.",
						},
					},
				},
				{
					name: "storj_ix", api: "STORJ_IX", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `STORJ_IX`.",
					children: []*node{
						{
							name: "access_key_id", api: "access_key_id", path: "provider.access_key_id",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "Storj S3-compatible access key ID for authentication.",
						},
						{
							name: "secret_access_key", api: "secret_access_key", path: "provider.secret_access_key",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "Storj S3-compatible secret access key for authentication.",
						},
						{
							name: "endpoint", api: "endpoint", path: "provider.endpoint",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "https://gateway.storjshare.io/",
							description: "Storj gateway endpoint URL for S3-compatible access. Defaults to `\"https://gateway.storjshare.io/\"`.",
						},
					},
				},
				{
					name: "openstack_swift", api: "OPENSTACK_SWIFT", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `OPENSTACK_SWIFT`.",
					children: []*node{
						{
							name: "user", api: "user", path: "provider.user",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "Swift username for authentication.",
						},
						{
							name: "key", api: "key", path: "provider.key",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "Swift password or API key for authentication.",
						},
						{
							name: "auth", api: "auth", path: "provider.auth",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "Swift authentication URL endpoint.",
						},
						{
							name: "user_id", api: "user_id", path: "provider.user_id",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Swift user ID for authentication. Defaults to `\"\"`.",
						},
						{
							name: "domain", api: "domain", path: "provider.domain",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Swift domain name for authentication. Defaults to `\"\"`.",
						},
						{
							name: "tenant", api: "tenant", path: "provider.tenant",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Swift tenant name for multi-tenancy. Defaults to `\"\"`.",
						},
						{
							name: "tenant_id", api: "tenant_id", path: "provider.tenant_id",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Swift tenant ID for multi-tenancy. Defaults to `\"\"`.",
						},
						{
							name: "tenant_domain", api: "tenant_domain", path: "provider.tenant_domain",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Swift tenant domain name. Defaults to `\"\"`.",
						},
						{
							name: "region", api: "region", path: "provider.region",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Swift region name for geographic distribution. Defaults to `\"\"`.",
						},
						{
							name: "storage_url", api: "storage_url", path: "provider.storage_url",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Swift storage URL endpoint. Defaults to `\"\"`.",
						},
						{
							name: "auth_token", api: "auth_token", path: "provider.auth_token",
							kind: kindString, role: roleOptional,
							sensitive: true, writeOnly: true, readable: true,
							description: "Swift authentication token for pre-authenticated access.",
						},
						{
							name: "application_credential_id", api: "application_credential_id", path: "provider.application_credential_id",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Swift application credential ID for authentication. Defaults to `\"\"`.",
						},
						{
							name: "application_credential_name", api: "application_credential_name", path: "provider.application_credential_name",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Swift application credential name for authentication. Defaults to `\"\"`.",
						},
						{
							name: "application_credential_secret", api: "application_credential_secret", path: "provider.application_credential_secret",
							kind: kindString, role: roleOptional,
							sensitive: true, writeOnly: true, readable: true,
							description: "Swift application credential secret for authentication.",
						},
						{
							name: "auth_version", api: "auth_version", path: "provider.auth_version",
							kind: kindInt, role: roleRequired,
							nullable: true, readable: true,
							description: "Swift authentication API version.\n\n* `0`: Legacy auth\n* `1`: TempAuth\n* `2`: Keystone v2.0\n* `3`: Keystone v3\n* `null`: Auto-detect",
						},
						{
							name: "endpoint_type", api: "endpoint_type", path: "provider.endpoint_type",
							kind: kindString, role: roleRequired,
							nullable: true, readable: true,
							description: "Swift endpoint type to use.\n\n* `public`: Public endpoint (default)\n* `internal`: Internal network endpoint\n* `admin`: Administrative endpoint\n* `null`: Use default",
						},
					},
				},
				{
					name: "webdav", api: "WEBDAV", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `WEBDAV`.",
					children: []*node{
						{
							name: "url", api: "url", path: "provider.url",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "WebDAV server URL endpoint.",
						},
						{
							name: "vendor", api: "vendor", path: "provider.vendor",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "WebDAV server vendor type for compatibility optimizations.\n\n* `NEXTCLOUD`: Nextcloud server\n* `OWNCLOUD`: ownCloud server\n* `SHAREPOINT`: Microsoft SharePoint\n* `OTHER`: Generic WebDAV server",
						},
						{
							name: "user", api: "user", path: "provider.user",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "WebDAV username for authentication.",
						},
						{
							name: "pass", api: "pass", path: "provider.pass",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "WebDAV password for authentication.",
						},
					},
				},
				{
					name: "yandex", api: "YANDEX", path: "provider",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `type` is `YANDEX`.",
					children: []*node{
						{
							name: "client_id", api: "client_id", path: "provider.client_id",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Yandex OAuth application client ID. Defaults to `\"\"`.",
						},
						{
							name: "client_secret", api: "client_secret", path: "provider.client_secret",
							kind: kindString, role: roleOptional,
							sensitive: true, writeOnly: true, readable: true,
							description: "Yandex OAuth application client secret.",
						},
						{
							name: "token", api: "token", path: "provider.token",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "Yandex OAuth access token for API authentication.",
						},
					},
				},
			},
		},
		{
			name: "storage_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `storage` again. Terraform never stores `storage`.",
		},
	},
}

// NewCloudsyncCredentials returns the cloudsync_credentials resource.
func NewCloudsyncCredentials() resource.Resource {
	return &cloudsyncCredentialsResource{crudResource{model: modelCloudsyncCredentials}}
}

type cloudsyncCredentialsResource struct{ crudResource }

// NewCloudsyncCredentialsDataSource returns the cloudsync_credentials data source.
func NewCloudsyncCredentialsDataSource() datasource.DataSource {
	return &dataSource{model: modelCloudsyncCredentials, attrs: dataAttrs(modelCloudsyncCredentials)}
}

// NewCloudsyncCredentialsList returns the cloudsync_credentials list resource, for terraform query.
func NewCloudsyncCredentialsList() list.ListResource {
	return &listResource{crudResource{model: modelCloudsyncCredentials}}
}

func (r *cloudsyncCredentialsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages credentials for cloud sync tasks. Set exactly one storage type under `storage`. Secret fields are write-only: they never enter Terraform state. Bump `storage_wo_version` to resend them.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Human-readable name for the cloud credential.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"storage": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Cloud provider configuration including type and authentication details.",
				Attributes: map[string]schema.Attribute{
					"azureblob": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `AZUREBLOB`.",
						Attributes: map[string]schema.Attribute{
							"account": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Azure Blob Storage account name for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"key": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Azure Blob Storage access key for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"endpoint": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Custom Azure Blob Storage endpoint URL. Empty string for default endpoints. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
						},
					},
					"b2": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `B2`.",
						Attributes: map[string]schema.Attribute{
							"account": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Backblaze B2 account ID for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"key": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Backblaze B2 application key for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"box": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `BOX`.",
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Box OAuth application client ID. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"client_secret": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Box OAuth application client secret.",
							},
							"token": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Box OAuth access token for API authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"dropbox": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `DROPBOX`.",
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Dropbox OAuth application client ID. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"client_secret": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Dropbox OAuth application client secret.",
							},
							"token": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Dropbox OAuth access token for API authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"ftp": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `FTP`.",
						Attributes: map[string]schema.Attribute{
							"host": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "FTP server hostname or IP address.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"port": schema.NumberAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "FTP server port number. Defaults to `21`.",
								Validators:          []validator.Number{numberIsInteger{}},
								Default:             numberdefault.StaticBigFloat(big.NewFloat(21)),
							},
							"user": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "FTP username for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"pass": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "FTP password for authentication.",
							},
						},
					},
					"google_cloud_storage": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `GOOGLE_CLOUD_STORAGE`.",
						Attributes: map[string]schema.Attribute{
							"service_account_credentials": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "JSON service account credentials for Google Cloud Storage authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"google_drive": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `GOOGLE_DRIVE`.",
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "OAuth client ID for Google Drive API access. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"client_secret": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OAuth client secret for Google Drive API access.",
							},
							"token": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OAuth access token for Google Drive authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"team_drive": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Google Drive team drive ID or empty string for personal drive. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
						},
					},
					"google_photos": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `GOOGLE_PHOTOS`.",
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "OAuth client ID for Google Photos API access. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"client_secret": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OAuth client secret for Google Photos API access.",
							},
							"token": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OAuth access token for Google Photos authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"http": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `HTTP`.",
						Attributes: map[string]schema.Attribute{
							"url": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "HTTP URL for file access.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1), stringvalidator.LengthAtMost(2083)},
							},
						},
					},
					"hubic": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `HUBIC`.",
						Attributes: map[string]schema.Attribute{
							"token": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OAuth access token for Hubic authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"mega": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `MEGA`.",
						Attributes: map[string]schema.Attribute{
							"user": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "MEGA username for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"pass": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "MEGA password for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"onedrive": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `ONEDRIVE`.",
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "OAuth client ID for OneDrive API access. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"client_secret": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OAuth client secret for OneDrive API access.",
							},
							"token": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OAuth access token for OneDrive authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"drive_type": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Type of OneDrive to access.",
								Validators:          []validator.String{stringvalidator.OneOf("PERSONAL", "BUSINESS", "DOCUMENT_LIBRARY")},
							},
							"drive_id": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "OneDrive drive identifier.",
							},
						},
					},
					"pcloud": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `PCLOUD`.",
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "OAuth client ID for pCloud API access. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"client_secret": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OAuth client secret for pCloud API access.",
							},
							"token": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OAuth access token for pCloud authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"hostname": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "pCloud hostname or empty string for default. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
						},
					},
					"s3": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `S3`.",
						Attributes: map[string]schema.Attribute{
							"access_key_id": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "S3 access key ID for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"secret_access_key": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "S3 secret access key for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"endpoint": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "S3-compatible endpoint URL or empty string for AWS S3. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"region": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "S3 region or empty string for default. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"skip_region": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Whether to skip region validation. Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"signatures_v2": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Whether to use AWS Signature Version 2. Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"max_upload_parts": schema.NumberAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Maximum number of parts for multipart uploads. Defaults to `10000`.",
								Validators:          []validator.Number{numberIsInteger{}},
								Default:             numberdefault.StaticBigFloat(big.NewFloat(10000)),
							},
						},
					},
					"sftp": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `SFTP`.",
						Attributes: map[string]schema.Attribute{
							"host": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "SFTP server hostname or IP address.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"port": schema.NumberAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "SFTP server port number. Defaults to `22`.",
								Validators:          []validator.Number{numberIsInteger{}},
								Default:             numberdefault.StaticBigFloat(big.NewFloat(22)),
							},
							"user": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "SFTP username for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"pass": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "SFTP password for authentication or `null` for key-based auth.",
							},
							"private_key": schema.NumberAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "SSH private key ID for authentication or `null` for password auth.",
								Validators:          []validator.Number{numberIsInteger{}},
							},
						},
					},
					"storj_ix": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `STORJ_IX`.",
						Attributes: map[string]schema.Attribute{
							"access_key_id": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Storj S3-compatible access key ID for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"secret_access_key": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Storj S3-compatible secret access key for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"endpoint": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Storj gateway endpoint URL for S3-compatible access. Defaults to `\"https://gateway.storjshare.io/\"`.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1), stringvalidator.LengthAtMost(2083)},
								Default:             stringdefault.StaticString("https://gateway.storjshare.io/"),
							},
						},
					},
					"openstack_swift": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `OPENSTACK_SWIFT`.",
						Attributes: map[string]schema.Attribute{
							"user": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Swift username for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"key": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Swift password or API key for authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"auth": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Swift authentication URL endpoint.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"user_id": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Swift user ID for authentication. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"domain": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Swift domain name for authentication. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"tenant": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Swift tenant name for multi-tenancy. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"tenant_id": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Swift tenant ID for multi-tenancy. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"tenant_domain": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Swift tenant domain name. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"region": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Swift region name for geographic distribution. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"storage_url": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Swift storage URL endpoint. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"auth_token": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Swift authentication token for pre-authenticated access.",
							},
							"application_credential_id": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Swift application credential ID for authentication. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"application_credential_name": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Swift application credential name for authentication. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"application_credential_secret": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Swift application credential secret for authentication.",
							},
							"auth_version": schema.NumberAttribute{
								Required:            true,
								MarkdownDescription: "Swift authentication API version.\n\n* `0`: Legacy auth\n* `1`: TempAuth\n* `2`: Keystone v2.0\n* `3`: Keystone v3\n* `null`: Auto-detect",
								Validators:          []validator.Number{numberIsInteger{}},
							},
							"endpoint_type": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Swift endpoint type to use.\n\n* `public`: Public endpoint (default)\n* `internal`: Internal network endpoint\n* `admin`: Administrative endpoint\n* `null`: Use default",
								Validators:          []validator.String{stringvalidator.OneOf("public", "internal", "admin")},
							},
						},
					},
					"webdav": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `WEBDAV`.",
						Attributes: map[string]schema.Attribute{
							"url": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "WebDAV server URL endpoint.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1), stringvalidator.LengthAtMost(2083)},
							},
							"vendor": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "WebDAV server vendor type for compatibility optimizations.\n\n* `NEXTCLOUD`: Nextcloud server\n* `OWNCLOUD`: ownCloud server\n* `SHAREPOINT`: Microsoft SharePoint\n* `OTHER`: Generic WebDAV server",
								Validators:          []validator.String{stringvalidator.OneOf("NEXTCLOUD", "OWNCLOUD", "SHAREPOINT", "OTHER")},
							},
							"user": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "WebDAV username for authentication.",
							},
							"pass": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "WebDAV password for authentication.",
							},
						},
					},
					"yandex": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `type` is `YANDEX`.",
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Yandex OAuth application client ID. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
							"client_secret": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Yandex OAuth application client secret.",
							},
							"token": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Yandex OAuth access token for API authentication.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
				},
			},
			"storage_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `storage` again. Terraform never stores `storage`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
