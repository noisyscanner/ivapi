#!/bin/sh

export DB_HOST=$(ip route | awk '/default/ { print $3 }')
