# Shadow_Web

sudo service mysql start

sudo service mysql status

sudo mysql_secure_installation
sudo mysql -u root -p

CREATE DATABASE powerdns;
CREATE USER 'pdns'@'localhost' IDENTIFIED BY 'root';
GRANT ALL PRIVILEGES ON powerdns.* TO 'pdns'@'localhost';
FLUSH PRIVILEGES;
EXIT;

