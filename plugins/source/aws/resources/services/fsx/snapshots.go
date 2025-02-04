package fsx

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/fsx"
	"github.com/aws/aws-sdk-go-v2/service/fsx/types"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/cq-provider-sdk/provider/diag"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

//go:generate cq-gen --resource snapshots --config snapshots.hcl --output .
func Snapshots() *schema.Table {
	return &schema.Table{
		Name:         "aws_fsx_snapshots",
		Description:  "A snapshot of an Amazon FSx for OpenZFS volume",
		Resolver:     fetchFsxSnapshots,
		Multiplex:    client.ServiceAccountRegionMultiplexer("fsx"),
		IgnoreError:  client.IgnoreAccessDeniedServiceDisabled,
		DeleteFilter: client.DeleteAccountRegionFilter,
		Options:      schema.TableCreationOptions{PrimaryKeys: []string{"arn"}},
		Columns: []schema.Column{
			{
				Name:        "account_id",
				Description: "The AWS Account ID of the resource.",
				Type:        schema.TypeString,
				Resolver:    client.ResolveAWSAccount,
			},
			{
				Name:        "region",
				Description: "The AWS Region of the resource.",
				Type:        schema.TypeString,
				Resolver:    client.ResolveAWSRegion,
			},
			{
				Name:        "creation_time",
				Description: "The time that the resource was created, in seconds (since 1970-01-01T00:00:00Z), also known as Unix time",
				Type:        schema.TypeTimestamp,
			},
			{
				Name:        "lifecycle",
				Description: "The lifecycle status of the snapshot  * PENDING - Amazon FSx hasn't started creating the snapshot  * CREATING - Amazon FSx is creating the snapshot  * DELETING - Amazon FSx is deleting the snapshot  * AVAILABLE - The snapshot is fully available",
				Type:        schema.TypeString,
			},
			{
				Name:        "lifecycle_transition_reason_message",
				Description: "A detailed error message",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("LifecycleTransitionReason.Message"),
			},
			{
				Name:        "name",
				Description: "The name of the snapshot",
				Type:        schema.TypeString,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for a given resource",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("ResourceARN"),
			},
			{
				Name:        "snapshot_id",
				Description: "The ID of the snapshot",
				Type:        schema.TypeString,
			},
			{
				Name:        "tags",
				Description: "A list of Tag values, with a maximum of 50 elements",
				Type:        schema.TypeJSON,
				Resolver:    resolveSnapshotsTags,
			},
			{
				Name:        "volume_id",
				Description: "The ID of the volume that the snapshot is of",
				Type:        schema.TypeString,
			},
		},
	}
}

// ====================================================================================================================
//                                               Table Resolver Functions
// ====================================================================================================================

func fetchFsxSnapshots(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	cl := meta.(*client.Client)
	svc := cl.Services().FSX
	input := fsx.DescribeSnapshotsInput{MaxResults: aws.Int32(1000)}
	for {
		result, err := svc.DescribeSnapshots(ctx, &input)
		if err != nil {
			return diag.WrapError(err)
		}
		res <- result.Snapshots
		if aws.ToString(result.NextToken) == "" {
			break
		}
		input.NextToken = result.NextToken
	}
	return nil
}
func resolveSnapshotsTags(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	return diag.WrapError(resource.Set(c.Name, client.TagsToMap(resource.Item.(types.Snapshot).Tags)))
}
