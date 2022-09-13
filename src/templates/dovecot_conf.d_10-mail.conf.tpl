##
## Mailbox locations and namespaces
##

mail_location = maildir:/var/mail/vhosts/%d/%n/

namespace inbox {
  inbox = yes
}

mail_privileged_group = mail
protocol !indexer-worker {
}