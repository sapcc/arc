package main

import (
	log "github.com/Sirupsen/logrus"

	"gitHub.***REMOVED***/monsoon/arc/api-server/models"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
	"gitHub.***REMOVED***/monsoon/arc/transport"
	"gitHub.***REMOVED***/monsoon/arc/transport/fake"
)

/*
 * Returns a transport connection
 * Remember to disconnect when not any more in use. Use the Disconnect() method
 */
func arcNewConnection(config arc_config.Config) (transport.Transport, error) {
	// get transport
	tp, err := transport.New(config, false)
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

	concurrencySafe := true

	for {
		select {
		case registry := <-regChan:
			log.Infof("Got registration from %q with id %q and data %q", registry.Sender, registry.RegistrationID, registry.Payload)

			agent := models.Agent{}
			err := agent.ProcessRegistration(db, registry, tp.IdentityInformation()["identity"], concurrencySafe)
			if err == models.RegistrationExistsError {
				log.Info(models.RegistrationExistsError, " Registration id ", registry.RegistrationID)
			} else if err != nil {
				log.Errorf("Error updating registration %q. Got %q", registry.RegistrationID, err.Error())
			}

			// send done signal (for testing)
			ftp, ok := tp.(*fake.FakeClient)
			if ok {
				ftp.DoneSignal()
			}
		case reply := <-msgChan:
			log.Infof("Got reply with id %q and status %q", reply.RequestID, reply.State)

			// add log
			err := models.ProcessLogReply(db, reply, tp.IdentityInformation()["identity"], concurrencySafe)
			if err == models.ReplyExistsError {
				log.Info(models.ReplyExistsError, " Reply id ", reply.RequestID)
			} else if err != nil {
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
