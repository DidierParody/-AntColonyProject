# Salir si hay un error
set -e

# Construir el cliente de React
cd client
npm install
npm run build