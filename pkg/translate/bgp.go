package translate

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/cmmarslender/edgefig/pkg/config"
)

type bgpDir string

const (
	bgpDirTo   = bgpDir("To")
	bgpDirFrom = bgpDir("From")
)

// getBGPGroupName returns the peer specific group name to use for prefix-lists and route-maps
func getBGPGroupName(bgpPeer config.BGPPeer, direction bgpDir) string {
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%d", bgpPeer.ASN)))
	hasher.Write(bgpPeer.IP.AsSlice())
	hashBytes := hasher.Sum(nil)
	hash := hex.EncodeToString(hashBytes)

	return fmt.Sprintf("BGP-%s-%s", hash, direction)
}
