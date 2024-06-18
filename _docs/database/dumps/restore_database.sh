DB_HOST=localhost
DB_PORT=27017
DB_DATABASE=teos
DB_USERNAME=
DB_PASSWORD=
DB_AUTH=
DUMP_FILE=$1

if [ -z "$DUMP_FILE" ]
then
    echo "ERR: path to database dump directory required"
fi

if [ -z "$DB_USERNAME" ] || [ -z "$DB_PASSWORD" ]
then
    cmd="mongorestore --host=$DB_HOST --port=$DB_PORT --db=$DB_DATABASE -v $DUMP_FILE"
else
    cmd="mongorestore --host=$DB_HOST --port=$DB_PORT --db=$DB_DATABASE --username=$DB_USERNAME --password=$DB_PASSWORD --authenticationDatabase=$DB_AUTH -v $DUMP_FILE"
fi

echo "Restoring database '$DB_DATABASE' to '$DB_HOST':'$DB_PORT' from '$DUMP_FILE', wait..."
$cmd
if [ $? -eq 0 ] 
then
    echo "Database restore success."
else
    echo "ERR: database restore error!"
fi