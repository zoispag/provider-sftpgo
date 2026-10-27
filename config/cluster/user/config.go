package user

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// unusedBackends are the filesystem backends we don't use. Dropping them keeps their
// credentials out of the CRD, including the observed (status) copies.
var unusedBackends = []string{"azblobconfig", "cryptconfig", "ftpconfig", "gcsconfig", "httpconfig", "sftpconfig"}

// Configure configures the sftpgo_user resource.
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("sftpgo_user", func(r *ujconfig.Resource) {
		r.ShortGroup = "user"
		r.Kind = "User"
		// password and password_wo conflict; expose only password, which becomes a secret reference.
		delete(r.TerraformResource.Schema, "password_wo")
		delete(r.TerraformResource.Schema, "password_wo_version")
		pruneBackends(r.TerraformResource.Schema["filesystem"])
		// A virtual folder's filesystem is read back from SFTPGo, including S3 secret envelopes
		// and session tokens. It's never set from the User, so leave it out of the status.
		if vf, ok := r.TerraformResource.Schema["virtual_folders"].Elem.(*schema.Resource); ok {
			delete(vf.Schema, "filesystem")
		}
		// Upstream doesn't mark the S3 session token sensitive; do it so it becomes a secret reference.
		if fs, ok := r.TerraformResource.Schema["filesystem"].Elem.(*schema.Resource); ok {
			if s3, ok := fs.Schema["s3config"].Elem.(*schema.Resource); ok {
				if st, ok := s3.Schema["session_token"]; ok {
					st.Sensitive = true
				}
			}
		}
	})
}

func pruneBackends(fs *schema.Schema) {
	if fs == nil {
		return
	}
	if res, ok := fs.Elem.(*schema.Resource); ok {
		for _, b := range unusedBackends {
			delete(res.Schema, b)
		}
	}
}
