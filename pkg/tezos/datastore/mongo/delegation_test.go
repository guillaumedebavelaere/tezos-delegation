package mongo_test

import (
	"context"
	"errors"
	"time"

	"github.com/guillaumedebavelaere/tezos-delegation/pkg/tezos/datastore/model"
)

var errBulkWrite = errors.New("must provide at least one element in input slice")

func (suite *MongoTestSuite) TestDatastore_StoreDelegations() {
	cases := []struct {
		name    string
		init    func(ctx context.Context)
		want    []*model.Delegation
		wantErr error
	}{
		{
			name: "Success create",
			init: func(ctx context.Context) {},
			want: []*model.Delegation{
				{
					Timestamp: time.Date(2023, 12, 10, 11, 1, 1, 0, time.UTC),
					Amount:    124428330,
					Delegator: "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx",
					Block:     "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG",
				},
			},
		},
		{
			name: "Success update",
			init: func(ctx context.Context) {
				err := suite.mongoSvc.StoreDelegations(ctx, []*model.Delegation{
					{
						Timestamp: time.Date(2023, 12, 10, 11, 1, 1, 0, time.UTC),
						Amount:    124428330,
						Delegator: "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx",
						Block:     "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG",
					},
					{
						Timestamp: time.Date(2023, 12, 10, 11, 0, 1, 0, time.UTC),
						Amount:    499836,
						Delegator: "tz1NqVXDBf8fZNomacychFPSK1trQbi14PvA",
						Block:     "BLxQGrPcAPAwKaeCdivBVw45Choicesen6wrmdm3NBeGsCnkLKv",
					},
				})
				suite.Require().Nil(err)
			},
			want: []*model.Delegation{
				{
					Timestamp: time.Date(2023, 12, 10, 11, 1, 1, 0, time.UTC),
					Amount:    124428330,
					Delegator: "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx",
					Block:     "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG",
				},
			},
		},
		{
			name:    "Error BulkWrite",
			init:    func(ctx context.Context) {},
			want:    nil,
			wantErr: errBulkWrite,
		},
	}

	for _, c := range cases {
		suite.Run(c.name, func() {
			suite.SetupTest()
			defer suite.TearDownTest()

			ctx := context.Background()

			c.init(ctx)

			err := suite.mongoSvc.StoreDelegations(ctx, c.want)
			if c.wantErr != nil {
				suite.Require().Equal(c.wantErr, err)
			} else {
				suite.Require().Nil(err)

				latestDelegation, err := suite.mongoSvc.GetLatestDelegation(ctx)
				suite.Require().Equal(c.want[0], latestDelegation)
				suite.Require().Nil(err)
			}
		})
	}
}

func (suite *MongoTestSuite) TestDatastore_GetDelegations() {
	cases := []struct {
		name string
		init func(ctx context.Context)
		want []*model.Delegation
		year int
	}{
		{
			name: "Success empty",
			init: func(ctx context.Context) {},
			want: nil,
		},
		{
			name: "Success",
			init: func(ctx context.Context) {
				err := suite.mongoSvc.StoreDelegations(ctx, []*model.Delegation{
					{
						Timestamp: time.Date(2023, 12, 10, 11, 1, 1, 0, time.UTC),
						Amount:    124428330,
						Delegator: "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx",
						Block:     "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG",
					},
					{
						Timestamp: time.Date(2021, 12, 10, 11, 0, 1, 0, time.UTC),
						Amount:    499836,
						Delegator: "tz1NqVXDBf8fZNomacychFPSK1trQbi14PvA",
						Block:     "BLxQGrPcAPAwKaeCdivBVw45Choicesen6wrmdm3NBeGsCnkLKv",
					},
				})
				suite.Require().Nil(err)
			},
			want: []*model.Delegation{
				{
					Timestamp: time.Date(2023, 12, 10, 11, 1, 1, 0, time.UTC),
					Amount:    124428330,
					Delegator: "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx",
					Block:     "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG",
				},
				{
					Timestamp: time.Date(2021, 12, 10, 11, 0, 1, 0, time.UTC),
					Amount:    499836,
					Delegator: "tz1NqVXDBf8fZNomacychFPSK1trQbi14PvA",
					Block:     "BLxQGrPcAPAwKaeCdivBVw45Choicesen6wrmdm3NBeGsCnkLKv",
				},
			},
		},
		{
			name: "Success with year",
			init: func(ctx context.Context) {
				err := suite.mongoSvc.StoreDelegations(ctx, []*model.Delegation{
					{
						Timestamp: time.Date(2023, 12, 10, 11, 1, 1, 0, time.UTC),
						Amount:    124428330,
						Delegator: "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx",
						Block:     "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG",
					},
					{
						Timestamp: time.Date(2021, 12, 10, 11, 0, 1, 0, time.UTC),
						Amount:    499836,
						Delegator: "tz1NqVXDBf8fZNomacychFPSK1trQbi14PvA",
						Block:     "BLxQGrPcAPAwKaeCdivBVw45Choicesen6wrmdm3NBeGsCnkLKv",
					},
				})
				suite.Require().Nil(err)
			},
			year: 2021,
			want: []*model.Delegation{
				{
					Timestamp: time.Date(2021, 12, 10, 11, 0, 1, 0, time.UTC),
					Amount:    499836,
					Delegator: "tz1NqVXDBf8fZNomacychFPSK1trQbi14PvA",
					Block:     "BLxQGrPcAPAwKaeCdivBVw45Choicesen6wrmdm3NBeGsCnkLKv",
				},
			},
		},
	}

	for _, c := range cases {
		suite.Run(c.name, func() {
			suite.SetupTest()
			defer suite.TearDownTest()

			ctx := context.Background()

			c.init(ctx)

			result, err := suite.mongoSvc.GetDelegations(ctx, c.year)
			suite.Require().Equal(c.want, result)
			suite.Require().Nil(err)
		})
	}
}

func (suite *MongoTestSuite) TestDatastore_GetLatestDelegation() {
	cases := []struct {
		name string
		init func(ctx context.Context)
		want *model.Delegation
	}{
		{
			name: "Success",
			init: func(ctx context.Context) {
				err := suite.mongoSvc.StoreDelegations(ctx, []*model.Delegation{
					{
						Timestamp: time.Date(2023, 12, 10, 11, 1, 1, 0, time.UTC),
						Amount:    124428330,
						Delegator: "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx",
						Block:     "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG",
					},
					{
						Timestamp: time.Date(2021, 12, 10, 11, 0, 1, 0, time.UTC),
						Amount:    499836,
						Delegator: "tz1NqVXDBf8fZNomacychFPSK1trQbi14PvA",
						Block:     "BLxQGrPcAPAwKaeCdivBVw45Choicesen6wrmdm3NBeGsCnkLKv",
					},
				})
				suite.Require().Nil(err)
			},
			want: &model.Delegation{
				Timestamp: time.Date(2023, 12, 10, 11, 1, 1, 0, time.UTC),
				Amount:    124428330,
				Delegator: "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx",
				Block:     "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG",
			},
		},
	}

	for _, c := range cases {
		suite.Run(c.name, func() {
			suite.SetupTest()
			defer suite.TearDownTest()

			ctx := context.Background()

			c.init(ctx)

			result, err := suite.mongoSvc.GetLatestDelegation(ctx)
			suite.Require().Equal(c.want, result)
			suite.Require().Nil(err)
		})
	}
}
