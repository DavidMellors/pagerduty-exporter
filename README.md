PagerDuty Exporter
==================

[![license](https://img.shields.io/github/license/DavidMellors/pagerduty-exporter.svg)](https://github.com/DavidMellors/pagerduty-exporter/blob/master/LICENSE)
[![Docker](https://img.shields.io/badge/docker-wahooga%2Fpagerduty--exporter-blue.svg?longCache=true&style=flat&logo=docker)](https://hub.docker.com/r/wahooga/pagerduty-exporter/)
[![Docker Build Status](https://img.shields.io/docker/build/webdevops/pagerduty-exporter.svg)](https://hub.docker.com/r/wahooga/pagerduty-exporter/)

Prometheus exporter for PagerDuty informations (users, teams, schedules, oncalls, incidents...)
A fork of https://github.com/webdevops/pagerduty-exporter with the addition of log entries metrics

Configuration
-------------

Normally no configuration is needed but can be customized using environment variables.

| Environment variable                   | DefaultValue                | Description                                                               |
|-----------------------------------------|-----------------------------|--------------------------------------------------------------------------|
| `SCRAPE_TIME`                           | `5m`                        | Time (time.Duration) for general informations                            |
| `SCRAPE_TIME_LIVE`                      | `1m`                        | Time (time.Duration) for live metrics (incidents, oncall)                |
| `SERVER_BIND`                           | `:8080`                     | IP/Port binding                                                          |
| `PAGERDUTY_AUTH_TOKEN`                  | none                        | PagerDuty auth token                                                     |
| `PAGERDUTY_SCHEDULE_OVERRIDE_TIMEFRAME` | `48h`                       | PagerDuty schedule override list timeframe                               |
| `PAGERDUTY_SCHEDULE_ENTRY_TIMEFRAME`    | `72h`                       | PagerDuty schedule rendered list timeframe                               |
| `PAGERDUTY_SCHEDULE_ENTRY_TIMEFORMAT`   | `Mon, 02 Jan 15:04 MST`     | PagerDuty schedule entry timeformat (label)                              |
| `PAGERDUTY_INCIDENT_TIMEFORMAT`         | `Mon, 02 Jan 15:04 MST`     | PagerDuty incident entry timeformat (label)                              |
| `PAGERDUTY_DISABLE_TEAMS`               | `false`                     | Boolean (set to 'true' to skip collecting "team" data)                   |
| `PAGERDUTY_TEAM_FILTER`                 | none                        | Comma delimited list of Team IDs to pass to list options when applicable |
| `PAGERDUTY_MAX_CONNECTIONS`             | `4`                         | Maximum numbers of HTTP connections to PagerDuty API                     |
| `PAGERDUTY_LOG_ENTRIES_TIMEZONE`        | `UTC`                | Time zone in which dates in the result will be rendered.|
| `PAGERDUTY_LOG_ENTRIES_SINCE`           | none                 | The start of the date range over which you want to search.|
| `PAGERDUTY_LOG_ENTRIES_UNTIL`           | none                 | The end of the date range over which you want to search. |
| `PAGERDUTY_LOG_ENTRIES_ISOVERVIEW`      | `false`              | If true, will return a subset of log entries that show only the most important changes to the incident. |

Metrics
-------

| Metric                                | Scraper            | Description                                                                           |
|---------------------------------------|--------------------|---------------------------------------------------------------------------------------|
| `pagerduty_stats`                     | Collector          | Collector stats                                                                       |
| `pagerduty_api_counter`               | Collector          | PagerDuty api call counter                                                            |
| `pagerduty_team_info`                 | Team               | Team informations                                                                     |
| `pagerduty_user_info`                 | User               | User informations                                                                     |
| `pagerduty_service_info`              | Service            | Service (per team) informations                                                       |
| `pagerduty_maintenancewindow_info`    | MaintanaceWindows  | Maintenance window informations                                                       |
| `pagerduty_maintenancewindow_status`  | Maintenance window | status (start and endtime)                                         |
| `pagerduty_schedule_info`             | Schedule           | Schedule informations                                                                 |
| `pagerduty_schedule_layer_info`       | Schedule           | Schedule layer informations                                                           |
| `pagerduty_schedule_layer_entry`      | Schedule           | Schedule layer schedule entries                                                       |
| `pagerduty_schedule_layer_coverage`   | Schedule           | Schedule layer schedule coverage                                                      |
| `pagerduty_schedule_final_entry`      | Schedule           | Schedule final (rendered) schedule entries                                            |
| `pagerduty_schedule_final_coverage`   | Schedule           | Schedule final (rendered) schedule coverage                                           |
| `pagerduty_schedule_override`         | Schedule           | Schedule override informations                                                        |
| `pagerduty_schedule_oncall`           | Oncall             | Schedule oncall informations                                                          |
| `pagerduty_incident_info`             | Incident           | Incident informations                                                                 |
| `pagerduty_incident_status`           | Incident           | Incident status informations (acknowledgement, assignment)                            |
| `pagerduty_log_entries`               | Log_entries        | Log entries |

Prometheus queries
------------------

Current oncall person
```
pagerduty_schedule_oncall{scheduleID="$SCHEDULEID",type="startTime"}
* on (userID) group_left(userName) (pagerduty_user_info)
```

Next shift
```
bottomk(1,
  min by (userName, time) (
    pagerduty_schedule_final_entry{scheduleID="$SCHEDULEID",type="startTime"}
    * on (userID) group_left(userName) (pagerduty_user_info) 
  ) - time() > 0
)
```
