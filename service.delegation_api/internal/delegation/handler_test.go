package delegation_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	datastoremock "github.com/guillaumedebavelaere/tezos-delegation/pkg/tezos/datastore/mock"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/tezos/datastore/model"
	"github.com/guillaumedebavelaere/tezos-delegation/service.delegation_api/internal/delegation"
)

type underTest struct {
	mockCtrl      *gomock.Controller
	mockDatastore *datastoremock.MockDatastorer
	apiHandler    *delegation.APIHandler
}

func setupTest(t *testing.T) *underTest {
	t.Helper()

	ut := &underTest{}

	ut.mockCtrl = gomock.NewController(t)

	ut.mockDatastore = datastoremock.NewMockDatastorer(ut.mockCtrl)

	ut.apiHandler = delegation.New(ut.mockDatastore)

	return ut
}

var (
	errGetDelegations   = errors.New("error getting delegations")
	errCountDelegations = errors.New("error count delegations")
)

func TestDelegation_GetDelegationsHandler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		request        func() *http.Request
		init           func(*underTest)
		want           []*model.Delegation
		wantErr        error
		wantStatusCode int
	}{
		{
			name: "Success",
			request: func() *http.Request {
				// Create a request without a year parameter
				req, err := http.NewRequestWithContext(
					context.Background(),
					http.MethodGet,
					"/delegations?page=1&size=100",
					nil,
				)
				require.NoError(t, err, "Error creating request")

				return req
			},
			init: func(ut *underTest) {
				ut.mockDatastore.EXPECT().GetDelegations(
					gomock.Any(),
					gomock.Eq(1),
					gomock.Eq(100),
					gomock.Eq(0),
				).Return([]*model.Delegation{
					{
						Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						Amount:    57800,
						Delegator: "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6",
						Block:     "123456",
					},
					{
						Timestamp: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
						Amount:    157800,
						Delegator: "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK7",
						Block:     "56897",
					},
				}, nil)

				ut.mockDatastore.EXPECT().GetDelegationsCount(
					gomock.Any(),
					gomock.Eq(0),
				).Return(2, nil)
			},
			want: []*model.Delegation{
				{
					Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					Amount:    57800,
					Delegator: "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6",
					Block:     "123456",
				},
				{
					Timestamp: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					Amount:    157800,
					Delegator: "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK7",
					Block:     "56897",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Success empty",
			request: func() *http.Request {
				// Create a request without a year parameter
				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/delegations", nil)
				require.NoError(t, err, "Error creating request")

				return req
			},
			init: func(ut *underTest) {
				ut.mockDatastore.EXPECT().GetDelegations(
					gomock.Any(),
					gomock.Eq(1),
					gomock.Eq(100),
					gomock.Eq(0),
				).Return([]*model.Delegation{}, nil)

				ut.mockDatastore.EXPECT().GetDelegationsCount(
					gomock.Any(),
					gomock.Eq(0),
				).Return(0, nil)
			},
			want:           []*model.Delegation{},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Success with year parameter",
			request: func() *http.Request {
				// Create a request without a year parameter
				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/delegations?year=2020", nil)
				require.NoError(t, err, "Error creating request")

				return req
			},
			init: func(ut *underTest) {
				ut.mockDatastore.EXPECT().GetDelegations(
					gomock.Any(),
					gomock.Eq(1),
					gomock.Eq(100),
					gomock.Eq(2020),
				).Return([]*model.Delegation{
					{
						Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						Amount:    57800,
						Delegator: "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6",
						Block:     "123456",
					},
					{
						Timestamp: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
						Amount:    157800,
						Delegator: "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK7",
						Block:     "56897",
					},
				}, nil)

				ut.mockDatastore.EXPECT().GetDelegationsCount(
					gomock.Any(),
					gomock.Eq(2020),
				).Return(2, nil)
			},
			want: []*model.Delegation{
				{
					Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					Amount:    57800,
					Delegator: "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6",
					Block:     "123456",
				},
				{
					Timestamp: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					Amount:    157800,
					Delegator: "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK7",
					Block:     "56897",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Error year parameter",
			request: func() *http.Request {
				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/delegations?year=deux", nil)
				require.NoError(t, err, "Error creating request")

				return req
			},
			init:           func(ut *underTest) {},
			wantErr:        errors.New("Bad Request: couldn't parse value year for query parameter deux\n"), //nolint:revive
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Error GetDelegations from datastore",
			request: func() *http.Request {
				// Create a request without a year parameter
				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/delegations", nil)
				require.NoError(t, err, "Error creating request")

				return req
			},
			init: func(ut *underTest) {
				ut.mockDatastore.EXPECT().GetDelegations(
					gomock.Any(),
					gomock.Eq(1),
					gomock.Eq(100),
					gomock.Eq(0),
				).Return(nil, errGetDelegations)
			},
			wantErr:        errors.New("Internal Server Error\n"), //nolint:revive
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Error GetDelegationsCount from datastore",
			request: func() *http.Request {
				// Create a request without a year parameter
				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/delegations", nil)
				require.NoError(t, err, "Error creating request")

				return req
			},
			init: func(ut *underTest) {
				ut.mockDatastore.EXPECT().GetDelegations(
					gomock.Any(),
					gomock.Eq(1),
					gomock.Eq(100),
					gomock.Eq(0),
				).Return([]*model.Delegation{
					{
						Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						Amount:    57800,
						Delegator: "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6",
						Block:     "123456",
					},
					{
						Timestamp: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
						Amount:    157800,
						Delegator: "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK7",
						Block:     "56897",
					},
				}, nil)

				ut.mockDatastore.EXPECT().GetDelegationsCount(
					gomock.Any(),
					gomock.Eq(0),
				).Return(0, errCountDelegations)
			},
			wantErr:        errors.New("Internal Server Error\n"), //nolint:revive
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			ut := setupTest(t)
			c.init(ut)

			responseRecorder := httptest.NewRecorder()
			ut.apiHandler.GetDelegationsHandler(responseRecorder, c.request())

			// Check the status code
			assert.Equal(t, c.wantStatusCode, responseRecorder.Code)

			if c.wantErr != nil {
				assert.Equal(t, c.wantErr.Error(), responseRecorder.Body.String())
			} else {
				// Parse the JSON response
				var result []*model.Delegation
				err := json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				require.NoError(t, err, "Error parsing JSON response")
				assert.Equal(t, c.want, result)
			}
		})
	}
}
