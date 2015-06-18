package main

import (
	log "github.com/Sirupsen/logrus"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/transport"
)

/*
 * Returns a transport connection
 * Remember to disconnect when not any more in use. Use the Disconnect() method
 */
func arcNewConnection(config arc.Config) (transport.Transport, error) {
	// get transport
	tp, err := transport.New(config)
	if err != nil {
		return nil, err
	}

	// conect
	if err := tp.Connect(); err != nil {
		return nil, err
	}

	return tp, nil
}

func arcSubscribeReplies(tp transport.Transport) error {
	regChan, cancelRegSubscription := tp.Subscribe("registry")
	defer cancelRegSubscription()

	msgChan, cancelRepliesSubscription := tp.SubscribeReplies()
	defer cancelRepliesSubscription()

	for {
		select {
		case registry := <-regChan:
			log.Infof("Got registry from %q with data %q", registry.Sender, registry.Payload)
			err := models.UpdateFact(db, registry)
			if err != nil {
				log.Errorf("Error updating fact %q. Got %q", registry, err.Error())
				continue
			}
		case reply := <-msgChan:
			log.Infof("Got reply with id %q and status %q", reply.RequestID, reply.State)

			// update job
			err := models.UpdateJob(db, reply)
			if err != nil {
				log.Errorf("Error updating job %q. Got %q", reply.RequestID, err.Error())
				continue
			}

			// add log
			err = models.SaveLog(db, reply)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

func arcSendRequest(req *arc.Request) error {
	// send request
	log.Infof("Sending request %s", req.RequestID)
	tp.Request(req)

	return nil
}
