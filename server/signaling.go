package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap //definimos dos variables globales de tipo AllRooms y  RoomMap
// AllRooms es una variable global que contiene todas las salas creadas, es una map de Rooms

//Recuerda que un Handler es una funcion que responde una solicitud  http

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	// w es el objeto que representa la respuesta http que el servidor le envia al cliente
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.HHeader() devuelve un mapa de encabezados HTTP que acompanan la respuesta que el servidor envia al cliente
	// Acces-Control-Allow-Origin pertenece a las politicas de CORS
	//controla que origen puede hacer solicitudes a el servidor
	//en resumen se permite que cualer origen puedea hacer solicitudes a el servidor
	// r contiene toda la informacion que envio el cliente
	// w (wire o escribir) es lo que se le envia al cliente y r (leer o read) es lo que se recibe del cliente
	roomID := AllRooms.CreateRoom() //recuerda que esto retonrna el roomID y crea la sala

	//usamos el paquete de json para convertir nuestro roomID a formato json
	type resp struct {
		RoomID string `json:"room_id"`
	}

	// imprime los mapasa de las salas
	log.Println(AllRooms.Map)
	//enviamos una resopuesta (w) y la data de esta drespuesta es lo que pongamos en el struct
	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

//recuerda que el upgrader es para que podamos convertir la conexion de http a websocket

// upgrader permite actualizar una conexion http a una websocket
var upgrader = websocket.Upgrader{
	// CheckOrigin verifica que la  conexion sea a la que se le dio permisos de origen
	// retunr true porque permitimos cualquiier origen
	CheckOrigin: func(r *http.Request) bool {
		return true
	}, //la coma esta por buenas practicas de go
}

// un broadcast se usa para "transmitir a todos"
// en este caso broadcastMsg sirve para enviar un mensaje del servidor a todos los clientes conectados

type broadcastMsg struct {
	Message map[string]interface{} // mapa string y recibe cualcquier tipo de dato, porque no se sabe que datos tendra el mensaje
	// se usa comunmente map[string]interface{} cuando no sabemos que tipo de datos contrendra el mensaje
	RoomID string          //el identificador de la sala
	Client *websocket.Conn // es el puntero a la conexion websocket del cliente que envio el mensaje
	//se usa para saber quien lo envio, evitar erenviar el mesmo mensjae y para cerrar su conexion en caso de error
}

var broadcast = make(chan broadcastMsg) // se cre un canal de tipo broadcastMsg

func broadcaster() {
	//es un bucle infinito que escucha constantemene el canal broadcast
	for {
		//recibe l mensaje y lo guarda en msg
		msg := <-broadcast
		for _, client := range AllRooms.Map[msg.RoomID] {
			// el _, significa que queremos que ignore el key y solo se fije en el value de el map
			//recorre todos los clientes conectados en esa sala en especifico
			if client.Conn != msg.Client {
				//esta condicion evalua si el cliente que envio el mensaje es diferente al cliente que lo recibe
				err := client.Conn.WriteJSON(msg.Message)
				// si todo esta bien, entonces se envia el mensaje a los demas clientes
				// se toan lkos datos y se convierten a json
				//se enmvian a traves del socket
				//devuelve un error si algo sale mal
				//por defecto revuelve un nil
				if err != nil {
					//si err no devolvio un nil significa que hubo un error
					log.Fatal(err)
					client.Conn.Close()

				}

			}
		}
	}
}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	// hace una busqueda de la roomID en los parametros de la url
	// ok es un boolean que verifica si existe informacion dentro de la solicitud
	// esto hace una consulta en la url de la solicitud y lo transforma en un map con toda la informacion
	roomID, ok := r.URL.Query()["roomID"]

	//si no existe el roomID  dentro de la url
	if !ok {
		log.Println("RoomID missing in URL Parameters")
		return //pausa la ejecucion de la funcion
	}
	// se encarga de tranformar con el upgrader la onexion de http a websocket
	// ws, sirve como el retorno de la funcion updrader.Upgrade y representa el websocket
	// err, tambien tiene un error para verificar que todo salio bien
	ws, err := upgrader.Upgrade(w, r, nil)

	// si hubo un error
	if err != nil {
		log.Fatal("Web Socket Upgrade Error", err)
	}

	AllRooms.InsertIntoRoom(roomID[0], false, ws)
	// insertamos al cliente a la sala,
	// el roomID es 0 porque es el primer elemento del slice que se nos retorna al hacer la consulta a la url, y en este caso ese es el que nos interesa

	go broadcaster()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("Read Error:", err)
		}

		msg.Client = ws
		msg.RoomID = roomID[0]

		log.Println(msg.Message)

		broadcast <- msg

	}

}
