package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

// Custom lists the resources written by hand, because the API does not describe them in a way the
// generator can follow: an app's lifecycle, a filesystem ACL, and a file's contents.
var Custom = []func() resource.Resource{
	NewApp,
	NewFile,
	NewFilesystemACL,
}
