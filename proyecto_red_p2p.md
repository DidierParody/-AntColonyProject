# 📖 Tesis del Proyecto

La idea central es construir un **sistema web semidescentralizado** donde, a diferencia de la arquitectura tradicional cliente-servidor (donde toda la carga recae en el servidor), los usuarios mismos (navegadores) se convierten en **nodos activos de una red P2P (peer-to-peer)**.

En este modelo:

* Cuantos más usuarios se conecten, **más capacidad y robustez gana la red**.
* La carga del servidor disminuye porque el tráfico y los datos se distribuyen entre los propios usuarios.
* El servidor solo actúa como **punto de arranque** o respaldo cuando la red no puede resolver algo por sí sola.

**Finalidad:**

* Escalabilidad sin aumentar exponencialmente los costos de infraestructura.
* Robustez: la red puede sobrevivir a la caída de nodos individuales gracias a los **nodos de contingencia**.
* Seguridad: cada nodo participa solo con datos encriptados y minimiza el intercambio con el servidor.

---

# 🌐 Dinámica abstracta de la red de nodos

Imagina la red como una **colonia de hormigas**:

* El **servidor** es como la entrada al hormiguero, necesaria al inicio para organizar, pero no controla a todas las hormigas.
* Cada **nodo (usuario)** es una hormiga que lleva parte de la información y puede comunicarse con otras hormigas.
* Cuando una hormiga (nodo) cae, otra hormiga **contingente** hereda su responsabilidad para que la colonia siga funcionando.
* La **DHT** es el “mapa” que usan las hormigas para saber hacia dónde ir cuando buscan información.

De esta forma, el sistema crece con más usuarios y se vuelve más fuerte y distribuido.

---

# ⚙️ Explicación técnica del flujo

### 1. Ingreso del usuario

* El usuario abre la aplicación web desde su navegador (`https://miapp.com`).
* Hasta este punto, el navegador no sabe nada de la red P2P.

---

### 2. Signaling Server (fase de arranque)

* Se abre un canal **WebSocket (Gorilla en Go)** entre el navegador y el **signaling server**.
* Función: **teléfono temporal** para que dos nodos intercambien información de conexión.
* Aquí se negocian:

  * **SDP (Session Description Protocol):** describe cómo quiere comunicarse cada nodo.
  * **ICE Candidates:** posibles rutas de red (IP local, IP pública, relay).

👉 El signaling server no transmite datos de usuario ni controla la red. Solo **facilita el primer apretón de manos**.

---

### 3. Establecimiento de conexión P2P con WebRTC

* Con la info (SDP + ICE), los navegadores abren un canal **WebRTC DataChannel**.
* Ahora se comunican directamente sin pasar por el signaling server.
* El signaling server puede apagarse para este par de nodos: **misión cumplida**.

👉 **WebRTC** = autopista directa y segura entre dos nodos, optimizada para baja latencia.

---

### 4. Integración a la red libp2p

* Ahora que el navegador tiene un canal directo, entra en la **red P2P con libp2p**.
* Libp2p le permite:

  * **Conectarse con múltiples peers** (no solo el primero).
  * Usar la **DHT (Distributed Hash Table)** para descubrir nodos.

👉 **DHT** = como un “DNS distribuido”: cada nodo guarda una parte del mapa de la red y sabe cómo enrutar consultas hasta llegar al nodo que tiene lo que buscas.

---

### 5. Registro en la DHT

* El nuevo nodo se registra en la DHT publicando:

  * Su **PeerID único**.
  * Sus **direcciones accesibles** (ejemplo: su WebRTC transport).
* Ahora puede ser encontrado por cualquier otro nodo que consulte la DHT.

---

### 6. 📌 Entrada de los nodos de contingencia

* Una vez registrado, la red le asigna un **nodo contingente**.
* Función del nodo contingente:

  * Respaldar la información del nuevo nodo en caso de caída.
  * Intermediar consultas mientras el nodo gana conexiones directas.
  * Reasignar responsabilidades si el nuevo nodo desaparece.

👉 Piensa en el contingente como un **hermano mayor** que protege al nodo más joven.

#### Ejemplo con nodos A, B y C

* **Nodo A** tiene como contingente al **Nodo B**.  
* **Nodo B** tiene como contingente al **Nodo C**.  

Cuando **B se desconecta**:

1. B avisa a A que ya no estará disponible.  
2. B le dice a A que su **nuevo contingente será C**.  
3. La información que tenía B de respaldo de A se transfiere a C.  
4. Ahora A → contingente = C, y la red mantiene continuidad sin perder datos.  

---

### 7. Intercambio de datos

* Cuando un nodo necesita información:

  1. Consulta a la **DHT** → “¿Quién tiene este dato o quién es este PeerID?”
  2. La DHT responde con la ruta más cercana.
  3. El nodo establece un canal directo vía WebRTC.
* Los datos se transfieren en **MessagePack**:

  * Similar a JSON pero binario, más rápido y compacto.

---

# ✅ Resumen del flujo

1. Usuario entra → carga frontend.
2. Se abre canal WebSocket → Signaling Server.
3. Se negocian SDP + ICE → WebRTC crea canal P2P.
4. Libp2p adopta WebRTC como transporte.
5. Nodo se registra en la DHT → publica PeerID + direcciones.
6. Nodo obtiene un **nodo de contingencia**.
7. Intercambio de datos P2P en **MessagePack**.
8. El servidor solo interviene si no hay nodo disponible o como fallback.

---

# 📊 Diagrama del flujo (versión simplificada)

```
[ Usuario entra al sitio ]
           |
           v
 [ Signaling Server ]
   (WebSocket Gorilla)
           |
           v
  < Intercambio SDP + ICE >
           |
           v
 [ WebRTC Canal Directo ]
           |
           v
 [ libp2p integra el nodo ]
     -> usa DHT para:
         - Registro PeerID
         - Descubrimiento de peers
           |
           v
 [ Nodo obtiene contingente ]
   (ejemplo: A→B, B→C,
    si B cae: A→C)
           |
           v
 [ Operación normal ]
 - Datos viajan P2P
 - Codificados en MessagePack
 - DHT enruta consultas
 - Contingente respalda
```
