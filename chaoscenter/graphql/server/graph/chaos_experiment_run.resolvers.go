package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/authorization"
	data_store "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/data-store"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) ChaosExperimentRun(ctx context.Context, request model.ExperimentRunRequest) (string, error) {
	return r.chaosExperimentRunHandler.ChaosExperimentRunEvent(request)
}

func (r *mutationResolver) RunChaosExperiment(ctx context.Context, experimentID string, projectID string) (*model.RunChaosExperimentResponse, error) {
	logFields := logrus.Fields{
		"projectId":         projectID,
		"chaosExperimentId": experimentID,
	}

	logrus.WithFields(logFields).Info("request received to run chaos experiment")
	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.CreateChaosWorkFlow],
		model.InvitationAccepted.String())
	if err != nil {
		return nil, err
	}

	query := bson.D{
		{"experiment_id", experimentID},
		{"is_removed", false},
	}

	experiment, err := r.chaosExperimentHandler.GetDBExperiment(query)
	if err != nil {
		return nil, errors.New("could not get experiment run, error: " + err.Error())
	}

	var uiResponse *model.RunChaosExperimentResponse

	uiResponse, err = r.chaosExperimentRunHandler.RunChaosWorkFlow(ctx, projectID, experiment, data_store.Store)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return nil, err
	}

	return &model.RunChaosExperimentResponse{NotifyID: uiResponse.NotifyID}, err
}

func (r *queryResolver) GetExperimentRun(ctx context.Context, projectID string, experimentRunID *string, notifyID *string) (*model.ExperimentRun, error) {
	logFields := logrus.Fields{
		"projectId":            projectID,
		"chaosExperimentRunId": experimentRunID,
	}
	logrus.WithFields(logFields).Info("request received to fetch chaos experiment run")
	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.GetWorkflowRun],
		model.InvitationAccepted.String())
	if err != nil {
		return nil, err
	}

	expRunResponse, err := r.chaosExperimentRunHandler.GetExperimentRun(ctx, projectID, experimentRunID, notifyID)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return nil, err
	}
	return expRunResponse, err
}

func (r *queryResolver) ListExperimentRun(ctx context.Context, projectID string, request model.ListExperimentRunRequest) (*model.ListExperimentRunResponse, error) {
	logFields := logrus.Fields{
		"projectId":             projectID,
		"chaosExperimentIds":    request.ExperimentIDs,
		"chaosExperimentRunIds": request.ExperimentRunIDs,
	}
	logrus.WithFields(logFields).Info("request received to list chaos experiment run")

	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.ListWorkflowRuns],
		model.InvitationAccepted.String())
	if err != nil {
		return nil, err
	}
	uiResponse, err := r.chaosExperimentRunHandler.ListExperimentRun(projectID, request)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return nil, err
	}
	return uiResponse, err
}

func (r *queryResolver) GetExperimentRunStats(ctx context.Context, projectID string) (*model.GetExperimentRunStatsResponse, error) {
	logFields := logrus.Fields{
		"projectId": projectID,
	}
	logrus.WithFields(logFields).Info("request received to get chaos experiment run stats")
	err := authorization.ValidateRole(ctx, projectID,
		authorization.MutationRbacRules[authorization.ListWorkflowRuns],
		model.InvitationAccepted.String())
	if err != nil {
		return nil, err
	}

	uiResponse, err := r.chaosExperimentRunHandler.GetExperimentRunStats(ctx, projectID)
	if err != nil {
		logrus.WithFields(logFields).Error(err)
		return nil, err
	}
	return uiResponse, err
}
