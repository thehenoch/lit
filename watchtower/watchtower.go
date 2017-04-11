package watchtower

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/btcsuite/btcd/wire"
	"github.com/mit-dci/lit/lnutil"
)

// The main watchtower struct
type WatchTower struct {
	Path    string   // where the DB goes?  needed?
	WatchDB *bolt.DB // DB with everything in it

	Accepting bool // true if new channels and sigs are allowed in
	Watching  bool // true if there are txids to watch for

	SyncHeight int32 // last block we've sync'd to.  Not needed?

	OutBox chan *wire.MsgTx // where the tower sends its justice txs
}

// 2 structs used in the DB: IdxSigs and ChanStatic

// IdxSig is what we save in the DB for each txid
type IdxSig struct {
	PKHIdx   uint32   // Who
	StateIdx uint64   // When
	Sig      [64]byte // What
}

func (w *WatchTower) HandleMessage(lm lnutil.LitMsg) error {
	fmt.Printf("got message from %x\n", lm.Peer())

	switch lm.MsgType() {
	case lnutil.MSGID_WATCH_DESC:
		fmt.Printf("new channel to watch\n")
		desc, err := lnutil.NewWatchDescMsgFromBytes(lm.Bytes(), lm.Peer())
		if err != nil {
			return err
		} else {
			return w.AddNewChannel(*desc)
		}

	case lnutil.MSGID_WATCH_COMMSG:
		fmt.Printf("new commsg\n")
		commsg, err := lnutil.NewComMsgFromBytes(lm.Bytes(), lm.Peer())
		if err != nil {
			return err
		} else {
			return w.AddState(*commsg)
		}

	case lnutil.MSGID_WATCH_DELETE:
		fmt.Printf("delete message\n")
		// delete not yet implemented
	default:
		fmt.Printf("unknown message type %x\n", lm.MsgType())
	}
	return nil
}

func (w *WatchTower) JusticeOutbox() chan *wire.MsgTx {
	return w.OutBox
}
