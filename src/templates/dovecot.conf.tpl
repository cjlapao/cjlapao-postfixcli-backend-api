!include_try /usr/share/dovecot/protocols.d/*.protocol
protocols = imap pop3 lmtp


dict {
  #quota = mysql:/etc/dovecot/dovecot-dict-sql.conf.ext
  #expire = sqlite:/etc/dovecot/dovecot-dict-sql.conf.ext
}

!include conf.d/*.conf

!include_try local.conf