package main

import (
	ingressNginx "k8s-addons/ingress-nginx"
	cockroachdb "k8s-addons/cockroachdb"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Get stack name
		stackName := ctx.Stack()

		err := ingressNginx.DeployIngressNginx(ctx, stackName)
		if err != nil {
			return err
		}

		err = cockroachdb.DeployCockroachDB(ctx, stackName)
		if err != nil {
			return err
		}
		
		return nil
	})

}
