#!/bin/bash

scp /tmp/xandrtools xandrtools:xandrtools.new

ssh xandrtools '
		mv --backup=numbered xandrtools xandrtools.old
		mv xandrtools.new xandrtools
		sudo service restart xandrtools
'
