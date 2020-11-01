// This package will contain logic for processing rules results

package process

import (
	"sync"

	"github.com/cwlogsalert/model"
	"github.com/hashicorp/go-memdb"
	"github.com/rs/zerolog/log"
)

//Results -
func Results(r chan *model.ResultItem, n chan *model.NotificationItem, db *memdb.MemDB) error {

	var wg sync.WaitGroup
	for k := range r {
		log.Debug().Msgf("Results worker will process [%s]", k.Rule.Name)
		wg.Add(1)
		go EvalRule(&wg, k, n, db)
	}

	wg.Wait()
	return nil
}

//EvalRule - evaluate rule against results
func EvalRule(wg *sync.WaitGroup, result *model.ResultItem, notifications chan *model.NotificationItem, db *memdb.MemDB) error {

	defer wg.Done()

	log.Debug().Msgf("Evaluate rule against log results for [%s]", result.Rule.Name)

	status, err := GetAlertStatus(result.Rule.Name, db)

	if err != nil {
		log.Error().Msgf("Cannot get Alert state for [%s]: %s", result.Rule.Name, err)
	}

	// alerting
	if int(*result.Result.Statistics.RecordsMatched) >= int(result.Rule.NumEvents) {
		//create record and send notification
		if status == "" {
			_, err := CreateAlertRecord(result.Rule.Name, db)
			if err != nil {
				return err
			}
			log.Info().Msgf("[%s] is in Alert state, sending notification", result.Rule.Name)
			notifications <- &model.NotificationItem{Rule: result.Rule, Result: result.Result, State: "Alert"}
		} else { //rule is already alerting do nothing
			log.Info().Msgf("[%s] is already in Alert state, skipping...", result.Rule.Name)
			return nil
		}
	} else {
		//rule is OK, delete alert record if exist and send recovery notification
		if status == "Alert" {
			log.Info().Msgf("[%s] has recovered and is no longer in alert state", result.Rule.Name)

			err := DeleteAlertRecord(result.Rule.Name, db)
			if err != nil {
				return err
			}
			if result.Rule.NotifyOk == false {
				log.Info().Msgf("[%s] recovery notification is disabled, skipping...", result.Rule.Name)
				return nil
			}
			notifications <- &model.NotificationItem{Rule: result.Rule, Result: result.Result, State: "Ok"}
			return nil
		}

		log.Info().Msgf("[%s] is OK", result.Rule.Name)

	}

	return nil
}

//GetAlertStatus - check if rule is already alerting
func GetAlertStatus(rule string, db *memdb.MemDB) (string, error) {
	log.Debug().Msgf("GetAlertStatus for [%s]", rule)

	// Create read-only transaction
	txn := db.Txn(false)
	defer txn.Abort()

	// Lookup by rule
	raw, err := txn.First("alerts", "id", rule)

	if err != nil {
		log.Error().Msgf("Cannot get alert state for rule [%s] in memdb: %s", rule, err)
		return "", err
	}

	if raw == nil {
		return "", nil
	} else {
		return "Alert", nil
	}

}

//CreateAlertRecord - create record for alerting rule
func CreateAlertRecord(rule string, db *memdb.MemDB) (bool, error) {
	log.Info().Msgf("Creating alert record for [%s]", rule)

	t := db.Txn(true)
	defer t.Abort()

	alert := &model.Alert{RuleName: rule, Alerting: true}

	if err := t.Insert("alerts", alert); err != nil {
		log.Error().Msgf("Cannot insert in memdb: %s", err)
	}
	// Commit the transaction
	t.Commit()
	return true, nil
}

//DeleteAlertRecord - delete record from memdb
func DeleteAlertRecord(rule string, db *memdb.MemDB) error {
	log.Info().Msgf("Deleting alert record for [%s]", rule)

	t := db.Txn(true)
	defer t.Abort()

	alert := &model.Alert{RuleName: rule, Alerting: true}

	if err := t.Delete("alerts", alert); err != nil {
		log.Error().Msgf("Cannot delete from memdb: %s", err)
		return err
	}
	// Commit the transaction
	t.Commit()
	return nil
}
