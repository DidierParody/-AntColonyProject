package server

import (
	"log"
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
)

type Participant struct {
	Host bool
	Conn *websocket.Conn
}

type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant //map de slices de participantes
}

// utilizamos init para inicializar el map de RoomMap
func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant) // se crea un map vacio de slices(es una lista como en pytho)  de participant
}

// utilizamos Get para obtener los participantes de una sala especifica	por medio de un roomID string
func (r *RoomMap) Get(roomID string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Map[roomID]
}

// CreateRoom
func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	//rand.Seed(time.Now().UnixNano()) ete metodo quedo obsoleto
	// letters es un slice de runes (las cuales son una especie de tabla ascci)
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 8) //se crea un slice de runes de 8 posiciones (caracteres en este caso)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	} //con este for insertamos letras al azar de el slice letters en el slice b

	roomID := string(b) //qui obtenemos el id unico para el room

	r.Map[roomID] = []Participant{} //inseta una key para una lista de participantes

	return roomID //devolvemos el identificador de nuestra sala
}

// utilizamos InsertIntoRoom para insertar un participante en una sala especifica
func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{host, conn}

	log.Println("Inserting into Room with RoomID: ", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)
}

// DeleteRoom elimina una sala y todos sus participantes
func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)

}
