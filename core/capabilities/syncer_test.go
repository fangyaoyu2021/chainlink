package capabilities_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	commonMocks "github.com/smartcontractkit/chainlink-common/pkg/types/mocks"
	coreCapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
)

func TestSyncer_CleanStartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	registry := commonMocks.NewCapabilitiesRegistry(t)
	registry.On("Add", mock.Anything, mock.Anything).Return(nil)
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("SetReceiver", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	onchainRegistry := coreCapabilities.NewOnchainCapabilityRegistry(types.MustEIP55Address("0x0000000000000000000000000000000000000001").Address(), lggr)

	syncer := coreCapabilities.NewRegistrySyncer(wrapper, registry, dispatcher, lggr, onchainRegistry)
	require.NoError(t, syncer.Start(ctx))
	require.NoError(t, syncer.Close())
}
