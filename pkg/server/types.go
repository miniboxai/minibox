package server

import (
	"database/sql"
	"strings"
	"time"

	"minibox.ai/pkg/api/v1/types"
	"minibox.ai/pkg/models"
	"minibox.ai/pkg/server/internal/storage"
)

// 类型转换表, types 到 models 的转换
func prjType2Model(typ *types.Project) *models.Project {
	var createdAt time.Time
	if typ.CreatedAt != nil {
		createdAt = *typ.CreatedAt
	}

	var prj = &models.Project{
		Model: models.Model{
			ID:        typ.ID,
			CreatedAt: createdAt,
		},
		Name:      typ.Name,
		Private:   typ.Private,
		Namespace: typ.Namespace,
	}

	if typ.Author != nil {
		prj.UserID = uint(typ.Author.ID)
	}

	return prj
}

func prjModel2Type(mdl *models.Project) *types.Project {
	var prj = &types.Project{
		ID:        mdl.ID,
		Name:      mdl.Name,
		Namespace: mdl.Namespace,
		Private:   mdl.Private,
		CreatedAt: Time(mdl.CreatedAt),
	}

	prj.Author = usrModel2Type(&mdl.User)

	return prj
}

func usrType2Model(typ *types.User) *models.User {
	var createdAt, updatedAt time.Time

	if typ.CreatedAt != nil {
		createdAt = *typ.CreatedAt
	}

	if typ.UpdatedAt != nil {
		updatedAt = *typ.UpdatedAt
	}
	var usr = &models.User{
		ModelByInt: models.ModelByInt{
			ID:        uint(typ.ID),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Name:      typ.Name,
		Namespace: typ.Namespace,
	}

	return usr
}

func usrModel2Type(mdl *models.User) *types.User {
	var usr = &types.User{
		ID:        uint32(mdl.ID),
		Name:      mdl.Name,
		Namespace: mdl.Namespace,
		CreatedAt: Time(mdl.CreatedAt),
		UpdatedAt: Time(mdl.UpdatedAt),
	}

	if len(mdl.Projects) > 0 {
		usr.Projects = make([]*types.Project, len(mdl.Projects))

		for i, prj := range mdl.Projects {
			usr.Projects[i] = prjModel2Type(&prj)
		}
	}

	// usr.Author = &types.User{
	// 	ID: uint32(mdl.UserID),
	// }

	return usr
}
func orgType2Model(typ *types.Organization) *models.Organization {
	var createdAt, updatedAt time.Time
	if typ.CreatedAt != nil {
		createdAt = *typ.CreatedAt
	}

	if typ.UpdatedAt != nil {
		updatedAt = *typ.UpdatedAt
	}

	var ds = &models.Organization{
		ModelByInt: models.ModelByInt{
			ID:        uint(typ.ID),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Name:        typ.Name,
		Namespace:   typ.Namespace,
		Description: typ.Description,
		AvatarURL:   typ.AvatarUrl,
		SiteURL:     typ.SiteUrl,
		HomeURL:     typ.HomeUrl,
	}

	if typ.Owner != nil {
		ds.OwnerID = uint(typ.Owner.ID)
	}

	return ds
}

func dsType2Model(typ *types.Dataset) *models.Dataset {
	var createdAt, publishAt, modifyAt time.Time
	if typ.CreatedAt != nil {
		createdAt = *typ.CreatedAt
	}

	if typ.PublishAt != nil {
		publishAt = *typ.PublishAt
	}

	if typ.ModifyAt != nil {
		modifyAt = *typ.ModifyAt
	}

	var ds = &models.Dataset{
		ModelByInt: models.ModelByInt{
			ID:        uint(typ.ID),
			CreatedAt: createdAt,
			UpdatedAt: modifyAt,
		},
		Name:          typ.Name,
		Namespace:     typ.Namespace,
		Private:       sql.NullBool{Bool: typ.Private, Valid: true},
		Description:   typ.Description,
		Preprocessing: typ.Preprocessing,
		Instances:     int(typ.Instances),
		Format:        strings.Join(typ.Format, ","),
		Lisense:       typ.Lisense,
		PublishAt:     publishAt,
	}

	if typ.Author != nil {
		ds.UserID = uint(typ.Author.ID)
	}

	return ds
}

func dsModel2Type(mdl *models.Dataset) *types.Dataset {
	var ds = &types.Dataset{
		ID:            uint32(mdl.ID),
		Name:          mdl.Name,
		Namespace:     mdl.Namespace,
		Private:       mdl.Private.Bool,
		CreatedAt:     Time(mdl.CreatedAt),
		Description:   mdl.Description,
		Preprocessing: mdl.Preprocessing,
		Instances:     int32(mdl.Instances),
		Format:        strings.Split(mdl.Format, ","),
		Lisense:       mdl.Lisense,
		PublishAt:     Time(mdl.PublishAt),
		ModifyAt:      Time(mdl.UpdatedAt),
	}

	ds.Author = usrModel2Type(&mdl.User)

	return ds
}

func orgModel2Type(mdl *models.Organization) *types.Organization {
	var ds = &types.Organization{
		ID:          uint32(mdl.ID),
		Name:        mdl.Name,
		Namespace:   mdl.Namespace,
		CreatedAt:   Time(mdl.CreatedAt),
		UpdatedAt:   Time(mdl.UpdatedAt),
		Description: mdl.Description,
		SiteUrl:     mdl.SiteURL,
		HomeUrl:     mdl.SiteURL,
		AvatarUrl:   mdl.AvatarURL,
	}

	ds.Owner = usrModel2Type(&mdl.Owner)

	return ds
}

func init() {
	storage.RegisterConvert(&types.Project{}, &models.Project{}, prjType2Model, prjModel2Type)
	storage.RegisterConvert(&types.User{}, &models.User{}, usrType2Model, usrModel2Type)
	storage.RegisterConvert(&types.Organization{}, &models.Organization{}, orgType2Model, orgModel2Type)
	storage.RegisterConvert(&types.Dataset{}, &models.Dataset{}, dsType2Model, dsModel2Type)

}
