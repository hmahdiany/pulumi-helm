package ingressNginx

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

func DeployIngressNginx(ctx *pulumi.Context, stackName string) (err error) {

	var configValues Values

	cfg := config.New(ctx, "ingress-nginx")
	cfg.RequireObject("configs", &configValues)

	appLabels := pulumi.StringMap{
		"app":        pulumi.String("nginx-ingress"),
		"deployedBy": pulumi.String("pulumi"),
		"stack": pulumi.String(stackName),
	}

	// Create a new namespace (user supplies the name of the namespace)
	ingressNs, err := corev1.NewNamespace(ctx, configValues.Namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap(appLabels),
			Name:   pulumi.String(configValues.Namespace),
		},
	})
	if err != nil {
		return err
	}

	// Use Helm to install the Nginx ingress controller
	ingresscontroller, err := helmv3.NewRelease(ctx, "ingresscontroller", &helmv3.ReleaseArgs{
		Chart:     pulumi.String("ingress-nginx"),
		Namespace: ingressNs.Metadata.Name(),
		RepositoryOpts: &helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://kubernetes.github.io/ingress-nginx"),
		},
		Version:  pulumi.String("4.10.0"),
		Name:     pulumi.String(stackName + "-ingress-nginx"),
		SkipCrds: pulumi.Bool(true),
		Values: pulumi.Map{
			"nameOverride":     pulumi.String("ingress-nginx"),
			"fullnameOverride": pulumi.String("ingress-nginx"),
			"controller": pulumi.Map{
				"replicaCount":          pulumi.Int(2),
				"enableCustomResources": pulumi.Bool(true),
				"appprotect": pulumi.Map{
					"enable": pulumi.Bool(false),
				},
				"appprotectdos": pulumi.Map{
					"enable": pulumi.Bool(false),
				},
				"service": pulumi.Map{
					"type": pulumi.String("ClusterIP"),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	// Export some values for use elsewhere
	ctx.Export("name", ingresscontroller.Name)
	ctx.Export("name", ingresscontroller.Status)

	return nil
}
