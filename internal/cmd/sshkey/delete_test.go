package sshkey_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := sshkey.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sshKey := &hcloud.SSHKey{
		ID:   123,
		Name: "test",
	}

	fx.Client.SSHKeyClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(sshKey, nil, nil)
	fx.Client.SSHKeyClient.EXPECT().
		Delete(gomock.Any(), sshKey).
		Return(nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "SSH Key test deleted\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := sshkey.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	keys := []*hcloud.SSHKey{
		{
			ID:   123,
			Name: "test1",
		},
		{
			ID:   456,
			Name: "test2",
		},
		{
			ID:   789,
			Name: "test3",
		},
	}

	var names []string
	for _, key := range keys {
		names = append(names, key.Name)
		fx.Client.SSHKeyClient.EXPECT().
			Get(gomock.Any(), key.Name).
			Return(key, nil, nil)
		fx.Client.SSHKeyClient.EXPECT().
			Delete(gomock.Any(), key).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, names)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "SSH Keys test1, test2, test3 deleted\n", out)
}
