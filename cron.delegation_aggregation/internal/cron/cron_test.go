package cron_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/cron"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/tezos"
	tezosmock "github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/tezos/mock"
	datastoremock "github.com/guillaumedebavelaere/tezos-delegation/pkg/tezos/datastore/mock"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/tezos/datastore/model"
)

type underTest struct {
	mockCtrl         *gomock.Controller
	mockTezosService *tezosmock.MockAPI
	mockDatastore    *datastoremock.MockDatastorer
	cron             *cron.Cron
}

func setupTest(t *testing.T) *underTest {
	t.Helper()

	ut := &underTest{}

	ut.mockCtrl = gomock.NewController(t)

	ut.mockTezosService = tezosmock.NewMockAPI(ut.mockCtrl)
	ut.mockDatastore = datastoremock.NewMockDatastorer(ut.mockCtrl)

	ut.cron = cron.New(
		ut.mockTezosService,
		ut.mockDatastore,
	)

	return ut
}

func TestCron_New(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ut := setupTest(t)
		assert.NotNil(t, ut.cron)
	})
}

var errAny = errors.New("any error")

func TestCron_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		init    func(*underTest)
		wantErr error
	}{
		{
			name: "Success first run",
			init: func(ut *underTest) {
				// No delegation in datastore
				getLatestDelegation := ut.mockDatastore.EXPECT().GetLatestDelegation(gomock.Any()).
					Return(nil, nil)
				listDelegations := ut.mockTezosService.EXPECT().ListDelegations(
					gomock.Any(),
					gomock.Nil(),
				).After(getLatestDelegation).Return([]*tezos.Delegation{
					{
						Timestamp: time.Date(2023, 1, 1, 17, 0, 0, 0, time.UTC),
						Amount:    100,
						Block:     "block2",
						Sender: tezos.Sender{
							Address: "tz2",
						},
					},
					{
						Timestamp: time.Date(2023, 1, 1, 16, 0, 0, 0, time.UTC),
						Amount:    100,
						Block:     "block1",
						Sender: tezos.Sender{
							Address: "tz1",
						},
					},
				}, nil)
				ut.mockDatastore.EXPECT().StoreDelegations(
					gomock.Any(),
					gomock.Eq(
						[]*model.Delegation{
							{
								Delegator: "tz2",
								Block:     "block2",
								Amount:    100,
								Timestamp: time.Date(
									2023,
									1,
									1,
									17,
									0,
									0,
									0,
									time.UTC,
								),
							},
							{
								Delegator: "tz1",
								Block:     "block1",
								Amount:    100,
								Timestamp: time.Date(
									2023,
									1,
									1,
									16,
									0,
									0,
									0,
									time.UTC,
								),
							},
						}),
				).After(listDelegations).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "Success second run",
			init: func(ut *underTest) {
				lastTimestamp := time.Date(
					2023,
					1,
					1,
					16,
					0,
					0,
					0,
					time.UTC,
				)

				// One delegation in datastore
				getLatestDelegation := ut.mockDatastore.EXPECT().GetLatestDelegation(gomock.Any()).
					Return(&model.Delegation{
						Delegator: "tz1",
						Block:     "block1",
						Amount:    100,
						Timestamp: lastTimestamp,
					}, nil)
				listDelegations := ut.mockTezosService.EXPECT().ListDelegations(
					gomock.Any(),
					gomock.Eq(&lastTimestamp),
				).After(getLatestDelegation).Return([]*tezos.Delegation{
					{
						Timestamp: time.Date(2023, 1, 1, 17, 0, 0, 0, time.UTC),
						Amount:    100,
						Block:     "block2",
						Sender: tezos.Sender{
							Address: "tz2",
						},
					},
				}, nil)
				ut.mockDatastore.EXPECT().StoreDelegations(
					gomock.Any(),
					gomock.Eq(
						[]*model.Delegation{
							{
								Delegator: "tz2",
								Block:     "block2",
								Amount:    100,
								Timestamp: time.Date(
									2023,
									1,
									1,
									17,
									0,
									0,
									0,
									time.UTC,
								),
							},
						}),
				).After(listDelegations).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "Success second run with no results from tezos service",
			init: func(ut *underTest) {
				lastTimestamp := time.Date(
					2023,
					1,
					1,
					16,
					0,
					0,
					0,
					time.UTC,
				)

				// One delegation in datastore
				getLatestDelegation := ut.mockDatastore.EXPECT().GetLatestDelegation(gomock.Any()).
					Return(&model.Delegation{
						Delegator: "tz1",
						Block:     "block1",
						Amount:    100,
						Timestamp: lastTimestamp,
					}, nil)
				ut.mockTezosService.EXPECT().ListDelegations(
					gomock.Any(),
					gomock.Eq(&lastTimestamp),
				).After(getLatestDelegation).Return([]*tezos.Delegation{}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Error GetLatestDelegation",
			init: func(ut *underTest) {
				ut.mockDatastore.EXPECT().GetLatestDelegation(gomock.Any()).
					Return(nil, errAny)
			},
			wantErr: errAny,
		},
		{
			name: "Error ListDelegations",
			init: func(ut *underTest) {
				latestDelegation := ut.mockDatastore.EXPECT().GetLatestDelegation(gomock.Any()).
					Return(nil, nil)
				ut.mockTezosService.EXPECT().ListDelegations(
					gomock.Any(),
					gomock.Nil(),
				).After(latestDelegation).Return(nil, errAny)
			},
			wantErr: errAny,
		},
		{
			name: "Error StoreDelegations",
			init: func(ut *underTest) {
				latestDelegation := ut.mockDatastore.EXPECT().GetLatestDelegation(gomock.Any()).
					Return(nil, nil)
				listDelegations := ut.mockTezosService.EXPECT().ListDelegations(
					gomock.Any(),
					gomock.Nil(),
				).After(latestDelegation).Return([]*tezos.Delegation{
					{
						Timestamp: time.Date(2023, 1, 1, 17, 0, 0, 0, time.UTC),
						Amount:    100,
						Block:     "block2",
						Sender: tezos.Sender{
							Address: "tz2",
						},
					},
					{
						Timestamp: time.Date(2023, 1, 1, 16, 0, 0, 0, time.UTC),
						Amount:    100,
						Block:     "block1",
						Sender: tezos.Sender{
							Address: "tz1",
						},
					},
				}, nil)
				ut.mockDatastore.EXPECT().StoreDelegations(
					gomock.Any(),
					gomock.Eq(
						[]*model.Delegation{
							{
								Delegator: "tz2",
								Block:     "block2",
								Amount:    100,
								Timestamp: time.Date(
									2023,
									1,
									1,
									17,
									0,
									0,
									0,
									time.UTC,
								),
							},
							{
								Delegator: "tz1",
								Block:     "block1",
								Amount:    100,
								Timestamp: time.Date(
									2023,
									1,
									1,
									16,
									0,
									0,
									0,
									time.UTC,
								),
							},
						}),
				).After(listDelegations).Return(errAny)
			},
			wantErr: errAny,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			ut := setupTest(t)
			c.init(ut)

			assert.Equal(t, c.wantErr, ut.cron.Run())
		})
	}
}
