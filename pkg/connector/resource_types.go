package connector

import v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"

var (
	resourceTypeUser = &v2.ResourceType{
		Id:          "user",
		DisplayName: "User",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
		Annotations: annotationsForUserResourceType(),
	}
	resourceTypeOrg = &v2.ResourceType{
		Id:          "org",
		DisplayName: "Org",
		Annotations: v1AnnotationsForResourceType("org"),
	}
	resourceTypeVault = &v2.ResourceType{
		Id:          "vault",
		DisplayName: "Vault",
		Annotations: v1AnnotationsForResourceType("vault"),
	}
)
