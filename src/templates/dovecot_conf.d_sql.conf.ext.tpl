# Database driver: mysql, pgsql, sqlite
driver = mysql
connect = host={{ .Sql.ServerName }} dbname={{ .Sql.DatabaseName }} user={{ .Sql.Username }} password={{ .Sql.Password }}
default_pass_scheme = SHA512.b64
password_query = SELECT email as user, password FROM virtual_users WHERE email='%u';