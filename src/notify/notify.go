// Package notify contains logic for notifications
// Currently supported notifications are only via SNS topics
// Keeps track of sent notifications
// Sends OK when is no longer alerting
package notify

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cwlogsalert/model"
	"github.com/rs/zerolog/log"
)

//Notification keep track of alerting rules

//ProcessNotifications - notifies subscriber of rule evaluation was succesfful
func ProcessNotifications(notifications chan *model.NotificationItem, tmpl string) error {
	var wg sync.WaitGroup
	log.Debug().Msgf("processing [%d] channel notifications", len(notifications))
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sns := sns.New(sess)

	for n := range notifications {
		log.Debug().Msgf("Notification worker will process [%s]", n.Rule.Name)
		switch n.State {
		case "Alert":
			wg.Add(1)
			go SendAlert(&wg, n, tmpl, sns)
		case "Ok":
			wg.Add(1)
			go SendOk(&wg, n, tmpl, sns)
		}

	}
	wg.Wait()
	return nil
}

//SendAlert - alerting
func SendAlert(wg *sync.WaitGroup, n *model.NotificationItem, t string, client *sns.SNS) error {
	defer wg.Done()
	log.Info().Msgf("Sending Alert for [%s]", n.Rule.Name)
	msg, err := RenderMessageTemplate(n, t)
	if err != nil {
		log.Error().Msgf("Cannot send alert %s", err)
		return err
	}

	log.Debug().Msgf("Message to send is: %s", msg)
	input := &sns.PublishInput{}
	subject := fmt.Sprintf("[ALERT] rule %s is in alert state", n.Rule.Name)

	input.SetSubject(subject)
	input.SetMessage(msg)
	input.SetTopicArn(n.Rule.SnsTopic)

	result, err := client.Publish(input)
	if err != nil {
		log.Error().Msgf("SNS publish error:", err)
		return err
	}
	log.Debug().Msgf("SNS publish result for [%s]: %v+", n.Rule.Name, result)
	return nil
}

//SendOk - ok
func SendOk(wg *sync.WaitGroup, n *model.NotificationItem, t string, client *sns.SNS) error {
	wg.Done()
	log.Info().Msgf("Sending OK for [%s]", n.Rule.Name)
	msg, err := RenderMessageTemplate(n, t)
	if err != nil {
		log.Error().Msgf("Cannot send alert %s", err)
		return err
	}

	log.Debug().Msgf("Message to send is: %s", msg)

	input := &sns.PublishInput{}
	subject := fmt.Sprintf("[OK] rule %s has recovered from alert state", n.Rule.Name)

	input.SetSubject(subject)
	input.SetMessage(msg)
	input.SetTopicArn(n.Rule.SnsTopic)

	result, err := client.Publish(input)
	if err != nil {
		log.Error().Msgf("SNS publish error:", err)
		return err
	}
	log.Debug().Msgf("SNS publish result for [%s]: %v+", n.Rule.Name, result)
	return nil
}

//RenderMessageTemplate - render message template and return text message
func RenderMessageTemplate(notification *model.NotificationItem, t string) (string, error) {
	tmpl, err := template.New("sns").Parse(t)
	if err != nil {
		log.Panic().Err(err)
	}
	var b bytes.Buffer

	err = tmpl.Execute(&b, notification)
	if err != nil {
		log.Panic().Err(err)
	}

	return b.String(), err
}
