package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-vgs/pkg/client"
)

type orgResourceType struct {
	resourceType *v2.ResourceType
	client       *client.VGSClient
}

const (
	orgRoleMember = "member"
	orgRoleAdmin  = "admin"
)

var orgAccessLevels = []string{
	orgRoleMember,
	orgRoleAdmin,
}

func (o *orgResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return o.resourceType
}

// List returns all the organizations from the database as resource objects.
func (o *orgResourceType) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var ret []*v2.Resource
	orgs, err := o.client.ListOrganizations(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("vgs-connector: failed to fetch org: %w", err)
	}

	for _, org := range orgs {
		orgResource, err := rs.NewResource(
			org.Name,
			resourceTypeOrg,
			org.Id,
			rs.WithParentResourceID(parentResourceID),
			rs.WithAnnotation(
				&v2.ExternalLink{Url: org.Name},
				&v2.V1Identifier{Id: fmt.Sprintf("org:%s", org.Id)},
				&v2.ChildResourceType{ResourceTypeId: resourceTypeUser.Id},
			),
		)

		if err != nil {
			return nil, "", nil, err
		}

		ret = append(ret, orgResource)
	}

	return ret, "", nil, nil
}

func (o *orgResourceType) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	rv := make([]*v2.Entitlement, 0, len(orgAccessLevels))
	for _, level := range orgAccessLevels {
		rv = append(rv, ent.NewPermissionEntitlement(resource, level,
			ent.WithDisplayName(fmt.Sprintf("%s Organization %s", resource.DisplayName, titleCase(level))),
			ent.WithDescription(fmt.Sprintf("Access to %s organization in VGS", resource.DisplayName)),
			ent.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("org:%s:role:%s", resource.Id.Resource, level),
			}),
			ent.WithGrantableTo(resourceTypeUser),
		))
	}

	return rv, "", nil, nil
}

func (o *orgResourceType) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *orgResourceType) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	return nil, nil
}

func (o *orgResourceType) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	return nil, nil
}

func orgBuilder(c *client.VGSClient) *orgResourceType {
	return &orgResourceType{
		resourceType: resourceTypeOrg,
		client:       c,
	}
}
