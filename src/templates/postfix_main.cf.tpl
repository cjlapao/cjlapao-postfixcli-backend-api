# See /usr/share/postfix/main.cf.dist for a commented, more complete version

# Debian specific:  Specifying a file name will cause the first
# line of that file to be used as the name.  The Debian default
# is /etc/mailname.
#myorigin = /etc/mailname
compatibility_level=2
maillog_file = /var/log/mail.log
smtpd_banner = ESMTP {{ .LoadBalancer.Hostname }} (Ubuntu)
biff = no

# appending .domain is the MUA's job.
append_dot_mydomain = no

# Uncomment the next line to generate "delayed mail" warnings
#delay_warning_time = 4h

readme_directory = no

header_checks = regexp:/etc/postfix/header_checks

# TLS parameters
smtpd_tls_cert_file=/etc/mailserver/shared/tls/{{ .SubDomain }}.{{ .Domain }}.cer
smtpd_tls_CAfile=/etc/mailserver/shared/tls/ca.cer
smtpd_tls_key_file=/etc/mailserver/shared/tls/{{ .SubDomain }}.{{ .Domain }}.key
smtpd_use_tls=yes
smtpd_tls_auth_only = yes
smtp_tls_security_level = may
smtpd_tls_security_level = may
smtpd_tls_received_header = yes
smtpd_tls_loglevel = 1
smtpd_sasl_security_options = noanonymous, noplaintext
smtpd_sasl_tls_security_options = noanonymous

smtp_tls_mandatory_ciphers = high
smtp_tls_mandatory_protocols = >=TLSv1.2
smtp_tls_secure_cert_match = nexthop, dot-nexthop

# Authentication
smtpd_sasl_type = dovecot
smtpd_sasl_path = private/auth
smtpd_sasl_auth_enable = yes

# See /usr/share/doc/postfix/TLS_README.gz in the postfix-doc package for
# information on enabling SSL in the smtp client.

# Restrictions
smtpd_helo_restrictions =
        permit_mynetworks,
        permit_sasl_authenticated,
        reject_invalid_helo_hostname,
        reject_non_fqdn_helo_hostname,
        reject_unknown_helo_hostname

smtpd_recipient_restrictions = 
        permit_mynetworks,
        permit_sasl_authenticated,
        reject_unauth_destination,
#        check_policy_service unix:private/policyd-spf,
        check_policy_service inet:127.0.0.1:10023,
        reject_invalid_hostname,
        reject_non_fqdn_hostname,
        reject_non_fqdn_sender,
        reject_non_fqdn_recipient,
        reject_unknown_sender_domain,
        reject_rbl_client sbl.spamhaus.org,
        reject_rbl_client cbl.abuseat.org

smtpd_sender_restrictions =
        reject_non_fqdn_sender,
        reject_sender_login_mismatch,
        reject_unlisted_sender,
        reject_unauth_destination,
        reject_invalid_hostname,
        reject_unknown_sender_domain,
        permit_mynetworks,
        permit_sasl_authenticated

smtpd_relay_restrictions =
        permit_mynetworks,
        permit_sasl_authenticated,
        defer_unauth_destination

# See /usr/share/doc/postfix/TLS_README.gz in the postfix-doc package for
# information on enabling SSL in the smtp client.

myhostname = {{ .SubDomain }}.{{ .Domain }}
alias_maps = hash:/etc/aliases
alias_database = hash:/etc/aliases
mydomain = {{ .Domain }}
myorigin = $mydomain
mydestination = localhost
relayhost =
mynetworks = 127.0.0.0/8 10.0.0.0/8 172.16.0.0/12 192.168.0.0/16 [::ffff:127.0.0.0]/104 [::1]/128
mailbox_size_limit = 0
recipient_delimiter = +
inet_interfaces = all
inet_protocols = all

# Handing off local delivery to Dovecot's LMTP, and telling it where to store mail
virtual_transport = lmtp:unix:private/dovecot-lmtp

# Virtual domains, users, and aliases
virtual_mailbox_domains = mysql:/etc/postfix/mysql-virtual-mailbox-domains.cf
virtual_mailbox_maps = mysql:/etc/postfix/mysql-virtual-mailbox-maps.cf
virtual_alias_maps = mysql:/etc/postfix/mysql-virtual-alias-maps.cf,
        mysql:/etc/postfix/mysql-virtual-email2email.cf

# Even more Restrictions and MTA params
disable_vrfy_command = yes
strict_rfc821_envelopes = yes
#smtpd_etrn_restrictions = reject
#smtpd_reject_unlisted_sender = yes
#smtpd_reject_unlisted_recipient = yes
smtpd_delay_reject = yes
smtpd_helo_required = yes
smtp_always_send_ehlo = yes
#smtpd_hard_error_limit = 1
smtpd_timeout = 30s
smtp_helo_timeout = 15s
smtp_rcpt_timeout = 15s
smtpd_recipient_limit = 40
minimal_backoff_time = 180s
maximal_backoff_time = 3h

# Reply Rejection Codes
invalid_hostname_reject_code = 550
non_fqdn_reject_code = 550
unknown_address_reject_code = 550
unknown_client_reject_code = 550
unknown_hostname_reject_code = 550
unverified_recipient_reject_code = 550
unverified_sender_reject_code = 550

policyd-spf_time_limit = 3600

# Milter configuration
# OpenDKIM
milter_default_action = accept
# Postfix ≥ 2.6 milter_protocol = 6, Postfix ≤ 2.5 milter_protocol = 2
milter_protocol = 6
smtpd_milters = local:opendkim/opendkim.sock,local:opendmarc/opendmarc.sock
non_smtpd_milters = local:opendkim/opendkim.sock