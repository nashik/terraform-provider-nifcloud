package sslcertificate

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	res := describeResponses{}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		req := svc.DescribeSslCertificatesRequest(expandDescribeSSLCertificatesInput(d))
		res.describeSSLCertificatesResponse, err = req.Send(ctx)
		if err != nil {
			return fmt.Errorf("failed reading SSLCertificate: %s", err.Error())
		}
		return nil
	})

	eg.Go(func() error {
		var err error
		req := svc.DownloadSslCertificateRequest(expandDownloadSSLCertificateInputForCert(d))
		res.downloadSSLCertificateResponseForCert, err = req.Send(ctx)
		if err != nil {
			return checkNotFoundError(err)
		}
		return nil
	})

	eg.Go(func() error {
		var err error
		req := svc.DownloadSslCertificateRequest(expandDownloadSSLCertificateInputForKey(d))
		res.downloadSSLCertificateResponseForKey, err = req.Send(ctx)
		if err != nil {
			return checkNotFoundError(err)
		}
		return nil
	})

	if _, ok := d.GetOk("ca"); ok {
		eg.Go(func() error {
			var err error
			req := svc.DownloadSslCertificateRequest(expandDownloadSSLCertificateInputForCA(d))
			res.downloadSSLCertificateResponseForCA, err = req.Send(ctx)
			if err != nil {
				return checkNotFoundError(err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return diag.FromErr(err)
	}

	if err := flatten(d, &res); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func checkNotFoundError(err error) error {
	var awserr awserr.Error
	if errors.As(err, &awserr) && awserr.Code() == "Client.InvalidParameterNotFound.SslCertificate" {
		return nil
	}
	return fmt.Errorf("failed downloading certificate: %s", err.Error())
}
