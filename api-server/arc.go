package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/arc/api-server/models"
	"github.com/sapcc/arc/arc"
	arc_config "github.com/sapcc/arc/config"
	"github.com/sapcc/arc/transport"
	"github.com/sapcc/arc/transport/fake"
	"github.com/sapcc/arc/transport/helpers"
)

var (
	//8 buckets, starting from 0.01 and multiplying by 3 between each
	// 0.01, 0.03, 0.09, 0.27, 0.81, 2.43, 7.29, 21.87
	metricMessageDurationsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "arc_api_mqtt_message_duration_seconds",
			Help:    "Latency of processing MQTT message.",
			Buckets: prometheus.ExponentialBuckets(0.01, 3, 8),
		},
		[]string{"message"},
	)
)

func init() {
	prometheus.MustRegister(metricMessageDurationsHistogram)
}

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

	// connect
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

	concurrencySafe := false
	if tp.IdentityInformation().Transport == helpers.MQTT {
		log.Info("Concurrency safe mode on")
		concurrencySafe = true
	}

	for {
		select {
		case registry := <-regChan:
			log.Debugf("Got registration from %q with id %q and data %q", registry.Sender, registry.RegistrationID, registry.Payload)
			timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
				metricMessageDurationsHistogram.WithLabelValues("registration").Observe(v)
			}))

			// add registration
			err := models.ProcessRegistration(db, registry, tp.IdentityInformation().Identity, concurrencySafe)
			// save processing time
			timer.ObserveDuration()
			// check errors
			if _, ok := err.(models.RegistrationExistsError); ok {
				log.Debug(err.Error(), " Registration id ", registry.RegistrationID)
			} else {
				if err != nil {
					log.Errorf("error updating registration %q. Got %q", registry.RegistrationID, err.Error())
				}
			}

			// send done signal (for testing)
			ftp, ok := tp.(*fake.FakeClient)
			if ok {
				ftp.DoneSignal()
			}
		case reply := <-msgChan:
			log.Infof("Got reply with id %q and status %q", reply.RequestID, reply.State)
			timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
				metricMessageDurationsHistogram.WithLabelValues("reply").Observe(v)
			}))

			// add log
			err := models.ProcessLogReply(db, reply, tp.IdentityInformation().Identity, concurrencySafe)
			// save processing time
			timer.ObserveDuration()
			// check errors
			if _, ok := err.(models.ReplyExistsError); ok {
				log.Debug(err.Error(), " Reply id ", reply.RequestID)
			} else {
				if err != nil {
					log.Error(err)
					continue
				}
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
	return tp.Request(req)
}
