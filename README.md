 Arc automation agent
======================
The repository contains the agent running on VMs using the new Monsoon automation service.

Development setup
-----------------
To run a local mosquitto broker run the following docker container

    docker run -it --rm -p 1883:1883 --name mosquitto sapcc/mosquitto

to build and start the agent run:

    go build -o onos-agent .
    ./onos-agent -endpoint=tcp://$(boot2docker ip):1883 #when using bash
    ./onos-agent -endpoint=tcp://(boot2docker ip):1883 #when using fish

To submit messages to the MQTT broker you can reuse the mosquitto image from above:

    docker run -i --rm --link mosquitto:broker mosquitto mosquitto_pub -h broker -t [TOPIC] -s < payload.json

For your convience you can define an alias for the mosquitto_pub command:

    alias mosquitto_pub="docker run -i --rm --link mosquitto:broker mosquitto mosquitto_pub -h broker"

Using this alias you can publish messages just be executing

    mosquitto_pub -t [TOPIC] -s < payload.json