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
)

type InvitationView struct {
	class.View
	invitationGetter   invitation.GetterUsecase
	invitationAccepter invitation.AccepterUsecase
	invitationRejecter invitation.RejecterUsecase
}

func NewInvitationView(configuration internal.Configuration, logger core.Logger, invitationGetter invitation.GetterUsecase,
	invitationAccepter invitation.AccepterUsecase, invitationRejecter invitation.RejecterUsecase) *InvitationView {
	return &InvitationView{
		View:               *class.NewView(configuration, logger),
		invitationGetter:   invitationGetter,
		invitationAccepter: invitationAccepter,
		invitationRejecter: invitationRejecter,
	}
}

func (self *InvitationView) GetInvitationList(ctx echo.Context) error {
	requestUser := RequestUser(ctx)
	response := &payload.GetInvitationListResponse{Invitations: []payload.Invitation{}}
	return self.Handle(ctx, class.Endpoint{}, func() error {
		resInvitations, err := self.invitationGetter.List(ctx.Request().Context(), *requestUser)
		switch {
		case err == nil:
			for _, resInvitation := range resInvitations {
				response.Invitations = append(response.Invitations, payload.Invitation{
					ID:          resInvitation.ID,
					Phone:       resInvitation.Phone,
					CommunityID: resInvitation.CommunityID,
					Door:        resInvitation.Door,
					Role:        resInvitation.Role,
					CreatedAt:   resInvitation.CreatedAt,
				})
			}
			return ctx.JSON(http.StatusOK, response)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *InvitationView) PostInvitationAccept(ctx echo.Context) error {
	request := &payload.PostInvitationAcceptRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostInvitationAcceptResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		resMembership, err := self.invitationAccepter.Accept(ctx.Request().Context(), *requestUser, request.ID)
		switch {
		case err == nil:
			response.Membership = payload.Membership{
				ID:          resMembership.ID,
				UserID:      resMembership.UserID,
				CommunityID: resMembership.CommunityID,
				Door:        resMembership.Door,
				Role:        resMembership.Role,
				CreatedAt:   resMembership.CreatedAt,
			}
			return ctx.JSON(http.StatusOK, response)
		case invitation.ErrInvalid().Is(err), community.ErrInvalid().Is(err), community.ErrInvalidDoor().Is(err),
			community.ErrInvalidRole().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrAlreadyJoined().Is(err):
			return ExcUserAlreadyJoined.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *InvitationView) PostInvitationReject(ctx echo.Context) error {
	request := &payload.PostInvitationRejectRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostInvitationRejectResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		err := self.invitationRejecter.Reject(ctx.Request().Context(), *requestUser, request.ID)
		switch {
		case err == nil:
			return ctx.JSON(http.StatusOK, response)
		case invitation.ErrInvalid().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}
