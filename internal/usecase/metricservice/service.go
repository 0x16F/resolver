package metricservice

import (
	"context"

	"github.com/0x16f/vpn-resolver/internal/entity"
	"github.com/sirupsen/logrus"
)

type metricSrv interface {
	CreateMetric(ctx context.Context, req entity.MetricCreateReq) error
}

type Service struct {
	metricSrv metricSrv
}

func New(metricSrv metricSrv) *Service {
	return &Service{
		metricSrv: metricSrv,
	}
}

func (s *Service) CreateMetric(ctx context.Context, req entity.MetricCreateReq) error {
	err := s.metricSrv.CreateMetric(ctx, req)
	if err != nil {
		logrus.Errorf("failed to create metric: %v", err)

		return err
	}

	return nil
}
