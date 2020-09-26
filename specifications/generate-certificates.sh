# Create CA certificate and key
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt -subj "/CN=Admission Controller Demo CA"

# Generate the private key for the reverse proxy
openssl genrsa -out key.pem 2048

# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key key.pem -subj "/CN=az-fx-dac-rp.default.svc" |
    openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out cert.pem

# The API server requires the B64 encoded CA certificate to ensure that request is originating from the correct source.
openssl base64 -in ca.crt -out b64ca.crt

# The generated certificate has newline characters which need to be removed.
cat b64ca.crt | tr -d '\n' > b64ca-formatted.crt

# Store the certificates in a secret
kubectl delete secret/dac-rp-cert
kubectl create secret generic dac-rp-cert --from-file=cert.pem --from-file=key.pem
