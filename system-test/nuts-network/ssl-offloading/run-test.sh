#!/usr/bin/env bash
source ../../util.sh

echo "------------------------------------"
echo "Cleaning up running Docker containers and volumes, and key material..."
echo "------------------------------------"
docker-compose down
docker-compose rm -f -v
rm -rf ./node-*/data
mkdir -p ./node-A/data/keys
mkdir -p ./node-B/data/keys
touch ./node-A/data/keys/truststore.pem

echo "------------------------------------"
echo "Starting Docker containers..."
echo "------------------------------------"
docker-compose up -d
waitForDCService nodeA-backend
waitForDCService nodeB

echo "------------------------------------"
echo "Registering vendors..."
echo "------------------------------------"
# Register Vendor A
docker-compose exec -e NUTS_MODE=cli nodeA-backend nuts crypto selfsign-vendor-cert "Vendor A" /opt/nuts/keys/vendor_certificate.pem
docker-compose exec -e NUTS_MODE=cli nodeA-backend nuts registry register-vendor /opt/nuts/keys/vendor_certificate.pem

# Register Vendor B
docker-compose exec -e NUTS_MODE=cli nodeB nuts crypto selfsign-vendor-cert "Vendor B" /opt/nuts/keys/vendor_certificate.pem
docker-compose exec -e NUTS_MODE=cli nodeB nuts registry register-vendor /opt/nuts/keys/vendor_certificate.pem

# Since node B connects to A's gRPC server, so A needs to trust B's Vendor CA certificate since it's used to issue the client certificate
docker cp ./node-B/data/keys/vendor_certificate.pem $(docker-compose ps -q nodeA):/etc/nginx/ssl/truststore.pem
# This also means that B must trust A's server certificate (by trusting our custom Root CA)
docker cp ../../keys/ca-certificate.pem $(docker-compose ps -q nodeB):/usr/local/share/ca-certificates/rootca.crt
docker-compose exec nodeB update-ca-certificates

docker-compose restart

echo "------------------------------------"
echo "Waiting for services to restart..."
echo "------------------------------------"
waitForDCService nodeA-backend
waitForDCService nodeB

echo "------------------------------------"
echo "Performing assertions..."
echo "------------------------------------"
# Wait for Nuts Network nodes to build connections
sleep 5
# Assert that node A is connected to B and vice versa using diagnostics. It should look something like this:
# [P2P Network] Connected peers #: 1
#	[P2P Network] Connected peers: (ID=172.19.0.2:43882,NodeID=urn:oid:1.3.6.1.4.1.54851.4:00000002,Addr=172.19.0.2:43882)
RESPONSE=$(curl -s http://localhost:11323/status/diagnostics)
if echo $RESPONSE | grep -q "Connected peers #: 1"; then
  echo "Number of peers of node A is OK"
else
  echo "FAILED: Node A does not report 1 connected peer!" 1>&2
  echo $RESPONSE
  exit 1
fi
RESPONSE=$(curl -s http://localhost:21323/status/diagnostics)
if echo $RESPONSE | grep -q "Connected peers #: 1"; then
  echo "Number of peers of node B is OK"
else
  echo "FAILED: Node B does not report 1 connected peer!" 1>&2
  echo $RESPONSE
  exit 1
fi

echo "------------------------------------"
echo "Stopping Docker containers..."
echo "------------------------------------"
docker-compose stop