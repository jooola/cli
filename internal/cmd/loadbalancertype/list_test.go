package loadbalancertype_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancertype"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := loadbalancertype.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.LoadBalancerTypeClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.LoadBalancerTypeListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     nil, // Load Balancer Types do not support sorting
			},
		).
		Return([]*hcloud.LoadBalancerType{
			{
				ID:             123,
				Name:           "test",
				MaxServices:    12,
				MaxConnections: 100,
				MaxTargets:     5,
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   DESCRIPTION   MAX SERVICES   MAX CONNECTIONS   MAX TARGETS
123   test   -             12             100               5          
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
