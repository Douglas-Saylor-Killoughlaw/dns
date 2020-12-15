package dns

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qdm12/cloudflare-dns-server/internal/constants"
	"github.com/qdm12/golibs/files/mock_files"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/qdm12/golibs/network/mock_network"
)

func Test_DownloadRootHints(t *testing.T) { //nolint:dupl
	t.Parallel()
	tests := map[string]struct {
		content   []byte
		status    int
		clientErr error
		writeErr  error
		err       error
	}{
		"no data": {
			status: http.StatusOK,
		},
		"bad status": {
			status: http.StatusBadRequest,
			err:    fmt.Errorf("HTTP status code is 400 for https://raw.githubusercontent.com/qdm12/files/master/named.root.updated"),
		},
		"client error": {
			clientErr: fmt.Errorf("error"),
			err:       fmt.Errorf("error"),
		},
		"write error": {
			status:   http.StatusOK,
			writeErr: fmt.Errorf("error"),
			err:      fmt.Errorf("error"),
		},
		"data": {
			content: []byte("content"),
			status:  http.StatusOK,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ctx := context.Background()
			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Info("downloading root hints from %s", constants.NamedRootURL).Times(1)
			client := mock_network.NewMockClient(mockCtrl)
			client.EXPECT().Get(ctx, string(constants.NamedRootURL)).
				Return(tc.content, tc.status, tc.clientErr).Times(1)
			fileManager := mock_files.NewMockFileManager(mockCtrl)
			if tc.clientErr == nil && tc.status == http.StatusOK {
				fileManager.EXPECT().WriteToFile(
					string(constants.RootHints),
					tc.content,
					gomock.Any()).
					Return(tc.writeErr).Times(1)
			}
			c := &configurator{logger: logger, client: client, fileManager: fileManager}
			err := c.DownloadRootHints(ctx)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_DownloadRootKey(t *testing.T) { //nolint:dupl
	t.Parallel()
	tests := map[string]struct {
		content   []byte
		status    int
		clientErr error
		writeErr  error
		err       error
	}{
		"no data": {
			status: http.StatusOK,
		},
		"bad status": {
			status: http.StatusBadRequest,
			err:    fmt.Errorf("HTTP status code is 400 for https://raw.githubusercontent.com/qdm12/files/master/root.key.updated"),
		},
		"client error": {
			clientErr: fmt.Errorf("error"),
			err:       fmt.Errorf("error"),
		},
		"write error": {
			status:   http.StatusOK,
			writeErr: fmt.Errorf("error"),
			err:      fmt.Errorf("error"),
		},
		"data": {
			content: []byte("content"),
			status:  http.StatusOK,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ctx := context.Background()
			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Info("downloading root key from %s", constants.RootKeyURL).Times(1)
			client := mock_network.NewMockClient(mockCtrl)
			client.EXPECT().Get(ctx, string(constants.RootKeyURL)).
				Return(tc.content, tc.status, tc.clientErr).Times(1)
			fileManager := mock_files.NewMockFileManager(mockCtrl)
			if tc.clientErr == nil && tc.status == http.StatusOK {
				fileManager.EXPECT().WriteToFile(
					string(constants.RootKey),
					tc.content,
					gomock.Any(),
				).Return(tc.writeErr).Times(1)
			}
			c := &configurator{logger: logger, client: client, fileManager: fileManager}
			err := c.DownloadRootKey(ctx)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
