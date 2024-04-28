package cockroachdb

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type Values struct {
	Namespace string
}

func DeployCockroachDB(ctx *pulumi.Context, stackName string) (err error) {

	var configValues Values

	cfg := config.New(ctx, "cockroachdb")
	cfg.RequireObject("configs", &configValues)

	appLabels := pulumi.StringMap{
		"app":        pulumi.String("cockroachdb"),
		"deployedBy": pulumi.String("pulumi"),
		"stack": pulumi.String(stackName),
	}

	// Create a new namespace (user supplies the name of the namespace)
	cockroachdbNS, err := corev1.NewNamespace(ctx, configValues.Namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap(appLabels),
			Name:   pulumi.String(configValues.Namespace),
		},
	})
	if err != nil {
		return err
	}

	// Use Helm to install the Nginx ingress controller
	cockroachdb, err := helmv3.NewRelease(ctx, "cockroachdb", &helmv3.ReleaseArgs{
		Chart:     pulumi.String("cockroachdb"),
		Namespace: cockroachdbNS.Metadata.Name(),
		RepositoryOpts: &helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://charts.cockroachdb.com/"),
		},
		Version:  pulumi.String("12.0.3"),
		Name:     pulumi.String(stackName + "-cockroachdb"),
		SkipCrds: pulumi.Bool(true),
		Values: pulumi.Map{
			"nameOverride":     pulumi.String("cockroachdb"),
			"fullnameOverride": pulumi.String("cockroachdb"),
			"conf": pulumi.Map{
				"cluster-name": pulumi.String("mydbcluster"),
				"log": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
			},
			"statefulset": pulumi.Map{
				"replicas": pulumi.Int(1),
			},
			"ingress": pulumi.Map{
				"enabled": pulumi.Bool(false),
				"hosts": pulumi.StringArray{
					pulumi.String("mycockroach.local"),
				},
			},
			"prometheus": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"storage": pulumi.Map{
				"persistentVolume": pulumi.Map{
					"enabled": pulumi.Bool(false),
				},
			},
			"tls": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
		},
	})
	if err != nil {
		return err
	}

	// Export some values for use elsewhere
	ctx.Export("name", cockroachdb.Name)
	ctx.Export("name", cockroachdb.Status)

	return nil
}
