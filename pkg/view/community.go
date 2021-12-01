package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/payload"
	"github.com/neoxelox/odin/pkg/usecase/community"
	"github.com/neoxelox/odin/pkg/usecase/invitation"
	"github.com/neoxelox/odin/pkg/usecase/user"
)

type CommunityView struct {
	class.View
	communityCreator  community.CreatorUsecase
	communityGetter   community.GetterUsecase
	communityLeaver   community.LeaverUsecase
	invitationCreator invitation.CreatorUsecase
}

func NewCommunityView(configuration internal.Configuration, logger core.Logger, communityCreator community.CreatorUsecase,
	communityGetter community.GetterUsecase, communityLeaver community.LeaverUsecase, invitationCreator invitation.CreatorUsecase) *CommunityView {
	return &CommunityView{
		View:              *class.NewView(configuration, logger),
		communityCreator:  communityCreator,
		communityGetter:   communityGetter,
		communityLeaver:   communityLeaver,
		invitationCreator: invitationCreator,
	}
}

func (self *CommunityView) PostCommunity(ctx echo.Context) error {
	request := &payload.PostCommunityRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostCommunityResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		newCommunity, newMembership, err := self.communityCreator.Create(ctx.Request().Context(), request.Address, request.Name, request.Categories, requestUser)
		switch {
		case err == nil:
			response.Community = payload.Community{
				ID:         newCommunity.ID,
				Address:    newCommunity.Address,
				Name:       newCommunity.Name,
				Categories: newCommunity.Categories,
				PinnedIDs:  newCommunity.PinnedIDs,
				CreatedAt:  newCommunity.CreatedAt,
			}
			response.Membership = payload.Membership{
				ID:          newMembership.ID,
				UserID:      newMembership.UserID,
				CommunityID: newMembership.CommunityID,
				Door:        newMembership.Door,
				Role:        newMembership.Role,
				CreatedAt:   newMembership.CreatedAt,
			}
			return ctx.JSON(http.StatusOK, response)
		case community.ErrInvalid().Is(err), community.ErrInvalidAddress().Is(err), community.ErrInvalidName().Is(err),
			community.ErrInvalidDoor().Is(err), community.ErrInvalidRole().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrAlreadyJoined().Is(err):
			return ExcUserAlreadyJoined.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *CommunityView) GetCommunity(ctx echo.Context) error {
	request := &payload.GetCommunityRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.GetCommunityResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		resCommunity, resMembership, err := self.communityGetter.Get(ctx.Request().Context(), *requestUser, request.ID)
		switch {
		case err == nil:
			response.Community = payload.Community{
				ID:         resCommunity.ID,
				Address:    resCommunity.Address,
				Name:       resCommunity.Name,
				Categories: resCommunity.Categories,
				PinnedIDs:  resCommunity.PinnedIDs,
				CreatedAt:  resCommunity.CreatedAt,
			}
			response.Membership = payload.Membership{
				ID:          resMembership.ID,
				UserID:      resMembership.UserID,
				CommunityID: resMembership.CommunityID,
				Door:        resMembership.Door,
				Role:        resMembership.Role,
				CreatedAt:   resMembership.CreatedAt,
			}
			return ctx.JSON(http.StatusOK, response)
		case community.ErrInvalid().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrNotBelongs().Is(err):
			return ExcUserNotBelongs.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *CommunityView) GetCommunityList(ctx echo.Context) error {
	requestUser := RequestUser(ctx)
	response := &payload.GetCommunityListResponse{Communities: []payload.CommunityAndMembership{}}
	return self.Handle(ctx, class.Endpoint{}, func() error {
		resCommunities, resMemberships, err := self.communityGetter.List(ctx.Request().Context(), *requestUser)
		switch {
		case err == nil:
			for i := 0; i < len(resCommunities); i++ {
				response.Communities = append(response.Communities, payload.CommunityAndMembership{
					Community: payload.Community{
						ID:         resCommunities[i].ID,
						Address:    resCommunities[i].Address,
						Name:       resCommunities[i].Name,
						Categories: resCommunities[i].Categories,
						PinnedIDs:  resCommunities[i].PinnedIDs,
						CreatedAt:  resCommunities[i].CreatedAt,
					},
					Membership: payload.Membership{
						ID:          resMemberships[i].ID,
						UserID:      resMemberships[i].UserID,
						CommunityID: resMemberships[i].CommunityID,
						Door:        resMemberships[i].Door,
						Role:        resMemberships[i].Role,
						CreatedAt:   resMemberships[i].CreatedAt,
					},
				})
			}
			return ctx.JSON(http.StatusOK, response)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *CommunityView) PostCommunityInvite(ctx echo.Context) error {
	request := &payload.PostCommunityInviteRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostCommunityInviteResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		resInvitation, err := self.invitationCreator.Create(ctx.Request().Context(), *requestUser, request.ID, request.Phone, request.Door, request.Role)
		switch {
		case err == nil:
			response.Invitation = payload.Invitation{
				ID:          resInvitation.ID,
				Phone:       resInvitation.Phone,
				CommunityID: resInvitation.CommunityID,
				Door:        resInvitation.Door,
				Role:        resInvitation.Role,
				CreatedAt:   resInvitation.CreatedAt,
			}
			return ctx.JSON(http.StatusOK, response)
		case user.ErrInvalidPhone().Is(err), community.ErrInvalidDoor().Is(err), community.ErrInvalidRole().Is(err),
			invitation.ErrInvitingYourself().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrNotBelongs().Is(err):
			return ExcUserNotBelongs.Cause(err)
		case community.ErrNotPermission().Is(err):
			return ExcUserNotPermission.Cause(err)
		case community.ErrAlreadyJoined().Is(err):
			return ExcUserAlreadyJoined.Cause(err)
		case invitation.ErrAlreadyInvited().Is(err):
			return ExcUserAlreadyInvited.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *CommunityView) PostCommunityLeave(ctx echo.Context) error {
	request := &payload.PostCommunityLeaveRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostCommunityLeaveResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		err := self.communityLeaver.Leave(ctx.Request().Context(), *requestUser, request.ID)
		switch {
		case err == nil:
			return ctx.JSON(http.StatusOK, response)
		case community.ErrInvalid().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrNotBelongs().Is(err):
			return ExcUserNotBelongs.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}
