package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/server-may-cry/bubble-go/models"
	"github.com/server-may-cry/bubble-go/protocol"
	"github.com/server-may-cry/bubble-go/storage"
)

var i = 0
var upgrader = websocket.Upgrader{} // use default options

// Test used for debug
func Test(w http.ResponseWriter, r *http.Request) {
	i = i + 1

	jsonByte, _ := json.Marshal(models.Test{
		Ttt: i,
	})

	w.Write(jsonByte)
}

// Redis used for send "ping" and receive "pong" from redis
func Redis(w http.ResponseWriter, r *http.Request) {
	pong, err := storage.Redis.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	jsonByte, _ := json.Marshal(protocol.RedisResponse{
		Ping: pong,
	})

	w.Write(jsonByte)
}

// Echo conroll websocket with echo
func Echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	log.Println("closed")
}

// Home serve index page with websocket js template
func Home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "wss://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server,
"Send" to send a message to the server and "Close" to close the connection.
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
