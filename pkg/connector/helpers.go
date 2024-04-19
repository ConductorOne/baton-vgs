package connector

import (
	"fmt"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-vgs/pkg/client"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func annotationsForUserResourceType() annotations.Annotations {
	annos := annotations.Annotations{}
	annos.Update(&v2.SkipEntitlementsAndGrants{})
	return annos
}

func titleCase(s string) string {
	titleCaser := cases.Title(language.English)

	return titleCaser.String(s)
}

// Populate entitlement options for a 1password resource.
func PopulateOptions(displayName, permission, resource string) []ent.EntitlementOption {
	options := []ent.EntitlementOption{
		ent.WithGrantableTo(resourceTypeUser),
		ent.WithDescription(fmt.Sprintf("%s of VGS %s %s", permission, displayName, resource)),
		ent.WithDisplayName(fmt.Sprintf("%s %s %s", displayName, resource, permission)),
	}
	return options
}

func v1AnnotationsForResourceType(resourceTypeID string) annotations.Annotations {
	annos := annotations.Annotations{}
	annos.Update(&v2.V1Identifier{
		Id: resourceTypeID,
	})

	return annos
}

// splitFullName returns firstName and lastName.
func splitFullName(name string) (string, string) {
	names := strings.SplitN(name, " ", 2)
	var firstName, lastName string

	switch len(names) {
	case 1:
		firstName = names[0]
	case 2:
		firstName = names[0]
		lastName = names[1]
	}

	return firstName, lastName
}

func getUserResource(user *client.OrganizationUser) (*v2.Resource, error) {
	var userStatus v2.UserTrait_Status_Status = v2.UserTrait_Status_STATUS_ENABLED
	firstName, lastName := splitFullName(user.Name)
	profile := map[string]interface{}{
		"login":      user.Email,
		"first_name": firstName,
		"last_name":  lastName,
		"email":      user.Email,
	}

	userTraits := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
		rs.WithUserLogin(user.Email),
		rs.WithEmail(user.Email, true),
	}

	displayName := user.Name
	if user.Name == "" {
		displayName = user.Email
	}

	ret, err := rs.NewUserResource(displayName, resourceTypeUser, user.Id, userTraits)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
