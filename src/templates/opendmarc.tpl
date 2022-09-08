# This is a basic configuration that can easily be adapted to suit a standard
# installation. For more advanced options, see opendkim.conf(5) and/or
# /usr/share/doc/opendmarc/examples/opendmarc.conf.sample.

AuthservID OpenDMARC
IgnoreAuthenticatedClients true
RequiredHeaders true
SPFSelfValidate true
PidFile /var/run/opendmarc/opendmarc.pid
PublicSuffixList /usr/share/publicsuffix
RejectFailures false
Socket local:/var/spool/postfix/opendmarc/opendmarc.sock
Syslog true
TrustedAuthservIDs {{ .SubDomain }}.{{ .Domain }}
UMask 0002
UserID opendmarc
