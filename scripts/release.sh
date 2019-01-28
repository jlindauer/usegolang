#!/bin/bash

cd "$GOPATH/src/github.com/jlindauer/usegolang"

echo "==== Releasing usegolang.com ===="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm usegolang
echo "  Done!"

echo "  Deleting existing code..."
ssh root@167.99.229.109 "rm -rf /root/go/src/github.com/jlindauer/usegolang"
echo "  Code deleted successfully!"

echo "  Creating deployment directory..."
ssh root@167.99.229.109 "mkdir /root/go/src/github.com/jlindauer/usegolang"
echo "  Directory created successfully!"

echo "  Uploading code..."
rsync -avr --exclude '.git/*' --exclude 'tmp/*' \
  --exclude 'images/*' ./ \
  root@167.99.229.109:/root/go/src/github.com/jlindauer/usegolang/
echo "  Code uploaded successfully!"

echo "  Go getting deps..."
ssh root@167.99.229.109 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
ssh root@167.99.229.109 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/mux"
ssh root@167.99.229.109 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/schema"
ssh root@167.99.229.109 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/lib/pq"
ssh root@167.99.229.109 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/jinzhu/gorm"
ssh root@167.99.229.109 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/csrf"

echo "  Building the code on remote server..."
ssh root@167.99.229.109 'export GOPATH=/root/go; \
  cd /root/app; \
  /usr/local/go/bin/go build -o ./server \
    $GOPATH/src/github.com/jlindauer/usegolang/*.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh root@167.99.229.109 "cd /root/app; \
  cp -R /root/go/src/github.com/jlindauer/usegolang/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh root@167.99.229.109 "cd /root/app; \
  cp -R /root/go/src/github.com/jlindauer/usegolang/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh root@167.99.229.109 "cd /root/app; \
  cp /root/go/src/github.com/jlindauer/usegolang/Caddyfile ."
echo "  Caddyfile moved successfully!"

echo "  Restarting the server..."
ssh root@167.99.229.109 "sudo service usegolang.com restart"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@167.99.229.109 "sudo service caddy restart"
echo "  Caddy restarted successfully!"

echo "==== Done releasing usegolang.com ===="
