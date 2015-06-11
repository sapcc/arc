package main

import (
	log "github.com/Sirupsen/logrus"	

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"	
	"gitHub.***REMOVED***/monsoon/arc/transport"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

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
	msgChan, cancelSubscription := tp.SubscribeReplies()
	defer cancelSubscription()
	
	for {
		select {
		case reply := <-msgChan:
			log.Infof("Got reply with id %q and status %q", reply.RequestID, reply.State)
			
			affect, err := models.UpdateJob(db, reply)
			if err != nil {
				log.Errorf("Error updating job %q. Got %q", reply.RequestID, err.Error())
			}
			
			log.Infof("%v rows where updated with id %q", affect, reply.RequestID)
		}
	}
}

func arcSendRequest(req *arc.Request) error {	
	// send request
	log.Infof("Sending request %s", req.RequestID)
	tp.Request(req)
	
	return nil
}