package main

import (
	"compress/flate"
	"io"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsflate"
	"github.com/gobwas/ws/wsutil"
)

func main() {
	http.HandleFunc("/ws", simple)
	http.HandleFunc("/wsflate", wsflatef)

	// Run raw TCP listener to serve HTTP and handle upgrade manually
	log.Println("Gobwas echo server on :9004")
	log.Fatal(http.ListenAndServe(":9004", nil))
}

func simple(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	defer conn.Close()

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("read error: %v", err)
			return
		}

		err = wsutil.WriteServerMessage(conn, op, msg)
		if err != nil {
			log.Printf("write error: %v", err)
			return
		}
	}
}

func wsflatef(w http.ResponseWriter, r *http.Request) {
	e := wsflate.Extension{
		Parameters: wsflate.Parameters{
			ServerNoContextTakeover: true,
			ClientNoContextTakeover: true,
		},
	}
	u := ws.HTTPUpgrader{
		Negotiate: e.Negotiate,
	}
	conn, _, _, err := u.Upgrade(r, w)
	if err != nil {
		log.Printf("upgrade error: %s", err)
		return
	}
	defer conn.Close()

	if _, ok := e.Accepted(); !ok {
		log.Printf("no accepted extension")
		return
	}

	// Using nil as a destination io.Writer since we will Reset() it in the
	// loop below.
	fw := wsflate.NewWriter(nil, func(w io.Writer) wsflate.Compressor {
		// As flat.NewWriter() docs says:
		//   If level is in the range [-2, 9] then the error returned will
		//   be nil.
		f, _ := flate.NewWriter(w, 9)
		return f
	})
	// Using nil as a source io.Reader since we will Reset() it in the loop
	// below.
	fr := wsflate.NewReader(nil, func(r io.Reader) wsflate.Decompressor {
		return flate.NewReader(r)
	})

	// MessageState implements wsutil.Extension and is used to check whether
	// received WebSocket message is compressed. That is, it's generally
	// possible to receive uncompressed messaged even if compression extension
	// was negotiated.
	var msg wsflate.MessageState

	// Note that control frames are all written without compression.
	controlHandler := wsutil.ControlFrameHandler(conn, ws.StateServerSide)
	rd := wsutil.Reader{
		Source:         conn,
		State:          ws.StateServerSide | ws.StateExtended,
		CheckUTF8:      false,
		OnIntermediate: controlHandler,
		Extensions:     []wsutil.RecvExtension{&msg},
	}

	wr := wsutil.NewWriter(conn, ws.StateServerSide|ws.StateExtended, 0)
	wr.SetExtensions(&msg)

	for {
		h, err := rd.NextFrame()
		if err != nil {
			log.Printf("next frame error: %v", err)
			return
		}
		if h.OpCode.IsControl() {
			if err := controlHandler(h, &rd); err != nil {
				log.Printf("handle control frame error: %v", err)
				return
			}
			continue
		}

		wr.ResetOp(h.OpCode)

		var (
			src io.Reader = &rd
			dst io.Writer = wr
		)
		if msg.IsCompressed() {
			fr.Reset(src)
			fw.Reset(dst)
			src = fr
			dst = fw
		}
		// Copy incoming bytes right into writer, probably through decompressor
		// and compressor.
		if _, err = io.Copy(dst, src); err != nil {
			log.Fatal(err)
		}
		if msg.IsCompressed() {
			// Flush the flate writer.
			if err = fw.Close(); err != nil {
				log.Fatal(err)
			}
		}
		// Flush WebSocket fragment writer. We could send multiple fragments
		// for large messages.
		if err = wr.Flush(); err != nil {
			log.Fatal(err)
		}
	}
}
