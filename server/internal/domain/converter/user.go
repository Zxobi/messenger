package converter

import (
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/cutils"
	protocolv1 "github.com/dvid-messanger/protos/gen/protocol"
)

func UserToDTO(usr *model.User) *protocolv1.User {
	return &protocolv1.User{
		Id:    usr.Id,
		Email: usr.Email,
		Bio:   usr.Bio,
	}
}

func UserFromDTO(usr *protocolv1.User) *model.User {
	return &model.User{
		Id:    usr.Id,
		Email: usr.Email,
		Bio:   usr.Bio,
	}
}

func UsersToDTO(users []model.User) []*protocolv1.User {
	return cutils.Map(users, func(usr model.User) *protocolv1.User {
		return UserToDTO(&usr)
	})
}

func UsersFromDTO(users []*protocolv1.User) []model.User {
	return cutils.Map(users, func(usr *protocolv1.User) model.User {
		return *UserFromDTO(usr)
	})
}
