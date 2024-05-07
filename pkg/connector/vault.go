package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-vgs/pkg/client"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type vaultResourceType struct {
	resourceType *v2.ResourceType
	client       *client.VGSClient
}

const (
	vaultRoleWrite = "write"
	vaultRoleAdmin = "admin"
)

var vaultAccessLevels = []string{
	vaultRoleWrite,
	vaultRoleAdmin,
}

func (v *vaultResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return v.resourceType
}

// List returns all the vaults from the database as resource objects.
func (v *vaultResourceType) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var ret []*v2.Resource
	vaults, err := v.client.ListVaults(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("vgs-connector: failed to fetch vault: %w", err)
	}

	for _, vault := range vaults {
		vaultResource, err := rs.NewResource(
			vault.Name,
			resourceTypeVault,
			vault.Id,
			rs.WithParentResourceID(parentResourceID),
			rs.WithAnnotation(
				&v2.ExternalLink{Url: vault.Name},
				&v2.V1Identifier{Id: fmt.Sprintf("vault:%s", vault.Id)},
				&v2.ChildResourceType{ResourceTypeId: resourceTypeVault.Id},
			),
		)

		if err != nil {
			return nil, "", nil, err
		}

		ret = append(ret, vaultResource)
	}

	return ret, "", nil, nil
}

func (v *vaultResourceType) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	rv := make([]*v2.Entitlement, 0, len(vaultAccessLevels))
	for _, level := range vaultAccessLevels {
		rv = append(rv, ent.NewPermissionEntitlement(resource, level,
			ent.WithDisplayName(fmt.Sprintf("%s Vault %s", resource.DisplayName, titleCase(level))),
			ent.WithDescription(fmt.Sprintf("Access to %s vault in VGS", resource.DisplayName)),
			ent.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("vault:%s:role:%s", resource.Id.Resource, level),
			}),
			ent.WithGrantableTo(resourceTypeVault),
		))
	}

	return rv, "", nil, nil
}

func (v *vaultResourceType) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var (
		err error
		rv  []*v2.Grant
	)
	users, err := v.client.ListVaultUsers(ctx, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, err
	}

	for _, usr := range users {
		userCopy := &client.OrganizationUser{
			Id:    usr.Attributes.Id,
			Name:  usr.Attributes.Email,
			Type:  "users",
			Email: usr.Attributes.Email,
		}
		ur, err := getUserResource(userCopy, resource.Id)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error creating user resource for role %s: %w", resource.Id.Resource, err)
		}

		gr := grant.NewGrant(resource, usr.Attributes.Role, ur.Id)
		rv = append(rv, gr)
	}

	return rv, "", nil, nil
}

func (v *vaultResourceType) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	var role string
	l := ctxzap.Extract(ctx)
	if principal.Id.ResourceType != resourceTypeUser.Id {
		l.Warn(
			"baton-vgs: only users can be granted role membership",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
		return nil, fmt.Errorf("baton-vgs: only users can be granted role membership")
	}

	_, parts, err := parseEntitlementID(entitlement.Id)
	if err != nil {
		return nil, err
	}

	role = parts[len(parts)-1]
	err = v.client.UpdateVault(ctx,
		entitlement.Resource.Id.Resource,
		principal.Id.Resource,
		role)
	if err != nil {
		return nil, err
	}

	l.Warn("Role Membership has been added.",
		zap.String("vaultIdentifier", entitlement.Resource.Id.Resource),
		zap.String("userId", principal.Id.Resource),
	)

	return nil, nil
}

func (v *vaultResourceType) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	return nil, nil
}

func vaultBuilder(c *client.VGSClient) *vaultResourceType {
	return &vaultResourceType{
		resourceType: resourceTypeVault,
		client:       c,
	}
}
