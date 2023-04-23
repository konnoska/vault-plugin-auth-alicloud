package alicloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAliCloud,
		},
		Fields: map[string]*framework.FieldSchema{
			"ram_alias": {
				Type:    framework.TypeString,
				Default: defaultRamAlias,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: "auth-configuration",
				},
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "auth",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "auth",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "delete",
					OperationSuffix: "auth-configuration",
				},
			},
		},
		ExistenceCheck: b.pathConfigExistenceCheck,

		HelpSynopsis:    confHelpSyn,
		HelpDescription: confHelpDesc,
	}
}

type alicloudConfig struct {
	RamAlias string `json:"ram_alias"`
}

func newAlicloudConfig() alicloudConfig {
	return alicloudConfig{
		RamAlias: defaultRamAlias,
	}
}

func (b *backend) config(ctx context.Context, s logical.Storage) (*alicloudConfig, error) {
	fmt.Println("------ENTERED cinfig")
	config := newAlicloudConfig()

	rawEntry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}
	if rawEntry != nil {
		err = rawEntry.DecodeJSON(&config)
		if err != nil {
			return nil, err
		}
	}

	//if err := entry.DecodeJSON(config); err != nil {
	//	return nil, err
	//}
	fmt.Printf("CONFIG TO RETURN %v\n", config)
	return &config, nil
}

func (b *backend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, _ *framework.FieldData) (bool, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return false, err
	}
	return config != nil, nil
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	config.RamAlias = defaultRamAlias
	ramAlias, ok := data.GetOk("ram_alias")
	if ok {
		config.RamAlias = ramAlias.(string)
	}

	entry, err := logical.StorageEntryJSON("config", config)
	fmt.Printf("------------------> Config Write: %+v\n", config)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	fmt.Printf("------------------> Config READ: %+v\n", config)

	resp := &logical.Response{
		Data: map[string]interface{}{
			"ram_alias": config.RamAlias,
		},
	}

	return resp, nil
}

func (b *backend) pathConfigDelete(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "config")

	return nil, err
}

func (b *backend) saveConfig(ctx context.Context, config *alicloudConfig, s logical.Storage) error {
	entry, err := logical.StorageEntryJSON(configStoragePath, config)
	if err != nil {
		return err
	}

	err = s.Put(ctx, entry)
	if err != nil {
		return err
	}

	return nil
}

const (
	defaultRamAlias   = "principalId"
	configStoragePath = "config"
	confHelpSyn       = `Configures the Alicloud authentication backend.`
	confHelpDesc      = ""
)
