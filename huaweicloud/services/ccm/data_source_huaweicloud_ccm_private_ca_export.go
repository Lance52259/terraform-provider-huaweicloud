// Generated by PMS #277
package ccm

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourcePrivateCaExport() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateCaExportRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"ca_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the CA certificate you want to export.`,
			},
			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The certificate content.`,
			},
			"certificate_chain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The content of the certificate chain.`,
			},
		},
	}
}

type PrivateCaExportDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newPrivateCaExportDSWrapper(d *schema.ResourceData, meta interface{}) *PrivateCaExportDSWrapper {
	return &PrivateCaExportDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourcePrivateCaExportRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newPrivateCaExportDSWrapper(d, meta)
	expCerAutCerRst, err := wrapper.ExportCertificateAuthorityCertificate()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.exportCertificateAuthorityCertificateToSchema(expCerAutCerRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCM POST /v1/private-certificate-authorities/{ca_id}/export
func (w *PrivateCaExportDSWrapper) ExportCertificateAuthorityCertificate() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ccm")
	if err != nil {
		return nil, err
	}

	uri := "/v1/private-certificate-authorities/{ca_id}/export"
	uri = strings.ReplaceAll(uri, "{ca_id}", w.Get("ca_id").(string))
	return httphelper.New(client).
		Method("POST").
		URI(uri).
		Request().
		Result()
}

func (w *PrivateCaExportDSWrapper) exportCertificateAuthorityCertificateToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("certificate", body.Get("certificate").Value()),
		d.Set("certificate_chain", body.Get("certificate_chain").Value()),
	)
	return mErr.ErrorOrNil()
}
