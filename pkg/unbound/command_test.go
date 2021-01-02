package unbound

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/command/mock_command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Start(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	const unboundEtcDir = "/unbound"
	const unboundPath = "/usr/sbin/unbound"
	commander := mock_command.NewMockCommander(mockCtrl)
	commander.EXPECT().Start(context.Background(), unboundPath, "-d", "-c", "/unbound/unbound.conf", "-vv").
		Return(nil, nil, nil, nil).Times(1)
	c := &configurator{
		commander:     commander,
		unboundEtcDir: unboundEtcDir,
		unboundPath:   unboundPath,
	}
	stdout, waitFn, err := c.Start(context.Background(), 2)
	assert.Nil(t, stdout)
	assert.Nil(t, waitFn)
	assert.NoError(t, err)
}

func Test_Version(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		runOutput string
		runErr    error
		version   string
		err       error
	}{
		"no data": {
			err: fmt.Errorf(`unbound version was not found in ""`),
		},
		"2 lines with version": {
			runOutput: "Version  \nVersion 1.0-a hello\n",
			version:   "1.0-a",
		},
		"run error": {
			runErr: fmt.Errorf("error"),
			err:    fmt.Errorf("unbound version: error"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			commander := mock_command.NewMockCommander(mockCtrl)
			ctx := context.Background()

			const unboundEtcDir = "/unbound"
			const unboundPath = "/usr/sbin/unbound"
			commander.EXPECT().Run(ctx, unboundPath, "-V").
				Return(tc.runOutput, tc.runErr).Times(1)
			c := &configurator{
				commander:     commander,
				unboundEtcDir: unboundEtcDir,
				unboundPath:   unboundPath,
			}
			version, err := c.Version(ctx)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.version, version)
		})
	}
}
