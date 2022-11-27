#!/bin/sh
set -e

# containers on linux share file permissions with hosts.
# assigning the same uid/gid from the host user
# ensures that the files can be read/write from both sides
if ! id taro > /dev/null 2>&1; then
  USERID=${USERID:-1000}
  GROUPID=${GROUPID:-1000}

  echo "adding user taro ($USERID:$GROUPID)"
  addgroup -g $GROUPID taro
  adduser -u $USERID -D -G taro taro
fi

# if [ $(echo "$1" | cut -c1) = "-" ]; then
#   echo "$0: assuming arguments for taro"
#   set -- tarod "$@"
# fi
sleep 1
if [ "$1" = "tarod" ] || [ "$1" = "tarocli" ]; then
  echo "Running as taro user: $@"
  exec su-exec taro "$@"
fi

echo "$@"
exec "$@"