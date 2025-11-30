package mappings

import (
	"github.com/reoden/go-NFT/pkg/mapper"
	userService "github.com/reoden/go-NFT/user/internal/shared/grpc/genproto"
	datamodel "github.com/reoden/go-NFT/user/internal/user/data/datamodels"
	dtoV1 "github.com/reoden/go-NFT/user/internal/user/dtos/v1"
	"github.com/reoden/go-NFT/user/internal/user/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConfigureUserMappings() error {
	err := mapper.CreateMap[*models.User, *dtoV1.UserDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*dtoV1.UserDto, *models.User]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*datamodel.UserDataModel, *models.User]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.User, *datamodel.UserDataModel]()
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap[*dtoV1.UserDto, *userService.User](
		func(user *dtoV1.UserDto) *userService.User {
			if user == nil {
				return nil
			}
			return &userService.User{
				Id:        user.Id,
				UserId:    user.UserId.String(),
				Nickname:  user.Nickname,
				Phone:     user.Phone,
				CreatedAt: timestamppb.New(user.CreatedAt),
				UpdatedAt: timestamppb.New(user.UpdatedAt),
			}
		},
	)
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap(
		func(user *models.User) *userService.User {
			return &userService.User{
				Id:        user.Id,
				UserId:    user.UserId.String(),
				Nickname:  user.Nickname,
				Phone:     user.Phone,
				CreatedAt: timestamppb.New(user.CreatedAt),
				UpdatedAt: timestamppb.New(user.UpdatedAt),
			}
		},
	)

	return nil
}
