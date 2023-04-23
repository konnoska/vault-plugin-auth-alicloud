// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// operationPrefixAliCloud is used as a prefix for OpenAPI operation id's.
const operationPrefixAliCloud = "ali-cloud"

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	client := cleanhttp.DefaultClient()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	b := newBackend(client)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// newBackend exists for testability. It allows us to inject a fake client.
func newBackend(client *http.Client) *backend {
	b := &backend{
		identityClient: client,
	}
	b.Backend = &framework.Backend{
		AuthRenew: b.pathLoginRenew,
		Help:      backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},
		Paths: []*framework.Path{
			pathConfig(b),
			pathLogin(b),
			pathListRole(b),
			pathListRoles(b),
			pathRole(b),
		},
		BackendType: logical.TypeCredential,
	}
	return b
}

type backend struct {
	*framework.Backend

	identityClient *http.Client
}

func (b *backend) getAliasName(ctx context.Context, data logical.Storage, arn *arn, identity *sts.GetCallerIdentityResponse) (string, error) {
	fmt.Println("-------------ENTERED getAliasName")
	defer fmt.Println("-------------EXITED getAliasName")
	config, err := b.config(ctx, data)

	if err != nil {
		return "", fmt.Errorf("unable to retrieve backend configuration: %w", err)
	}

	fmt.Println("-------------BEFORE SWITCH getAliasName")
	fmt.Printf("CONFIG:  %+v\n", config)
	switch config.RamAlias {
	case "roleArn":
		fmt.Println("The case is roleArn")
		return arn.RoleArn, nil
	default:
		fmt.Println("The case is default")
		return identity.PrincipalId, nil
	}

}

const backendHelp = `
That AliCloud RAM auth method allows entities to authenticate based on their
identity and pre-configured roles.
`
