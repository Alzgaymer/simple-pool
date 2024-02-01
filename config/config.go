package config

import "time"

const Address = ":8080"
const Network = "tcp"

var Timeout = time.Now().Add(10 * time.Second)
