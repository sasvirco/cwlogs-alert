package cwlog

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/cwlogsalert/config"
	"github.com/cwlogsalert/model"
	"github.com/rs/zerolog/log"
)

//ProcessQueries - process rule cloudwatch insights queries from config
func ProcessQueries(rules map[string]config.Rule, results chan *model.ResultItem) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cwl := cloudwatchlogs.New(sess)

	var wg sync.WaitGroup
	for k := range rules {
		wg.Add(1)
		go Query(&wg, cwl, rules[k], results)
	}
	wg.Wait()

}

//Query Gets cloudwatch insights query results
func Query(wg *sync.WaitGroup, cwl *cloudwatchlogs.CloudWatchLogs, rule config.Rule, results chan *model.ResultItem) (*cloudwatchlogs.GetQueryResultsOutput, error) {

	log.Info().Msgf("Query cloudwatchlog search for [%s]", rule.Name)
	defer wg.Done()

	now := model.MakeTimestamp()
	timeframe := model.ConvertDurationToMilliseconds(rule.Timeframe)
	query := &cloudwatchlogs.StartQueryInput{}

	query.SetStartTime(now - timeframe)
	query.SetEndTime(now)
	query.SetLimit(rule.NumEvents)
	query.SetLogGroupName(rule.LogGroup)
	query.SetQueryString(rule.Query)

	if err := query.Validate(); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	StartQueryOutput, err := cwl.StartQuery(query)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	getQueryResultInput := &cloudwatchlogs.GetQueryResultsInput{}
	getQueryResultInput.SetQueryId(*StartQueryOutput.QueryId)

	for {
		getQueryResultOutput, err := cwl.GetQueryResults(getQueryResultInput)
		if err != nil {
			log.Error().Err(err)
			return nil, err
		}
		time.Sleep(5 * time.Second)
		switch *getQueryResultOutput.Status {
		case "Running":
			if len(getQueryResultOutput.Results) < int(rule.NumEvents) {
				continue
			}
			stopQueryInput := &cloudwatchlogs.StopQueryInput{}
			stopQueryInput.SetQueryId(*StartQueryOutput.QueryId)
			stopResult, err := cwl.StopQuery(stopQueryInput)
			if err != nil {
				log.Error().Msgf("stop query [%s]: error=%s status=%v", rule.Name, err.Error(), stopResult)
				return nil, fmt.Errorf("stop query error=%s status=%v", err.Error(), stopResult)
			}
			return getQueryResultOutput, nil
		case "Scheduled":
			continue
		case "Failed":
			log.Error().Msgf("Query for  [%s] failed", rule.Name)
			return nil, errors.New("job failed")
		case "Cancelled":
			log.Error().Msgf("Query for [%s] canceled", rule.Name)
			return nil, errors.New("job cancelled")
		case "Complete":
			log.Info().Msgf("Query result for [%s]: %d matches", rule.Name, int(*getQueryResultOutput.Statistics.RecordsMatched))
			log.Debug().Msgf("Query result for [%s] : %v+", rule.Name, *getQueryResultOutput)
			results <- &model.ResultItem{Result: getQueryResultOutput, Rule: rule}
			return getQueryResultOutput, nil
		default:
			log.Error().Msgf("Query unknown status for [%s]: %s", rule.Name, *getQueryResultOutput.Status)
			return nil, fmt.Errorf("unknown status: %s", *getQueryResultOutput.Status)
		}
	}

}
