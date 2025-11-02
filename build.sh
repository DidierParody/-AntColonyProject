
# Salir si hay un error
set -e

# Listar todos los archivos para diagnosticar
echo "--- Listando archivos del proyecto ---"
ls -R
echo "------------------------------------"

# Construir el cliente de React
cd client
npm install
npm run build
