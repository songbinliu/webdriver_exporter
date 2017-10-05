#!/bin/sh

options=""
if [ "X$HOST" != "X" ] ; then
    # HOST="127.0.0.1"
    options="$options --listenIP=$HOST"
fi

if [ "X$PORT" != "X" ] ; then
    # PORT="9156"
    options="$options --port=$PORT"
fi

#1. start prometheus
echo "[`date`] begin to start prometheus."
prometheus=/bin/prometheus
aoptions="-config.file=/etc/prometheus/prometheus.yml"
aoptions="$aoptions -storage.local.path=/prometheus"
aoptions="$aoptions -web.console.libraries=/etc/prometheus/console_libraries"
aoptions="$aoptions -web.console.templates=/etc/prometheus/consoles"
echo "$prometheus $aoptions"
$prometheus $aoptions &

echo "[`date`] prometheus is running."

#2. start webdriver_exporter
echo "[`date`] begin to run webDriver. " 
webdriver=/usr/bin/webdriver
echo "$webdriver $options"
$webdriver $options 
echo "[`date`] webDriver is stopped. " 

