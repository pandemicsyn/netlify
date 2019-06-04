# netlify (WIP)

Still a work in progress - theres an WIP branch/pull request for the worker that I'll wrap up tomorrow. After that I need to do a pass to clean up logging/monitoring - and fix up docs. 

# general design/goal

![https://d.pr/free/i/jkYvJM+](https://d.pr/free/i/jkYvJM+)

# components

* etlfunctions - right now its just a single etl function (ChurnTransform) that gets invoked on csv upload to a cloud storage bucket
* transform - package to transform csv's to ChurnProfile json files
* pkg/events - the pubsub event the etlfunction emits to notify downstream subscribers a new json file is available
* pkg/profiles - the package for our ChurnProfiles and EnrichedProfiles as well as the supporting sources/stores
* services/enrichment - just a single worker (Enrichment) right now, that enriches ChurnProfiles with Churn Scores - and stores the resulting EnrichedProfiles in our primary db (postgres)

# deployment (wip)

# todo

- the HTTP GET /{customerID} endpoint
- cleanup logging
- cleanup/additional unit tests (probably a bit refactoring)
- deployment docs

# order of operations for this flow/pipeline

1. A new csv batch file arrives in a short term GCS bucket. Which triggers a CloudFunction on completion
2. The cloud function streams in the csv, and streams out json to our permanent GCS bucket. (DONE)
    1. Note: it drops the mock PII fields (gender, senior citizen)
3. The cloud function emits a `FileEvent` to down stream subscribers (DONE)

```
FileEvent
{
    "Bucket": "thebucket",
    "Object": "some/object.json",
    //our schema version for future reference - date based ones work well too
    "Version": 1,
    //right now created is our only status, later we might have purge (gdpr request?)/update/reprocess
    "Status": "created"
    //it might be worth adding a type field like "churnprofile", or "invoices" to future proof things
}
```

4. The Enrichment worker listens for new `FileEvents` and proceeds to consume the json records, enrich them with a mock churn score - and pushes them to permanent storage.
    1. On message receive it checks to make sure Status = created (TODO)
    2. We then use a cloud datastore transaction to create a `LogEntry` so we can track whether a file has been processed. If a record for a given file already exists we skip it, as its either a pubsub duplicate or another worker has processed it. I don't think google cloud pubsub support's exactly once single worker delivery without piping through cloud dataproc. So, i'm using these log entries to make sure that dupe events or workers won't cause duplicate work. They also provide a handy record of what needs to be reprocessed if something breaks. (DONE)
    3. Next it reads the json out of storage and turns each ChurnProfile into an EnrichedProfile that contains a ChurnScore field with a mock score. (DONE)
        1. Note: im using scanner.Scan() since im assuming we know the json lines are a reasonable size. If they end up too large I'd probably switch to Readline.
    4. It Bulk inserts the EnrichedProfiles and stores them in our postgres db (DONE)
        1. Note: im using Postgres Copy In (https://godoc.org/github.com/lib/pq#hdr-Bulk_imports)
    5. If everything was successful it updates the `LogEntry` in cloud datastore recording the file was processed successfully and when it was completed. (DONE)

A few important interfaces:

    1. The ProfileStore interface is used to obtain an io.Reader of wherever the raw .json files live. I used a Google Cloud Storage implementation.
    2. The EnrichedProfileStore interface is the abstraction of where the Enriched Profiles should live. I used a Postgres implementation.
    3. The LogEntry interface is the abstraction of tracking log entries that are created to track process files. I used a Google Cloud Datastore implementation.

# testing

go test right now requires a running instance of the go datastore emulator - if i had a bit more time I could mock that and abstract it as well. For now if you want to run go test something like the below should work:

```
#!/bin/bash
export DATASTORE_DATASET=netlify-242319
export DATASTORE_EMULATOR_HOST=localhost:8081
export DATASTORE_EMULATOR_HOST_PATH=localhost:8081/datastore
export DATASTORE_HOST=http://localhost:8081
export DATASTORE_PROJECT_ID=netlify-242319
gcloud beta emulators datastore start --no-store-on-disk &
go test ./...
```

Alternatively you could also remove services/enrichment/log_entry_test.go

# random note

Its been awhile since I wrote a lot of Go. Got off to a bit of a rough start but its getting there.

# Specific Docs

- etlfunctions/README.md
- services/README.md (TODO)