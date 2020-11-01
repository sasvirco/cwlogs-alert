# cwlogs-alert

Cwlogs-alert is a tool for alerting of patters of interest from data collected by Cloudwatch logs. 

# Overview

Cwlogs-alert works by executing LogInsights queries against cloudwatch log groups and sending an SNS notification if match is found.
It keeps an in-memory database where it tracks rules that are in Alert mode, and it sends a recovery notification when the rule no longer match.

# Running cwlogs-alert
## Install 
Simply run make in the root folder:
```
make
```
The resulting binary will be stored in bin/cwlogsalert
Current version does not support multiple aws accounts, so you need to run multiple versions of it for every account you want to monitor.
Cwlogs-alert will inherit your aws configuration, so before your run it you need to specify which aws profile it will use:

```
export AWS_PROFILE=my-profile
```

The following invocation can be used to run cwlogs-alert after compiling it:

```
$bin/cwlogsalert -config examples/config.toml
```

# Configuration
Configuration file contains 2 major sections general, and rules. The format is toml.
## General
```
run_interval = "60s" 

# golang text/template
sns_message_template = """
Rule {{.Rule.Name }} changed status to {{.State}}
At least {{ .Rule.NumEvents }} match(es) occured in the last {{ .Rule.Timeframe }}

Query: 
{{ .Rule.Query }}

Result: 
{{  .Result }}
"""
```
* run_interval = a string representing how long to wait before the query cycles (defaults to 60s ). 
* template - a golang text/template that is used to render the message that will be send.

## Rules
```
[rules]
#  A timeframe string is a sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms" or "2h45m". 
#  Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
  [rules.1]
    name = "appAlpha exceptions"
    log_group = "/aws/containerinsights/my-kubernetes-cluster/application"
    num_events = 1
    timeframe = "5m"
    query =  """
        fields @message
        | filter @message  like /Exception/ and `kubernetes.labels.app_kubernetes_io/name` like 'appAlpha' 
        and `kubernetes.namespace_name` == 'production'
      """
    sns_topic = "arn:aws:sns:us-west-2:1234565789:cwlogs-alert"
  [rules.2]
    name = "appBravo exceptions"
    log_group = "/aws/containerinsights/my-kubernetes-cluster/application"
    num_events = 1
    timeframe = "5m"
    query =  """
        fields @message
        | filter @message  like /Exception/ and `kubernetes.labels.app_kubernetes_io/name` like 'appBravo' 
        and `kubernetes.namespace_name` == 'production'
      """
    sns_topic = "arn:aws:sns:us-west-2:1234565789:cwlogs-alert"
  [rules.3]
    name = "appCharlie exceptions"
    log_group = "/aws/containerinsights/my-kubernetes-cluster/application"
    num_events = 1
    timeframe = "5m"
    query =  """
        fields @message
        | filter @message  like /Exception/ and `kubernetes.labels.app_kubernetes_io/name` like 'appCharlie' 
        and `kubernetes.namespace_name` == 'production'
      """
    sns_topic = "arn:aws:sns:us-west-2:1234565789:cwlogs-alert"
  [rule.X]
     ......
'''
* name - Name of the rule
* log_group - the cloudwatchlogs log group that contains the log streams we would like to search in
* num_events - the minimum number of events our query should produce to trigger the alert
* timeframe - how long back in time we want to search (Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h")
* query - cloudwatchlog insights query to execute (use the CloudWatch Logs Insights console to validate the query)
* sns_topic - an aws sns topic to send the notification
