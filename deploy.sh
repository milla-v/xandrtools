#!/bin/bash

GOOS=linux go build -o /tmp/xandrtools xandrtools/cmd/xandrtools

scp /tmp/xandrtools xandrtools:xandrtools.new

ssh xandrtools '
		mv --backup=numbered xandrtools xandrtools.old | head
		mv xandrtools.new xandrtools | head
		chmod +x xandrtools
		sudo service xandrtools restart
'
sleep 5
curl https://xandrtools.com
