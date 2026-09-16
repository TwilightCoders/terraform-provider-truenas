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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/big"
)

var modelCertificate = &model{
	typeName:     "certificate",
	namespace:    "certificate",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "certificate.create",
	updateMethod: "certificate.update",
	getMethod:    "certificate.get_instance",
	deleteMethod: "certificate.delete",
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
			description: "Certificate name.",
		},
		{
			name: "create_type", api: "create_type", path: "create_type",
			kind: kindString, role: roleRequired,
			createOnly:  true,
			description: "Type of certificate creation operation. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
		},
		{
			name: "add_to_trusted_store", api: "add_to_trusted_store", path: "add_to_trusted_store",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to add this certificate to the trusted certificate store. Defaults to `false`.",
		},
		{
			name: "certificate", api: "certificate", path: "certificate",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "PEM-encoded certificate to import or `null`. Changing this forces a new resource.",
		},
		{
			name: "privatekey", api: "privatekey", path: "privatekey",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, sensitive: true, readable: true,
			description: "PEM-encoded private key. TrueNAS generates one for `CERTIFICATE_CREATE_ACME` and `CERTIFICATE_CREATE_CSR`, and reports it on read, so it is stored in Terraform state for every certificate this resource manages. Protect state accordingly. Changing this forces a new resource.",
		},
		{
			name: "csr", api: "CSR", path: "CSR",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "PEM-encoded certificate signing request to import or `null`. Changing this forces a new resource.",
		},
		{
			name: "key_length", api: "key_length", path: "key_length",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "RSA key length in bits or `null`. Changing this forces a new resource.",
		},
		{
			name: "key_type", api: "key_type", path: "key_type",
			kind: kindString, role: roleOptionalComputed,
			replace: true, stable: true, readable: true,
			description: "Type of cryptographic key to generate. Changing this forces a new resource.",
		},
		{
			name: "ec_curve", api: "ec_curve", path: "ec_curve",
			kind: kindString, role: roleOptional,
			createOnly:  true,
			description: "Elliptic curve to use for EC keys. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
		},
		{
			name: "passphrase", api: "passphrase", path: "passphrase",
			kind: kindString, role: roleOptional,
			nullable: true, sensitive: true, writeOnly: true,
			description: "Passphrase to protect the private key or `null`.",
		},
		{
			name: "city", api: "city", path: "city",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "City or locality name for certificate subject or `null`. Changing this forces a new resource.",
		},
		{
			name: "common", api: "common", path: "common",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "Common name for certificate subject or `null`. Changing this forces a new resource.",
		},
		{
			name: "country", api: "country", path: "country",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "Country name for certificate subject or `null`. Changing this forces a new resource.",
		},
		{
			name: "email", api: "email", path: "email",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "Email address for certificate subject or `null`. Changing this forces a new resource.",
		},
		{
			name: "organization", api: "organization", path: "organization",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "Organization name for certificate subject or `null`. Changing this forces a new resource.",
		},
		{
			name: "organizational_unit", api: "organizational_unit", path: "organizational_unit",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "Organizational unit for certificate subject or `null`. Changing this forces a new resource.",
		},
		{
			name: "state", api: "state", path: "state",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "State or province name for certificate subject or `null`. Changing this forces a new resource.",
		},
		{
			name: "digest_algorithm", api: "digest_algorithm", path: "digest_algorithm",
			kind: kindString, role: roleOptional,
			createOnly:  true,
			description: "Hash algorithm for certificate signing. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
		},
		{
			name: "san", api: "san", path: "san",
			kind: kindList, role: roleOptionalComputed,
			replace: true, stable: true, readable: true,
			description: "Subject alternative names. Write them bare (`example.com`) or DNS-prefixed (`DNS:example.com`); both are accepted and mean the same name. Prefix an address with `IP:` to request an IP name. TrueNAS treats every other prefix, including `email:` and `URI:`, as part of a DNS name rather than as that name type. Changing this forces a new resource.",
			canonical:   CanonicalSAN,
			elem: &node{
				name: "", api: "", path: "san",
				kind: kindString, role: roleRequired,
				readable:  true,
				canonical: CanonicalSAN,
			},
		},
		{
			name: "cert_extensions", api: "cert_extensions", path: "cert_extensions",
			kind: kindObject, role: roleOptional,
			createOnly:  true,
			description: "Certificate extensions configuration. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
			children: []*node{
				{
					name: "basicconstraints", api: "BasicConstraints", path: "cert_extensions.BasicConstraints",
					kind: kindObject, role: roleOptional,
					description: "Basic Constraints extension configuration for certificate authority capabilities. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
					children: []*node{
						{
							name: "ca", api: "ca", path: "cert_extensions.BasicConstraints.ca",
							kind: kindBool, role: roleOptional,
							description: "Whether this certificate is authorized to sign other certificates as a Certificate Authority (CA). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "enabled", api: "enabled", path: "cert_extensions.BasicConstraints.enabled",
							kind: kindBool, role: roleOptional,
							description: "Whether the Basic Constraints X.509 extension is present in the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "path_length", api: "path_length", path: "cert_extensions.BasicConstraints.path_length",
							kind: kindInt, role: roleOptional,
							nullable:    true,
							description: "Maximum number of intermediate CA certificates that may follow this certificate in a valid certificate chain.     `null` indicates no path length constraint. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "extension_critical", api: "extension_critical", path: "cert_extensions.BasicConstraints.extension_critical",
							kind: kindBool, role: roleOptional,
							description: "Whether the Basic Constraints extension is marked as critical. If `true`, applications that do not understand     this extension must reject the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
					},
				},
				{
					name: "extendedkeyusage", api: "ExtendedKeyUsage", path: "cert_extensions.ExtendedKeyUsage",
					kind: kindObject, role: roleOptional,
					description: "Extended Key Usage extension configuration specifying certificate purposes. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
					children: []*node{
						{
							name: "usages", api: "usages", path: "cert_extensions.ExtendedKeyUsage.usages",
							kind: kindList, role: roleOptional,
							description: "Array of Extended Key Usage (EKU) purposes that define what the certificate may be used for     (e.g., 'SERVER_AUTH', 'CLIENT_AUTH', 'CODE_SIGNING'). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							elem: &node{
								name: "", api: "", path: "cert_extensions.ExtendedKeyUsage.usages",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
						{
							name: "enabled", api: "enabled", path: "cert_extensions.ExtendedKeyUsage.enabled",
							kind: kindBool, role: roleOptional,
							description: "Whether the Extended Key Usage X.509 extension is present in the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "extension_critical", api: "extension_critical", path: "cert_extensions.ExtendedKeyUsage.extension_critical",
							kind: kindBool, role: roleOptional,
							description: "Whether the Extended Key Usage extension is marked as critical. If `true`, applications that do not understand     this extension must reject the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
					},
				},
				{
					name: "keyusage", api: "KeyUsage", path: "cert_extensions.KeyUsage",
					kind: kindObject, role: roleOptional,
					description: "Key Usage extension configuration defining permitted cryptographic operations. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
					children: []*node{
						{
							name: "enabled", api: "enabled", path: "cert_extensions.KeyUsage.enabled",
							kind: kindBool, role: roleOptional,
							description: "Whether the Key Usage X.509 extension is present in the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "digital_signature", api: "digital_signature", path: "cert_extensions.KeyUsage.digital_signature",
							kind: kindBool, role: roleOptional,
							description: "Whether the certificate may be used for digital signatures to verify identity or integrity. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "content_commitment", api: "content_commitment", path: "cert_extensions.KeyUsage.content_commitment",
							kind: kindBool, role: roleOptional,
							description: "Whether the certificate may be used for non-repudiation (proving content commitment). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "key_encipherment", api: "key_encipherment", path: "cert_extensions.KeyUsage.key_encipherment",
							kind: kindBool, role: roleOptional,
							description: "Whether the certificate's public key may be used for encrypting symmetric keys. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "data_encipherment", api: "data_encipherment", path: "cert_extensions.KeyUsage.data_encipherment",
							kind: kindBool, role: roleOptional,
							description: "Whether the certificate's public key may be used for directly encrypting raw data. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "key_agreement", api: "key_agreement", path: "cert_extensions.KeyUsage.key_agreement",
							kind: kindBool, role: roleOptional,
							description: "Whether the certificate's public key may be used for key agreement protocols (e.g., Diffie-Hellman). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "key_cert_sign", api: "key_cert_sign", path: "cert_extensions.KeyUsage.key_cert_sign",
							kind: kindBool, role: roleOptional,
							description: "Whether the certificate may be used to sign other certificates (CA functionality). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "crl_sign", api: "crl_sign", path: "cert_extensions.KeyUsage.crl_sign",
							kind: kindBool, role: roleOptional,
							description: "Whether the certificate may be used to sign Certificate Revocation Lists (CRLs). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "encipher_only", api: "encipher_only", path: "cert_extensions.KeyUsage.encipher_only",
							kind: kindBool, role: roleOptional,
							description: "Whether the public key may only be used for encryption when `key_agreement` is also set. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "decipher_only", api: "decipher_only", path: "cert_extensions.KeyUsage.decipher_only",
							kind: kindBool, role: roleOptional,
							description: "Whether the public key may only be used for decryption when `key_agreement` is also set. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
						{
							name: "extension_critical", api: "extension_critical", path: "cert_extensions.KeyUsage.extension_critical",
							kind: kindBool, role: roleOptional,
							description: "Whether the Key Usage extension is marked as critical. If `true`, applications that do not understand     this extension must reject the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						},
					},
				},
			},
		},
		{
			name: "acme_directory_uri", api: "acme_directory_uri", path: "acme_directory_uri",
			kind: kindString, role: roleOptional,
			nullable: true, createOnly: true,
			description: "ACME directory URI to be used for ACME certificate creation. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
		},
		{
			name: "csr_id", api: "csr_id", path: "csr_id",
			kind: kindInt, role: roleOptional,
			nullable: true, createOnly: true,
			description: "CSR to be used for ACME certificate creation. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
		},
		{
			name: "tos", api: "tos", path: "tos",
			kind: kindBool, role: roleOptional,
			nullable: true, createOnly: true,
			description: "Set this when creating an ACME certificate to accept terms of service of the ACME service. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
		},
		{
			name: "dns_mapping", api: "dns_mapping", path: "dns_mapping",
			kind: kindMap, role: roleOptional,
			createOnly:  true,
			description: "A mapping of domain to ACME DNS Authenticator ID for each domain listed in SAN or common name of the CSR. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
			elem: &node{
				name: "", api: "", path: "dns_mapping",
				kind: kindInt, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "renew_days", api: "renew_days", path: "renew_days",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "10",
			description: "Number of days before the certificate expiration date to attempt certificate renewal. If certificate renewal     fails, renewal will be reattempted every day until expiration. Defaults to `10`.",
		},
		{
			name: "type", api: "type", path: "type",
			kind: kindInt, role: roleComputed,
			readable:    true,
			description: "Internal certificate type identifier used to determine certificate capabilities.",
		},
		{
			name: "acme_uri", api: "acme_uri", path: "acme_uri",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "ACME directory server URI used for automated certificate management. `null` for non-ACME certificates.",
		},
		{
			name: "domains_authenticators", api: "domains_authenticators", path: "domains_authenticators",
			kind: kindAny, role: roleComputed,
			nullable: true, readable: true,
			description: "Mapping of domain names to ACME DNS authenticator IDs for domain validation. `null` for non-ACME     certificates.",
		},
		{
			name: "acme", api: "acme", path: "acme",
			kind: kindAny, role: roleComputed,
			nullable: true, readable: true,
			description: "ACME registration and account information used for certificate lifecycle management. `null` for     non-ACME certificates.",
		},
		{
			name: "root_path", api: "root_path", path: "root_path",
			kind: kindString, role: roleComputed,
			readable:    true,
			description: "Filesystem path where certificate-related files are stored.",
		},
		{
			name: "certificate_path", api: "certificate_path", path: "certificate_path",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "Filesystem path to the certificate file (.crt). `null` if no certificate is available.",
		},
		{
			name: "privatekey_path", api: "privatekey_path", path: "privatekey_path",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "Filesystem path to the private key file (.key). `null` if no private key is available.",
		},
		{
			name: "csr_path", api: "csr_path", path: "csr_path",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "Filesystem path to the certificate signing request file (.csr). `null` if no CSR is available.",
		},
		{
			name: "cert_type", api: "cert_type", path: "cert_type",
			kind: kindString, role: roleComputed,
			readable:    true,
			description: "Human-readable certificate type, typically 'CERTIFICATE' for standard certificates.",
		},
		{
			name: "cert_type_existing", api: "cert_type_existing", path: "cert_type_existing",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether this is an existing certificate (imported or generated).",
		},
		{
			name: "cert_type_csr", api: "cert_type_CSR", path: "cert_type_CSR",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether this entry represents a Certificate Signing Request (CSR) rather than a signed certificate.",
		},
		{
			name: "cert_type_ca", api: "cert_type_CA", path: "cert_type_CA",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether this certificate is a Certificate Authority (CA) certificate.",
		},
		{
			name: "chain_list", api: "chain_list", path: "chain_list",
			kind: kindList, role: roleComputed,
			readable:    true,
			description: "Array of PEM-encoded certificates in the certificate chain, starting with the leaf certificate.",
			elem: &node{
				name: "", api: "", path: "chain_list",
				kind: kindString, role: roleComputed,
				readable: true,
			},
		},
		{
			name: "dn", api: "DN", path: "DN",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "Distinguished Name (DN) of the certificate subject in RFC 2253 format. `null` if certificate parsing failed.",
		},
		{
			name: "subject_name_hash", api: "subject_name_hash", path: "subject_name_hash",
			kind: kindInt, role: roleComputed,
			nullable: true, readable: true,
			description: "Hash of the certificate subject name. `null` if certificate parsing failed.",
		},
		{
			name: "extensions", api: "extensions", path: "extensions",
			kind: kindAny, role: roleComputed,
			readable:    true,
			description: "X.509 certificate extensions parsed into a dictionary structure.",
		},
		{
			name: "lifetime", api: "lifetime", path: "lifetime",
			kind: kindInt, role: roleComputed,
			nullable: true, readable: true,
			description: "Certificate validity period in seconds. `null` if certificate parsing failed.",
		},
		{
			name: "from", api: "from", path: "from",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "Certificate validity start date in ISO 8601 format. `null` if certificate parsing failed.",
		},
		{
			name: "until", api: "until", path: "until",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "Certificate validity end date in ISO 8601 format. `null` if certificate parsing failed.",
		},
		{
			name: "serial", api: "serial", path: "serial",
			kind: kindInt, role: roleComputed,
			nullable: true, readable: true,
			description: "Certificate serial number. `null` if certificate parsing failed.",
		},
		{
			name: "chain", api: "chain", path: "chain",
			kind: kindBool, role: roleComputed,
			nullable: true, readable: true,
			description: "Whether this certificate has an associated certificate chain. `null` if unavailable.",
		},
		{
			name: "fingerprint", api: "fingerprint", path: "fingerprint",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "SHA-256 fingerprint of the certificate in hexadecimal format. `null` if certificate parsing failed.",
		},
		{
			name: "expired", api: "expired", path: "expired",
			kind: kindBool, role: roleComputed,
			nullable: true, readable: true,
			description: "Whether the certificate has expired. `null` if certificate parsing failed.",
		},
		{
			name: "parsed", api: "parsed", path: "parsed",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether the certificate data was successfully parsed and validated.",
		},
		{
			name: "passphrase_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `passphrase` again. Terraform never stores `passphrase`.",
		},
	},
}

// NewCertificate returns the certificate resource.
func NewCertificate() resource.Resource {
	return &certificateResource{crudResource{model: modelCertificate}}
}

type certificateResource struct{ crudResource }

// NewCertificateDataSource returns the certificate data source.
func NewCertificateDataSource() datasource.DataSource {
	return &dataSource{model: modelCertificate, attrs: dataAttrs(modelCertificate)}
}

// NewCertificateList returns the certificate list resource, for terraform query.
func NewCertificateList() list.ListResource {
	return &listResource{crudResource{model: modelCertificate}}
}

func (r *certificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a certificate. `create_type` selects how: `CERTIFICATE_CREATE_IMPORTED` (bring `certificate` and `privatekey`), `CERTIFICATE_CREATE_CSR` (TrueNAS generates the key and a CSR), `CERTIFICATE_CREATE_IMPORTED_CSR`, or `CERTIFICATE_CREATE_ACME` (issue from an existing CSR via `csr_id`, `dns_mapping` and `tos`). TrueNAS renews ACME certificates itself.\n\nTrueNAS generates the private key for the CSR and ACME types and reports it on read, so `privatekey` is stored in Terraform state. That is what lets an issued certificate be served elsewhere, and it means state holds key material: protect it accordingly.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Certificate name.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1), stringvalidator.LengthAtMost(120)},
			},
			"create_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of certificate creation operation. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
				Validators:          []validator.String{stringvalidator.OneOf("CERTIFICATE_CREATE_IMPORTED", "CERTIFICATE_CREATE_CSR", "CERTIFICATE_CREATE_IMPORTED_CSR", "CERTIFICATE_CREATE_ACME")},
			},
			"add_to_trusted_store": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to add this certificate to the trusted certificate store. Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"certificate": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "PEM-encoded certificate to import or `null`. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"privatekey": schema.StringAttribute{
				Optional: true, Computed: true, Sensitive: true,
				MarkdownDescription: "PEM-encoded private key. TrueNAS generates one for `CERTIFICATE_CREATE_ACME` and `CERTIFICATE_CREATE_CSR`, and reports it on read, so it is stored in Terraform state for every certificate this resource manages. Protect state accordingly. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"csr": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "PEM-encoded certificate signing request to import or `null`. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"key_length": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "RSA key length in bits or `null`. Changing this forces a new resource.",
				Validators:          []validator.Number{numberIsInteger{}},
				PlanModifiers:       []planmodifier.Number{numberplanmodifier.RequiresReplace(), numberplanmodifier.UseStateForUnknown()},
			},
			"key_type": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Type of cryptographic key to generate. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.OneOf("RSA", "EC")},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"ec_curve": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Elliptic curve to use for EC keys. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
				Validators:          []validator.String{stringvalidator.OneOf("SECP256R1", "SECP384R1", "SECP521R1", "ed25519")},
			},
			"passphrase": schema.StringAttribute{
				Optional: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "Passphrase to protect the private key or `null`.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"city": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "City or locality name for certificate subject or `null`. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"common": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Common name for certificate subject or `null`. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"country": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Country name for certificate subject or `null`. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"email": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Email address for certificate subject or `null`. Changing this forces a new resource.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"organization": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Organization name for certificate subject or `null`. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"organizational_unit": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Organizational unit for certificate subject or `null`. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"state": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "State or province name for certificate subject or `null`. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"digest_algorithm": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Hash algorithm for certificate signing. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
				Validators:          []validator.String{stringvalidator.OneOf("SHA224", "SHA256", "SHA384", "SHA512")},
			},
			"san": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Subject alternative names. Write them bare (`example.com`) or DNS-prefixed (`DNS:example.com`); both are accepted and mean the same name. Prefix an address with `IP:` to request an IP name. TrueNAS treats every other prefix, including `email:` and `URI:`, as part of a DNS name rather than as that name type. Changing this forces a new resource.",
				ElementType:         types.StringType,
				PlanModifiers:       []planmodifier.List{listplanmodifier.RequiresReplace(), listplanmodifier.UseStateForUnknown()},
			},
			"cert_extensions": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Certificate extensions configuration. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
				Attributes: map[string]schema.Attribute{
					"basicconstraints": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Basic Constraints extension configuration for certificate authority capabilities. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						Attributes: map[string]schema.Attribute{
							"ca": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether this certificate is authorized to sign other certificates as a Certificate Authority (CA). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"enabled": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the Basic Constraints X.509 extension is present in the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"path_length": schema.NumberAttribute{
								Optional:            true,
								MarkdownDescription: "Maximum number of intermediate CA certificates that may follow this certificate in a valid certificate chain.     `null` indicates no path length constraint. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
								Validators:          []validator.Number{numberIsInteger{}},
							},
							"extension_critical": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the Basic Constraints extension is marked as critical. If `true`, applications that do not understand     this extension must reject the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
						},
					},
					"extendedkeyusage": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Extended Key Usage extension configuration specifying certificate purposes. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						Attributes: map[string]schema.Attribute{
							"usages": schema.ListAttribute{
								Optional:            true,
								MarkdownDescription: "Array of Extended Key Usage (EKU) purposes that define what the certificate may be used for     (e.g., 'SERVER_AUTH', 'CLIENT_AUTH', 'CODE_SIGNING'). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
								ElementType:         types.StringType,
							},
							"enabled": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the Extended Key Usage X.509 extension is present in the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"extension_critical": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the Extended Key Usage extension is marked as critical. If `true`, applications that do not understand     this extension must reject the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
						},
					},
					"keyusage": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Key Usage extension configuration defining permitted cryptographic operations. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the Key Usage X.509 extension is present in the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"digital_signature": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the certificate may be used for digital signatures to verify identity or integrity. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"content_commitment": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the certificate may be used for non-repudiation (proving content commitment). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"key_encipherment": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the certificate's public key may be used for encrypting symmetric keys. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"data_encipherment": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the certificate's public key may be used for directly encrypting raw data. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"key_agreement": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the certificate's public key may be used for key agreement protocols (e.g., Diffie-Hellman). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"key_cert_sign": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the certificate may be used to sign other certificates (CA functionality). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"crl_sign": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the certificate may be used to sign Certificate Revocation Lists (CRLs). TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"encipher_only": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the public key may only be used for encryption when `key_agreement` is also set. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"decipher_only": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the public key may only be used for decryption when `key_agreement` is also set. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
							"extension_critical": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether the Key Usage extension is marked as critical. If `true`, applications that do not understand     this extension must reject the certificate. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
							},
						},
					},
				},
			},
			"acme_directory_uri": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "ACME directory URI to be used for ACME certificate creation. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"csr_id": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "CSR to be used for ACME certificate creation. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"tos": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Set this when creating an ACME certificate to accept terms of service of the ACME service. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
			},
			"dns_mapping": schema.MapAttribute{
				Optional:            true,
				MarkdownDescription: "A mapping of domain to ACME DNS Authenticator ID for each domain listed in SAN or common name of the CSR. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
				ElementType:         types.NumberType,
			},
			"renew_days": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Number of days before the certificate expiration date to attempt certificate renewal. If certificate renewal     fails, renewal will be reattempted every day until expiration. Defaults to `10`.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 30}},
				Default:             numberdefault.StaticBigFloat(big.NewFloat(10)),
			},
			"type": schema.NumberAttribute{
				Computed:            true,
				MarkdownDescription: "Internal certificate type identifier used to determine certificate capabilities.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"acme_uri": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ACME directory server URI used for automated certificate management. `null` for non-ACME certificates.",
			},
			"domains_authenticators": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Mapping of domain names to ACME DNS authenticator IDs for domain validation. `null` for non-ACME     certificates.",
			},
			"acme": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ACME registration and account information used for certificate lifecycle management. `null` for     non-ACME certificates.",
			},
			"root_path": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Filesystem path where certificate-related files are stored.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"certificate_path": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Filesystem path to the certificate file (.crt). `null` if no certificate is available.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"privatekey_path": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Filesystem path to the private key file (.key). `null` if no private key is available.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"csr_path": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Filesystem path to the certificate signing request file (.csr). `null` if no CSR is available.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"cert_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Human-readable certificate type, typically 'CERTIFICATE' for standard certificates.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"cert_type_existing": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this is an existing certificate (imported or generated).",
			},
			"cert_type_csr": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this entry represents a Certificate Signing Request (CSR) rather than a signed certificate.",
			},
			"cert_type_ca": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this certificate is a Certificate Authority (CA) certificate.",
			},
			"chain_list": schema.ListAttribute{
				Computed:            true,
				MarkdownDescription: "Array of PEM-encoded certificates in the certificate chain, starting with the leaf certificate.",
				ElementType:         types.StringType,
			},
			"dn": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Distinguished Name (DN) of the certificate subject in RFC 2253 format. `null` if certificate parsing failed.",
			},
			"subject_name_hash": schema.NumberAttribute{
				Computed:            true,
				MarkdownDescription: "Hash of the certificate subject name. `null` if certificate parsing failed.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"extensions": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "X.509 certificate extensions parsed into a dictionary structure.",
			},
			"lifetime": schema.NumberAttribute{
				Computed:            true,
				MarkdownDescription: "Certificate validity period in seconds. `null` if certificate parsing failed.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"from": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Certificate validity start date in ISO 8601 format. `null` if certificate parsing failed.",
			},
			"until": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Certificate validity end date in ISO 8601 format. `null` if certificate parsing failed.",
			},
			"serial": schema.NumberAttribute{
				Computed:            true,
				MarkdownDescription: "Certificate serial number. `null` if certificate parsing failed.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"chain": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this certificate has an associated certificate chain. `null` if unavailable.",
			},
			"fingerprint": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SHA-256 fingerprint of the certificate in hexadecimal format. `null` if certificate parsing failed.",
			},
			"expired": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the certificate has expired. `null` if certificate parsing failed.",
			},
			"parsed": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the certificate data was successfully parsed and validated.",
			},
			"passphrase_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `passphrase` again. Terraform never stores `passphrase`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
