##
## SSL settings
##

ssl = required

ssl_cert = </etc/mailserver/shared/tls/{{ .SubDomain }}.{{ .Domain }}.cer
ssl_key = </etc/mailserver/shared/tls/{{ .SubDomain }}.{{ .Domain }}.key
ssl_client_ca_dir = /etc/ssl/certs
