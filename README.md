# netlify (WIP)

Still a work in progress - theres an WIP branch/pull request for the worker that I'll wrap up tomorrow. After that I need to do a pass to clean up logging/monitoring - and fix up docs. 

# general design/goal

![https://d.pr/free/i/jkYvJM+](https://d.pr/free/i/jkYvJM+)

# components

* etlfunctions - right now its just a single etl function (ChurnTransform), additional etl functions would live here
* churnprofiles - package to work with the churn profile's and source files
* services - just a single worker (Enrichment) right now, that enriches Churn Profiles with Churn Scores - and stores the EnrichedProfiles in our primary datastore

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
    3. Next it reads the json out of storage and turns each ChurnProfile into an EnrichedProfile that contains a ChurnScore field with a mock score. (WIP)
    4. It batches out the EnrichedProfiles and stores them in our primary db (TODO)
    5. If everything was successful it updates the `LogEntry` in cloud datastore recording the file was processed successfully and when it was completed. (DONE) 

# random note

Its been awhile since I wrote a lot of Go. I'm finding it surprisingly difficult to come back and remember how to structure my application in away that makes it easy to unittest it while leveraging Cloud Services.

# Specific Docs

- etlfunctions/README.md
- services/README.md (TODO)