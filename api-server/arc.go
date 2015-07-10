package main

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	"gitHub.***REMOVED***/monsoon/arc/transport"
	"gitHub.***REMOVED***/monsoon/arc/transport/fake"
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
	regChan, cancelRegSubscription := tp.SubscribeRegistrations()
	defer cancelRegSubscription()

	msgChan, cancelRepliesSubscription := tp.SubscribeReplies()
	defer cancelRepliesSubscription()

	for {
		select {
		case registry := <-regChan:
			log.Infof("Got registry from %q with data %q", registry.Sender, registry.Payload)

			agent := models.Agent{}
			err := agent.ProcessRegistration(db, registry)
			if err != nil {
				log.Errorf("Error updating fact %q. Got %q", registry, err.Error())
				continue
			}

			// send done signal (for testing)
			ftp, ok := tp.(*fake.FakeClient)
			if ok {
				ftp.DoneSignal()
			}
		case reply := <-msgChan:
			log.Infof("Got reply with id %q and status %q", reply.RequestID, reply.State)

			// update job
			job := models.Job{Request: arc.Request{RequestID: reply.RequestID}, Status: reply.State, UpdatedAt: time.Now()}
			err := job.Update(db)
			if err != nil {
				log.Errorf("Error updating job %q. Got %q", reply.RequestID, err.Error())
				continue
			}

			// add log
			err = models.ProcessLogReply(db, reply)
			if err != nil {
				log.Error(err)
				continue
			}

			// send done signal (for testing)
			ftp, ok := tp.(*fake.FakeClient)
			if ok {
				ftp.DoneSignal()
			}
		}
	}
}

func arcSendRequest(req *arc.Request) error {
	// send request
	log.Infof("Sending request %s", req.RequestID)
	err := tp.Request(req)
	if err != nil {
		return err
	}
	return nil
}
