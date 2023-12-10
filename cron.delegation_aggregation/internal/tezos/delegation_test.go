package tezos_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/tezos"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/terrs"
)

func TestTezos_ListDelegations(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		init      func(ut *underTest)
		want      []*tezos.Delegation
		unmarshal func(data []byte, v any) error
		wantErr   error
	}{
		{
			name: "Success",
			init: func(ut *underTest) {
				ut.mockTransport.RegisterResponder(http.MethodGet,
					"https://api.tezos.test/v1/operations/delegations"+
						"?limit=100&select=timestamp%2Camount%2Csender%2Cblock&sort.desc=id",
					httpmock.NewStringResponder(http.StatusOK, `
						[
							{
								"timestamp": "2023-12-10T11:01:01Z",
								"amount": 124428330,
								"sender": {
									"address": "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx"
								},
								"block": "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG"
							},
							{
								"timestamp": "2023-12-10T11:00:01Z",
								"amount": 499836,
								"sender": {
									"alias": "The E Major",
									"address": "tz1NqVXDBf8fZNomacychFPSK1trQbi14PvA"
								},
								"block": "BLxQGrPcAPAwKaeCdivBVw45Choicesen6wrmdm3NBeGsCnkLKv"
							}]
					`))
			},
			want: []*tezos.Delegation{
				{
					Timestamp: time.Date(2023, 12, 10, 11, 1, 1, 0, time.UTC),
					Amount:    124428330,
					Sender: tezos.Sender{
						Address: "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx",
					},
					Block: "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG",
				},
				{
					Timestamp: time.Date(2023, 12, 10, 11, 0, 1, 0, time.UTC),
					Amount:    499836,
					Sender: tezos.Sender{
						Address: "tz1NqVXDBf8fZNomacychFPSK1trQbi14PvA",
					},
					Block: "BLxQGrPcAPAwKaeCdivBVw45Choicesen6wrmdm3NBeGsCnkLKv",
				},
			},
			unmarshal: nil,
			wantErr:   nil,
		},
		{
			name: "Error get",
			init: func(ut *underTest) {
				ut.mockTransport.RegisterResponder(http.MethodGet,
					"https://api.tezos.test/v1/operations/delegations"+
						"?limit=100&select=timestamp%2Camount%2Csender%2Cblock&sort.desc=id",
					func(req *http.Request) (*http.Response, error) {
						return nil, terrs.NewTestError()
					})
			},
			want:      nil,
			unmarshal: nil,
			wantErr: fmt.Errorf("couldn't list delegations from tezos api error: %w",
				&url.Error{
					Op: "Get",
					URL: "https://api.tezos.test/v1/operations/delegations" +
						"?limit=100&select=timestamp%2Camount%2Csender%2Cblock&sort.desc=id",
					Err: terrs.NewTestError(),
				},
			),
		},
		{
			name: "Error unmarshal",
			init: func(ut *underTest) {
				ut.mockTransport.RegisterResponder(http.MethodGet,
					"https://api.tezos.test/v1/operations/delegations"+
						"?limit=100&select=timestamp%2Camount%2Csender%2Cblock&sort.desc=id",
					httpmock.NewStringResponder(http.StatusOK, `
						[
							{
								"timestamp": "2023-12-10T11:01:01Z",
								"amount": 124428330,
								"sender": {
									"address": "tz1eZsUhWxawxDn5U24LGKiozLapYvAbw2yx"
								},
								"block": "BMWE6vssezoqBCSSJd9M24jmExjLRyVrAr4f7sWFywGKjD3TeSG"
							},
							{
								"timestamp": "2023-12-10T11:00:01Z",
								"amount": 499836,
								"sender": {
									"alias": "The E Major",
									"address": "tz1NqVXDBf8fZNomacychFPSK1trQbi14PvA"
								},
								"block": "BLxQGrPcAPAwKaeCdivBVw45Choicesen6wrmdm3NBeGsCnkLKv"
							}]
					`))
			},
			want: nil,
			unmarshal: func(data []byte, v any) error {
				return terrs.NewTestError()
			},
			wantErr: fmt.Errorf("couldn't list delegations from tezos api error: %w",
				terrs.NewTestError(),
			),
		},
		{
			name: "Error internal",
			init: func(ut *underTest) {
				ut.mockTransport.RegisterResponder(http.MethodGet,
					"https://api.tezos.test/v1/operations/delegations"+
						"?limit=100&select=timestamp%2Camount%2Csender%2Cblock&sort.desc=id",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(http.StatusInternalServerError, map[string]string{
							"code": "500",
							"msg":  "error",
						})
					})
			},
			want:      nil,
			unmarshal: nil,
			wantErr: fmt.Errorf(
				`couldn't list delegations from tezos api error: {"code":"500","msg":"error"}`,
			),
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			ut := setupTest(t, c.unmarshal)
			defer ut.mockTransport.Reset()

			c.init(ut)
			resp, err := ut.client.ListDelegations(context.Background(), nil)

			assert.Equal(t, c.want, resp)
			assert.Equal(t, c.wantErr, err)
		})
	}
}
