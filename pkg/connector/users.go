package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-vgs/pkg/client"
)

type userResourceType struct {
	resourceType *v2.ResourceType
	client       *client.VGSClient
}

func (u *userResourceType) ResourceType(ctx context.Context) *v2.ResourceType {
	return u.resourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (u *userResourceType) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var (
		pageToken string
		rv        []*v2.Resource
	)
	_, b, err := unmarshalSkipToken(pToken)
	if err != nil {
		return nil, "", nil, err
	}

	if b.Current() == nil {
		// Push onto stack in reverse
		b.Push(pagination.PageState{
			ResourceTypeID: "invites",
		})
		b.Push(pagination.PageState{
			ResourceTypeID: "users",
		})
	}

	switch b.Current().ResourceTypeID {
	case "users":
		users, err := u.client.ListUsers(ctx, u.client.GetOrganizationId(), u.client.GetVaultId())
		if err != nil {
			return nil, "", nil, fmt.Errorf("vgs-connector: failed to fetch users: %w", err)
		}

		for _, usr := range users {
			usrCopy := usr
			ur, err := getUserResource(&usrCopy, parentResourceID)
			if err != nil {
				return nil, "", nil, err
			}
			rv = append(rv, ur)
		}

		b.Pop()
		pageToken, err = marshalSkipToken(len(users), 0, b)
		if err != nil {
			return nil, "", nil, err
		}
	case "invites":
		userInvites, err := u.client.ListInvites(ctx, u.client.GetOrganizationId())
		if err != nil {
			return nil, "", nil, fmt.Errorf("vgs-connector: failed to fetch invites: %w", err)
		}

		for _, usr := range userInvites {
			usrCopy := usr
			ur, err := getUserResource(&usrCopy, parentResourceID)
			if err != nil {
				return nil, "", nil, err
			}
			rv = append(rv, ur)
		}

		return rv, "", nil, nil
	default:
		return nil, "", nil, fmt.Errorf("baton-vgs: unknown page state: %s", b.Current().ResourceTypeID)
	}

	return rv, pageToken, nil, nil
}

// Entitlements always returns an empty slice for users.
func (u *userResourceType) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (u *userResourceType) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func userBuilder(c *client.VGSClient) *userResourceType {
	return &userResourceType{
		resourceType: resourceTypeUser,
		client:       c,
	}
}
