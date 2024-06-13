DB_HOST=localhost
DB_PORT=27017
DB_DATABASE=teos
DB_USERNAME=
DB_PASSWORD=
DB_AUTH=
DUMP_FILE="./teos_$(date +"%Y-%m-%d")"

if [ -z "$DB_USERNAME" ] || [ -z "$DB_PASSWORD" ]
then
    cmd="mongodump --host=$DB_HOST --db=$DB_DATABASE --port=$DB_PORT -v --out=$DUMP_FILE"
else
    cmd="mongodump --host=$DB_HOST --db=$DB_DATABASE --port=$DB_PORT --username=$DB_USERNAME --password=$DB_PASSWORD --authenticationDatabase=$DB_AUTH -v --out=$DUMP_FILE"
fi

echo "Dumping database '$DB_DATABASE' from '$DB_HOST':'$DB_PORT' to '$DUMP_FILE', wait..."
$cmd
if [ $? -eq 0 ] 
then
    echo "Database dump success."
else
    echo "Database dump error!"
fi
