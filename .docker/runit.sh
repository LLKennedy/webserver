#!/bin/sh
exec /sbin/setuser webserver /usr/bin/webserver/webserver.o >> /var/log/webserver.log 2>&1