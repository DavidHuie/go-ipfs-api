package shell

import (
	"context"
	"encoding/json"
	"strconv"

	ma "github.com/multiformats/go-multiaddr"
)

// P2PListener describes a P2P listener.
type P2PListener struct {
	Protocol string
	Address  string
}

// P2POpenListener forwards P2P connections to a network multiaddr.
func (s *Shell) P2POpenListener(ctx context.Context, protocol, maddr string) (*P2PListener, error) {
	if _, err := ma.NewMultiaddr(maddr); err != nil {
		return nil, err
	}

	req := s.newRequest(ctx, "p2p/listener/open", protocol, maddr)
	resp, err := req.Send(s.httpcli)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.Error != nil {
		return nil, resp.Error
	}

	var response *P2PListener
	if err := json.NewDecoder(resp.Output).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

// P2PCloseListener closes one or all active P2P listeners.
func (s *Shell) P2PCloseListener(ctx context.Context, protocol string, closeAll bool) error {
	var args []string
	if protocol != "" {
		args = append(args, protocol)
	}

	req := s.newRequest(ctx, "p2p/listener/close", args...)
	req.Opts["all"] = strconv.FormatBool(closeAll)

	resp, err := req.Send(s.httpcli)
	if err != nil {
		return err
	}
	defer resp.Close()

	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

// P2PListenerList contains a slice of P2PListeners.
type P2PListenerList struct {
	Listeners []*P2PListener
}

// P2PListListeners lists all P2P listeners.
func (s *Shell) P2PListListeners(ctx context.Context) (*P2PListenerList, error) {
	req := s.newRequest(ctx, "p2p/listener/ls")
	resp, err := req.Send(s.httpcli)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.Error != nil {
		return nil, resp.Error
	}

	var response *P2PListenerList
	if err := json.NewDecoder(resp.Output).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

// P2PStream describes a P2P stream.
type P2PStream struct {
	Protocol string
	Address  string
}

// P2PStreamDial dials to a peer's P2P listener.
func (s *Shell) P2PStreamDial(ctx context.Context, peerID, protocol, listenerMaddr string) (*P2PStream, error) {
	if _, err := ma.NewMultiaddr(listenerMaddr); err != nil {
		return nil, err
	}

	args := []string{peerID, protocol}
	if listenerMaddr != "" {
		args = append(args, listenerMaddr)
	}

	req := s.newRequest(ctx, "p2p/stream/dial", args...)
	resp, err := req.Send(s.httpcli)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.Error != nil {
		return nil, resp.Error
	}

	var response *P2PStream
	if err := json.NewDecoder(resp.Output).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

// P2PCloseStream closes one or all active P2P streams.
func (s *Shell) P2PCloseStream(ctx context.Context, handlerID string, closeAll bool) error {
	var args []string
	if handlerID != "" {
		args = append(args, handlerID)
	}

	req := s.newRequest(ctx, "p2p/stream/close", args...)
	req.Opts["all"] = strconv.FormatBool(closeAll)

	resp, err := req.Send(s.httpcli)
	if err != nil {
		return err
	}
	defer resp.Close()

	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

// P2PStreamsList contains a slice of streams.
type P2PStreamsList struct {
	Streams []*struct {
		HandlerID     string
		Protocol      string
		LocalPeer     string
		LocalAddress  string
		RemotePeer    string
		RemoteAddress string
	}
}

// P2PListStreams lists all P2P streams.
func (s *Shell) P2PListStreams(ctx context.Context) (*P2PStreamsList, error) {
	req := s.newRequest(ctx, "p2p/stream/ls")
	req.Opts["headers"] = strconv.FormatBool(true)

	resp, err := req.Send(s.httpcli)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.Error != nil {
		return nil, resp.Error
	}

	var response *P2PStreamsList
	if err := json.NewDecoder(resp.Output).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}
