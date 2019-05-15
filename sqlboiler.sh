#!/usr/bin/env bash

set -e

# we create subdirs in here.
# we'll run a PG docker here
readonly _RANDO_PORT=5932

function main {
    case "${1}" in
        db) _make_db ;;
        gen) _gen ;;
        test) _run_tests ;;
        stop) _stop_db ;;
        *) 
            echo "Usage: sqlboiler.sh (db|gen|test|stop)" >&2
            exit 1
            ;;
    esac
}

function _make_db {
    docker run \
        --detach \
        --rm \
        --name sqlboiler-db \
        --publish ${_RANDO_PORT}:5432 \
        postgres

    echo "Giving PG 5 seconds to catch its breath" >&2
    sleep 5
 
    docker exec -i sqlboiler-db psql -Upostgres < sample.sql
}

function _stop_db {
    docker stop sqlboiler-db
}

function _run_tests {
    POSTGRES_DSN="postgres://postgres:postgres@localhost:${_RANDO_PORT}/?sslmode=disable" go test ./...
}

function _gen {
    trap 'rm -f /tmp/sqlboiler.*.toml' EXIT

    # get around mktemp weirdness on mac
    local path=/tmp/sqlboiler.$(cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f0-9' | head -c 10).toml
    if [[ -f "${path}" ]]; then
        echo "file exists" >&2
        exit 1
    fi

    cat <<EOF > "${path}"
[psql]
    host = "localhost"
    dbname = "postgres"
    user = "postgres"
    sslmode = "disable"
    port = ${_RANDO_PORT}
    
    schema = "public"
EOF

    if ! [ -x "$(command -v sqlboiler)" ]; then
        go get -u -t github.com/volatiletech/sqlboiler
        go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql
    fi
    sqlboiler psql \
        --config "${path}"  \
        --no-tests \
        --wipe
}

main "${1}"

